package animations

import "github.com/andew42/brightlight/controller"

type staticColour struct {
	colour controller.Rgb
}

func newStaticColour(colour controller.Rgb) *staticColour {
	var sc staticColour
	sc.colour = colour
	return &sc
}

func (r *staticColour) animateBegin(seg controller.Segment) {
	for s := uint(0); s < seg.Len(); s++ {
		seg.Set(s, r.colour)
	}
}

func (r *staticColour) animateNextFrame() {
}
