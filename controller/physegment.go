package controller

// A physical segment aggregates a number of LedStrips
type PhySegment struct {
	Strips []LedStrip
}

// Constructor
func NewPhySegment(strips []LedStrip) PhySegment {
	return PhySegment{Strips: strips}
}

// Number of LEDs in the segment
func (seg PhySegment) Len() uint {
	l := uint(0)
	for i := 0; i < len(seg.Strips); i++ {
		l += uint(len(seg.Strips[i].Leds))
	}
	return l
}

// Set a particular LED colour from the left of the strip
func (seg PhySegment) Set(pos uint, colour Rgb) {
	// Locate the strip for this position
	i := 0
	for ; i < len(seg.Strips); i++ {
		stripLen := uint(len(seg.Strips[i].Leds))
		// Is this the required strip?
		if pos < stripLen {
			break
		}
		pos -= stripLen
	}

	// Was the position out of range?
	if i == len(seg.Strips) {
		panic("position out of range")
	}

	// Set at position within strip
	if seg.Strips[i].LeftToRight {
		seg.Strips[i].Leds[pos] = colour
	} else {
		seg.Strips[i].Leds[uint(len(seg.Strips[i].Leds))-pos-1] = colour
	}
}
