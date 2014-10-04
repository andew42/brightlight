package animations

import "github.com/andew42/brightlight/controller"

type runner struct {
	seg    controller.Segment
	colour controller.Rgb
	pos    uint
}

func newRunner(seg controller.Segment, colour controller.Rgb) *runner {
	var r runner
	r.seg = seg
	r.colour = colour
	return &r
}

func (r *runner) animateNextFrame() {
	r.seg.Set(r.pos, controller.NewRgb(0, 0, 0))
	r.pos = nextPos(r.pos, r.seg.Len())
	r.seg.Set(r.pos, r.colour)
}
