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
	// All physical strips in a single segment
	allLeds                 controller.Segment

	// One segment for each physical strip (i.e. 8 per Teensy)
	physicalStrips            []controller.Segment

	// A segment representing the strips above the curtain
	curtainLeds                controller.Segment

	// Communicate next animation to driver
	animationDriverChan     chan []animator

	// Unknown animation error
	ErrInvalidAnimationName = errors.New("invalid animation name")
)

// High level request to play a static colour animation
func AnimateStaticColour(colour controller.Rgb) {
	animationDriverChan <- []animator{newStaticColour(allLeds, colour)}
}

// High level request to light the first stripLength LEDs of a physical strip
func AnimateStripLength(stripIndex uint, stripLength uint) {

	animations := make([]animator, 1, 2)

	// Turn off all leds
	animations[0] = newStaticColour(allLeds, controller.NewRgb(0, 0, 0))

	// Check request fits physical strip
	if stripIndex < uint(len(physicalStrips)) && stripLength <= physicalStrips[stripIndex].Len() {
		// Turn on test strip, if a strip is revers direction then this may not
		// show up on the virtual display which shows only the FIRST 20 LEDs
		animations = append(animations, newStaticColour(
				controller.NewSubSegment(physicalStrips[stripIndex], 0, stripLength),
				controller.NewRgb(128, 128, 128)))
	}

	animationDriverChan <- animations
}

// High level request to play an animation from web UI
func Animate(animationName string) error {

	animations := make([]animator, 0, 1)
	switch {
	case animationName == "runner":
		animations = append(animations, newRunner(allLeds, controller.NewRgb(0, 0, 255)))

	case animationName == "cylon":
		animations = append(animations, newStaticColour(allLeds, controller.NewRgbFromInt(0)))
		animations = append(animations, newCylon(NewLogSegment(allLeds, 8, 20)))
		animations = append(animations, newCylon(NewLogSegment(allLeds, controller.MaxLedStripLen, 20)))

	case animationName == "rainbow":
		animations = append(animations, newRainbow(curtainLeds, time.Second*5))

	case animationName == "sweetshop":
		animations = append(animations, newSweetshop(allLeds, time.Second*1))

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

	// TODO: Make this a config file
	// All frame buffer strips as a single long segment
	allLeds = controller.NewPhySegment(fb.Strips)

	// Each frame buffer strip as its own segment
	physicalStrips = make([]controller.Segment, len(fb.Strips))
	for i, _ := range fb.Strips {
		physicalStrips[i] = controller.NewPhySegment(fb.Strips[i:i+1])
	}

	// Two physical strips above curtains
	x := make([]controller.LedStrip, 2)
	x[0] = fb.Strips[3]
	x[1] = fb.Strips[7]
	curtainLeds = controller.NewPhySegment(x)

	// Start the animator go routine
	animationDriverChan = make(chan []animator)
	go animateDriver(animationDriverChan, fb, statistics)
}

// The animation go routine
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
