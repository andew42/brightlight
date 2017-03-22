package animations

import "github.com/andew42/brightlight/framebuffer"

type staticColour struct {
	colour framebuffer.Rgb
}

func newStaticColour(colour framebuffer.Rgb) *staticColour {

	return &staticColour{colour}
}

func (sc *staticColour) clone() animator {
	clone := *sc
	return &clone
}

func (sc *staticColour) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	for s := uint(0); s < frame.Len(); s++ {
		frame.Set(s, sc.colour)
	}
}
