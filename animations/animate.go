package animations

import (
	"errors"
	"github.com/andew42/brightlight/controller"
	"time"
)

const frameRateMs = 20

var (
	allLeds                 controller.Segment
	animationDriverChan     chan []animator
	ErrInvalidAnimationName = errors.New("invalid animation name")
)

// High level request to play a static colour animation
func AnimateStaticColour(colour controller.Rgb) {
	animations := make([]animator, 0)
	animations = append(animations, newStaticColour(allLeds, colour))
	animationDriverChan <- animations
}

// High level request to play an animation from web UI
func Animate(animationName string) error {
	animations := make([]animator, 0)
	switch {
	case animationName == "runner":
		animations = append(animations, newRunner(allLeds, controller.NewRgb(0, 0, 255)))

	case animationName == "cylon":
		animations = append(animations, newStaticColour(allLeds, controller.NewRgbFromInt(0)))
		animations = append(animations, newCylon(NewLogSegment(allLeds, 8, 20)))
		animations = append(animations, newCylon(NewLogSegment(allLeds, controller.MaxLedStripLen, 20)))

	case animationName == "rainbow":
		animations = append(animations, newRainbow(NewLogSegment(allLeds, 0, 28), 10000))

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
	// All 8 strips as a single long segment
	allLeds = controller.NewPhySegment(fb.Strips)
	animationDriverChan = make(chan []animator)
	go animateDriver(animationDriverChan, fb)
}

// The animation GO routine
func animateDriver(newAnimations chan []animator, fb *controller.FrameBuffer) {
	frameSync := time.Tick(frameRateMs * time.Millisecond)
	currentAnimations := make([]animator, 0)
	for {
		select {
		case <-frameSync:
			// Wait for a frame tick
			for _, value := range currentAnimations {
				value.animateNextFrame()
			}
			fb.Flush()

		case currentAnimations = <-newAnimations:
		}
	}
}
