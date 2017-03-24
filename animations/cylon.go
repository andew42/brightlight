package animations

import "github.com/andew42/brightlight/framebuffer"

type cylon struct {
}

func newCylon() *cylon {

	return &cylon{}
}

func (c *cylon) animateFrame(frameCount uint, frame framebuffer.Segment) {

	if frame.Len() == 0 {
		return;
	}

	// Get an incrementing position twice the frame length (forward then backwards)
	pos := uint(frameCount) % (frame.Len() * 2)
	if pos >= frame.Len() {
		// Backwards
		pos = 2*frame.Len() - pos - 1
	}
	frame.Set(pos, framebuffer.NewRgb(255, 0, 0))

	// TODO ADD Trail
}
