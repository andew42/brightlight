package animations

import "github.com/andew42/brightlight/controller"

type cylon struct {
	seg controller.Segment
	pos uint
	dir bool
}

func (r *cylon) animateBegin(seg controller.Segment) {
	r.seg = seg
	r.pos = 0
}

func (r *cylon) animateNextFrame() {
	r.seg.Set(r.pos, controller.NewRgb(0, 0, 0))

	if r.dir {
		if r.pos == r.seg.Len() - 1 {
			r.dir = false
		} else {
			r.pos++
		}
	} else {
		if r.pos == 0 {
			r.dir = true
		} else {
			r.pos--
		}
	}

	r.seg.Set(r.pos, controller.NewRgb(255, 0, 0))
}
