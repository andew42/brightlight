# Alexa voice control

Voice phrasing: **"Alexa, ask bedroom lights for rainbow"** (also "... to run
rainbow", or just "Alexa, ask bedroom lights" then say the setting name).
The spoken name is matched case-insensitively against the button names in
`backend/ui-config/user-buttons.json` (fallback `default-buttons.json`),
with small speech-to-text errors forgiven and spoken numbers matched to
digits ("two tone" finds "2 Tone").

```
Alexa -> custom skill -> AWS Lambda (this directory)
      -> HTTPS + bearer token + pinned cert -> brightlight /alexa/RunButton
```

See `docs/alexa-integration-plan.md` for the design discussion.

## 1. Generate a shared token and server certificate

On the brightlight box:

```sh
# 64 char random token
openssl rand -hex 32 > alexa-token.txt

# 10 year self-signed cert; CN/SAN must be the public hostname the Lambda
# will connect to (your DDNS name)
openssl req -x509 -newkey rsa:2048 -nodes -days 3650 \
  -keyout alexa-key.pem -out alexa-cert.pem \
  -subj "/CN=myhouse.example.com" \
  -addext "subjectAltName=DNS:myhouse.example.com"
```

Keep `alexa-key.pem` and the token private and outside the web content tree.

## 2. Configure brightlight

Set the environment variables before starting brightlight (the Alexa
listener is disabled when `BRIGHTLIGHT_ALEXA_TOKEN` is unset):

| Variable | Value |
|---|---|
| `BRIGHTLIGHT_ALEXA_TOKEN` | contents of alexa-token.txt |
| `BRIGHTLIGHT_ALEXA_CERT` | path to alexa-cert.pem |
| `BRIGHTLIGHT_ALEXA_KEY` | path to alexa-key.pem |
| `BRIGHTLIGHT_ALEXA_PORT` | optional, default 8443 |

Forward exactly one port on the router: external 8443 to the brightlight
box port 8443. Do not forward 8080.

Test from anywhere:

```sh
curl --cacert alexa-cert.pem \
  -H "Authorization: Bearer $(cat alexa-token.txt)" \
  https://myhouse.example.com:8443/alexa/Buttons

curl --cacert alexa-cert.pem \
  -H "Authorization: Bearer $(cat alexa-token.txt)" \
  -d '{"button":"rainbow"}' \
  https://myhouse.example.com:8443/alexa/RunButton
```

## 3. Build and deploy the Lambda

```sh
cd alexa/lambda
GOOS=linux GOARCH=arm64 go build -o bootstrap .
cp /path/to/alexa-cert.pem server-cert.pem
zip lambda.zip bootstrap server-cert.pem
```

In the AWS console create a Lambda: runtime **provided.al2023**,
architecture arm64, upload `lambda.zip`. Set environment variables:

| Variable | Value |
|---|---|
| `BRIGHTLIGHT_URL` | `https://myhouse.example.com:8443` |
| `BRIGHTLIGHT_ALEXA_TOKEN` | same token as the server |
| `ALEXA_SKILL_ID` | skill id from step 4 (add after creating the skill) |

Add an **Alexa Skills Kit** trigger with the skill id.

## 4. Create the skill

In the [Alexa developer console](https://developer.amazon.com/alexa/console/ask)
create a **Custom** skill named e.g. "Bedroom Lights", then in the JSON
editor paste `interaction-model.json` and build the model. Set the
endpoint to the Lambda's ARN. Copy the skill id into the Lambda
environment and trigger.

The skill stays in development mode, which is fine for personal use on
devices signed in to the same Amazon account.

## 5. Keeping button names recognisable

The `BUTTON_NAME` slot values are recognition hints, not a closed list.
When you add or rename buttons, fetch the current names with
`GET /alexa/Buttons` and paste them into the slot values in the developer
console for best recognition. Unlisted names usually still work if they
are ordinary English words.
