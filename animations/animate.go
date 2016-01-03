package animations

import (
	"errors"
	"time"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/stats"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

// Animation action to perform on a segment
type SegmentAction struct {
	SegmentId  string
	Action     string
	Params     string
}

// Frame rate as duration
// 20ms -> 50Hz
// * 25ms -> 40Hz
// 40ms -> 25Hz
// 50ms -> 20Hz
const frameRate time.Duration = 25 * time.Millisecond

var (
	// Definitions of named segments by segment id
	namedSegments map[string]segmentInfo

	// One segment for each physical strip (i.e. 8 per Teensy)
	physicalStrips []controller.Segment

	// Communicate next animation to driver
	animationDriverChan chan []animator

	// Bad parameters
	ErrInvalidParameter = errors.New("invalid parameters")
)

// Animate a bunch of segments supplied by a UI button press
func RunAnimations(segments []SegmentAction) {

	// Build a slice of animators
	animators := make([]animator, 1, 4)

	// Initial animation turns all lights off
	animators[0] = newStaticColour(namedSegments["s0"].Seg, controller.NewRgbFromInt(0))

	// Foreach supplied segment action
	for _, seg := range segments {

		// Lookup action's segment id
		segInfo, ok := namedSegments[seg.SegmentId]
		if !ok {
			log.WithField(seg.SegmentId, "SegmentId").Warn("Unknown segment id")
			continue
		}

		// Append animators for this segment action
		appendAnimatorsForAction(&animators, segInfo.Seg, seg.Action, seg.Params)
	}

	log.WithField("len(animators)", len(animators)).Info("RunAnimations")

	// Send the (possibly) new animation to driver
	animationDriverChan <- animators
}

// Append an animation specified as a string
func appendAnimatorsForAction(animators *[]animator, seg controller.Segment, action string, params string) {

	switch action {
	case "static":
		colour, err := strconv.ParseInt(params, 16, 32)
		if err != nil {
			log.WithFields(log.Fields{"params": params, "Error": err.Error()}).Warn("Bad animataion parameter")
		}
		*animators = append(*animators, newStaticColour(seg, controller.NewRgbFromInt(int(colour))))

	case "runner":
		*animators = append(*animators, newRunner(seg, controller.NewRgb(0, 0, 255)))

	case "cylon":
		*animators = append(*animators, newCylon(seg))

	case "rainbow": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, newRainbow(seg, time.Second*5))

	case "sweetshop": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, newSweetshop(seg, time.Second*1))

	case "candle": // TODO MAKE POSITION AND REPEAT PARAMETERS
		*animators = append(*animators, newCandle(NewLogSegment(seg, 8, 3)))

	case "christmas": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, newChristmas(seg, time.Second*1))

	default:
		log.WithField("action", action).Warn("Unknown animataion action")
	}
}

// High level request to light the first stripLength LEDs of a physical strip
func AnimateStripLength(stripIndex uint, stripLength uint) error {

	// Check request fits physical strip
	if stripIndex < uint(len(physicalStrips)) && stripLength <= physicalStrips[stripIndex].Len() {
		// Turn off all leds
		animations := make([]animator, 2)
		animations[0] = newStaticColour(namedSegments["s0"].Seg, controller.NewRgb(0, 0, 0))

		// Turn on test strip, if a strip is reverse direction then this may not
		// show up on the virtual display which shows only the FIRST 20 LEDs
		animations[1] = newStaticColour(
			controller.NewSubSegment(physicalStrips[stripIndex], 0, stripLength),
			controller.NewRgb(128, 128, 128))
		animationDriverChan <- animations
		return nil
	} else {
		// Turn all leds red (for error)
		animations := make([]animator, 1)
		animations[0] = newStaticColour(namedSegments["s0"].Seg, controller.NewRgb(128, 0, 0))
		animationDriverChan <- animations
		return ErrInvalidParameter
	}
}

// Start animate driver
func StartDriver(fb *controller.FrameBuffer, statistics *stats.Stats) {

	if animationDriverChan != nil {
		log.Panic("StartAnimateDriver called twice")
	}

	// Each frame buffer strip as its own segment
	physicalStrips = make([]controller.Segment, len(fb.Strips))
	for i, _ := range fb.Strips {
		physicalStrips[i] = controller.NewPhySegment(fb.Strips[i:i+1])
	}

	// Construct list of named segments
	namedSegments = NewNamedSegments(fb)

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
