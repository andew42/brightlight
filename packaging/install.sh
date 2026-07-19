#!/usr/bin/env bash
# Brightlight installer for Raspberry Pi.
#
# Downloads the latest CI build, installs it to /opt/brightlight and sets up
# a systemd service so it starts on boot. Safe to re-run to upgrade; user
# button edits (user-buttons.json) are preserved.
#
#   curl -fsSL https://github.com/andew42/brightlight/releases/latest/download/install.sh | sudo bash

set -euo pipefail

REPO="andew42/brightlight"
INSTALL_DIR="/opt/brightlight"
TARBALL_URL="https://github.com/${REPO}/releases/latest/download/brightlight-pi.tar.gz"
USER_BUTTONS="${INSTALL_DIR}/backend/ui-config/user-buttons.json"

if [ "$(id -u)" -ne 0 ]; then
    echo "This installer needs root. Re-run with: curl -fsSL ${TARBALL_URL%brightlight-pi.tar.gz}install.sh | sudo bash" >&2
    exit 1
fi

echo "Stopping brightlight service (if running)..."
systemctl stop brightlight 2>/dev/null || true

# Preserve user-edited buttons across upgrades
BUTTONS_BACKUP=""
if [ -f "${USER_BUTTONS}" ]; then
    BUTTONS_BACKUP="$(mktemp)"
    cp "${USER_BUTTONS}" "${BUTTONS_BACKUP}"
    echo "Preserving existing user-buttons.json"
fi

echo "Downloading ${TARBALL_URL} ..."
mkdir -p "${INSTALL_DIR}"
curl -fsSL "${TARBALL_URL}" | tar -xz -C "${INSTALL_DIR}"
chmod +x "${INSTALL_DIR}/brightlight"

if [ -n "${BUTTONS_BACKUP}" ]; then
    mv "${BUTTONS_BACKUP}" "${USER_BUTTONS}"
fi

echo "Installing systemd service..."
cp "${INSTALL_DIR}/brightlight.service" /etc/systemd/system/brightlight.service
systemctl daemon-reload
systemctl enable brightlight
systemctl start brightlight

IP="$(hostname -I 2>/dev/null | awk '{print $1}')"
echo
echo "Brightlight installed and running."
echo "  Web UI:  http://${IP:-<pi-address>}:8080"
echo "  Status:  systemctl status brightlight"
echo "  Logs:    journalctl -u brightlight -f"
