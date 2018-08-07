package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"time"
	"github.com/andew42/brightlight/segment"
)

type sweetshop struct {
	period        time.Duration
	brightness    int
	minSaturation int
	changeTime    time.Time
	pixies        []framebuffer.Rgb
}

func newSweetshop(period time.Duration, brightness int, minSaturation int) *sweetshop {

	return &sweetshop{period: period, brightness: brightness, minSaturation: minSaturation}
}

func (s *sweetshop) animateFrame(frameCount uint, frame segment.Segment) {

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
			s.pixies[i] = framebuffer.NewRgbFromHsl(randBetween(0, 359), randBetween(s.minSaturation, 100), uint(s.brightness))
		}
	}

	// Update the frame buffer from the pixies
	for i := uint(0); i < frame.Len(); i++ {
		frame.Set(i, s.pixies[i])
	}
}
