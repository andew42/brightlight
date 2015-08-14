package animations

import (
	"github.com/andew42/brightlight/controller"
	log "github.com/Sirupsen/logrus"
)

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

// Set a particular LED colour from the left of the segment
func (seg LogSegment) Set(pos uint, colour controller.Rgb) {

	// Is position out of range?
	if pos >= seg.len {
		log.Panic("position out of range")
	}
	// Set at position within segment
	seg.baseSeg.Set(seg.start+pos, colour)
}
