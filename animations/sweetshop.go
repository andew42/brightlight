package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"math/rand"
	"time"
)

type sweetshop struct {
	period     time.Duration
	changeTime time.Time
}

func newSweetshop(period time.Duration) *sweetshop {

	return &sweetshop{period: period}
}

func (s *sweetshop) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	if time.Now().Sub(s.changeTime) > 0 {

		// Set next change time
		s.changeTime = time.Now().Add(s.period)

		// Refresh the random colours
		for i := uint(0); i < frame.Len(); i++ {
			frame.Set(i, framebuffer.NewRgbFromInt(rand.Int()&(1<<24-1)))
		}
	}
}
