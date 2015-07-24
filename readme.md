## Brightlight

Brightlight is a lighting controller for pixel addressable LED strips,
intended for domestic mood lighting.

It has two parts. A low level controller responsible for generating the LED strip waveforms,
written in Arduino C using the Teensy 3.x controller.

https://www.pjrc.com/teensy/td_libs_OctoWS2811.html

A web site interface to control one or more Teensy boards via USB. This is written in GO and
runs on a raspberry pi or similar.

## Scripts

// **OSX:** Connecting to raspberry pi fom OSX terminal session
ssh pi@192.168.0.44

// **pi:** clean
sudo killall -q -9 brightlight
rm -f -r /home/pi/go
rm -f /home/pi/brightlight.log
mkdir -p /home/pi/go/src

// **OSX:** copy latest source
scp -r /Users/andrew/Dropbox/go/src pi@192.168.0.44:/home/pi/go

// **pi:** build
export GOPATH=/home/pi/go
cd /home/pi/go/src/github.com/andew42/brightlight
go install

// **pi:** run
stty -F /dev/ttyACM0 raw
export GOPATH=/home/pi/go
/home/pi/go/bin/brightlight

## Dev Environment
**go 1.4** https://golang.org/dl/
**WebStorm 9.02** https://www.jetbrains.com/webstorm/
**Go language plugin 0.9.15.3** Newer alpha version have trouble with import paths
External Libraries breaks in project when new go version is installed:
Create a new project
Copy misc.xml from hidden .idea directory to same directory in brightlight project

##To Do
Button edits should persist
Support editing button names
Support editing animations
Doesn't support landscape well (as home screen app)
Doesn't support on call bar well (as home screen app)
Colour picker doesn't display colour and doesn't adjust the room lights
