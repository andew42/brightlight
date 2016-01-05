package animations

import "github.com/andew42/brightlight/controller"

type staticColour struct {
	seg    controller.Segment
	colour controller.Rgb
	done   bool
}

func newStaticColour(seg controller.Segment, colour controller.Rgb) *staticColour {

	return &staticColour{seg, colour, false}
}

func (sc *staticColour) animateNextFrame() {

	if !sc.done {
		for s := uint(0); s < sc.seg.Len(); s++ {
			sc.seg.Set(s, sc.colour)
		}
		sc.done = true
	}
}
