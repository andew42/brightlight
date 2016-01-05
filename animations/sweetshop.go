package animations

import (
	"github.com/andew42/brightlight/controller"
	"math/rand"
	"time"
)

type sweetshop struct {
	seg        controller.Segment
	period     time.Duration
	changeTime time.Time
}

func newSweetshop(seg controller.Segment, period time.Duration) *sweetshop {

	return &sweetshop{seg: seg, period: period}
}

func (s *sweetshop) animateNextFrame() {

	if time.Now().Sub(s.changeTime) > 0 {
		s.changeTime = time.Now().Add(s.period)
		for i := uint(0); i < s.seg.Len(); i++ {
			s.seg.Set(i, controller.NewRgbFromInt(rand.Int()&(1<<24-1)))
		}
	}
}
