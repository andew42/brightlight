# Brightlight

A lighting controller for pixel-addressable WS2811 LED strips intended for domestic mood lighting.

## Architecture

- **Teensy 3.x microcontroller** — generates LED strip waveforms (Arduino C, OctoWS2811)
- **Go server on Raspberry Pi** — drives the Teensy over USB serial, serves the web UI
- **React web app** — controls animations and button configuration

## Repository layout

```
frontend/   React 19 + Vite 8 web app
backend/    Go server
firmware/   Arduino/Teensy C firmware (OctoWS2811)
deploy/     Build output staging area for Pi deployment (see deploy/README.md)
```

## Building

Run `full-build.bat` from the repo root (Windows). This:
1. Builds the React frontend (`frontend/`) with Vite
2. Cross-compiles the Go backend for Linux ARM (Raspberry Pi)
3. Copies artefacts to `deploy/`

## Setting up a new Pi

1. Flash Raspberry Pi OS Lite, enable SSH
2. SSH in: `ssh pi@192.168.0.XXX`
3. Configure autostart in `/etc/rc.local`:
   ```bash
   export BRIGHTLIGHT=/home/pi
   /home/pi/brightlight > /dev/null 2>&1 &
   ```
4. Run `full-build.bat` on your PC
5. Deploy to Pi:
   ```bash
   scp ./deploy/brightlight pi@192.168.0.XXX:/home/pi/brightlight
   scp -r ./deploy/frontend/build pi@192.168.0.XXX:/home/pi/frontend/build
   ssh pi@192.168.0.XXX "sudo chmod +x /home/pi/brightlight && sudo reboot"
   ```

## Development environment

- Go 1.17+ — https://golang.org/dl/
- Node.js 20+ — https://nodejs.org/
- JetBrains GoLand / IntelliJ — open `frontend/` and `backend/` as separate projects

### Frontend dev server
```bash
cd frontend
npm install
npm start    # dev server at http://localhost:5173, proxies API to device
```

### Kill running instance on Pi
```bash
sudo killall -q -9 brightlight
```

## Setting up Arduino / Teensy environment

- [Arduino 1.6.3](https://www.arduino.cc/en/Main/OldSoftwareReleases#previous)
- [Teensyduino 1.26](https://www.pjrc.com/teensy/td_download.html)
- Or use [PlatformIO](https://platformio.org/)

## To Do

- Candle animation
- Fairground light chasers
- Static string (Gazebo lights)
- Clock animation
- Fade between animations
- Support for Alexa
