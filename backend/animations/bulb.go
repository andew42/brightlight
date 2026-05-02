package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
)

type bulb struct {
	colour framebuffer.Rgb
	offset uint
	width  uint
}

// A single 'bulb' of width LEDs wide of specified colour
// Typically used with a repeater to create a string of bulbs
func newBulb(colour framebuffer.Rgb, offset uint, width uint) *bulb {

	return &bulb{colour, offset, width}
}

func (b *bulb) animateFrame(_ uint, frame segment.Segment) {

	// Return if bulb won't fit
	if b.offset+b.width > frame.Len() {
		return
	}

	// Light the bulb
	for s := b.offset; s < b.offset+b.width; s++ {
		frame.Set(s, b.colour)
	}
}
