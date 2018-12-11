package animations

import (
	"time"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/segment"
	log "github.com/Sirupsen/logrus"
)

type linearFade struct {
	animators       []animator
	framesPerPeriod uint
	reverseOnRepeat bool
	forward         bool
}

// Fade between two or more animations using linear interpolation, if more than
// two animations are provided then a chain is formed 1 -> 2 -> 3 once the last
// animation as played the sequence it either starts again or is reversed 3 ->
// 2 -> 1 (depending on reverseOnRepeat). The animation always repeats forever.
func newLinearFade(period time.Duration, reverseOnRepeat bool, animators ...animator) *linearFade {

	// Must be at least two animations
	if len(animators) < 2 {
		log.WithField("animators", len(animators)).
			Fatal("LinearFade animation requires at least two animations")
	}

	// Calculate how many frames require to complete animation from -> to
	framesPerPeriod := uint(float32(period) / float32(config.FramePeriodMs))
	if framesPerPeriod <= 0 {
		// Must be at least two steps
		framesPerPeriod = 2
	}

	return &linearFade{animators, framesPerPeriod, reverseOnRepeat, !reverseOnRepeat}
}

// Assume two animations with 3 frames per period
// Frame 0 from = 100% to = 0%
// Frame 1 from = 50% to = 50%
// Frame 2 from = 0% to = 100%
func (lf *linearFade) animateFrame(frameCount uint, frame segment.Segment) {

	// How far into the entire animation chain (all segments) are we?
	index := frameCount % lf.framesPerPeriod

	// Determine direction for this frame, this is done in a pure deterministic
	// way based only on frame count as the animation may be reused multiple
	// times with the same frame count if it is embedded in a repeater animation
	if lf.reverseOnRepeat {
		lf.forward = (frameCount/lf.framesPerPeriod)%2 == 0
	} else {
		lf.forward = true
	}

	// Flip the index if we are moving backwards
	if !lf.forward {
		index = (lf.framesPerPeriod - index) - 1
	}

	// Segment count per period depends on the reverseOnRepeat value
	// Consider an animation with 3 segments with reverseOnRepeat:
	// true  - 0->1->2->1->0 i.e. two segments forward (then two segments reversed)
	// false - 0->1->2->0    i.e. three segments forward
	segmentsPerPeriod := len(lf.animators)
	if lf.reverseOnRepeat {
		segmentsPerPeriod--
	}

	// Frames per segment in the chain of animations
	framesPerSegment := float32(lf.framesPerPeriod) / float32(segmentsPerPeriod)

	// Determine current from -> to animations
	animationIndex := float32(index) / framesPerSegment
	from := lf.animators[uint(animationIndex)]
	to := lf.animators[uint(animationIndex+1)%uint(len(lf.animators))]

	// Work out what percentage of from and to animations to show
	segmentIndex := float32(index) - float32(uint(animationIndex))*framesPerSegment
	toPercent := segmentIndex / framesPerSegment
	fromPercent := 1 - toPercent

	// Make a copy of the frame buffer to a temporary in memory buffer
	tmpSegment := segment.NewMemSegment(frame.Len())
	for i := uint(0); i < frame.Len(); i++ {
		tmpSegment.Set(i, frame.Get(i))
	}

	// Apply the first animation to the frame buffer and scale with fromPercent
	from.animateFrame(frameCount, frame)
	if fromPercent == 100 {
		return
	}
	scaleFrameBuffer(frame, fromPercent)

	// Apply the second animation to the temporary in memory buffer and scale with toPercent
	to.animateFrame(frameCount, tmpSegment)
	scaleFrameBuffer(tmpSegment, toPercent)

	// Finally add the tmp segment into the frame buffer, we should never
	// overflow a byte unless rounding causes a problem...
	for i := uint(0); i < frame.Len(); i++ {
		frame.Set(i, frame.Get(i).Add(tmpSegment.Get(i)))
	}
}

// Scale each colour in the segment by f between 0 -> 1
func scaleFrameBuffer(segment segment.Segment, f float32) {

	for i := uint(0); i < segment.Len(); i++ {
		segment.Set(i, segment.Get(i).ScaleRgb(f))
	}
}
