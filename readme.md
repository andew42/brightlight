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
* [Build latest raspbian card](https://www.raspberrypi.org/downloads/)
* [Download go tar ball 1.5.2](http://dave.cheney.net/unofficial-arm-tarballs)
* Update profile in /etc to include /usr/local/go/bin in path (sudo nano /etc/profile)
* Create **clean**, **build** and **run** script and chmod 744 than to make them executable
* sudo nano /etc/rc.local and add:
```bash
# Start bright light
export GOPATH=/home/pi/go
/home/pi/go/bin/brightlight > /dev/null 2>&1 &
```

### Scripts

###### OSX: Connecting to raspberry pi fom OSX terminal session
```bash
ssh pi@192.168.0.46
```

###### Pi: clean
```bash
sudo killall -q -9 brightlight
rm -f -r /home/pi/go
rm -f /home/pi/brightlight.log
mkdir -p /home/pi/go/src
```

###### OSX: copy latest source
```bash
rsync -rav -e ssh --exclude='.git' \
/Users/andrew/Dropbox/go/src/ \
pi@192.168.0.46:/home/pi/go/src/
```

###### Pi: build
```bash
export GOPATH=/home/pi/go
cd /home/pi/go/src/github.com/andew42/brightlight
go install
```

###### Pi: run
```bash
export GOPATH=/home/pi/go
/home/pi/go/bin/brightlight
```

### Dev Environment
#### go 1.5.2
https://golang.org/dl/

go get "github.com/Sirupsen/logrus"

#### WebStorm 10.04
https://www.jetbrains.com/webstorm/

#### Go language plugin 0.9.748
https://github.com/go-lang-plugin-org

#### Grep Console plugin
With ANSI terminal emulator enabled for logrus

#### To Do
#####UI
* Does not support landscape well (as home screen app)
* Does not support on call bar well (as home screen app)
* Indicate selected button
#####Engine
* Indicate which button is pressed
* Allow same segment to appear multiple times
* Candle
* Fair ground light chasers
* Static string (Gazebo lights)
* Overlay pieces (corner + mirror)
* Clock
* Security indicators
* Fade between animations
* Cylon stacking
* Twinkle colours
* Plasma charging bulls
* Meteors
* Fireworks
* Support for alexa

#### Building for pi
env GOOS=linux GOARCH=arm go build -o ./pi
