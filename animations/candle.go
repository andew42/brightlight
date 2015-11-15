package animations

import (
	"math/rand"
	"github.com/andew42/brightlight/controller"
)

type candle struct {

	seg controller.Segment
	framesPerUpdate uint
}

// Typical candle is a 3 LED segment
func newCandle(seg controller.Segment) *candle {

	return &candle{seg, 0}
}

func (c *candle) animateNextFrame() {

	if c.framesPerUpdate == 0 {
		c.framesPerUpdate = 6;
		// Yellow
		r := byte(rand.Intn(120) + 135);
		c.seg.Set(0, controller.NewRgb(r, r, 0))
		// Red
		r = byte(rand.Intn(120) + 135);
		c.seg.Set(1, controller.NewRgb(r, 0, 0))
		// Yellow
		r = byte(rand.Intn(120) + 135);
		c.seg.Set(2, controller.NewRgb(r, r, 0))
	} else {
		c.framesPerUpdate--
	}
}
