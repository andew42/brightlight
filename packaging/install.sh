#!/usr/bin/env bash
# Brightlight installer for Raspberry Pi.
#
# Downloads the latest CI build, installs it to /opt/brightlight and sets up
# a systemd service so it starts on boot. Safe to re-run to upgrade; user
# button edits (user-buttons.json) are preserved.
#
#   curl -fsSL https://github.com/andew42/brightlight/releases/latest/download/install.sh | sudo bash
#
# The hardware site layout (which house/room wiring the server drives) can be
# selected with --site titania|bedroom (or a BRIGHTLIGHT_SITE env var):
#
#   curl -fsSL .../install.sh | sudo bash -s -- --site bedroom
#
# The choice is stored in a systemd drop-in and kept across upgrades; when no
# site is given the existing choice (or the built-in default, titania) is kept.

set -euo pipefail

main() {
    REPO="andew42/brightlight"
    INSTALL_DIR="/opt/brightlight"
    TARBALL_URL="https://github.com/${REPO}/releases/latest/download/brightlight-pi.tar.gz"
    USER_BUTTONS="${INSTALL_DIR}/backend/ui-config/user-buttons.json"
    DROPIN_DIR="/etc/systemd/system/brightlight.service.d"

    SITE="${BRIGHTLIGHT_SITE:-}"
    while [ $# -gt 0 ]; do
        case "$1" in
            --site) SITE="${2:-}"; shift 2 ;;
            --site=*) SITE="${1#*=}"; shift ;;
            *) echo "Unknown option: $1" >&2; exit 1 ;;
        esac
    done
    if [ -n "${SITE}" ] && [ "${SITE}" != "titania" ] && [ "${SITE}" != "bedroom" ]; then
        echo "Invalid --site '${SITE}' (expected titania or bedroom)" >&2
        exit 1
    fi

    if [ "$(id -u)" -ne 0 ]; then
        echo "This installer needs root. Re-run with: curl -fsSL ${TARBALL_URL%brightlight-pi.tar.gz}install.sh | sudo bash" >&2
        exit 1
    fi

    WORK="$(mktemp -d)"
    trap 'rm -rf "${WORK}"' EXIT

    # Download before touching anything so a slow or failed download can
    # never leave the service stopped or the install half-replaced
    echo "Downloading ${TARBALL_URL} ..."
    curl -fsSL --retry 4 --retry-delay 2 -o "${WORK}/brightlight-pi.tar.gz" "${TARBALL_URL}"

    # Preserve user-edited buttons across upgrades
    if [ -f "${USER_BUTTONS}" ]; then
        cp "${USER_BUTTONS}" "${WORK}/user-buttons.json"
        echo "Preserving existing user-buttons.json"
    fi

    # Stop the service, but never let it stall the install: units installed
    # before TimeoutStopSec was added can hang `systemctl stop` for minutes
    # if the binary is wedged, so give it 15s then SIGKILL and move on
    echo "Stopping brightlight service (if running)..."
    if ! timeout 15 systemctl stop brightlight 2>/dev/null; then
        echo "Service did not stop within 15s; killing it"
        systemctl kill -s SIGKILL brightlight 2>/dev/null || true
        sleep 1
    fi

    # --unlink-first so extraction still succeeds if the old binary is
    # somehow running (overwriting a running executable in place fails
    # with "Text file busy")
    echo "Installing to ${INSTALL_DIR} ..."
    mkdir -p "${INSTALL_DIR}"
    tar -xz --unlink-first -f "${WORK}/brightlight-pi.tar.gz" -C "${INSTALL_DIR}"
    chmod +x "${INSTALL_DIR}/brightlight"

    if [ -f "${WORK}/user-buttons.json" ]; then
        cp "${WORK}/user-buttons.json" "${USER_BUTTONS}"
    fi

    echo "Installing systemd service..."
    cp "${INSTALL_DIR}/brightlight.service" /etc/systemd/system/brightlight.service

    if [ -n "${SITE}" ]; then
        echo "Setting site layout to '${SITE}'"
        mkdir -p "${DROPIN_DIR}"
        printf '[Service]\nEnvironment=BRIGHTLIGHT_SITE=%s\n' "${SITE}" > "${DROPIN_DIR}/site.conf"
    elif [ -f "${DROPIN_DIR}/site.conf" ]; then
        echo "Keeping existing site layout: $(grep -o 'BRIGHTLIGHT_SITE=.*' "${DROPIN_DIR}/site.conf")"
    else
        echo "No --site given; using built-in default (titania)"
    fi

    systemctl daemon-reload
    systemctl enable brightlight
    systemctl restart brightlight

    IP="$(hostname -I 2>/dev/null | awk '{print $1}')"
    echo
    echo "Brightlight installed and running."
    echo "  Web UI:  http://${IP:-<pi-address>}:8080"
    echo "  Status:  systemctl status brightlight"
    echo "  Logs:    journalctl -u brightlight -f"
}

# Wrapped in a function and invoked at the very end so that when the script
# is streamed via `curl | bash` nothing executes until it has downloaded in
# full — a dropped connection part-way can't run half an install
main "$@"
