# Alexa Light Control — Integration Plan

Goal: say something like *"Alexa, tell bedroom lights rainbow"* and have brightlight run
the button of that name, exactly as if it had been pressed in the UI.

Guiding principles: **minimal new code**, **reuse the existing button/animation pipeline**,
**never expose the existing UI/config endpoints to the internet**.

## Architecture overview

```
 "Alexa, tell bedroom lights rainbow"
        │
        ▼
 Alexa Custom Skill ("bedroom lights")
   RunButtonIntent, slot {buttonName}
        │
        ▼
 POST https://bedroom-lights.elms.andrewandlaura.com/alexa/skill
   Alexa request envelope, signed by Amazon
        │
        ▼
 Caddy (already internet-facing for the site, terminates TLS)
   reverse_proxy → pi:8443, /alexa/skill only
        │
        ▼
 brightlight (own listener, separate mux, port 8443)
   servers/alexasignature.go
     - validate SignatureCertChainUrl, chain-verify the cert
     - check it is issued for echo-api.amazon.com
     - verify Signature-256 over the raw body
   servers/alexa.go
     - check the application id is our skill
     - check the timestamp is within 150s
     - look up button by name in backend/ui-config/user-buttons.json
       (fallback default-buttons.json)
     - animations.RunAnimations(button.Segments)
     - updateActiveButtonKey(button.Key)   → UIs stay in sync
     - speak the matched name back
```

## Part 1 — brightlight (Go)

`servers/alexa.go` and `servers/alexasignature.go`, plus a few lines in
`main.go`.

### New endpoints

| Route | Method | Purpose |
|---|---|---|
| `/alexa/skill` (port 8443) | POST | The skill endpoint: verified Alexa envelope in, speech out |
| `/api/AlexaButtons` (port 8080) | GET | List button names for keeping the skill slot list in sync |

The skill route is alone on its own mux and port so a reverse-proxy
misconfiguration cannot expose anything else. The button listing is on the
LAN server precisely because it does not need to be public.

Implementation notes:

* The ui2 button file (`user-buttons.json` / `default-buttons.json`) already has the exact
  shape of `animations.Button` (`key`, `name`, `segments`), so the handler is: read file,
  unmarshal `{ "buttons": []animations.Button }`, find name (case-insensitive, trimmed),
  call `animations.RunAnimations(...)` and `updateActiveButtonKey(...)` — the same two
  calls `RunAnimationsHandler` makes today. No new animation code at all.
* Reading the file per request (like the config handler does) keeps it simple and always
  picks up button edits — no cache invalidation logic needed. It's one small file read per
  voice command; performance is irrelevant here.

### Security

* **Separate listener**: `http.ListenAndServe` on its own port (8443) with its own
  `http.ServeMux` containing *only* the skill route. The existing HTTP server (UI,
  config PUT, SSE streams) remains LAN-only and untouched. Caddy proxies only the one
  path from the internet.
* **TLS at Caddy**: the site's existing internet-facing Caddy terminates TLS with a
  Let's Encrypt certificate on `bedroom-lights.elms.andrewandlaura.com` and forwards
  plain HTTP over the LAN. Amazon requires a trusted-CA certificate for an HTTPS skill
  endpoint, which this satisfies with no extra plumbing and nothing to renew by hand.
* **Amazon request signature**: every request carries `SignatureCertChainUrl` and
  `Signature-256`. We validate the URL points at Amazon's own certificate location
  (scheme, host, port, and path after `path.Clean` so traversal cannot disguise it),
  chain-verify the fetched certificate against the system roots, require the
  `echo-api.amazon.com` SAN, then check the signature over the *raw* body — so the body
  is read and verified before anything is unmarshalled from it. Chains are cached but
  re-verified per request so expiry still bites.
