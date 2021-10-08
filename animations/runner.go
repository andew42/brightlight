package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
)

type runner struct {
	colour framebuffer.Rgb
}

func newRunner(colour framebuffer.Rgb) *runner {

	return &runner{colour}
}

func (r *runner) animateFrame(frameCount uint, frame segment.Segment) {

	frame.Set(frameCount%frame.Len(), r.colour)
}
