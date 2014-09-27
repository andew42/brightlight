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

	// TODO: Make this more dynamic?
	fb.Strips = make([]LedStrip, 0, 2)
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 10))
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 8))

	return &fb
}

// Signal the frame buffer has changed to any listeners
func (fb *FrameBuffer) Flush() {
	fb.Mutex.Lock()
	fb.Cond.Broadcast()
	fb.Mutex.Unlock()
}
