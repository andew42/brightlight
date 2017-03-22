package animations

import "github.com/andew42/brightlight/framebuffer"

type cylon struct {
	pos uint
	dir bool
}

func newCylon() *cylon {

	return &cylon{0, false}
}

func (c *cylon) clone() animator {
	clone := *c
	return &clone
}

func (c *cylon) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	frame.Set(c.pos, framebuffer.NewRgb(0, 0, 0))

	if c.dir {
		if c.pos == frame.Len()-1 {
			c.dir = false
		} else {
			c.pos++
		}
	} else {
		if c.pos == 0 {
			c.dir = true
		} else {
			c.pos--
		}
	}

	frame.Set(c.pos, framebuffer.NewRgb(255, 0, 0))
}
