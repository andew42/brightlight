package framebuffer

import (
	log "github.com/Sirupsen/logrus"
)

// A logical segment is a slice of another segment
type LogSegment struct {
	baseSeg Segment
	start   uint
	len     uint
}

// Constructor
func NewLogSegment(baseSeg Segment, start uint, len uint) LogSegment {

	baseLen := baseSeg.Len()
	if start >= baseLen {
		log.Panic("invalid segment start")
	}
	if start+len > baseLen {
		log.Panic("invalid segment length")
	}
	return LogSegment{baseSeg, start, len}
}

// Number of LEDs in the segment
func (seg LogSegment) Len() uint {

	return seg.len
}

// Get a particular LED colour from the left of the segment
func (seg LogSegment) Get(pos uint) Rgb {

	// Is position out of range?
	if pos >= seg.len {
		log.Panic("position out of range")
	}
	// Get at position within segment
	return seg.baseSeg.Get(seg.start + pos)
}

// Set a particular LED colour from the left of the segment
func (seg LogSegment) Set(pos uint, colour Rgb) {

	// Is position out of range?
	if pos >= seg.len {
		log.Panic("position out of range")
	}
	// Set at position within segment
	seg.baseSeg.Set(seg.start+pos, colour)
}
