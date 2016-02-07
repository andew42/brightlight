package framebuffer

import (
	log "github.com/Sirupsen/logrus"
)

// Two segments connected together
type CombinedSegment struct {
	Seg1 Segment
	Seg2 Segment
}

// Constructor
func NewCombinedSegment(seg1 Segment, seg2 Segment) CombinedSegment {

	return CombinedSegment{seg1, seg2}
}

// Number of LEDs in the segment
func (s CombinedSegment) Len() uint {

	return s.Seg1.Len() + s.Seg2.Len()
}

// Set a particular LED colour 0 is 0 in seg1
func (s CombinedSegment) Set(pos uint, colour Rgb) {

	// Is position out of range?
	if pos >= s.Len() {
		log.Panic("position out of range")
	}

	// Set in seg 1 or seg 2
	if pos >= s.Seg1.Len() {
		s.Seg2.Set(pos-s.Seg1.Len(), colour)
	} else {
		s.Seg1.Set(pos, colour)
	}
}
