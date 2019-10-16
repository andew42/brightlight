package animations

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
	"github.com/andew42/brightlight/stats"
	"strconv"
	"time"
)

// Button request to animate
type Button struct {
	Key      int
	Name     string
	Segments segments
}

// Animation action to perform on a segment (from UI)
type SegmentAction struct {
	Name      string
	Base      string
	Start     uint
	Length    uint
	Animation string
	Params    segmentParams
}

type segmentParam struct {
	Type  string
	Value interface{}
}

type segments []SegmentAction
type segmentParams []segmentParam
type propertyMap map[string]interface{}

func (params segmentParams) asColour(index int) (framebuffer.Rgb, error) {

	if params == nil || len(params) <= index {
		return framebuffer.Rgb{}, errors.New("no parameter at index: " + strconv.Itoa(index))
	}
	p := params[index]
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
	p := params[index]
	if p.Type != "range" {
		return 0, errors.New("parameter type 'range' expected but '" + p.Type + "' provided")
	}
	if val, ok := p.Value.(float64); ok {
		return int(val), nil
	}
	return 0, errors.New("unexpected parameter type")
}

func (params segmentParams) asCheckbox(index int) (bool, error) {

	if params == nil || len(params) <= index {
		return false, errors.New("no parameter at index: " + strconv.Itoa(index))
	}
	p := params[index]
	if p.Type != "checkbox" {
		return false, errors.New("parameter type 'checkbox' expected but '" + p.Type + "' provided")
	}
	if val, ok := p.Value.(bool); ok {
		return bool(val), nil
	}
	return false, errors.New("unexpected parameter type")
}

func (m *propertyMap) IsValidNumber(i string) bool {
	var v interface{}
	var ok bool
	if v, ok = (*m)[i]; !ok {
		return false
	}
	_, ok = v.(float64)
	return ok
}

var animationChanged = make(chan []SegmentAction)

// Internal type describes a named segment with it animator
type segActionAndAnimator struct {
	segAction SegmentAction
	animator  animator
}

// RunAnimations animates a bunch of segments supplied by a UI button press
func RunAnimations(segments []SegmentAction) {
	animationChanged <- segments
}

func buildAnimatorList(segments []SegmentAction) []segActionAndAnimator {
	// Build a slice of animators with segment names
	animators := make([]segActionAndAnimator, 0, 4)

	// Foreach supplied segment action
	for _, seg := range segments {
		appendAnimatorsForAction(&animators, seg)
	}
	return animators
}

// Append an animation specified as a string
func appendAnimatorsForAction(animators *[]segActionAndAnimator, seg SegmentAction) {

	switch seg.Animation {
	case "Off":
		// Currently used by hue animations to ignore a segment

	case "Static":
		if colour, err := seg.Params.asColour(0); err == nil {
			*animators = append(*animators,
				segActionAndAnimator{seg, newStaticColour(colour)})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animation parameter")
		}

	case "Runner":
		*animators = append(*animators, segActionAndAnimator{seg,
			newRunner(framebuffer.NewRgb(0, 0, 255))})

	case "Cylon":
		*animators = append(*animators, segActionAndAnimator{seg, newCylon()})

	case "Rainbow":
		var err error
		duration, err := seg.Params.asRange(0)
		var brightness int
		if err == nil {
			brightness, err = seg.Params.asRange(1)
		}
		if err == nil {
			*animators = append(*animators,
				segActionAndAnimator{seg,
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
			*animators = append(*animators, segActionAndAnimator{seg,
				newSweetshop(config.FramePeriodMs*time.Duration(duration), brightness, minSaturation)})
		}
	case "Twinkle":
		*animators = append(*animators, segActionAndAnimator{seg, newTwinkle()})

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
				segActionAndAnimator{seg, newRepeater(
					newRainbow(time.Second*time.Duration(duration), brightness), uint(length))})
		} else {
			log.WithFields(log.Fields{"params": seg.Params, "Error": err.Error()}).Warn("Bad animation parameter")
		}

	case "Christmas":
		*animators = append(*animators, segActionAndAnimator{seg, newRepeater(
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
		*animators = append(*animators, segActionAndAnimator{seg, newRepeater(
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
			*animators = append(*animators, segActionAndAnimator{seg, newRepeater(
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
		var autoRepeat bool
		if err == nil {
			autoRepeat, err = seg.Params.asCheckbox(3)
		}
		if err == nil {
			*animators = append(*animators, segActionAndAnimator{seg,
				newLife(colour, uint(duration), rule, autoRepeat)})
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
		segments = append(segments, SegmentAction{"All", "", 0, 0, "Static", nil}) // TODO

		// Special pXX:YY segment id to address physical strip XX of length YY
		// NOTE: if a strip is reverse direction then this may not show up on
		// the virtual display which shows only the FIRST 20 LEDs

		segId := "p" + strconv.Itoa(int(stripIndex)) + ":" + strconv.Itoa(int(stripLength))
		segments = append(segments, SegmentAction{segId, "", 0, 0, "Static", nil}) // TODO 808080
	} else {
		// Invalid request, light all LEDs red
		segments = append(segments, SegmentAction{"All", "", 0, 0, "Static", nil}) // TODO
	}

	// Perform the animation
	animationChanged <- segments
}

// Handle predefined named segments and user defined segments
func resolveSegment(sa segActionAndAnimator, fb *framebuffer.FrameBuffer) (segment.Segment, error) {

	// User defined segments specify a base segment (start and length)
	if len(sa.segAction.Base) != 0 {

		// Lookup base segment (which should be a predefined named segment)
		ns, err := segment.GetNamedSegment(sa.segAction.Base)
		if err != nil {
			return nil, err
		}
		baseSeg := ns.GetSegment(fb)

		// Validate start falls within the segment
		l := baseSeg.Len()
		if sa.segAction.Start >= l {
			return nil, errors.New("segment start not within base segment")
		}

		// Clip length if it exceeds end of base segment
		requestedLength := sa.segAction.Length
		if sa.segAction.Start+requestedLength > l {
			requestedLength = l - sa.segAction.Start
		}

		return segment.NewSubSegment(baseSeg, sa.segAction.Start, requestedLength), nil
	}
	// Here we are dealing with a predefined named segment
	ns, err := segment.GetNamedSegment(sa.segAction.Name)
	if err != nil {
		return nil, err
	}
	return ns.GetSegment(fb), nil
}

// Start animate driver
func StartDriver(renderer chan *framebuffer.FrameBuffer) {
	// Start the animator go routine
	go func() {
		// The animations in play from the UI (default all off)
		var animators = make([]segActionAndAnimator, 1)
		animators[0] = segActionAndAnimator{SegmentAction{Name: "All"}, newStaticColour(framebuffer.NewRgbFromInt(0))}
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
					if seg, err := resolveSegment(v, fb); err == nil {
						v.animator.animateFrame(frameCounter, seg)
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
