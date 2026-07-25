package controller

import (
	"sync/atomic"
	"time"

	"log/slog"

	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/stats"
)

var teensyDriverStarted bool

// StartTeensyDriver Run driver or two as a go routine
func StartTeensyDriver() {

	if teensyDriverStarted {
		slog.Error("Teensy driver started twice")
		panic("Teensy driver started twice")
	}
	teensyDriverStarted = true

	// Start 2 sub drivers (16 channels)
	shutdownWaitGroup.Add(2)
	go teensyDriver(0)
	go teensyDriver(1)
}

// TeensyConnections Report connection state, atomics as the driver go
// routines update the state concurrently with HTTP handlers reading it
func TeensyConnections() []bool {

	return []bool{teensyConnections[0].Load(), teensyConnections[1].Load()}
}

var teensyConnections [2]atomic.Bool

// How long a frame write to the Teensy may take before we treat the
// connection as dead. A healthy frame takes a few ms at 12Mb USB speed;
// a Teensy that has crashed or stopped draining its serial input would
// otherwise block the write forever with a full kernel tty queue —
// exactly the state that wedges the Pi's USB controller on process exit
const teensyWriteTimeout = 2 * time.Second

// Monitors changes to frame buffer and update Teensy via USB
func teensyDriver(driverIndex int) {

	defer shutdownWaitGroup.Done()

	port := getPortName(teensyPortMappings, driverIndex)
	if port == "" {
		slog.Warn("teensyDriver unknown port name", "driverIndex", driverIndex)
		return
	}

	for {
		teensyConnections[driverIndex].Store(false)
		f := openUsbPortWithRetry(port)
		if f == nil {
			// Shutdown requested while waiting for the port
			return
		}
		teensyConnections[driverIndex].Store(true)

		// Allocate buffer once to avoid garbage collections in loop
		var data = make([]byte, 4+config.MaxLedStripLen*8*4+4)

		// Request frame buffer updates
		src, done := framebuffer.AddListener(port, true)

		// Push frame buffer changes to Teensy
		for {
			var fb *framebuffer.FrameBuffer
			select {
			case fb = <-src:
			case <-shutdownChan:
				// Quiesced close: no writes in progress, queue flushed
				closePortQuiesced(f)
				return
			}

			// Skip if the frame buffer has no strips for this Teensy
			// (e.g. a single Teensy configuration with a second port present)
			startStrip := driverIndex * config.StripsPerTeensy
			if startStrip+config.StripsPerTeensy > len(fb.Strips) {
				continue
			}

			started := time.Now()
			// Build the frame buffer, start with header of 4 * 0xff
			i := 0
			for z := 0; z < 4; z++ {
				data[i] = 0xff
				i++
			}
			var checksum int32 = 0
			// Buffer is send 8*LED1, 8*LED2 ... 8*(LEDS_PER_STRIP - 1)
			for l := 0; l < config.MaxLedStripLen; l++ {
				for s := startStrip; s < startStrip+config.StripsPerTeensy; s++ {
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

			// Bound the write; ignore the error as not all platforms
			// support deadlines on serial ports (Linux does)
			_ = f.SetWriteDeadline(time.Now().Add(teensyWriteTimeout))
			if _, err := f.Write(data); err != nil {
				slog.Warn("teensyDriver send failed", "error", err.Error())
				closePortQuiesced(f)

				// Close down listener then try and reconnect
				done <- src
				break
			}
			stats.AddSerialSendTimeSample(port, time.Since(started))
		}
	}
}
