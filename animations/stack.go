package animations

import (
	"github.com/andew42/brightlight/segment"
)

// Takes existing animations and stacks them left to right
// Space divide equally with any remaining space left empty
type stack struct {
	animators []animator
}

func newStack(args ...animator) *stack {

	var s stack
	s.animators = args
	return &s
}

func (s *stack) animateFrame(frameCount uint, frame segment.Segment) {

	// Number of pixels per animation
	segLength := frame.Len() / uint(len(s.animators))

	// Repeat the animation repeat times
	for i := uint(0); i < uint(len(s.animators)); i++ {
		// Create a logical segment for this animation
		seg := segment.NewSubSegment(frame, i*segLength, segLength)
		// Animate over the logical segment
		s.animators[i].animateFrame(frameCount, seg)
	}
}
