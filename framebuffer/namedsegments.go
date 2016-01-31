package framebuffer

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"strings"
)

// The list of names segment used to build scene in the UI
// SHOULD MATCH UI LIST IN UI/CONFIG/STATIC.JS
// A named segment is a single LogSegment (logical segment)
// A LogSegment consists of one or more segments (logical or physical) plus a start and end led
// A PhySegment (physical segment) consists of one or more led strips
// A LedStrip represents the physical strip of LEDS connected to a particular controller pin

// Construct a named segment
func GetNamedSegment(fb *FrameBuffer, name string) (Segment, error) {

	// Handle special physical segment pxx:yy segment index : length
	if len(name) > 3 && name[0] == 'p' && strings.IndexByte(name, ':') != -1 {
		colonIndex := strings.IndexByte(name, ':')
		if stripIndex, err := strconv.Atoi(name[1:colonIndex]); err == nil {
			if len, err := strconv.Atoi(name[colonIndex+1:]); err == nil {
				return NewSubSegment(NewPhySegment(fb.Strips[stripIndex:stripIndex+1]), 0, uint(len)), nil
			}
		}
	}

	switch name {
	case "All":
		// All physical strips, ceiling then wall
		return NewPhySegment([]LedStrip{
			fb.Strips[3], fb.Strips[7], fb.Strips[6], fb.Strips[14], fb.Strips[13], fb.Strips[9], fb.Strips[11], fb.Strips[4],
			fb.Strips[5], fb.Strips[15], fb.Strips[10], fb.Strips[8], fb.Strips[2]}),
		nil

	case "Ceiling":
		// Ceiling Strip starting at bed's corner
		return NewPhySegment([]LedStrip{
			fb.Strips[3], fb.Strips[7], fb.Strips[6], fb.Strips[14], fb.Strips[13], fb.Strips[9], fb.Strips[11], fb.Strips[4]}),
		nil

	case "Wall":
		// Ceiling Strip starting at bed's corner
		return NewPhySegment([]LedStrip{
			fb.Strips[5], fb.Strips[15], fb.Strips[10], fb.Strips[8], fb.Strips[2]}),
		nil

	case "Bedroom":
		// Crude representation of bedroom just controller 1 (TODO)
		return NewPhySegment([]LedStrip{
				fb.Strips[2], fb.Strips[3], fb.Strips[4], fb.Strips[5], fb.Strips[6], fb.Strips[7]}),
			nil

	case "Bathroom":
		// Crude representation of bathroom just controller 2 (TODO)
		return NewPhySegment([]LedStrip{
				fb.Strips[8], fb.Strips[9], fb.Strips[10], fb.Strips[11], fb.Strips[13], fb.Strips[14], fb.Strips[15]}),
			nil

	case "Curtains":
		// Strip above curtains
		return NewPhySegment([]LedStrip{fb.Strips[3], fb.Strips[7]}),
			nil

	case "Test 4":
		return NewPhySegment([]LedStrip{fb.Strips[4]}),
			nil
	case "Test 5":
		return NewPhySegment([]LedStrip{fb.Strips[5]}),
			nil
	default:
		log.WithField("name", name).Warn("Unknown named segment")
		return nil, errors.New("Unknown named segment")
	}
}
