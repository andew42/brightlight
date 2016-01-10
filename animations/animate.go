package animations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/framebuffer"
	"strconv"
	"time"
)

// Animation action to perform on a segment (from UI)
type SegmentAction struct {
	SegmentId string
	Action    string
	Params    string
}

var animationChanged = make(chan []SegmentAction)

// Internal type describes a named segment with it animator
type segNameAndAnimator struct {
	namedSegment string
	animator     animator
}

// Animate a bunch of segments supplied by a UI button press
func RunAnimations(segments []SegmentAction) {
	animationChanged <- segments
}

func buildAnimatorList(segments []SegmentAction) []segNameAndAnimator {
	// Build a slice of animators with segment names
	animators := make([]segNameAndAnimator, 1, 4)

	// Initial animation turns all (s0) lights off
	animators[0] = segNameAndAnimator{"s0", newStaticColour(framebuffer.NewRgbFromInt(0))}

	// Foreach supplied segment action
	for _, seg := range segments {
		appendAnimatorsForAction(&animators, seg)
	}
	return animators
}

// Append an animation specified as a string
func appendAnimatorsForAction(animators *[]segNameAndAnimator, seg SegmentAction) {

	switch seg.Action {
	case "static":
		if colour, err := strconv.ParseInt(seg.Params, 16, 32); err == nil {
			*animators = append(*animators, segNameAndAnimator{seg.SegmentId,
				newStaticColour(framebuffer.NewRgbFromInt(int(colour)))})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animataion parameter")
		}

	case "runner":
		*animators = append(*animators, segNameAndAnimator{seg.SegmentId,
			newRunner(framebuffer.NewRgb(0, 0, 255))})

	case "cylon":
		*animators = append(*animators, segNameAndAnimator{seg.SegmentId, newCylon()})

	case "rainbow": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, segNameAndAnimator{seg.SegmentId, newRainbow(time.Second * 5)})

	case "sweetshop": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, segNameAndAnimator{seg.SegmentId, newSweetshop(time.Second * 1)})

	case "candle": // TODO MAKE POSITION AND REPEAT PARAMETERS
		*animators = append(*animators, segNameAndAnimator{seg.SegmentId, newCandle()})

	case "christmas": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, segNameAndAnimator{seg.SegmentId, newChristmas(time.Second * 1)})

	default:
		log.WithField("action", seg.Action).Warn("Unknown animataion action")
	}
}

// High level request to light the first stripLength LEDs of a physical strip
func AnimateStripLength(stripIndex uint, stripLength uint) {
	// Create a new frame buffer as a model for physical strip lengths
	fb := framebuffer.NewFrameBuffer()

	// Check request is for a valid strip and length
	segments := make([]SegmentAction, 0)
	if stripIndex < uint(len(fb.Strips)) && stripLength <= uint(len(fb.Strips[stripIndex].Leds)) {
		// Clear all lights
		segments = append(segments, SegmentAction{"s0", "static", "0"})

		// Special pXX:YY segment id to address physical strip XX of length YY
		// NOTE: if a strip is reverse direction then this may not show up on
		// the virtual display which shows only the FIRST 20 LEDs

		segId := "p" + strconv.Itoa(int(stripIndex)) + ":" + strconv.Itoa(int(stripLength))
		segments = append(segments, SegmentAction{segId, "static", "808080"})
	} else {
		// Invalid request, light all LEDS red
		segments = append(segments, SegmentAction{"s0", "static", "800000"})
	}

	// Perform the animation
	animationChanged <- segments
}

// Start animate driver new version
func StartDriver(renderer chan *framebuffer.FrameBuffer) {
	// Start the animator go routine
	go func() {
		// The animations in play from the UI (default all off)
		var animators []segNameAndAnimator = make([]segNameAndAnimator, 1)
		animators[0] = segNameAndAnimator{"s0", newStaticColour(framebuffer.NewRgbFromInt(0))}
		for {
			select {
			// Request to render a frame buffer
			case fb := <-renderer:
				// Animate and return update frame buffer
				for _, v := range animators {
					// Resolve the segment to animate, based on string name, and animate it
					if seg, err := framebuffer.GetNamedSegment(fb, v.namedSegment); err == nil {
						v.animator.animateNextFrame(seg.Seg)
					}
				}
				renderer <- fb

			// Request animation update
			case currentAnimations := <-animationChanged:
				animators = buildAnimatorList(currentAnimations)
			}
		}
	}()
}
