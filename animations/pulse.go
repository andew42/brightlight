package animations

import "github.com/andew42/brightlight/framebuffer"

type pulse struct {
	up  bool
	val byte
}

func newPulse() *pulse {

	return &pulse{}
}

func (p *pulse) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	if p.up {
		if p.val != 255 {
			p.val++
		} else {
			p.up = false
			p.val--
		}
	} else {
		if p.val != 0 {
			p.val--
		} else {
			p.up = true
			p.val++
		}
	}

	colour := framebuffer.NewRgb(p.val, p.val, p.val)
	for s := uint(0); s < frame.Len(); s++ {
		frame.Set(s, colour)
	}
}
