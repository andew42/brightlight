package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"time"
	"github.com/andew42/brightlight/config"
	log "github.com/Sirupsen/logrus"
)

type stepFade struct {
	animators       []animator
	framesPerPeriod uint
	reverseOnRepeat bool
	forward         bool
}

// Steps between two or more animations. Once the last animation has played it either starts
// again or is reversed (depending on reverseOnRepeat). The animation always repeats forever.
func newStepFade(period time.Duration, reverseOnRepeat bool, animators ...animator) *stepFade {

	// Must be at least two animations
	if len(animators) < 2 {
		log.WithField("animators", len(animators)).
			Fatal("StepFade animation requires at least two animations")
	}

	// Calculate how many frames require to complete entire animation
	framesPerPeriod := uint(float32(period) / float32(config.FramePeriodMs))
	if framesPerPeriod <= 0 {
		// Must be at least two steps
		framesPerPeriod = 2
	}

	return &stepFade{animators, framesPerPeriod, reverseOnRepeat, !reverseOnRepeat}
}

// Animate current frame
func (sf *stepFade) animateFrame(frameCount uint, frame framebuffer.Segment) {

	// How far into the entire animation chain (all segments) are we?
	index := frameCount % sf.framesPerPeriod

	// Determine direction for this frame, this is done in a pure deterministic
	// way based only on frame count as the animation may be reused multiple
	// times with the same frame count if it is embedded in a repeater animation
	if sf.reverseOnRepeat {
		sf.forward = (frameCount/sf.framesPerPeriod)%2 == 0
	} else {
		sf.forward = true
	}

	// Flip the index if we are moving backwards
	if !sf.forward {
		index = (sf.framesPerPeriod - index) - 1
	}

	// Frames per segment in the chain of animations
	framesPerSegment := float32(sf.framesPerPeriod) / float32(len(sf.animators))

	// Determine current animation
	animationIndex := float32(index) / framesPerSegment

	// Apply the animation to the frame buffer
	sf.animators[uint(animationIndex)].animateFrame(frameCount, frame)
}
