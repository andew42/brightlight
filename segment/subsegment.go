package segment

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/framebuffer"
)

type SubSegment struct {
	baseSeg Segment
	start   uint
	len     uint
}

// A sub segment is a slice of another segment
func NewSubSegment(baseSeg Segment, start uint, len uint) SubSegment {

	baseLen := baseSeg.Len()
	if start >= baseLen {
		log.Panic("invalid segment start")
	}
	if start+len > baseLen {
		log.Panic("invalid segment length")
	}
	return SubSegment{baseSeg, start, len}
}

// Number of LEDs in the segment
func (seg SubSegment) Len() uint {

	return seg.len
}

// Get a particular LED colour from the left of the segment
func (seg SubSegment) Get(pos uint) framebuffer.Rgb {

	// Is position out of range?
	if pos >= seg.len {
		log.Panic("position out of range")
	}
	// Get at position within segment
	return seg.baseSeg.Get(seg.start + pos)
}

// Set a particular LED colour from the left of the segment
func (seg SubSegment) Set(pos uint, colour framebuffer.Rgb) {

	// Is position out of range?
	if pos >= seg.len {
		log.Panic("position out of range")
	}
	// Set at position within segment
	seg.baseSeg.Set(seg.start+pos, colour)
}