* **Skill id check**: a valid signature only proves the request came from Amazon. Anyone
  can create a skill and point it at this URL, and Amazon will sign those requests too,
  so `BRIGHTLIGHT_ALEXA_SKILL_ID` is what actually binds a request to our skill. It is
  mandatory: unset means the listener does not start (fail closed), unlike the old
  optional check in the Lambda where the bearer token was the real gate.
* **Replay window**: requests with a timestamp more than 150 seconds from now are
  refused, per Amazon's requirement.
* **Failures logged** via slog with the reason, returning 401 with no detail in the body.

Worst case if all of this were beaten: an attacker can change your light colours. No
config writes, no file access, no UI exposure.

## Part 2 — Alexa skill (`alexa/` directory in this repo)

* **Custom skill** (not Smart Home — see alternatives below), invocation name
  `bedroom lights`. One intent `RunButtonIntent` with a custom slot type `ButtonName`
  seeded with the current button names, plus the mandatory built-ins (Stop/Cancel/Help).
  Utterances: "{buttonName}", "run {buttonName}", "set {buttonName}".
* **HTTPS endpoint, no Lambda**: the skill's endpoint is set to the brightlight URL
  directly, so there is no AWS account, no deploy step and no `aws-lambda-go` dependency
  — the backend module stays stdlib-only. The envelope handling the Lambda used to do
  (LaunchRequest / RunButtonIntent / Help / Stop / SessionEnded → speech) moves into
  `servers/alexa.go`.
* **Configuration**: `BRIGHTLIGHT_ALEXA_SKILL_ID` (and optionally
  `BRIGHTLIGHT_ALEXA_PORT`) in a systemd drop-in, so it survives installer upgrades the
  same way `BRIGHTLIGHT_SITE` does.
* Directory contents: `alexa/interaction-model.json` (skill definition for the Alexa
  developer console) and `alexa/readme.md` (skill, systemd and Caddy setup steps).

## Voice-to-button-name mapping

Buttons are arbitrarily user-entered, so the mapping is done in two stages:

**Stage 1 — Alexa speech → slot text.** The `ButtonName` custom slot's value list is
training data, not a closed enum: Alexa biases recognition toward listed values but still
returns its best-guess transcription for anything else. Keeping the list in sync:

1. *Manual (v1)*: paste names into the developer console when buttons change;
   `GET /api/AlexaButtons` makes this a copy-paste job.
2. *Dynamic entities (v2)*: answer the LaunchRequest with a
   `Dialog.UpdateDynamicEntities` directive carrying the live button names (≤100 values)
   — recognition would then track button edits automatically; ~30 extra lines in
   `servers/alexa.go`.
3. *SMAPI model rebuilds*: overkill.

**Stage 2 — slot text → button.** In `servers/alexa.go`:

1. Normalise both sides: lowercase, trim, collapse whitespace, strip punctuation
   ("sweet shop" ↔ "Sweet Shop", "off" ↔ "OFF").
2. Exact match on normalised form.
3. Fuzzy fallback for speech quirks: digit/word equivalence ("two" ↔ "2"), plus a small
   prefix/edit-distance pass ("rainbows" → "Rainbow"). ~20 lines, no library.
4. No match → the reply speaks back what it heard so misrecognitions are audible.

Caveat: unpronounceable names ("My 3", "BR2-Low") will always be unreliable — the fix is
speakable button names, not code.

## Prerequisites (not code)

* A DNS record for `bedroom-lights.elms.andrewandlaura.com` pointing at the same public
  address as the rest of the site, and a Caddy site block proxying `/alexa/skill` to the
  Pi. No new router port-forward — this rides on the 443 Caddy already serves.
* An Amazon developer account for the skill. No AWS account.

## Alternatives considered

1. **Smart Home skill (scenes)** — would allow the most natural phrasing
   ("Alexa, turn on rainbow") by exposing each button as a scene, but *requires* OAuth2
   account linking, device discovery responses, and the Smart Home API schema. Far more
   moving parts for the same outcome. Custom skill phrasing costs us the word
   "tell/ask": "Alexa, tell bedroom lights rainbow". Recommend custom skill for v1;
   nothing in the brightlight endpoint would need to change if we upgraded later.
