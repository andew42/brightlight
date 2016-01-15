package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/stats"
	"io"
	"os"
	"runtime"
	"time"
)

var driverStarted bool

// Run driver or two as a go routine
func StartDriver() {

	if driverStarted {
		log.Panic("Teensy driver started twice")
	}
	driverStarted = true

	// Start 2 sub drivers (16 channels)
	go teensyDriver(0)
	go teensyDriver(1)
}

func IsDriverConnected() bool {

	return usbConnected
}

var usbConnected bool

// Monitors changes to frame buffer and update Teensy via USB
func teensyDriver(driverIndex int) {

	port := getPortName(driverIndex)
	if port == "" {
		log.WithField("driverIndex", driverIndex).Warn("teensyDriver unknown port name")
		return
	}

	for {
		usbConnected = false
		f := openUsbPort(port)
		usbConnected = true

		// Allocate buffer once to avoid garbage collections in loop
		var data = make([]byte, 4+config.MaxLedStripLen*8*4+4)

		// Request frame buffer updates
		src, done := framebuffer.AddListener(port, true)

		// Push frame buffer changes to Teensy
		for {
			select {
			case fb := <-src:
				started := time.Now()
				// Build the frame buffer, start with header of 4 * 0xff
				i := 0
				for z := 0; z < 4; z++ {
					data[i] = 0xff
					i++
				}
				startStrip := driverIndex * 8
				var checksum int32 = 0
				// Buffer is send 8*LED1, 8*LED2 ... 8*(LEDS_PER_STRIP - 1)
				for l := 0; l < config.MaxLedStripLen; l++ {
					for s := startStrip; s < startStrip+8; s++ {
						if l >= len(fb.Strips[s].Leds) {
							// Pad frame buffer with zeros as strip is < MaxLedStripLen
							for z := 0; z < 4; z++ {
								data[i] = 0
								i++
							}
						} else {
							// Colours are sent as 4 bytes with leading 0x00
							rgb := fb.Strips[s].Leds[l]
							data[i] = 0
							i++
							data[i] = rgb.Red
							i++
							data[i] = rgb.Green
							i++
							data[i] = rgb.Blue
							i++
							// Update the checksum
							checksum += ((int32(rgb.Red) << 16) + (int32(rgb.Green) << 8) + int32(rgb.Blue))
						}
					}
				}

				// Append checksum MSB first
				for z := 3; z >= 0; z-- {
					data[i] = byte((checksum >> (8 * uint(z))) & 0xff)
					i++
				}

				// Send the frame buffer
				n, err := f.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				if err != nil {
					log.WithField("err", err).Warn("teensyDriver")
					f.Close()

					// Close down listener
					done <- src

					// Try and reconnect
					break
				}
				stats.AddSerialSendTimeSample(port, time.Since(started))
			}
		}
	}
}

// Retry port open until it succeeds
func openUsbPort(port string) *os.File {

	errorLogged := false
	for {
		f, err := os.Create(port)
		if err == nil {
			log.WithField("port", port).Info("openUsbPort connected")
			return f
		}

		if !errorLogged {
			log.WithField("err", err).Warn("openUsbPort")
			errorLogged = true
		}

		// Try again in a second
		time.Sleep(1000 * time.Millisecond)
	}
}

// Determine port name based on index and OS
func getPortName(index int) string {

	if runtime.GOOS == "darwin" {
		// OSX
		switch index {
		case 0:
			return "/dev/cu.usbmodem288181"
			// Teensy 3.0 "/dev/cu.usbmodem103721"
			// Teensy 3.1 "/dev/cu.usbmodem103101"
		}
	} else if runtime.GOOS == "windows" {
		// Windows
		switch index {
		case 0:
			return "COM3"
		case 1:
			return "COM4"
		}
	} else {
		// Raspberry pi
		switch index {
		case 0:
			return "/dev/ttyACM0"
		case 1:
			return "/dev/ttyACM1"
		}
	}
	return ""
}
