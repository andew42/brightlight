package controller

import (
	"sync"
)

// Frame buffer is a slice of strips
// A strip is a slice of LEDs with a direction
type FrameBuffer struct {
	sync.Mutex
	*sync.Cond
	Strips [] LedStrip
}

// Create a frame buffer
func NewFrameBuffer() *FrameBuffer {
	var fb FrameBuffer
	fb.Cond = sync.NewCond(&fb.Mutex)

	// TODO: Make this more dynamic?
	fb.Strips = make([] LedStrip, 0, 2)

	// TODO: Strip 1 and 2
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 10))
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 8))

	return &fb
}

// Signal the frame buffer has changed to any listeners
func (fb *FrameBuffer) SignalChanged() {
	fb.Mutex.Lock()
	fb.Cond.Broadcast()
	fb.Mutex.Unlock()
}

// Set all frame buffer entries to a particular colour
func (fb *FrameBuffer) SetColour(colour int64) {
	red := byte(colour >> 16)
	green := byte(colour >> 8)
	blue := byte(colour)

	for s := 0; s < len(fb.Strips); s++ {
		for l := 0; l < len(fb.Strips[s].Leds); l++ {
			led := &fb.Strips[s].Leds[l]
			led.Red = red
			led.Green = green
			led.Blue = blue
		}
	}
}
