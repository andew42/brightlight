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
│   ├── public/        Static assets including ui-config/ (button config JSON)
│   ├── build/         Vite build output — served by Go server (gitignored)
│   ├── package.json   npm deps with overrides for React 19 peer dep compat
│   ├── vite.config.js Vite 8 config: outDir=build, esbuild CSS minify
│   └── ui-asset-source/ Pixelmator source for apple-touch-icon
├── backend/           Go server + Arduino firmware
│   ├── main.go        Entry point; serves frontend/build/ as static files
│   ├── go.mod         module github.com/andew42/brightlight, Go 1.17
│   ├── animations/    LED animation effects (14 animations)
│   ├── config/        Network config, presets, static data
│   ├── controller/    Teensy USB serial driver + relay driver
│   ├── framebuffer/   In-memory LED frame buffer
│   ├── segment/       Named/physical/combined LED segment abstractions
│   ├── servers/       HTTP handlers (animations, config, websockets)
│   ├── stats/         Performance statistics
│   └── firmware/      Arduino/Teensy C firmware (OctoWS2811)
├── deploy/            Staging area for Pi deployment artefacts (generated)
│   └── README.md      How to build and deploy
├── full-build.bat     Windows build script: builds frontend + backend → deploy/
├── readme.md          Original project README with Pi setup instructions
└── CLAUDE.md          This file
```

## Tech stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 19, React Router 7, Semantic UI React 2, Vite 8 |
| Backend | Go 1.17, logrus, golang.org/x/net/websocket |
| Hardware | Raspberry Pi (Linux ARM), Teensy 3.x, WS2811 LED strips |
| Build | Vite (frontend), cross-compiled Go GOARCH=arm GOARM=5 (backend) |

## Runtime configuration

The Go binary uses `BRIGHTLIGHT` env var as the base path for static files.
With `BRIGHTLIGHT=/home/pi` it looks for:
- Frontend: `/home/pi/frontend/build/` (index.html, assets/)
- Button config: `/home/pi/frontend/build/ui-config/`

## How to build locally (development)

```bash
# Frontend dev server (proxies /RunAnimations to device at 192.168.0.68:8080)
cd frontend && npm start

# Build for deployment
full-build.bat   # Windows — outputs to deploy/
```

## HTTP API (backend → frontend contract)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/ui-config/*.json` | GET/PUT | Button/segment config JSON files |
| `/RunAnimations/` | POST | Run animation (JSON button payload) |
| `/StripLength/` | POST | Set strip length on room lights |
| `/ButtonState` | WebSocket | Push active button key + config version |
| `/FrameBuffer` | WebSocket | Push raw LED frame data (virtual display) |
| `/Stats` | WebSocket | Push performance stats |
| `/option/` | POST | Set server options |

## Key gotchas and past decisions

### React 19 upgrade (April 2026)
- `ReactDOM.render()` removed → use `createRoot()` in `frontend/src/index.jsx`
- React Router v4 → v7: `Switch`→`Routes`, `render`/`component` props → `element`
  prop, `history.push()` → `useNavigate()` hook. Class components can't use
  hooks directly so `ButtonPadRoute` and `ButtonEditorRoute` wrapper components
  inject `navigate`/`location` as a history-compatible object.
- `semantic-ui-react` peer deps cap at React 18 — use npm `overrides` in
  `package.json` to force React 19. The library itself has no runtime
  incompatibilities (verified: zero uses of `findDOMNode` or removed APIs).

### Vite 8 / Rolldown upgrade (April 2026)
- **All source files use `.jsx` extension** — Vite 8 uses Rolldown which will
  not parse JSX in `.js` files. Do not rename them back to `.js`.
- `semantic-ui-css` 2.5.0 has an invalid CSS selector that lightningcss (Vite 8
  default minifier) rejects. Fixed with `css.lightningcss.errorRecovery: true`
  in `vite.config.js`. Do not remove this.
- `esbuild`/`optimizeDeps.esbuildOptions` Vite 6 workarounds are gone; the
  oxc/rolldown pipeline handles everything natively now.

### Frontend dev proxy
`vite.config.js` proxies `/RunAnimations` to `http://192.168.0.68:8080` during
`npm start`. Change the IP if the device address changes.

### Old UI removed
The original `ui/` directory (plain HTML/JS with Ractive.js) has been deleted.
The `/ui/` HTTP route and `/config/` handler have been removed from `main.go`.

### Repo structure (April 2026)
Reorganised from flat layout into `frontend/`, `backend/`, `deploy/` to allow
opening each part as a separate IntelliJ project. Previously `ui2/` was the
frontend and all Go files lived at the repo root.
