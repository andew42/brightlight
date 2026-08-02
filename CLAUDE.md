# Brightlight — Project Context for Claude

## What this project is

A domestic LED lighting controller. A Go web server runs on a Raspberry Pi and
drives WS2811 pixel-addressable LED strips via a Teensy 3.x microcontroller
connected over USB serial. A React web app (served by the Go server) lets the
user choose and configure lighting animations. An optional Alexa custom skill
calling the server directly provides voice control of the configured buttons.

## Repository layout

```
brightlight/
├── .github/workflows/
│   └── build.yml      CI: backend + frontend builds, Pi bundle, rolling release
├── frontend/          React web app (Vite + React 19)
│   ├── src/           JSX source — all files use .jsx extension
│   ├── public/        Static assets (icons, manifest, segment-icons)
│   ├── build/         Vite build output — served by Go server (gitignored)
│   ├── package.json   npm deps (react, react-dom, react-router-dom, vite)
│   ├── vite.config.js Vite 8 config: outDir=build, /api dev proxy
│   └── ui-asset-source/ Pixelmator source for apple-touch-icon
├── backend/           Go server
│   ├── main.go        Entry point; serves frontend/build/ as static files
│   ├── go.mod         module github.com/andew42/brightlight, Go 1.26
│   ├── animations/    LED animation effects (rainbow, plasma, twinkle, …)
│   ├── config/        Network config, presets, static data
│   ├── controller/    Teensy USB serial driver + relay driver
│   ├── framebuffer/   In-memory LED frame buffer
│   ├── segment/       Named/physical/combined LED segment abstractions
│   ├── servers/       HTTP handlers (animations, config, SSE streams, Alexa)
│   ├── stats/         Performance statistics
│   └── ui-config/     Button/segment config JSON (served via /api/ui-config)
├── alexa/             Alexa voice control (see alexa/readme.md)
│   └── interaction-model.json  Custom skill interaction model
├── firmware/          Arduino/Teensy C firmware (OctoWS2811)
├── packaging/         Pi install: systemd unit + install.sh (bundled by CI)
├── deploy/            Staging area for Pi deployment artefacts (generated)
│   └── README.md      How to build and deploy
├── docs/              Design notes (Alexa integration plan)
├── full-build.sh      Linux/macOS build script: frontend + backend → deploy/
├── full-build.bat     Windows equivalent
├── readme.md          Original project README with Pi setup instructions
└── CLAUDE.md          This file
```

## Tech stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 19, React Router 7, Vite 8 (semantic-ui removed) |
| Backend | Go 1.26, stdlib only (log/slog, net/http with SSE) — zero external deps |
| Hardware | Raspberry Pi 2B (Linux ARMv7), Teensy 3.x, WS2811 LED strips |
| Voice | Alexa custom skill → HTTPS endpoint on the Go server (no Lambda, no AWS) |
| Build | Vite (frontend), cross-compiled Go GOOS=linux GOARCH=arm GOARM=7 (backend) |
| CI/CD | GitHub Actions (`.github/workflows/build.yml`) → rolling `latest` GitHub release |

## Runtime configuration

`BRIGHTLIGHT_SITE` env var selects the hardware site layout at runtime:
`titania` (default) or `bedroom`. It switches the frame-buffer strip layout
and named segments (`config.Site` / `config.Titania`), and makes the config
server serve `static-data-<site>.json` / `default-buttons-<site>.json` when
the UI requests the generic names. On a Pi it is set via a systemd drop-in
written by `packaging/install.sh --site <site>`.

The Go binary uses `BRIGHTLIGHT` env var as the base path for static files.
The Pi installer sets `BRIGHTLIGHT=/opt/brightlight` (in the systemd unit),
under which it looks for:
- Frontend: `/opt/brightlight/frontend/build/` (index.html, assets/)
- Button config: `/opt/brightlight/backend/ui-config/`

When `BRIGHTLIGHT` is unset (development) static content is not served — use
the Vite dev server, which proxies `/api` to `http://localhost:8080` — and
button config is read from `ui-config/` relative to the working directory.

The Alexa endpoint is opt-in (see `alexa/readme.md`): it starts only when
`BRIGHTLIGHT_ALEXA_SKILL_ID` is set, and listens on its own plain-HTTP port
(`BRIGHTLIGHT_ALEXA_PORT`, default 8443) with TLS terminated by the Caddy
reverse proxy in front of it. Fail-closed is deliberate: the skill id is what
binds a request to our skill, so without it there is nothing to check.

## How to build

```bash
# Frontend dev server (proxies /api to backend at localhost:8080)
cd frontend && npm start

# Backend dev server
cd backend && go run .

# Build for deployment — outputs to deploy/ (same steps/layout as CI)
full-build.sh    # Linux/macOS
full-build.bat   # Windows
```

CI (`.github/workflows/build.yml`) runs on every push/PR to `master`:
it vets and cross-compiles the backend for Raspberry Pi 2B
(GOARM=7), builds the frontend with Vite (`npm ci && npm run build`), then on
push assembles `brightlight-pi.tar.gz` (binary, frontend build, ui-config,
systemd unit, installer) and republishes it as the rolling `latest` GitHub
release. A Pi installs/upgrades with one command — see `deploy/README.md`.

