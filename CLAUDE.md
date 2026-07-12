# Brightlight — Project Context for Claude

## What this project is

A domestic LED lighting controller. A Go web server runs on a Raspberry Pi and
drives WS2811 pixel-addressable LED strips via a Teensy 3.x microcontroller
connected over USB serial. A React web app (served by the Go server) lets the
user choose and configure lighting animations.

## Repository layout

```
brightlight/
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
│   ├── animations/    LED animation effects (14 animations)
│   ├── config/        Network config, presets, static data
│   ├── controller/    Teensy USB serial driver + relay driver
│   ├── framebuffer/   In-memory LED frame buffer
│   ├── segment/       Named/physical/combined LED segment abstractions
│   ├── servers/       HTTP handlers (animations, config, SSE streams)
│   ├── stats/         Performance statistics
│   └── ui-config/     Button/segment config JSON (served via /api/ui-config)
├── firmware/          Arduino/Teensy C firmware (OctoWS2811)
├── deploy/            Staging area for Pi deployment artefacts (generated)
│   └── README.md      How to build and deploy
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
| Hardware | Raspberry Pi (Linux ARM), Teensy 3.x, WS2811 LED strips |
| Build | Vite (frontend), cross-compiled Go GOARCH=arm GOARM=5 (backend) |

## Runtime configuration

The Go binary uses `BRIGHTLIGHT` env var as the base path for static files.
With `BRIGHTLIGHT=/home/pi` it looks for:
- Frontend: `/home/pi/frontend/build/` (index.html, assets/)
- Button config: `/home/pi/backend/ui-config/`

When `BRIGHTLIGHT` is unset (development) static content is not served — use
the Vite dev server, which proxies `/api` to `http://localhost:8080` — and
button config is read from `ui-config/` relative to the working directory.

## How to build locally (development)

```bash
# Frontend dev server (proxies /api to backend at localhost:8080)
cd frontend && npm start

# Build for deployment — outputs to deploy/
full-build.sh    # Linux/macOS
full-build.bat   # Windows
```

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
