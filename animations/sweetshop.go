package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"math/rand"
	"time"
)

type sweetshop struct {
	period           time.Duration
	changeTime       time.Time
	currentSelection []framebuffer.Rgb
}

func newSweetshop(period time.Duration) *sweetshop {

	return &sweetshop{period: period}
}

func (s *sweetshop) animateNextFrame(seg framebuffer.Segment) {

	if time.Now().Sub(s.changeTime) > 0 {
		// Set next change time
		s.changeTime = time.Now().Add(s.period)

		// Refresh the random colours
		s.currentSelection = make([]framebuffer.Rgb, seg.Len())
		for i := range s.currentSelection {
			s.currentSelection[i] = framebuffer.NewRgbFromInt(rand.Int() & (1<<24 - 1))
		}
	}
	for i := range s.currentSelection {
		seg.Set(uint(i), s.currentSelection[i])
	}
}
