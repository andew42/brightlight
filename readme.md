# Brightlight

A lighting controller for pixel-addressable WS2811 LED strips intended for domestic mood lighting.

## Architecture

- **Teensy 3.x microcontroller** — generates LED strip waveforms (Arduino C, OctoWS2811)
- **Go server on Raspberry Pi** — drives the Teensy over USB serial, serves the web UI
- **React web app** — controls animations and button configuration
- **Alexa voice control (optional)** — custom skill + AWS Lambda calling the server over HTTPS

## Repository layout

```
frontend/   React 19 + Vite 8 web app
backend/    Go server (stdlib only)
alexa/      Alexa custom skill: interaction model + AWS Lambda (see alexa/readme.md)
firmware/   Arduino/Teensy C firmware (OctoWS2811)
packaging/  systemd unit + Pi installer bundled by CI
deploy/     Build output staging area for Pi deployment (see deploy/README.md)
docs/       Design notes
```

## Building

### CI (GitHub Actions)

`.github/workflows/build.yml` runs on every push and pull request to `master`
or `develop`:

1. Vets and cross-compiles the Go backend for Raspberry Pi 2B
   (`GOOS=linux GOARCH=arm GOARM=7`)
2. Builds the React frontend with Vite (`npm ci && npm run build`)
3. On push, packages the binary, frontend build, ui-config, systemd unit and
   installer into `brightlight-pi.tar.gz` and republishes it as the rolling
   `latest` GitHub release

### Local build

Run `full-build.sh` (Linux/macOS) or `full-build.bat` (Windows) from the repo
root. It performs the same steps as CI and stages the artefacts in `deploy/`
(see `deploy/README.md` for deploying a local build manually).

## Setting up a new Pi

1. Flash Raspberry Pi OS Lite, enable SSH
2. SSH in and run the installer:
   ```bash
   curl -fsSL https://github.com/andew42/brightlight/releases/latest/download/install.sh | sudo bash
   ```
   Append `-s -- --site bedroom` to select the hardware site layout
   (default `titania`; the choice is remembered across upgrades).

This downloads the latest CI build, installs it to `/opt/brightlight` and sets
up the `brightlight` systemd service so it starts on boot. Re-run the same
command to upgrade; user button edits (`user-buttons.json`) are preserved.

Useful commands on the Pi:

```bash
systemctl status brightlight       # service state
journalctl -u brightlight -f       # follow logs
sudo systemctl restart brightlight
```

The web UI is served on port 8080.

## Development environment

- Go 1.26+ — https://golang.org/dl/
- Node.js 20+ — https://nodejs.org/
- JetBrains GoLand / IntelliJ — open `frontend/` and `backend/` as separate projects

### Backend dev server
```bash
cd backend
go run .     # API at http://localhost:8080 (static content needs BRIGHTLIGHT set)
```

### Frontend dev server
```bash
cd frontend
npm install
npm start    # dev server at http://localhost:5173, proxies /api to localhost:8080
```

## Setting up Arduino / Teensy environment

- [Arduino 1.6.3](https://www.arduino.cc/en/Main/OldSoftwareReleases#previous)
- [Teensyduino 1.26](https://www.pjrc.com/teensy/td_download.html)
- Or use [PlatformIO](https://platformio.org/)

## To Do

- Candle animation
- Clock animation
- Fade between animations
