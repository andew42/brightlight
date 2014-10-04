package animations

import "github.com/andew42/brightlight/controller"

// A logical segment is a slice of another segment
type LogSegment struct {
	baseSeg controller.Segment
	start   uint
	len     uint
}

// Constructor
func NewLogSegment(baseSeg controller.Segment, start uint, len uint) LogSegment {
	baseLen := baseSeg.Len()
	if start >= baseLen {
		panic("invalid segment start")
	}
	if start+len > baseLen {
		panic("invalid segment length")
	}
	var x LogSegment
	x.baseSeg = baseSeg
	x.start = start
	x.len = len
	return x
}

// Number of LEDs in the segment
func (seg LogSegment) Len() uint {
	return seg.len
}

// Set a particular LED colour from the left of the segment
func (seg LogSegment) Set(pos uint, colour controller.Rgb) {
	// Is position out of range?
	if pos >= seg.len {
		panic("position out of range")
	}
	// Set at position within segment
	seg.baseSeg.Set(seg.start+pos, colour)
}
