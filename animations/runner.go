package animations

import "github.com/andew42/brightlight/framebuffer"

type runner struct {
	colour framebuffer.Rgb
}

func newRunner(colour framebuffer.Rgb) *runner {

	return &runner{colour}
}

func (r *runner) animateFrame(frameCount uint, frame framebuffer.Segment) {

	frame.Set(uint(frameCount)%frame.Len(), r.colour)
}
