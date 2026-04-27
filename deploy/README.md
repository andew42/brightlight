# Deploy

This directory is the staging area for artefacts deployed to the Raspberry Pi controller.

## Contents (generated, not committed)

| Path | Source | Description |
|------|--------|-------------|
| `brightlight` | `backend/` Go build | Linux ARM binary (GOARCH=arm GOARM=5) |
| `frontend/build/` | `frontend/` Vite build | React app static files |

Run `full-build.bat` (Windows) from the repo root to populate this directory.

## Deploying to the Pi

```bash
# Copy binary
scp ./deploy/brightlight pi@192.168.0.XXX:/home/pi/brightlight

# Copy frontend assets
scp -r ./deploy/frontend/build pi@192.168.0.XXX:/home/pi/frontend/build

# Make binary executable (run on Pi)
sudo chmod +x /home/pi/brightlight
```

## Pi setup

The binary reads the `BRIGHTLIGHT` environment variable to locate static files.
Set it in `/etc/rc.local` so it persists across reboots:

```bash
export BRIGHTLIGHT=/home/pi
/home/pi/brightlight > /dev/null 2>&1 &
```

With `BRIGHTLIGHT=/home/pi` the server expects:
- Binary at `/home/pi/brightlight`
- Frontend at `/home/pi/frontend/build/`
- Button config at `/home/pi/frontend/build/ui-config/`
