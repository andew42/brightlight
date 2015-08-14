package controller

import (
	log "github.com/Sirupsen/logrus"
)

// Part of another segment
type SubSegment struct {

	baseSeg Segment
	start   uint
	len     uint
}

// Constructor start and length within baseSeg
func NewSubSegment(baseSeg Segment, start uint, len uint) SubSegment {

	baseLen := baseSeg.Len()
	if start > baseLen || start+len > baseLen {
		log.Panic("invalid start or end")
	}
	return SubSegment{baseSeg, start, len}
}

// Number of LEDs in the segment
func (seg SubSegment) Len() uint {

	return seg.len
}

// Set an LED relative to sub segment start
func (seg SubSegment) Set(pos uint, colour Rgb) {

	baseSegPos := seg.start + pos
	if baseSegPos >= seg.Len() {
		log.Panic("position out of range")
	}
	seg.baseSeg.Set(baseSegPos, colour)
}
