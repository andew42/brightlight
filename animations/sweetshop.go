package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"time"
)

type sweetshop struct {
	period     time.Duration
	changeTime time.Time
	pixies     []framebuffer.Rgb
}

func newSweetshop(period time.Duration) *sweetshop {

	return &sweetshop{period: period}
}

func (s *sweetshop) animateFrame(frameCount uint, frame framebuffer.Segment) {

	// Create the pixies?
	if uint(len(s.pixies)) != frame.Len() {
		s.pixies = make([]framebuffer.Rgb, frame.Len())
	}

	// Time to change the pixies?
	if time.Now().Sub(s.changeTime) > 0 {

		// Set next change time
		s.changeTime = time.Now().Add(s.period)

		// Refresh the random colours
		for i := uint(0); i < frame.Len(); i++ {
			s.pixies[i] = framebuffer.NewRgbFromHsl(randBetween(0, 359), randBetween(50, 100),50)
		}
	}

	// Update the frame buffer from the pixies
	for i := uint(0); i < frame.Len(); i++ {
		frame.Set(i, s.pixies[i])
	}
}
