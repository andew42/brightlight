package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"time"
)

// Simulates traditional christmas lights where each light position is
// a fixed colour and light positions alternate around specified colours
// Each light can be represented by more than one adjacent LED
type christmas struct {
	period       time.Duration
	changeTime   time.Time
	lightSize    uint
	lightColours []framebuffer.Rgb
	nextColour   uint
}

func newChristmas(period time.Duration) *christmas {

	return &christmas{period: period, lightSize: 3,
		lightColours: []framebuffer.Rgb{framebuffer.NewRgb(255, 0, 0), framebuffer.NewRgb(0, 255, 0), framebuffer.NewRgb(0, 0, 255)}}
}

func (s *christmas) animateNextFrame(seg framebuffer.Segment) {

	// Time to change lights
	if time.Now().Sub(s.changeTime) > 0 {
		s.changeTime = time.Now().Add(s.period)

		// Move on to next light colour
		if s.nextColour == uint(len(s.lightColours))-1 {
			s.nextColour = 0
		} else {
			s.nextColour++
		}

		// Set each LED appropriately
		off := framebuffer.NewRgbFromInt(0)
		for i := uint(0); i < seg.Len(); i++ {
			// Which colour index should this be
			c := (i / s.lightSize) % uint(len(s.lightColours))

			// Either turn colour on or turn LED off
			if c == s.nextColour {
				seg.Set(i, s.lightColours[s.nextColour])
			} else {
				seg.Set(i, off)
			}
		}
	}
}
