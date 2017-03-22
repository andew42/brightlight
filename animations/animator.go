package animations

import "github.com/andew42/brightlight/framebuffer"

type animator interface {
	// clone returns a deep copy of animator (used by repeater)
	clone() animator

	// animateNextFrame fills frame with next frame
	animateNextFrame(frameCount int, frame framebuffer.Segment)
}
