## Brightlight

Brightlight is a lighting controller for pixel addressable LED strips, intended for domestic mood lighting.

It has two parts. A low level controller responsible for generating the LED strip waveforms,
written in Arduino C using the [Teensy 3.x controller](https://www.pjrc.com/teensy/td_libs_OctoWS2811.html).

A website interface to control one or more Teensy boards via USB. This is written in GO and
runs on a raspberry Pi or similar.

### Setting up Arduino environment

* [Install Arduino 1.6.3](https://www.arduino.cc/en/Main/OldSoftwareReleases#previous)
* [Install Teensyduino 1.26](https://www.pjrc.com/teensy/td_download.html)
* [Alternatively use PlatformIO](https://platformio.org/)

### Setting up new Pi
* [Download RASPBIAN STRETCH LITE](https://www.raspberrypi.org/downloads/raspbian/)
* Create a (4G) SD card with [etcher](https://etcher.io/)
* Create a new file called "ssh" on the SD card's FAT boot partition. This will enable the SSH daemon immediately after the first boot.
* SSH onto pi from Windows or Mac (default password raspberry)
 ```bash
 ssh pi@192.168.0.XXX
 ``` 
* On pi set up auto start
```bash
sudo nano /etc/rc.local
```
* On pi in nano add
```bash
export BRIGHTLIGHT=/home/pi
/home/pi/brightlight > /dev/null 2>&1 &
```
* On PC build website and brightlight executable
```bash
full-build.bat
```
* On PC copy executable to pi
```bash
scp  ./brightlight pi@192.168.0.XXX:/home/pi
```
* On pi make executable run-able
```bash
sudo chmod +x ./brightlight 
```
* On pi create folder
```bash
mkdir ui2 
```
* On PC copy ui files to pi
```bash
scp -r ./ui pi@192.168.0.XXX:/home/pi
scp -r ./ui2/build pi@192.168.0.XXX:/home/pi/ui2
```
* Reboot
```bash
sudo reboot
```
* Killing running copies
```bash
sudo killall -q -9 brightlight
```

### Dev Environment
#### go 1.17.1
https://golang.org/dl/

#### Goland 2021.2.3
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
