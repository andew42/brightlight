package animations

import "github.com/andew42/brightlight/framebuffer"

type runner struct {
	colour framebuffer.Rgb
	pos    uint
}

func newRunner(colour framebuffer.Rgb) *runner {

	return &runner{colour, 0}
}

func (r *runner) animateNextFrame(seg framebuffer.Segment) {

	seg.Set(r.pos, framebuffer.NewRgb(0, 0, 0))
	r.pos = nextPos(r.pos, seg.Len())
	seg.Set(r.pos, r.colour)
}