## HTTP API (backend → frontend contract)

All API routes are under `/api`. The push channels use Server-Sent Events
(SSE, `EventSource`) — the WebSocket implementation was replaced.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/ui-config/*.json` | GET/PUT | Button/segment config JSON files (PUT allowed only for user-buttons.json) |
| `/api/RunAnimations/` | POST | Run animation (JSON button payload) |
| `/api/StripLength/` | POST | Set strip length on room lights |
| `/api/ButtonState` | SSE | Push active button key + config version |
| `/api/FrameBuffer` | SSE | Push raw LED frame data (virtual display) |
| `/api/Stats` | SSE | Push performance stats |
| `/api/option/` | POST | Set server options |

`/api/AlexaButtons` (GET) lists button names for pasting into the skill's slot
values; it is on the LAN server so it is never exposed to the internet.

The skill endpoint itself lives on a separate opt-in listener (default port
8443 — see `backend/servers/alexa.go`), which serves exactly one route:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/alexa/skill` | POST | Alexa request envelope in, speech response out |

## Key gotchas and past decisions

### React 19 upgrade (April 2026)
- `ReactDOM.render()` removed → use `createRoot()` in `frontend/src/index.jsx`
- React Router v4 → v7: `Switch`→`Routes`, `render`/`component` props → `element`
  prop, `history.push()` → `useNavigate()` hook. Class components can't use
  hooks directly so `ButtonPadRoute` and `ButtonEditorRoute` wrapper components
  inject `navigate`/`location` as a history-compatible object.

### Semantic UI removed
`semantic-ui-react` / `semantic-ui-css` have been removed entirely along with
the npm `overrides` and the lightningcss `errorRecovery` workaround they
required. UI components are now hand-rolled (see `frontend/src/dialog/`).

### Vite 8 / Rolldown upgrade (April 2026)
- **All source files use `.jsx` extension** — Vite 8 uses Rolldown which will
  not parse JSX in `.js` files. Do not rename them back to `.js`.
- `esbuild`/`optimizeDeps.esbuildOptions` Vite 6 workarounds are gone; the
  oxc/rolldown pipeline handles everything natively now.

### Frontend dev proxy
`vite.config.js` proxies `/api` to `http://localhost:8080` during `npm start`
(run the Go backend locally). Change the target if you develop against the
device instead.

### WebSockets → SSE
Push channels were converted from WebSockets to Server-Sent Events; the
frontend wrapper (`frontend/src/server-proxy/webSocket.jsx`) keeps its old
`OpenWebSocket` name but opens an `EventSource` and reconnects when the page
becomes visible again (fixes stale UI after switching apps on a phone). This
removed the last external Go dependency (`golang.org/x/net`).

### Old UI removed
The original `ui/` directory (plain HTML/JS with Ractive.js) has been deleted.
The `/ui/` HTTP route and `/config/` handler have been removed from `main.go`.

### Repo structure (April 2026)
Reorganised from flat layout into `frontend/`, `backend/`, `deploy/` to allow
opening each part as a separate IntelliJ project. Previously `ui2/` was the
frontend and all Go files lived at the repo root.

### Alexa voice control (July 2026, revised August 2026)
Voice control goes Alexa → custom skill → HTTPS straight to
`bedroom-lights.elms.andrewandlaura.com/alexa/skill`, where the site's Caddy
reverse proxy terminates TLS and forwards that one path to the Pi.
The original design put an AWS Lambda in front holding a bearer token and a
pinned self-signed cert; that existed only because the box had no trusted-CA
certificate, so Caddy made it redundant and `alexa/lambda/` was deleted along
with `aws-lambda-go` — the repo has no external Go dependencies again.

Two checks replace the bearer token and **both** are required. Amazon's request
signature (`backend/servers/alexasignature.go`: cert chain URL validation,
chain verify, `echo-api.amazon.com` SAN, `Signature-256` over the raw body,
150s replay window) proves a request came from Amazon; `BRIGHTLIGHT_ALEXA_SKILL_ID`
proves it came from *our* skill, since anyone can point their own skill at the
URL and Amazon signs those too. Verify before unmarshalling — the signature
covers the exact bytes sent — and don't let the proxy rewrite the body.

Spoken names are fuzzy-matched against button names in `user-buttons.json`
(fallback `default-buttons.json`). Setup steps in `alexa/readme.md`; design
discussion and the removal rationale in `docs/alexa-integration-plan.md`.

### CI pipeline and one-command Pi install (July 2026)
`.github/workflows/build.yml` builds backend (ARMv7, GOARM=7 — the target is a
Pi 2B; the old scripts used GOARM=5) and frontend, then packages a Pi bundle
and republishes it as the rolling `latest` GitHub release on every push to
`master`. A Pi installs or upgrades with
`curl -fsSL .../releases/latest/download/install.sh | sudo bash`, which
installs to `/opt/brightlight` and sets up a systemd service (`packaging/`),
replacing the old manual scp + `/etc/rc.local` flow. `full-build.sh`/`.bat`
mirror the CI build steps for a local build staged into `deploy/`.
