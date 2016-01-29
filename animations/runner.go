package animations

import "github.com/andew42/brightlight/framebuffer"

type runner struct {
	colour framebuffer.Rgb
	pos    uint
}

func newRunner(colour framebuffer.Rgb) *runner {

	return &runner{colour, 0}
}

func (r *runner) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	frame.Set(r.pos, framebuffer.NewRgb(0, 0, 0))
	r.pos = nextPos(r.pos, frame.Len())
	frame.Set(r.pos, r.colour)
}
