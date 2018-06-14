package animations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
	"github.com/andew42/brightlight/stats"
	"strconv"
	"time"
	"errors"
)

// Animation action to perform on a segment (from UI)
type SegmentAction struct {
	Name      string
	Animation string
	Params    segmentParams
}

type segmentParam struct {
	Type  string
	Value interface{}
}

type segmentParams []segmentParam
type propertyMap map[string]interface{}

func (params segmentParams) asColour(index int) (framebuffer.Rgb, error) {

	if params == nil || len(params) <= index {
		return framebuffer.Rgb{}, errors.New("no parameter at index: " + strconv.Itoa(index))
	}
	p := params[index];
	if p.Type != "colour" {
		return framebuffer.Rgb{}, errors.New("parameter type 'colour' expected but '" + p.Type + "' provided")
	}
	var valueMap propertyMap
	var ok bool
	if valueMap, ok = p.Value.(map[string]interface{}); !ok {
		return framebuffer.Rgb{}, errors.New("unexpected parameter type")
	}
	if !valueMap.IsValidNumber("r") {
		return framebuffer.Rgb{}, errors.New("missing or incorrectly typed red colour value")
	}
	if !valueMap.IsValidNumber("g") {
		return framebuffer.Rgb{}, errors.New("missing or incorrectly typed green colour value")
	}
	if !valueMap.IsValidNumber("b") {
		return framebuffer.Rgb{}, errors.New("missing or incorrectly typed blue colour value")
	}
	return framebuffer.NewRgb(
		byte(valueMap["r"].(float64)),
		byte(valueMap["g"].(float64)),
		byte(valueMap["b"].(float64))), nil
}

func (m *propertyMap) IsValidNumber(i string) bool {
	var v interface{}
	var ok bool
	if v, ok = (*m)[i]; !ok {
		return false;
	}
	_, ok = v.(float64);
	return ok
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
		if colour, err := seg.Params.asColour(0); err == nil {
			*animators = append(*animators,
				segNameAndAnimator{seg.Name, newStaticColour(colour)})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animation parameter")
		}

	case "Runner":
		*animators = append(*animators, segNameAndAnimator{seg.Name,
			newRunner(framebuffer.NewRgb(0, 0, 255))})

	case "Cylon":
		*animators = append(*animators, segNameAndAnimator{seg.Name, newCylon()})

	case "Rainbow": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, segNameAndAnimator{seg.Name, newRainbow(time.Second * 15)})

	case "Sweet Shop": // TODO MAKE TIME A PARAMETER
		*animators = append(*animators, segNameAndAnimator{seg.Name, newSweetshop(time.Second * 1)})

	case "Twinkle":
		*animators = append(*animators, segNameAndAnimator{seg.Name, newTwinkle()})

	case "BabyBows":
		*animators = append(*animators, segNameAndAnimator{seg.Name, newRepeater(
			newRainbow(time.Second*8), 15)})

	case "Christmas":
		*animators = append(*animators, segNameAndAnimator{seg.Name, newRepeater(
			newLinearFade(
				time.Duration(10000*time.Millisecond),
				false,
				newBulb(framebuffer.Rgb{255, 0, 0}, 0, 1),
				newBulb(framebuffer.Rgb{255, 122, 0}, 1, 1),
				newBulb(framebuffer.Rgb{0, 255, 0}, 2, 1),
				newBulb(framebuffer.Rgb{178, 0, 255}, 3, 1),
				newBulb(framebuffer.Rgb{0, 0, 255}, 4, 1)),
			5)})

	case "Fairground":
		*animators = append(*animators, segNameAndAnimator{seg.Name, newRepeater(
			newStepFade(
				time.Duration(640*time.Millisecond),
				false,
				newBulb(framebuffer.Rgb{175, 107, 1}, 0, 1),
				newBulb(framebuffer.Rgb{181, 145, 0}, 1, 1),
				newBulb(framebuffer.Rgb{175, 107, 1}, 2, 1),
				newBulb(framebuffer.Rgb{181, 145, 0}, 3, 1)),
			4)})

	case "Discrete":
		*animators = append(*animators, segNameAndAnimator{seg.Name, newRepeater(
			newBulb(framebuffer.NewRgb(255, 255, 255), 0, 1), 15)})

	default:
		log.WithField("action", seg.Animation).Warn("Unknown animation action")
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
		segments = append(segments, SegmentAction{"All", "Static", nil}) // TODO

		// Special pXX:YY segment id to address physical strip XX of length YY
		// NOTE: if a strip is reverse direction then this may not show up on
		// the virtual display which shows only the FIRST 20 LEDs

		segId := "p" + strconv.Itoa(int(stripIndex)) + ":" + strconv.Itoa(int(stripLength))
		segments = append(segments, SegmentAction{segId, "Static", nil}) // TODO 808080
	} else {
		// Invalid request, light all LEDs red
		segments = append(segments, SegmentAction{"All", "Static", nil}) // TODO
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
					if seg, err := segment.GetNamedSegment(v.namedSegment); err == nil {
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
