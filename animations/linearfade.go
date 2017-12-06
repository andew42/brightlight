package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"time"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/segments"
)

type linearFade struct {
	from            animator
	to              animator
	framesPerPeriod uint
}

// TODO Take a slice of animators
// TODO Reverse animation when it hits end

// Fade between two animations using linear interpolation
func newLinearFade(from animator, to animator, period time.Duration) *linearFade {

	// Calculate how many frames require to complete animation from -> to
	framesPerPeriod := uint(float32(period) / float32(config.FramePeriodMs))
	if framesPerPeriod <= 0 {
		// Must be at least two steps
		framesPerPeriod = 2
	}
	return &linearFade{from, to, framesPerPeriod}
}

// Assume 3 frames per period
// Frame 0 from = 100% to = 0%
// Frame 1 from = 50% to = 50%
// Frame 2 from = 0% to = 100%
func (lf *linearFade) animateFrame(frameCount uint, frame framebuffer.Segment) {

	// How far into this animation are we?
	index := frameCount % lf.framesPerPeriod

	// Work out what percentage of from and to animations to show
	toPercent := float32(index) / float32(lf.framesPerPeriod-1)
	fromPercent := 1 - toPercent

	// Apply the first animation to the frame buffer and scale with fromPercent
	lf.from.animateFrame(frameCount, frame)
	if fromPercent == 100 {
		return
	}
	scaleFrameBuffer(frame, fromPercent)

	// Apply the second animation to a temporary in memory buffer and scale with toPercent
	tmpSegment := segments.NewMemSegment(frame.Len())
	lf.to.animateFrame(frameCount, tmpSegment)
	scaleFrameBuffer(tmpSegment, toPercent)

	// Finally add the tmp segment into the frame buffer, we should never
	// overflow a byte unless rounding causes a problem...
	for i := uint(0); i < frame.Len(); i++ {
		frame.Set(i, frame.Get(i).Add(tmpSegment.Get(i)))
	}
}

// Scale each colour in the segment by f between 0 -> 1
func scaleFrameBuffer(segment framebuffer.Segment, f float32) {

	for i := uint(0); i < segment.Len(); i++ {
		segment.Set(i, segment.Get(i).ScaleRgb(f))
	}
}
