package animations

import (
	"errors"
	"time"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/stats"
)

// Frame rate as duration 20ms -> 50Hz
const frameRate time.Duration = 20 * time.Millisecond

var (
	allLeds                 controller.Segment
	curtainLeds				controller.Segment
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
		animations = append(animations, newRainbow(curtainLeds, time.Second * 5))

	case animationName == "sweetshop":
		animations = append(animations, newSweetshop(allLeds, time.Second * 1))

	default:
		return ErrInvalidAnimationName
	}

	// Send the (possibly) new animation to driver
	animationDriverChan <- animations
	return nil
}

// Start animate driver
func StartDriver(fb *controller.FrameBuffer, statistics *stats.Stats) {
	if animationDriverChan != nil {
		panic("StartAnimateDriver called twice")
	}
	// All 8 strips as a single long segment
	// TODO: Make this a config file
	allLeds = controller.NewPhySegment(fb.Strips)
	x := make([]controller.LedStrip, 2)
	x[0] = fb.Strips[3]
	x[1] = fb.Strips[7]
	curtainLeds = controller.NewPhySegment(x)

	animationDriverChan = make(chan []animator)
	go animateDriver(animationDriverChan, fb, statistics)
}

// The animation GO routine
func animateDriver(newAnimations chan []animator, fb *controller.FrameBuffer, statistics *stats.Stats) {
	frameSync := time.Tick(frameRate)
	currentAnimations := make([]animator, 0)
	nextFrameTime := time.Now().Add(frameRate)
	for {
		select {
		case <-frameSync:
			// Wait for a frame tick
			started := time.Now()
			jitter := started.Sub(nextFrameTime)
			nextFrameTime = started.Add(frameRate)
			for _, value := range currentAnimations {
				value.animateNextFrame()
			}
			fb.Flush()
			statistics.AddAnimation(time.Since(started), jitter)

		case currentAnimations = <-newAnimations:
			statistics.Reset()
		}
	}
}
