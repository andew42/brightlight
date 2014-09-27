package animations

import "github.com/andew42/brightlight/controller"

type runner struct {
	seg controller.Segment
	pos uint
}

func (r *runner) animateBegin(seg controller.Segment) {
	r.seg = seg
	r.pos = 0
}

func (r *runner) animateNextFrame() {
	r.seg.Set(r.pos, controller.NewRgb(0, 0, 0)) // OFF then next RED
	r.pos = nextPos(r.pos, r.seg.Len())
	r.seg.Set(r.pos, controller.NewRgb(255, 0, 0)) // OFF then next RED
}
