package controller

import (
	"sync"
	"strconv"
	log "github.com/Sirupsen/logrus"
)

// Frame buffer is a slice of strips A Mutex
// and Cond are used to broadcast changes
type FrameBuffer struct {

	sync.Mutex
	*sync.Cond
	Strips []LedStrip
}

// Create a frame buffer
func NewFrameBuffer() *FrameBuffer {

	var fb FrameBuffer
	fb.Cond = sync.NewCond(&fb.Mutex)

	// TODO: Make this more dynamic from config file?
	fb.Strips = make([]LedStrip, 0, StripsPerTeensy)
	// 0, 1 Unused strips
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 0))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 0))

	// 2 Bed Wall
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 168))

	// 3 Bed Curtains
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 164))

	// 4 Bed Ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 165))

	// 5 Dressing Table Wall
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 85))

	// 6 Dressing Table Ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 80))

	// 7 Dressing Table Curtain
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 162))

	// 8, 15 Unused strips
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))

	// Sanity check
	numberOfStrips := len(fb.Strips)
	if numberOfStrips <= 0 || numberOfStrips%StripsPerTeensy != 0 {
		log.WithField("StripsPerTeensy", strconv.Itoa(StripsPerTeensy)).Panic("framebuffer strips must be multiple of")
	}
	return &fb
}

// Signal the frame buffer has changed to any listeners
func (fb *FrameBuffer) Flush() {

	fb.Mutex.Lock()
	fb.Cond.Broadcast()
	fb.Mutex.Unlock()
}