2. **AWS IoT Core MQTT (outbound-only)** — brightlight connects *out* to AWS IoT over
   mutual-TLS and subscribes to a command topic; a Lambda publishes to it. No open ports,
   no DDNS, survives IP changes. It is the more security-conservative design, but it
   adds an MQTT client dependency, X.509 device certificates, and an always-on
   connection to maintain (reconnect logic). More code and more AWS setup than option 1.
   Worth switching to if opening a port is unacceptable.
3. **Reusing `/RunAnimations/` directly** — rejected: it would mean exposing the existing
   server (including config PUT, which writes files) to the internet, and Alexa knows
   button *names*, not the full segment payload that endpoint requires.
4. **Lambda in front of the endpoint** (the original v1, since removed) — a Go Lambda
   held the bearer token and the pinned self-signed certificate, so brightlight only had
   to compare a token. It existed because there was no trusted-CA certificate for the
   box; once Caddy provided one it was a second deployment artefact, an AWS account and
   an external Go dependency buying nothing. Its cost was moved, not avoided: signature
   verification now lives in brightlight instead.

## Decisions (agreed 2026-07-12)

Inbound HTTPS with port-forward; ui2 `user-buttons.json` (fallback
`default-buttons.json`) is the only button source; simple Go Lambda with manual slot
sync (no dynamic entities); "ask" as the connecting word ("Alexa, ask bedroom lights
for rainbow"); fuzzy matching includes digit/word equivalence ("2" = "two").
Implemented in `servers/alexa.go` and `alexa/`; deployment steps in `alexa/readme.md`.

## Revision (agreed 2026-08-01) — Lambda removed

The site's Caddy reverse proxy already terminates TLS with a trusted certificate, which
is the one thing an HTTPS skill endpoint needs and the only reason the Lambda and the
pinned self-signed certificate existed. So:

* The skill now calls `https://bedroom-lights.elms.andrewandlaura.com/alexa/skill`
  directly. `alexa/lambda/` is deleted along with the `aws-lambda-go` dependency — the
  repo has no external Go dependencies again.
* `BRIGHTLIGHT_ALEXA_TOKEN`, `BRIGHTLIGHT_ALEXA_CERT` and `BRIGHTLIGHT_ALEXA_KEY` are
  replaced by a single mandatory `BRIGHTLIGHT_ALEXA_SKILL_ID`; the listener is plain HTTP
  behind Caddy.
* `/alexa/RunButton` and `/alexa/Buttons` are gone. The public listener serves only
  `/alexa/skill`; button names moved to `/api/AlexaButtons` on the LAN server.
* Authentication is Amazon's request signature *plus* the skill id. The signature alone
  is not sufficient — Amazon signs requests from anyone's skill, so the id is what makes
  a request ours. Trading a shared secret for signature verification is the real cost of
  this change, which is why it has its own file and test suite.
* The router port-forward of 8443 is no longer needed.

## Open questions

1. Single room for now? If Titania/other rooms need their own skill later, we can add a
   second invocation name and skill id, or a room slot — but v1 assumes one brightlight
   instance.
2. Should the installer write the `BRIGHTLIGHT_ALEXA_SKILL_ID` drop-in (as it does for
   `--site`) rather than leaving it a manual step in `alexa/readme.md`?

## Size as built

* `servers/alexa.go`: ~360 lines (file load, name match, envelope handling, handler)
* `servers/alexasignature.go`: ~225 lines (certificate and signature verification)
* `main.go`: ~5 lines (start listener if the skill id is configured)
* Tests: ~530 lines across `alexa_test.go` and `alexasignature_test.go`
* Skill JSON + readmes: no runtime code
