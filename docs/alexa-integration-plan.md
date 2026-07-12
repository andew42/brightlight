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
 AWS Lambda (alexa/ in this repo)
   POST https://<home-ddns>:8443/alexa/RunButton
   Authorization: Bearer <shared secret>
   TLS: pinned self-signed cert
        │  (router forwards one port only)
        ▼
 brightlight (new HTTPS listener, separate mux, port 8443)
   servers/alexa.go
     - verify bearer token (constant-time)
     - look up button by name in ui2/build/ui-config/user-buttons.json
       (fallback default-buttons.json)
     - animations.RunAnimations(button.Segments)
     - updateActiveButtonKey(button.Key)   → UIs stay in sync
```

## Part 1 — brightlight (Go)

One new file, `servers/alexa.go`, plus a few lines in `main.go`.

### New endpoints (on a *separate* mux/port, never the existing one)

| Route | Method | Purpose |
|---|---|---|
| `/alexa/RunButton` | POST `{"button":"rainbow"}` | Case-insensitive button name match → run it |
| `/alexa/Buttons` | GET | List button names (debugging / keeping the skill slot list in sync) |

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

* **Separate listener**: `http.ListenAndServeTLS` on its own port (e.g. 8443) with its own
  `http.ServeMux` containing *only* the two Alexa routes. The existing HTTP server (UI,
  config PUT, websockets) remains LAN-only and untouched. The router forwards only 8443.
* **Bearer token**: a long random secret loaded from an environment variable
  (e.g. `BRIGHTLIGHT_ALEXA_TOKEN`) or a file *outside* the served content tree. Compared
  with `crypto/subtle.ConstantTimeCompare`. Missing token at startup → Alexa listener is
  simply not started (feature is opt-in).
* **TLS with a pinned self-signed cert**: generate a long-lived self-signed cert once;
  the Lambda bundles the public cert and uses it as its only trust root. This gives
  encryption *and* server authentication without a domain, Let's Encrypt, or renewal
  plumbing — and the Lambda will refuse to talk to anything that isn't our box.
* **Failures logged** via the existing logrus setup (bad token, unknown button), returning
  401/404 with no detail in the body.

Worst case if the token ever leaked: an attacker can change your light colours. No config
writes, no file access, no UI exposure.

## Part 2 — Alexa skill + Lambda (`alexa/` directory in this repo)

* **Custom skill** (not Smart Home — see alternatives below), invocation name
  `bedroom lights`. One intent `RunButtonIntent` with a custom slot type `ButtonName`
  seeded with the current button names, plus the mandatory built-ins (Stop/Cancel/Help).
  Utterances: "{buttonName}", "run {buttonName}", "set {buttonName}".
* **Lambda in Go** (`aws-lambda-go`), keeping the repo single-language. An Alexa custom
  skill is just JSON in/out over Lambda — no SDK needed for one intent; the handler is
  ~100 lines: parse intent, extract slot, POST to brightlight with the bearer token and
  pinned cert, speak back "OK, rainbow" or a friendly error.
* **Configuration via Lambda environment variables**: `BRIGHTLIGHT_URL`,
  `BRIGHTLIGHT_ALEXA_TOKEN` (Lambda env vars are encrypted at rest). Skill ID verification
  in the handler so only our skill can invoke it.
* Directory contents: `alexa/main.go` (handler), `alexa/skill.json` +
  `alexa/interaction-model.json` (skill definition for the Alexa developer console),
  `alexa/readme.md` (build/deploy steps: `GOOS=linux go build`, zip, upload).

## Voice-to-button-name mapping

Buttons are arbitrarily user-entered, so the mapping is done in two stages:

**Stage 1 — Alexa speech → slot text.** The `ButtonName` custom slot's value list is
training data, not a closed enum: Alexa biases recognition toward listed values but still
returns its best-guess transcription for anything else. Keeping the list in sync:

1. *Manual (v1)*: paste names into the developer console when buttons change;
   `GET /alexa/Buttons` makes this a copy-paste job.
2. *Dynamic entities (v2)*: Lambda fetches `/alexa/Buttons` at session start and injects
   live names via `Dialog.UpdateDynamicEntities` (≤100 values) — recognition then tracks
   button edits automatically; ~30 extra Lambda lines.
3. *SMAPI model rebuilds*: overkill.

**Stage 2 — slot text → button.** In `servers/alexa.go`:

1. Normalise both sides: lowercase, trim, collapse whitespace, strip punctuation
   ("sweet shop" ↔ "Sweet Shop", "off" ↔ "OFF").
2. Exact match on normalised form.
3. Fuzzy fallback for speech quirks: digit/word equivalence ("two" ↔ "2"), plus a small
   prefix/edit-distance pass ("rainbows" → "Rainbow"). ~20 lines, no library.
4. No match → 404; the Lambda speaks back what it heard so misrecognitions are audible.

Caveat: unpronounceable names ("My 3", "BR2-Low") will always be unreliable — the fix is
speakable button names, not code.

## Prerequisites (not code)

* A stable way to reach home: DDNS hostname or static IP, and one router port-forward
  (external 8443 → brightlight box 8443).
* An Amazon developer account for the skill and an AWS account for the Lambda
  (free tier covers this comfortably).

## Alternatives considered

1. **Smart Home skill (scenes)** — would allow the most natural phrasing
   ("Alexa, turn on rainbow") by exposing each button as a scene, but *requires* OAuth2
   account linking, device discovery responses, and the Smart Home API schema. Far more
   moving parts for the same outcome. Custom skill phrasing costs us the word
   "tell/ask": "Alexa, tell bedroom lights rainbow". Recommend custom skill for v1;
   nothing in the brightlight endpoint would need to change if we upgraded later.
2. **AWS IoT Core MQTT (outbound-only)** — brightlight connects *out* to AWS IoT over
   mutual-TLS and subscribes to a command topic; Lambda publishes to it. No open ports,
   no DDNS, survives IP changes. It is the more security-conservative design, but it
   adds an MQTT client dependency, X.509 device certificates, and an always-on
   connection to maintain (reconnect logic). More code and more AWS setup than option 1.
   Worth switching to if opening a port is unacceptable.
3. **Reusing `/RunAnimations/` directly** — rejected: it would mean exposing the existing
   server (including config PUT, which writes files) to the internet, and Alexa knows
   button *names*, not the full segment payload that endpoint requires.

## Open questions to agree before coding

1. Inbound HTTPS + port-forward (this plan) vs outbound MQTT (alternative 2)?
2. Is `ui2/build/ui-config/user-buttons.json` the authoritative button set, or should the
   old UI's `user.json` also be searched?
3. Single room for now? If Titania/other rooms need their own skill later, we can add a
   second invocation name pointing at the same Lambda with a different target URL, or a
   room slot — but v1 assumes one brightlight instance.
4. Lambda language: Go (proposed, single-language repo) vs Node/Python with the ASK SDK.
5. Port number and where the token/cert files should live on the brightlight box.

## Estimated size

* `servers/alexa.go`: ~120 lines (auth middleware, file load, name match, two handlers)
* `main.go`: ~10 lines (start TLS listener if token configured)
* `alexa/main.go`: ~100 lines
* Skill JSON + readmes: no runtime code
