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
#
# The install itself runs detached (setsid) once the download completes, so a
# dropped SSH connection can't kill it part-way and leave a half-applied
# upgrade. Progress is logged to /var/log/brightlight-install.log — check it
# after reconnecting if the session drops.

set -euo pipefail

main() {
    REPO="andew42/brightlight"
    TARBALL_URL="https://github.com/${REPO}/releases/latest/download/brightlight-pi.tar.gz"
    LOG="/var/log/brightlight-install.log"

    export INSTALL_DIR="/opt/brightlight"
    export USER_BUTTONS_NAME="backend/ui-config/user-buttons.json"
    export DROPIN_DIR="/etc/systemd/system/brightlight.service.d"

    export SITE="${BRIGHTLIGHT_SITE:-}"
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

    # Until the detached stage takes over, clean the work dir on any failure
    export WORK
    WORK="$(mktemp -d)"
    trap 'rm -rf "${WORK}"' EXIT

    # Download before touching anything so a slow or failed download can
    # never leave the service stopped or the install half-replaced
    echo "Downloading ${TARBALL_URL} ..."
    curl -fsSL --retry 4 --retry-delay 2 -o "${WORK}/brightlight-pi.tar.gz" "${TARBALL_URL}"

    # Preserve user-edited buttons across upgrades
    if [ -f "${INSTALL_DIR}/${USER_BUTTONS_NAME}" ]; then
        cp "${INSTALL_DIR}/${USER_BUTTONS_NAME}" "${WORK}/user-buttons.json"
        echo "Preserving existing user-buttons.json"
    fi

    write_stage2 > "${WORK}/stage2.sh"

    # Hand work-dir ownership to the detached stage (it cleans up itself)
    trap - EXIT
    : > "${LOG}"
    setsid bash "${WORK}/stage2.sh" </dev/null >>"${LOG}" 2>&1 &
    STAGE2_PID=$!

    echo "Installing (detached; survives a dropped SSH session)..."
    echo "Log: ${LOG}"
    tail --pid="${STAGE2_PID}" -n +1 -f "${LOG}" 2>/dev/null || true
    if wait "${STAGE2_PID}"; then
        exit 0
    else
        echo "Install failed — see ${LOG}" >&2
        exit 1
    fi
}

# The part that must not die with the SSH session: replacing the files and
# restarting the service. Runs via setsid with all state passed in the
# environment (WORK, INSTALL_DIR, USER_BUTTONS_NAME, DROPIN_DIR, SITE).
#
# Ordering matters: install and sync ALL files while the old service is
# still running, and only restart it once the upgrade is complete on disk.
# That way an interrupted restart leaves the new version installed and
# starting on the next boot, rather than a half-applied upgrade.
write_stage2() {
    cat <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
trap 'rm -rf "${WORK}"' EXIT

# Install everything first, with the old service still running. The old
# binary is unlinked (not overwritten) so the running process is untouched;
# it keeps executing the old inode until we restart it below
echo "Installing to ${INSTALL_DIR} ..."
mkdir -p "${INSTALL_DIR}"
rm -f "${INSTALL_DIR}/brightlight"
tar -xzf "${WORK}/brightlight-pi.tar.gz" -C "${INSTALL_DIR}"
chmod +x "${INSTALL_DIR}/brightlight"

if [ -f "${WORK}/user-buttons.json" ]; then
    cp "${WORK}/user-buttons.json" "${INSTALL_DIR}/${USER_BUTTONS_NAME}"
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

# Everything the new install needs is now on disk — flush it (log included)
# before restarting, so the upgrade is durable whatever the restart does
echo "Files installed. Restarting the service..."
sync

# Bound the restart so the install can't sit here indefinitely. The service
# is already enabled and the new files are in place, so a restart that does
# not complete still leaves the Pi running the new version after a reboot
if ! timeout 30 systemctl restart brightlight; then
    echo "Service restart did not complete within 30s."
    echo "The upgrade is installed; reboot the Pi to start the new version:"
    echo "  sudo reboot"
    exit 1
fi

IP="$(hostname -I 2>/dev/null | awk '{print $1}')"
echo
echo "Brightlight installed and running."
echo "  Web UI:  http://${IP:-<pi-address>}:8080"
echo "  Status:  systemctl status brightlight"
echo "  Logs:    journalctl -u brightlight -f"
EOF
}

# Wrapped in functions invoked at the very end so that when the script is
# streamed via `curl | bash` nothing executes until it has downloaded in
# full — a dropped connection part-way can't run half an install
main "$@"
