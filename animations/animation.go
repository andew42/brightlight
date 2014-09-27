package animations

import "github.com/andew42/brightlight/controller"

type animation struct {
	animator animator
	Segment  controller.Segment
}

func newAnimation(animator animator, segment controller.Segment) animation {
	var a animation
	a.animator = animator
	a.Segment = segment
	return a
}
