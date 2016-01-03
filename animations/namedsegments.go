package animations

import (
	"github.com/andew42/brightlight/controller"
)

// The list of names segment used to build scene in the UI
// SHOULD MATCH UI LIST IN UI/CONFIG/STATIC.JS
// A named segment is a single LogSegment (logical segment)
// A LogSegment consists of one or more segments (logical or physical) plus a start and end led
// A PhySegment (physical segment) consists of one or more led strips
// A LedStrip represents the physical strip of LEDS connected to a particular controller pin

// Unmarshal JSON into typed slice
type segmentInfo struct {
	Name string
	Seg  controller.Segment
}

// Constructs the NamedSegments map indexed by segment id
func NewNamedSegments(fb *controller.FrameBuffer) map[string]segmentInfo {

	var namedSegments = make(map[string]segmentInfo)

	namedSegments["s0"] = segmentInfo {"All", controller.NewPhySegment(fb.Strips)}

	// TODO
	namedSegments["s1"] = segmentInfo {"Bedroom", controller.NewPhySegment(fb.Strips)}
	namedSegments["s2"] = segmentInfo {"Bathroom", controller.NewPhySegment(fb.Strips)}

	namedSegments["s3"] = segmentInfo {"Curtains",
		controller.NewPhySegment([]controller.LedStrip{fb.Strips[3], fb.Strips[7]})}

	namedSegments["s4"] = segmentInfo {"Test 4", controller.NewPhySegment([]controller.LedStrip{fb.Strips[4]})}
	namedSegments["s5"] = segmentInfo {"Test 5", controller.NewPhySegment([]controller.LedStrip{fb.Strips[5]})}

	return namedSegments
}
