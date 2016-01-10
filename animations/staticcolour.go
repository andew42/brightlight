package animations

import "github.com/andew42/brightlight/framebuffer"

type staticColour struct {
	colour framebuffer.Rgb
}

func newStaticColour(colour framebuffer.Rgb) *staticColour {

	return &staticColour{colour}
}

func (sc *staticColour) animateNextFrame(seg framebuffer.Segment) {

	for s := uint(0); s < seg.Len(); s++ {
		seg.Set(s, sc.colour)
	}
}
