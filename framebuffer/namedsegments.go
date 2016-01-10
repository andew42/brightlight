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

// Un-marshal JSON into typed slice
type SegmentInfo struct {
	Name string
	Seg  Segment
}

// Construct a named segment
func GetNamedSegment(fb *FrameBuffer, name string) (si SegmentInfo, err error) {
	err = nil

	// Handle special physical segment pxx:yy segment index : length
	if len(name) > 3 && name[0] == 'p' && strings.IndexByte(name, ':') != -1 {
		colonIndex := strings.IndexByte(name, ':')
		if stripIndex, err := strconv.Atoi(name[1:colonIndex]); err == nil {
			if len, err := strconv.Atoi(name[colonIndex+1:]); err == nil {
				return SegmentInfo{name, NewSubSegment(
						NewPhySegment(fb.Strips[stripIndex:stripIndex+1]), 0, uint(len))},
					nil
			}
		}
	}

	switch name {
	case "s0":
		// All physical strips in no particular order
		si = SegmentInfo{"All", NewPhySegment(fb.Strips)}
	case "s1":
		// Crude representation of bedroom just controller 1 (TODO)
		si = SegmentInfo{"Bedroom", NewPhySegment([]LedStrip{
			fb.Strips[2], fb.Strips[3], fb.Strips[4], fb.Strips[5], fb.Strips[6], fb.Strips[7]})}
	case "s2":
		// Crude representation of bathroom just controller 2 (TODO)
		si = SegmentInfo{"Bathroom", NewPhySegment([]LedStrip{
			fb.Strips[8], fb.Strips[9], fb.Strips[10], fb.Strips[11], fb.Strips[13], fb.Strips[14], fb.Strips[15]})}
	case "s3":
		// Strip above curtains
		si = SegmentInfo{"Curtains", NewPhySegment([]LedStrip{fb.Strips[3], fb.Strips[7]})}
	case "s4":
		si = SegmentInfo{"Test 4", NewPhySegment([]LedStrip{fb.Strips[4]})}
	case "s5":
		si = SegmentInfo{"Test 5", NewPhySegment([]LedStrip{fb.Strips[5]})}
	default:
		log.WithField("name", name).Warn("Unknown named segment")
		err = errors.New("Unknown named segment")
	}
	return
}
