package animations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segments"
	"github.com/andew42/brightlight/stats"
	"strconv"
	"time"
)

// Animation action to perform on a segment (from UI)
type SegmentAction struct {
	Segment   string
	Animation string
	Params    string
}

var animationChanged = make(chan []SegmentAction)

// Internal type describes a named segment with it animator
type segNameAndAnimator struct {
	namedSegment string
	animator     animator
}

// RunAnimations animates a bunch of segments supplied by a UI button press
func RunAnimations(segments []SegmentAction) {
	animationChanged <- segments
}

func buildAnimatorList(segments []SegmentAction) []segNameAndAnimator {
	// Build a slice of animators with segment names
	animators := make([]segNameAndAnimator, 0, 4)

	// Foreach supplied segment action
	for _, seg := range segments {
		appendAnimatorsForAction(&animators, seg)
	}
	return animators
}

// Append an animation specified as a string
func appendAnimatorsForAction(animators *[]segNameAndAnimator, seg SegmentAction) {

	switch seg.Animation {
	case "Off":
		// Currently used by hue animations to ignore a segment

	case "Static":
		if colour, err := strconv.ParseInt(seg.Params, 16, 32); err == nil {
			*animators = append(*animators, segNameAndAnimator{seg.Segment,
				newStaticColour(framebuffer.NewRgbFromInt(int(colour)))})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animataion parameter")
		}

	case "Runner":
		*animators = append(*animators, segNameAndAnimator{seg.Segment,
			newRunner(framebuffer.NewRgb(0, 0, 255))})

	case "Cylon":
		*animators = append(*animators, segNameAndAnimator{seg.Segment, newCylon()})

	case "Rainbow": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, segNameAndAnimator{seg.Segment, newRainbow(time.Second * 15)})

	case "Sweet Shop": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, segNameAndAnimator{seg.Segment, newSweetshop(time.Second * 1)})

	case "Twinkle":
		*animators = append(*animators, segNameAndAnimator{seg.Segment, newTwinkle()})

	case "BabyBows":
		*animators = append(*animators, segNameAndAnimator{seg.Segment, newRepeater(
			newRainbow(time.Second*8), 15)})

	default:
		log.WithField("action", seg.Animation).Warn("Unknown animataion action")
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
		segments = append(segments, SegmentAction{"All", "Static", "0"})

		// Special pXX:YY segment id to address physical strip XX of length YY
		// NOTE: if a strip is reverse direction then this may not show up on
		// the virtual display which shows only the FIRST 20 LEDs

		segId := "p" + strconv.Itoa(int(stripIndex)) + ":" + strconv.Itoa(int(stripLength))
		segments = append(segments, SegmentAction{segId, "Static", "808080"})
	} else {
		// Invalid request, light all LEDs red
		segments = append(segments, SegmentAction{"All", "Static", "800000"})
	}

	// Perform the animation
	animationChanged <- segments
}

// Start animate driver new version
func StartDriver(renderer chan *framebuffer.FrameBuffer) {
	// Start the animator go routine
	go func() {
		// The animations in play from the UI (default all off)
		var animators = make([]segNameAndAnimator, 1)
		animators[0] = segNameAndAnimator{"All", newStaticColour(framebuffer.NewRgbFromInt(0))}
		frameCounter := uint(0)
		for {
			select {
			// Request to render a frame buffer
			case fb := <-renderer:
				renderStartTime := time.Now()

				fb = framebuffer.NewFrameBuffer()

				// Animate and return updated frame buffer
				for _, v := range animators {
					// Resolve the segment to animate, based on string name
					if seg, err := segments.GetNamedSegment(v.namedSegment); err == nil {
						v.animator.animateFrame(frameCounter, seg.GetSegment(fb))
					}
				}

				// Report render time and send buffer
				stats.AddFrameRenderTimeSample(time.Since(renderStartTime))
				renderer <- fb
				frameCounter++

			// Request animation update
			case currentAnimations := <-animationChanged:
				animators = buildAnimatorList(currentAnimations)
				frameCounter = 0
			}
		}
	}()
}
