Brightlight is a lighting controller for pixel addressable LED strips.

It has two parts, a low level controller responsible for generating the LED strip waveforms,
written in Arduino C and using the Teensy 3.0 controller.

https://www.pjrc.com/teensy/td_libs_OctoWS2811.html

A web site interface to control one or more Teensy boards via USB. This is written in GO for
running on a raspberry pi or similar

// OSX: Connecting to raspberry pi fom OSX terminal session
ssh pi@192.168.0.44

// pi: clean
sudo killall -q -9 brightlight
rm -f -r /home/pi/go
rm -f /home/pi/brightlight.log
mkdir -p /home/pi/go/src

// OSX: copy latest source
scp -r /Users/andrew/Dropbox/go/src pi@192.168.0.44:/home/pi/go

// pi: build
export GOPATH=/home/pi/go
cd /home/pi/go/src/github.com/andew42/brightlight
go install

// pi: run
stty -F /dev/ttyACM0 raw
export GOPATH=/home/pi/go
/home/pi/go/bin/brightlight
