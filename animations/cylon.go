package animations

import "github.com/andew42/brightlight/controller"

type cylon struct {

	seg controller.Segment
	pos uint
	dir bool
}

func newCylon(seg controller.Segment) *cylon {

	return &cylon{seg, 0, false}
}

func (c *cylon) animateNextFrame() {

	c.seg.Set(c.pos, controller.NewRgb(0, 0, 0))

	if c.dir {
		if c.pos == c.seg.Len()-1 {
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

	c.seg.Set(c.pos, controller.NewRgb(255, 0, 0))
}
