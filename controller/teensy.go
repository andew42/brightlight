package controller

import (
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

// Monitors changes to the frame buffer and update Teensy via USB
func TeensyDriver(fb *FrameBuffer) {

	for {
		// Open the port
		f := openUsbPort()

		// Push frame buffer changes to Teensy
		for {
			fb.Mutex.Lock()

			// Send the frame buffer TODO:3
			var data []byte = make([]byte, 0)
			data = append(data, 0x20, 0x20, 0x20, 0x20)
			for l := 0; l < 3; l++ {
				rgb := fb.Strips[0].Leds[l]
				data = append(data, rgb.Red, rgb.Green, rgb.Blue)
			}
			// fmt.Print(data)
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
		// Open the serial port
		port := "/dev/ttyACM0"
		if runtime.GOOS == "darwin" {
			port = "/dev/cu.usbmodem103101"
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
