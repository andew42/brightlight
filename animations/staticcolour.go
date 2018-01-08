package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
)

type staticColour struct {
	colour framebuffer.Rgb
}

func newStaticColour(colour framebuffer.Rgb) *staticColour {

	return &staticColour{colour}
}

func (sc *staticColour) animateFrame(frameCount uint, frame segment.Segment) {

	for s := uint(0); s < frame.Len(); s++ {
		frame.Set(s, sc.colour)
	}
}
