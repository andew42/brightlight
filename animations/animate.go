package animations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
	"github.com/andew42/brightlight/stats"
	"strconv"
	"time"
	"errors"
	"github.com/andew42/brightlight/config"
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

func (params segmentParams) asRange(index int) (int, error) {

	if params == nil || len(params) <= index {
		return 0, errors.New("no parameter at index: " + strconv.Itoa(index))
	}
	p := params[index];
	if p.Type != "range" {
		return 0, errors.New("parameter type 'range' expected but '" + p.Type + "' provided")
	}
	if val, ok := p.Value.(float64); ok {
		return int(val), nil
	}
	return 0, errors.New("unexpected parameter type")
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

	case "Rainbow":
		var err error
		duration, err := seg.Params.asRange(0)
		var brightness int
		if err == nil {
			brightness, err = seg.Params.asRange(1)
		}
		if err == nil {
			*animators = append(*animators,
				segNameAndAnimator{seg.Name,
					newRainbow(time.Second*time.Duration(duration), brightness)})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animation parameter")
		}

	case "Sweet Shop":
		var err error
		duration, err := seg.Params.asRange(0)
		var brightness int
		if err == nil {
			brightness, err = seg.Params.asRange(1)
		}
		var minSaturation int
		if err == nil {
			minSaturation, err = seg.Params.asRange(2)
		}
		if err == nil {
			*animators = append(*animators, segNameAndAnimator{seg.Name,
				newSweetshop(config.FramePeriodMs*time.Duration(duration), brightness, minSaturation)})
		}
	case "Twinkle":
		*animators = append(*animators, segNameAndAnimator{seg.Name, newTwinkle()})

	case "Baby Bows":
		var err error
		length, err := seg.Params.asRange(0)
		var duration int
		if err == nil {
			duration, err = seg.Params.asRange(1)
		}
		var brightness int
		if err == nil {
			brightness, err = seg.Params.asRange(2)
		}
		if err == nil {
			*animators = append(*animators,
				segNameAndAnimator{seg.Name, newRepeater(
					newRainbow(time.Second*time.Duration(duration), brightness), uint(length))})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animation parameter")
		}

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
		var err error
		colour, err := seg.Params.asColour(0)
		var width int
		if err == nil {
			width, err = seg.Params.asRange(1)
		}
		var repeat int
		if err == nil {
			repeat, err = seg.Params.asRange(2)
		}
		if err == nil && (width < 1 || repeat < 1) {
			err = errors.New("bad width or repeat parameters")
		}
		if err == nil {
			*animators = append(*animators, segNameAndAnimator{seg.Name, newRepeater(
				newBulb(colour, 0, uint(width)), uint(repeat))})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animation parameter")
		}

	case "Life":
		var err error
		colour, err := seg.Params.asColour(0)
		var duration int
		if err == nil {
			duration, err = seg.Params.asRange(1)
		}
		var rule int
		if err == nil {
			rule, err = seg.Params.asRange(2)
		}
		if err == nil {
			*animators = append(*animators, segNameAndAnimator{seg.Name,
				newLife(colour, uint(duration), rule)})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animation parameter")
		}

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
