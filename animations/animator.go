package animations

import "github.com/andew42/brightlight/framebuffer"

type animator interface {
	animateNextFrame(seg framebuffer.Segment)
}
