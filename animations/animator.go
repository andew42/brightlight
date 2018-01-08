package animations

import (
	"github.com/andew42/brightlight/segment"
)

type animator interface {
	// animateFrame fills frame with content for frameCount
	// frame is initially clear (or contains contents of a
	// lower layer animation (i.e. no need to initially clear)
	animateFrame(frameCount uint, frame segment.Segment)
}
