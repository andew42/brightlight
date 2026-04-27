package segment

import (
	"errors"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	"strconv"
	"strings"
)

// Named segments are used to build scenes in the UI
// TODO: SHOULD MATCH UI LIST IN UI/CONFIG/STATIC.JS
// A SubSegment (logical segment) consists of one or more segments (logical or physical)
// plus a start and end led position
// A PhySegment (physical segment) consists of one or more led strips
// LedStrip represents a physical strip of LEDS connected to a particular controller pin

type NamedSegment struct {
	Name       string
	GetSegment func(fb *framebuffer.FrameBuffer) Segment
}

// GetNamedSegment Return a particular named segment
func GetNamedSegment(name string) (NamedSegment, error) {

	// Handle special physical segment pxx:yy segment index : length
	// This feels a little hacky, it's used by the strip length ui
	if len(name) > 3 && name[0] == 'p' && strings.IndexByte(name, ':') != -1 {
		colonIndex := strings.IndexByte(name, ':')
		if stripIndex, err := strconv.Atoi(name[1:colonIndex]); err == nil {
			if length, err := strconv.Atoi(name[colonIndex+1:]); err == nil {
				return NamedSegment{
					Name: name,
					GetSegment: func(fb *framebuffer.FrameBuffer) Segment {
						return NewSubSegment(
							NewPhySegment(fb.Strips[stripIndex:stripIndex+1]),
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

// GetAllNamedSegmentNames Return collection of all named segment names
func GetAllNamedSegmentNames() []string {
	names := make([]string, len(allNamedSegments))
	for i := 0; i < len(names); i++ {
		names[i] = allNamedSegments[i].Name
	}
	return names
}

// Slice off all named segments, use a function scope to hide internal variables
var allNamedSegments = func() []NamedSegment {

	if config.Titania {
		// *** TITANIA ***

		p1 := func(fb *framebuffer.FrameBuffer) Segment {
			return NewPhySegment([]framebuffer.LedStrip{fb.Strips[0]})
		}
		p2 := func(fb *framebuffer.FrameBuffer) Segment {
			return NewPhySegment([]framebuffer.LedStrip{fb.Strips[1]})
		}
		p3 := func(fb *framebuffer.FrameBuffer) Segment {
			return NewPhySegment([]framebuffer.LedStrip{fb.Strips[2]})
		}

		return []NamedSegment{
			{
				"All",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewPhySegment([]framebuffer.LedStrip{fb.Strips[0], fb.Strips[2], fb.Strips[1]})
				},
			},
			{
				"Left",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewSubSegment(p1(fb), 0, 140)
				},
			},
			{
				"Right",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewSubSegment(p3(fb), 88, 140)
				},
			},
			{
				"Back",
				p2,
			},
			{
				"Front",
				func(fb *framebuffer.FrameBuffer) Segment {
					t := NewPhySegment([]framebuffer.LedStrip{
						fb.Strips[0], fb.Strips[2]})
					return NewSubSegment(t, 140, 177)
				},
			},
			{
				"One",
				p1,
			},
			{
				"Two",
				p2,
			},
			{
				"Three",
				p3,
			},
		}
	} else {
		// *** BEDROOM ***
		// Bedroom Ceiling
		bedroomCeiling := func(fb *framebuffer.FrameBuffer) Segment {
			t := NewPhySegment([]framebuffer.LedStrip{
				fb.Strips[11], fb.Strips[4],
				fb.Strips[3], fb.Strips[7],
				fb.Strips[6], fb.Strips[14],
				fb.Strips[13]})
			return NewSubSegment(t, 135, t.Len()-135)
		}

		// Bedroom Wall
		bedroomWall := func(fb *framebuffer.FrameBuffer) Segment {
			t := NewPhySegment([]framebuffer.LedStrip{
				fb.Strips[8], fb.Strips[2], fb.Strips[5], fb.Strips[15], fb.Strips[10]})
			return NewSubSegment(t, 144, t.Len()-(144+241))
		}

		// Bathroom Ceiling
		bathroomCeiling := func(fb *framebuffer.FrameBuffer) Segment {
			t := NewPhySegment([]framebuffer.LedStrip{
				fb.Strips[9], fb.Strips[11]})
			return NewSubSegment(t, 0, t.Len()-27)
		}

		// Bathroom Wall
		bathroomWall := func(fb *framebuffer.FrameBuffer) Segment {
			t := NewPhySegment([]framebuffer.LedStrip{
				fb.Strips[10], fb.Strips[8]})
			return NewSubSegment(t, 50, t.Len()-(28+50))
		}

		// Build the segment list
		return []NamedSegment{
			{
				"All",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewPhySegment([]framebuffer.LedStrip{
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
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewPhySegment([]framebuffer.LedStrip{
						fb.Strips[3], fb.Strips[7],
						fb.Strips[6], fb.Strips[14],
						fb.Strips[13], fb.Strips[9],
						fb.Strips[11], fb.Strips[4]})
				},
			},
			{
				"All Wall",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewPhySegment([]framebuffer.LedStrip{
						fb.Strips[5], fb.Strips[15],
						fb.Strips[10], fb.Strips[8],
						fb.Strips[2]})
				},
			},
			{
				"Bedroom",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewCombinedSegment(bedroomCeiling(fb), bedroomWall(fb))
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
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewCombinedSegment(bathroomCeiling(fb), bathroomWall(fb))
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
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewPhySegment([]framebuffer.LedStrip{
						fb.Strips[3], fb.Strips[7]})
				},
			},
			{
				"Strip Three",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewSubSegment(NewPhySegment(
						[]framebuffer.LedStrip{fb.Strips[3]}), 0, 20)
				},
			},
			{
				"Strip Five",
				func(fb *framebuffer.FrameBuffer) Segment {
					return NewSubSegment(NewPhySegment(
						[]framebuffer.LedStrip{fb.Strips[5]}), 0, 19)
				},
			},
		}
	}
}()
