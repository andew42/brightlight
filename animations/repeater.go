package animations

import "github.com/andew42/brightlight/framebuffer"

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

func (r *repeater) clone() animator {
	// TODO: DEEP COPY?
	panic("Not Implemented")
}

func (r *repeater) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	// Number of repeats and remaining pixels
	repeat := frame.Len() / r.repeatLength
	// TODO HOW TO HANDLE REMAINDER (SEPARATE CLASSES OR REMAINDER MODE)
	//remainder := frame.Len() % r.repeatLength
	startOffset := uint(0)

	// Repeat the animation repeat times
	for i := uint(0); i < repeat; i++ {
		// Create a logical segment for this animation
		seg := framebuffer.NewLogSegment(frame, startOffset+i*r.repeatLength, r.repeatLength)
		// Animate over the logical segment using a clone except for the last animation
		// this ensures we get the same result for each repeat but an update for last
		if i+1 == repeat {
			// Last in repeat sequence
			r.animator.animateNextFrame(frameCount, seg)
		} else {
			r.animator.clone().animateNextFrame(frameCount, seg)
		}
	}
}
