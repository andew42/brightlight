# Alexa voice control

Voice phrasing: **"Alexa, ask bedroom lights for rainbow"** (also "... to run
rainbow", or just "Alexa, ask bedroom lights" then say the setting name).
The spoken name is matched case-insensitively against the button names in
`backend/ui-config/user-buttons.json` (fallback `default-buttons.json`),
with small speech-to-text errors forgiven and spoken numbers matched to
digits ("two tone" finds "2 Tone").

```
Alexa -> custom skill -> HTTPS to bedroom-lights.elms.andrewandlaura.com
      -> Caddy (TLS) -> brightlight :8443 /alexa/skill
```

There is no AWS account and no Lambda: the skill calls brightlight directly.
Caddy already terminates TLS for the site, so the skill gets a certificate
from a trusted CA without any extra plumbing. See
`docs/alexa-integration-plan.md` for the design discussion.

## What protects the endpoint

`/alexa/skill` is the only route on the Alexa listener and the only part of
brightlight reachable from the internet. Two checks gate every request and
both are needed:

* **The Amazon signature** proves the request came from the Alexa service.
  Requests carry `SignatureCertChainUrl` and `Signature-256`; brightlight
  validates the URL is Amazon's, chain-verifies the certificate, checks it
  was issued for `echo-api.amazon.com`, verifies the body signature and
  rejects timestamps more than 150 seconds old.
* **The skill id** proves the request came from *our* skill. A signature
  alone is not enough — anyone can create their own skill and point it at
  this URL, and Amazon will sign those requests too.

Because the skill id is what binds requests to your skill, the listener does
not start at all unless `BRIGHTLIGHT_ALEXA_SKILL_ID` is set.

## 1. Create the skill

In the [Alexa developer console](https://developer.amazon.com/alexa/console/ask)
create a **Custom** skill named e.g. "Bedroom Lights", then in the JSON
editor paste `interaction-model.json` and build the model.

Copy the skill id (`amzn1.ask.skill.…`) from the skill's page — it is needed
in step 2, and can be pasted into the endpoint config in step 4.

The skill stays in development mode, which is fine for personal use on
devices signed in to the same Amazon account. A development-stage skill is
not in the skill store and cannot be enabled by anyone else.

## 2. Configure brightlight

Set the skill id in a systemd drop-in on the Pi (alongside the site drop-in
written by the installer, so it survives upgrades):

```sh
sudo mkdir -p /etc/systemd/system/brightlight.service.d
printf '[Service]\nEnvironment=BRIGHTLIGHT_ALEXA_SKILL_ID=amzn1.ask.skill.your-skill-id\n' \
  | sudo tee /etc/systemd/system/brightlight.service.d/alexa.conf
sudo systemctl daemon-reload && sudo systemctl restart brightlight
```

| Variable | Value |
|---|---|
| `BRIGHTLIGHT_ALEXA_SKILL_ID` | skill id from step 1; unset disables the endpoint |
| `BRIGHTLIGHT_ALEXA_PORT` | optional, default 8443 |

Check the log says `serving Alexa endpoint`:

```sh
journalctl -u brightlight -n 20
```

## 3. Point Caddy at it

Add a site block to the Caddyfile. Only `/alexa/skill` is proxied — the
listener has no other route, but being explicit means a future change to the
LAN server can't accidentally become public:

```caddyfile
bedroom-lights.elms.andrewandlaura.com {
	handle /alexa/skill {
		reverse_proxy 192.168.1.50:8443
	}
	handle {
		abort
	}
}
```

Replace `192.168.1.50` with the Pi's LAN address. Add a DNS record for
`bedroom-lights.elms.andrewandlaura.com` pointing at the same public address
as the rest of the site; no new router port-forward is needed because this
uses the 443 Caddy already listens on.

The signature covers the exact bytes Amazon sent, so do not add anything to
this route that rewrites the request body. A plain `reverse_proxy` passes the
body through untouched.

Reload Caddy and check the route answers — an unsigned request should be
refused, which confirms the path is reachable *and* that the checks work:

```sh
curl -i -X POST https://bedroom-lights.elms.andrewandlaura.com/alexa/skill -d '{}'
```

Expect `401 unauthorized`. Anything else (404, 502, a TLS error) is a Caddy
or DNS problem rather than brightlight.

## 4. Set the skill endpoint

Back in the developer console, under **Endpoint** choose **HTTPS** and set:

| Field | Value |
|---|---|
| Default Region | `https://bedroom-lights.elms.andrewandlaura.com/alexa/skill` |
| Certificate type | *My development endpoint has a certificate from a trusted certificate authority* |

Caddy's Let's Encrypt certificate satisfies the trusted-CA option. Save and
build the model, then test from the console's **Test** tab or a real Echo.

## 5. Keeping button names recognisable

The `BUTTON_NAME` slot values are recognition hints, not a closed list.
When you add or rename buttons, fetch the current names from the LAN and
paste them into the slot values in the developer console for best
recognition. Unlisted names usually still work if they are ordinary English
words.

```sh
curl http://<pi-address>:8080/api/AlexaButtons
```

This listing is deliberately on the LAN server rather than the Alexa
listener, so it is not exposed to the internet.

## Troubleshooting

Every rejection is logged with a reason:

```sh
journalctl -u brightlight -f | grep -i alexa
```

| Log | Cause |
|---|---|
| `failed signature check` | request did not come from Alexa, or Caddy altered the body |
| `from unexpected skill` | `BRIGHTLIGHT_ALEXA_SKILL_ID` does not match the calling skill |
| `timestamp rejected` | the Pi's clock is wrong, or the request was replayed |
| `no matching button` | the spoken name did not match a button (the reply says what was heard) |
