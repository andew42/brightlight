package segments

import (
	"errors"
	"github.com/andew42/brightlight/framebuffer"
	"strconv"
	"strings"
)

// TODO: MOVE OTHER SEGMENTS TO THIS NAME SPACE

// Named segments are used to build scenes in the UI
// TODO: SHOULD MATCH UI LIST IN UI/CONFIG/STATIC.JS
// A LogSegment (logical segment) consists of one or more segments (logical or physical)
// plus a start and end led position
// A PhySegment (physical segment) consists of one or more led strips
// LedStrip represents a physical strip of LEDS connected to a particular controller pin

type NamedSegment struct {
	Name       string
	GetSegment func(fb *framebuffer.FrameBuffer) framebuffer.Segment
}

// Return a particular named segment
func GetNamedSegment(name string) (NamedSegment, error) {

	// Handle special physical segment pxx:yy segment index : length
	// This feels a little hacky, its used by the strip length ui
	if len(name) > 3 && name[0] == 'p' && strings.IndexByte(name, ':') != -1 {
		colonIndex := strings.IndexByte(name, ':')
		if stripIndex, err := strconv.Atoi(name[1:colonIndex]); err == nil {
			if length, err := strconv.Atoi(name[colonIndex+1:]); err == nil {
				return NamedSegment{
					Name: name,
					GetSegment: func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
						return framebuffer.NewLogSegment(
							framebuffer.NewPhySegment(fb.Strips[stripIndex:stripIndex+1]),
							0, uint(length))
					},
				}, nil
			}
		}
	}

	// Search name segment collection
	for i := range allNamedSegments {
		if allNamedSegments[i].Name == name {
			return allNamedSegments[i], nil
		}
	}

	// Error if named segment not found
	return NamedSegment{}, errors.New("Segment named '" + name + "' not found")
}

// Return collection of all named segment names
func GetAllNamedSegmentNames() []string {
	names := make([]string, len(allNamedSegments))
	for i := 0; i < len(names); i++ {
		names[i] = allNamedSegments[i].Name
	}
	return names
}

// Slice off all named segments, use a function scope to hide internal variables
var allNamedSegments = func() []NamedSegment {

	// Bedroom Ceiling
	bedroomCeiling := func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
		t := framebuffer.NewPhySegment([]framebuffer.LedStrip{
			fb.Strips[11], fb.Strips[4],
			fb.Strips[3], fb.Strips[7],
			fb.Strips[6], fb.Strips[14],
			fb.Strips[13]})
		return framebuffer.NewLogSegment(t, 135, t.Len()-135)
	}

	// Bedroom Wall
	bedroomWall := func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
		t := framebuffer.NewPhySegment([]framebuffer.LedStrip{
			fb.Strips[8], fb.Strips[2], fb.Strips[5], fb.Strips[15], fb.Strips[10]})
		return framebuffer.NewLogSegment(t, 144, t.Len()-(144+241))
	}

	// Bathroom Ceiling
	bathroomCeiling := func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
		t := framebuffer.NewPhySegment([]framebuffer.LedStrip{
			fb.Strips[9], fb.Strips[11]})
		return framebuffer.NewLogSegment(t, 0, t.Len()-27)
	}

	// Bathroom Wall
	bathroomWall := func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
		t := framebuffer.NewPhySegment([]framebuffer.LedStrip{
			fb.Strips[10], fb.Strips[8]})
		return framebuffer.NewLogSegment(t, 50, t.Len()-(28+50))
	}

	// Build the segment list
	return []NamedSegment{
		{
			"All",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewPhySegment([]framebuffer.LedStrip{
					fb.Strips[3], fb.Strips[7],
					fb.Strips[6], fb.Strips[14],
					fb.Strips[13], fb.Strips[9],
					fb.Strips[11], fb.Strips[4],
					fb.Strips[5], fb.Strips[15],
					fb.Strips[10], fb.Strips[8],
					fb.Strips[2]})
			},
		},
		{
			"All Ceiling",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewPhySegment([]framebuffer.LedStrip{
					fb.Strips[3], fb.Strips[7],
					fb.Strips[6], fb.Strips[14],
					fb.Strips[13], fb.Strips[9],
					fb.Strips[11], fb.Strips[4]})
			},
		},
		{
			"All Wall",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewPhySegment([]framebuffer.LedStrip{
					fb.Strips[5], fb.Strips[15],
					fb.Strips[10], fb.Strips[8],
					fb.Strips[2]})
			},
		},
		{
			"Bedroom",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewCombinedSegment(bedroomCeiling(fb), bedroomWall(fb))
			},
		},
		{
			"Bedroom Ceiling",
			bedroomCeiling,
		},
		{
			"Bedroom Wall",
			bedroomWall,
		},
		{
			"Bathroom",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewCombinedSegment(bathroomCeiling(fb), bathroomWall(fb))
			},
		},
		{
			"Bathroom Ceiling",
			bathroomCeiling,
		},
		{
			"Bathroom Wall",
			bathroomWall,
		},
		{
			"Curtains",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewPhySegment([]framebuffer.LedStrip{
					fb.Strips[3], fb.Strips[7]})
			},
		},
		{
			"Strip Three",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewLogSegment(framebuffer.NewPhySegment(
					[]framebuffer.LedStrip{fb.Strips[3]}), 0, 20)
			},
		},
		{
			"Strip Five",
			func(fb *framebuffer.FrameBuffer) framebuffer.Segment {
				return framebuffer.NewLogSegment(framebuffer.NewPhySegment(
					[]framebuffer.LedStrip{fb.Strips[5]}), 0, 19)
			},
		},
	}
}()