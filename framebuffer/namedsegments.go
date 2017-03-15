package framebuffer

import (
	"errors"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
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
			if length, err := strconv.Atoi(name[colonIndex+1:]); err == nil {
				return NewLogSegment(NewPhySegment(fb.Strips[stripIndex:stripIndex+1]), 0, uint(length)), nil
			}
		}
	}

	// Bedroom Ceiling
	t := NewPhySegment([]LedStrip{
		fb.Strips[11], fb.Strips[4], fb.Strips[3], fb.Strips[7], fb.Strips[6], fb.Strips[14], fb.Strips[13]})
	bedroomCeiling := NewLogSegment(t, 135, t.Len()-135)

	// Bedroom Wall
	t = NewPhySegment([]LedStrip{
		fb.Strips[8], fb.Strips[2], fb.Strips[5], fb.Strips[15], fb.Strips[10]})
	bedroomWall := NewLogSegment(t, 144, t.Len()-(144+241))

	// Bathroom Ceiling
	t = NewPhySegment([]LedStrip{fb.Strips[9], fb.Strips[11]})
	bathroomCeiling := NewLogSegment(t, 0, t.Len()-27)

	// Bathroom Wall
	t = NewPhySegment([]LedStrip{fb.Strips[10], fb.Strips[8]})
	bathroomWall := NewLogSegment(t, 50, t.Len()-(28+50))

	switch name {
	case "All":
		// All physical strips, ceiling then wall
		return NewPhySegment([]LedStrip{
			fb.Strips[3], fb.Strips[7], fb.Strips[6], fb.Strips[14], fb.Strips[13], fb.Strips[9], fb.Strips[11], fb.Strips[4],
			fb.Strips[5], fb.Strips[15], fb.Strips[10], fb.Strips[8], fb.Strips[2]}),
			nil
	case "All Ceiling":
		// Ceiling Strip starting at bed's corner
		return NewPhySegment([]LedStrip{
			fb.Strips[3], fb.Strips[7], fb.Strips[6], fb.Strips[14], fb.Strips[13], fb.Strips[9], fb.Strips[11], fb.Strips[4]}),
			nil
	case "All Wall":
		// Ceiling Strip starting at bed's corner
		return NewPhySegment([]LedStrip{
			fb.Strips[5], fb.Strips[15], fb.Strips[10], fb.Strips[8], fb.Strips[2]}),
			nil

	case "Bedroom":
		return NewCombinedSegment(bedroomCeiling, bedroomWall), nil
	case "Bedroom Ceiling":
		return bedroomCeiling, nil
	case "Bedroom Wall":
		return bedroomWall, nil

	case "Bathroom":
		return NewCombinedSegment(bathroomCeiling, bathroomWall), nil
	case "Bathroom Ceiling":
		return bathroomCeiling, nil
	case "Bathroom Wall":
		return bathroomWall, nil

	case "Curtains":
		// Strip above curtains
		return NewPhySegment([]LedStrip{fb.Strips[3], fb.Strips[7]}),
			nil

	case "Strip Three":
		// Useful for testing in virtual mode (first 20 leds)
		return NewLogSegment(NewPhySegment([]LedStrip{fb.Strips[3]}), 0, 20), nil
	case "Strip Five":
		// Useful for testing in virtual mode (first 19 leds)
		return NewLogSegment(NewPhySegment([]LedStrip{fb.Strips[5]}), 0, 19), nil

	default:
		log.WithField("name", name).Warn("Unknown named segment")
		return nil, errors.New("Unknown named segment")
	}
}
