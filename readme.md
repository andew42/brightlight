## Brightlight

Brightlight is a lighting controller for pixel addressable LED strips, intended for domestic mood lighting.

It has two parts. A low level controller responsible for generating the LED strip waveforms,
written in Arduino C using the [Teensy 3.x controller](https://www.pjrc.com/teensy/td_libs_OctoWS2811.html).

A web site interface to control one or more Teensy boards via USB. This is written in GO and
runs on a raspberry Pi or similar.

### Setting up Arduino environment

* [Install Arduino 1.6.3](https://www.arduino.cc/en/Main/OldSoftwareReleases#previous)
* [Install Teensyduino 1.26](https://www.pjrc.com/teensy/td_download.html)

### Setting up new Pi
* [Download RASPBIAN STRETCH LITE](https://www.raspberrypi.org/downloads/raspbian/)
* Create a (4G) SD card with [etcher](https://etcher.io/)
* Create a new file called "ssh" on the SD card's FAT boot partition. This will enable the SSH daemon immediately after the first boot.
* SSH onto pi from Windows or Mac (password raspberry)
 ```bash
 ssh pi@192.168.0.46
 ``` 
* Set up auto start
```bash
sudo nano /etc/rc.local
```
* Add
```bash
export GOPATH=/home/pi
/home/pi/brightlight > /dev/null 2>&1 &
```
* Build brightlight on Windows or Mac for Arm
```bash
env GOOS=linux GOARCH=arm go build -o ./brightlight
```
* Copy executable to pi
```bash
 scp  ./brightlight pi@192.168.0.46:/home/pi
```
* Create ui folder
```bash
mkdir -p /home/pi/src/github.com/andew42/brightlight/ui2/build
```
* Copy ui files (build first)
```bash
scp -r . pi@192.168.0.46:/home/pi/src/github.com/andew42/brightlight/ui2/build
```
* Killing running copies
```bash
sudo killall -q -9 brightlight
```

### Dev Environment
#### go 1.11.1
https://golang.org/dl/

go get "github.com/sirupsen/logrus"

#### Goland 2018.2
https://www.jetbrains.com/go/

#### Grep Console plugin
With ANSI terminal emulator enabled for logrus

#### To Do
#####UI
* Scrolling slider for config setting is horrid

#####Engine
* Candle
* Fair ground light chasers
* Static string (Gazebo lights)
* Overlay pieces (corner + mirror)
* Clock
* Security indicators
* Fade between animations
* Cylon stacking
* Plasma charging bulls
* Meteors
* Fireworks
* Support for alexa
