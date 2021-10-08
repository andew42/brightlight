package segment

import (
	"github.com/andew42/brightlight/framebuffer"
	log "github.com/sirupsen/logrus"
)

type PhySegment struct {
	Strips []framebuffer.LedStrip
}

// NewPhySegment A physical segment aggregates a number of LedStrips
func NewPhySegment(strips []framebuffer.LedStrip) PhySegment {

	// Strip out zero length segments
	ps := PhySegment{}
	for _, s := range strips {
		if len(s.Leds) > 0 {
			ps.Strips = append(ps.Strips, s)
		}
	}
	return ps
}

// Len Number of LEDs in the segment
func (seg PhySegment) Len() uint {

	l := uint(0)
	for i := 0; i < len(seg.Strips); i++ {
		l += uint(len(seg.Strips[i].Leds))
	}
	return l
}

// Get a particular LED colour from the left of the strip
func (seg PhySegment) Get(pos uint) framebuffer.Rgb {

	stripIndex, stripPos := seg.locate(pos)
	return seg.Strips[stripIndex].Leds[stripPos]
}

// Set a particular LED colour from the left of the strip
func (seg PhySegment) Set(pos uint, colour framebuffer.Rgb) {

	stripIndex, stripPos := seg.locate(pos)
	seg.Strips[stripIndex].Leds[stripPos] = colour
}

// Takes an index (pos) from the left of the physical segment and calculates
// the index of the LedStrip (stripIndex) and position (stripPos) within the
// sub strip
func (seg PhySegment) locate(pos uint) (stripIndex int, stripPos uint) {

	// Locate the LedStrip for this position
	stripIndex = 0
	stripPos = pos
	for ; stripIndex < len(seg.Strips); stripIndex++ {
		stripLen := uint(len(seg.Strips[stripIndex].Leds))
		// Is this the required strip?
		if stripPos < stripLen {
			break
		}
		stripPos -= stripLen
	}

	// Was the position out of range?
	if stripIndex == len(seg.Strips) {
		log.Panic("position out of range")
	}

	// Transpose strip position if LedStrip is anti-clockwise
	if !seg.Strips[stripIndex].Clockwise {
		stripPos = uint(len(seg.Strips[stripIndex].Leds)) - stripPos - 1
	}
	return
}
