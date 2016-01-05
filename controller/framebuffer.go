package controller

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
	"sync"
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

	// 2 Bed wall
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 168))

	// 3 Bed curtains
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 164))

	// 4 Bed ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 165))

	// 5 Dressing table wall
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 85))

	// 6 Dressing table ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 80))

	// 7 Dressing table curtain
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 162))

	// 8 Bathroom mirror wall
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 172))

	// 9 Bath ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 226))

	// 10 Bath+ wall
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 291))

	// 11 Bathroom mirror ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 162))

	// 12 Unused
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 0))

	// 13 Left of door ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 88))

	// 14 Right of door ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 142))

	// 15 Right of door wall
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 122))

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
