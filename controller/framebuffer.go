package controller

import "sync"

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
	fb.Strips = make([]LedStrip, 0, 8)
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, MaxLedStripLen))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 164))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 165))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 85))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 80))
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 162))

	return &fb
}

// Signal the frame buffer has changed to any listeners
func (fb *FrameBuffer) Flush() {
	fb.Mutex.Lock()
	fb.Cond.Broadcast()
	fb.Mutex.Unlock()
}
