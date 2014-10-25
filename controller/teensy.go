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
	go teensyDriver(fb, statistics)
}

func IsDriverConnected() bool {
	return usbConnected
}

var usbConnected bool

// Monitors changes to frame buffer and update Teensy via USB
func teensyDriver(fb *FrameBuffer, statistics *stats.Stats) {
	for {
		usbConnected = false
		f := openUsbPort()
		usbConnected = true

		// Push frame buffer changes to Teensy
		for {
			fb.Mutex.Lock()

			started := time.Now()
			// Send the frame buffer TODO:Initial size
			var data []byte = make([]byte, 0)
			data = append(data, 0x20, 0x20, 0x20, 0x20)
			for s := 0; s < len(fb.Strips); s++ {
				for l := 0; l < MaxLedStripLen; l++ {
					if l >= len(fb.Strips[s].Leds) {
						// Pad frame buffer as strip is < MaxLedStripLen
						data = append(data, 0, 0, 0)
					} else {
						rgb := fb.Strips[s].Leds[l]
						data = append(data, rgb.Red, rgb.Green, rgb.Blue)
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
func openUsbPort() *os.File {
	errorLogged := false
	for {
		// Open the serial port (raspberry pi or OSX)
		port := "/dev/ttyACM0"
		if runtime.GOOS == "darwin" {
// TODO			port = "/dev/cu.usbmodem103101"
			port = "/dev/cu.usbmodem103721"
		}

		f, err := os.Create(port)
		if err == nil {
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
