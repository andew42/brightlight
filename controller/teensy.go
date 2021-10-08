package controller

import (
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/stats"
	log "github.com/sirupsen/logrus"
	"time"
)

var teensyDriverStarted bool

// StartTeensyDriver Run driver or two as a go routine
func StartTeensyDriver() {

	if teensyDriverStarted {
		log.Panic("Teensy driver started twice")
	}
	teensyDriverStarted = true

	// Start 2 sub drivers (16 channels)
	go teensyDriver(0)
	go teensyDriver(1)
}

func TeensyConnections() []bool {

	// TODO: 1 of 2 Returning here causes race detector to fail
	return teensyConnections
}

var teensyConnections = make([]bool, 2)

// Monitors changes to frame buffer and update Teensy via USB
func teensyDriver(driverIndex int) {

	port := getPortName(teensyPortMappings, driverIndex)
	if port == "" {
		log.WithField("driverIndex", driverIndex).Warn("teensyDriver unknown port name")
		return
	}

	for {
		// TODO: 2 of 2 Writing here causes race detector to fail
		teensyConnections[driverIndex] = false
		f := openUsbPortWithRetry(port)
		teensyConnections[driverIndex] = true

		// Allocate buffer once to avoid garbage collections in loop
		var data = make([]byte, 4+config.MaxLedStripLen*8*4+4)

		// Request frame buffer updates
		src, done := framebuffer.AddListener(port, true)

		// Push frame buffer changes to Teensy
	pushLoop:
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
							// Perform the output mapping here
							rgb := mapOutput(fb.Strips[s].Leds[l])
							// Colours are sent as 4 bytes with leading 0x00
							data[i] = 0
							i++
							data[i] = rgb.Red
							i++
							data[i] = rgb.Green
							i++
							data[i] = rgb.Blue
							i++
							// Update the checksum
							checksum += (int32(rgb.Red) << 16) + (int32(rgb.Green) << 8) + int32(rgb.Blue)
						}
					}
				}

				// Append checksum MSB first
				for z := 3; z >= 0; z-- {
					data[i] = byte((checksum >> (8 * uint(z))) & 0xff)
					i++
				}

				if _, err := f.Write(data); err != nil {
					log.WithField("error", err.Error()).Warn("teensyDriver send failed")
					f.Close()

					// Close down listener
					done <- src

					// Try and reconnect
					break pushLoop
				}
				stats.AddSerialSendTimeSample(port, time.Since(started))
			}
		}
	}
}
