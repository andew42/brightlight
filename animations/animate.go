package animations

import (
	"errors"
	"time"
	"github.com/andew42/brightlight/controller"
)

const frameRateMs = 20

var (
	allLeds                 controller.Segment
	animationDriverChan     chan []animation
	ErrInvalidAnimationName = errors.New("invalid animation name")
)

// High level request to play a static colour animation
func AnimateStaticColour(colour controller.Rgb) {
	animations := make([]animation, 0)
	animations = append(animations, newAnimation(newStaticColour(colour), allLeds))
	animationDriverChan <- animations
}

// High level request to play an animation from web UI
func Animate(animationName string) error {
	animations := make([]animation, 0)
	switch {
	case animationName == "runner":
		var r runner
		animations = append(animations, newAnimation(&r, allLeds))
	case animationName == "cylon":
		animations = append(animations, newAnimation(newStaticColour(controller.NewRgbFromInt(0)), allLeds))
		var c1 cylon
		seg := NewLogSegment(allLeds, 8, 20)
		animations = append(animations, newAnimation(&c1, seg))
		var c2 cylon
		seg = NewLogSegment(allLeds, 30, 20)
		animations = append(animations, newAnimation(&c2, seg))
	default:
		return ErrInvalidAnimationName
	}

	// Send the (possibly) new animation to driver
	animationDriverChan <- animations
	return nil
}

// Start animate driver
func StartDriver(fb *controller.FrameBuffer) {
	if animationDriverChan != nil {
		panic("StartAnimateDriver called twice")
	}
	allLeds = controller.NewPhySegment(fb.Strips)
	animationDriverChan = make(chan []animation)
	go animateDriver(animationDriverChan, fb)
}

// The animation GO routine
func animateDriver(newAnimations chan []animation, fb *controller.FrameBuffer) {
	frameSync := time.Tick(frameRateMs * time.Millisecond)
	currentAnimations := make([]animation, 0)
	for {
		select {
		case <-frameSync:
			// Wait for a frame tick
			for _, value := range currentAnimations {
				value.animator.animateNextFrame()
			}
			fb.Flush()
		case currentAnimations = <-newAnimations:
			// Wait for new animation
			for _, value := range currentAnimations {
				value.animator.animateBegin(value.Segment)
			}
		}
	}
}
