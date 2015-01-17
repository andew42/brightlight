package controller

import (
	"io"
	"log"
	"os"
	"runtime"
	"time"
	"github.com/andew42/brightlight/stats"
)

var driverStarted bool

// Run driver as a go routine
func StartDriver(fb *FrameBuffer, statistics *stats.Stats) {

	if driverStarted {
		panic("Teensy Driver started twice")
	}
	driverStarted = true

	// Start 2 drivers (16 channels)
	go teensyDriver(0, fb, statistics)
	go teensyDriver(1, fb, statistics)
}

func IsDriverConnected() bool {
	return usbConnected
}

var usbConnected bool

// Monitors changes to frame buffer and update Teensy via USB
func teensyDriver(driverIndex int, fb *FrameBuffer, statistics *stats.Stats) {

	port := getPortName(driverIndex)
	if port == "" {
		log.Printf("unknown port name for driver index %v", driverIndex)
		return
	}

	for {
		usbConnected = false
		f := openUsbPort(port)
		usbConnected = true

		// Push frame buffer changes to Teensy
		for {
			fb.Mutex.Lock()

			started := time.Now()

			// Send the frame buffer
			var data []byte = make([]byte, 0)
			data = append(data, 0xff, 0xff, 0xff, 0xff)
			startStrip := driverIndex * 8
			for s := startStrip; s < startStrip + 8 ; s++ {
				for l := 0; l < MaxLedStripLen; l++ {
					if l >= len(fb.Strips[s].Leds) {
						// Pad frame buffer as strip is < MaxLedStripLen
						data = append(data, 0, 0, 0, 0)
					} else {
						// Colours are sent as 4 bytes with leading 0x00
						rgb := fb.Strips[s].Leds[l]
						data = append(data, 0, rgb.Red, rgb.Green, rgb.Blue)
					}
				}
			}
			n, err := f.Write(data)
			if err == nil && n < len(data) {
				err = io.ErrShortWrite
			}
			if err != nil {
				fb.Mutex.Unlock()
				log.Printf(err.Error())
				f.Close()

				// Try and reconnect
				break
			}

			statistics.AddSerial(time.Since(started))
			// Wait for next frame buffer update
			fb.Cond.Wait()
			fb.Mutex.Unlock()
		}
	}
}

// Retry port open until it succeeds
func openUsbPort(port string) *os.File {

	errorLogged := false
	for {
		f, err := os.Create(port)
		if err == nil {
			log.Printf(port + " connected")
			return f
		}

		if !errorLogged {
			log.Printf(err.Error())
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
		case 0: return "/dev/cu.usbmodem103721"
		// Teensy 3.1  "/dev/cu.usbmodem103101"
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
