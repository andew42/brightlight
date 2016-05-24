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
* Doesn't support landscape well (as home screen app)
* Doesn't support on call bar well (as home screen app)

WARNING: DATA RACE
Write by goroutine 6:
  github.com/andew42/brightlight/controller.teensyDriver()
      /Users/andrew/Dropbox/go/src/github.com/andew42/brightlight/controller/teensy.go:43 +0x27d

Previous read by goroutine 13:
  github.com/andew42/brightlight/servers.RunAnimationsHandler()
      /Users/andrew/Dropbox/go/src/github.com/andew42/brightlight/servers/runanimations.go:41 +0x867
  net/http.HandlerFunc.ServeHTTP()
      /private/var/folders/vd/7l9ys5k57l91x63sh28wl_kc0000gn/T/workdir/go/src/net/http/server.go:1422 +0x47
  net/http.(*ServeMux).ServeHTTP()
      /private/var/folders/vd/7l9ys5k57l91x63sh28wl_kc0000gn/T/workdir/go/src/net/http/server.go:1699 +0x212
  net/http.serverHandler.ServeHTTP()
      /private/var/folders/vd/7l9ys5k57l91x63sh28wl_kc0000gn/T/workdir/go/src/net/http/server.go:1862 +0x206
  net/http.(*conn).serve()
      /private/var/folders/vd/7l9ys5k57l91x63sh28wl_kc0000gn/T/workdir/go/src/net/http/server.go:1361 +0x117c
INFO[0394] Framebuffer listener removed                  name={name:/dev/cu.usbmodem288181 isSerial:true}

Goroutine 6 (running) created at:
  github.com/andew42/brightlight/controller.StartTeensyDriver()
      /Users/andrew/Dropbox/go/src/github.com/andew42/brightlight/controller/teensy.go:22 +0x156
  main.main()
      /Users/andrew/Dropbox/go/src/github.com/andew42/brightlight/main.go:33 +0x4a9

Goroutine 13 (running) created at:
  net/http.(*Server).Serve()
      /private/var/folders/vd/7l9ys5k57l91x63sh28wl_kc0000gn/T/workdir/go/src/net/http/server.go:1910 +0x464
  net/http.(*Server).ListenAndServe()
      /private/var/folders/vd/7l9ys5k57l91x63sh28wl_kc0000gn/T/workdir/go/src/net/http/server.go:1877 +0x174
  net/http.ListenAndServe()
      /private/var/folders/vd/7l9ys5k57l91x63sh28wl_kc0000gn/T/workdir/go/src/net/http/server.go:1967 +0xe2
  main.main()
      /Users/andrew/Dropbox/go/src/github.com/andew42/brightlight/main.go:68 +0xa8e
==================
WARN[0399] openUsbPortWithRetry failed to open port      error=open /dev/cu.usbmodem288181: no such file or directory
