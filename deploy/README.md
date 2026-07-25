# Deploy

## Automated install on the Pi (recommended)

CI packages every push into a rolling `latest` GitHub release containing the
ARMv7 binary, the built frontend, the ui-config files, a systemd unit and an
installer. On the Pi run:

```bash
curl -fsSL https://github.com/andew42/brightlight/releases/latest/download/install.sh | sudo bash
```

This installs everything to `/opt/brightlight`, sets up the `brightlight`
systemd service so it starts on boot, and preserves `user-buttons.json`
across upgrades.

### Site selection

One build drives both installations. Select which hardware layout a Pi runs
with `--site` on first install (stored in a systemd drop-in and kept across
upgrades; default is `titania`):

```bash
curl -fsSL https://github.com/andew42/brightlight/releases/latest/download/install.sh | sudo bash -s -- --site bedroom
```

The server reads the `BRIGHTLIGHT_SITE` environment variable at startup: it
selects the strip/segment layout compiled into the binary and which
`static-data-<site>.json` / `default-buttons-<site>.json` variants are served
to the UI. Useful commands afterwards:

```bash
systemctl status brightlight     # service state
journalctl -u brightlight -f     # follow logs
sudo systemctl restart brightlight
```

The web UI is served on port 8080. The installer sources live in
`packaging/`; the bundle is assembled by `.github/workflows/build.yml`.

## Manual build (this directory)

This directory is the staging area for locally built artefacts. Its layout
matches the CI Pi bundle:

| Path | Source | Description |
|------|--------|-------------|
| `brightlight` | `backend/` Go build | Linux ARMv7 binary (GOARM=7) |
| `frontend/build/` | `frontend/` Vite build | React app static files |
| `backend/ui-config/` | `backend/ui-config/` | Button/segment config JSON |
| `brightlight.service` | `packaging/` | systemd unit |
| `install.sh` | `packaging/` | Pi installer script |

Run `full-build.sh` (Linux/macOS) or `full-build.bat` (Windows) from the repo
root to populate it — the scripts run the same steps as CI (`npm ci` +
Vite build, `go vet`, ARMv7 cross-compile) — then copy the entries above to
the Pi. With `BRIGHTLIGHT=<base>` the server expects:

- Binary anywhere (conventionally `<base>/brightlight`)
- Frontend at `<base>/frontend/build/`
- Button config at `<base>/backend/ui-config/`
