package animations

import "github.com/andew42/brightlight/controller"

type animator interface {
	animateBegin(seg controller.Segment)
	animateNextFrame()
}
