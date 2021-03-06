package animations

import (
	"github.com/andew42/brightlight/segment"
)

// Takes an existing animation and repeats it every n LEDs
type repeater struct {
	animator     animator
	repeatLength uint
}

func newRepeater(animator animator, repeatLength uint) *repeater {

	var r repeater
	r.animator = animator
	r.repeatLength = repeatLength
	return &r
}

func (r *repeater) animateFrame(frameCount uint, frame segment.Segment) {

	// Number of repeats and remaining pixels
	repeat := frame.Len() / r.repeatLength
	// TODO HOW TO HANDLE REMAINDER (SEPARATE CLASSES OR REMAINDER MODE)
	//remainder := frame.Len() % r.repeatLength
	startOffset := uint(0)

	// Repeat the animation repeat times
	for i := uint(0); i < repeat; i++ {
		// Create a logical segment for this animation
		seg := segment.NewSubSegment(frame, startOffset+i*r.repeatLength, r.repeatLength)
		// Animate over the logical segment
		r.animator.animateFrame(frameCount, seg)
	}
}
