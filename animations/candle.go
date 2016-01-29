package animations

import (
	"github.com/andew42/brightlight/framebuffer"
	"math/rand"
)

type candle struct {
	framesPerUpdate uint
}

// Typical candle is a 3 LED segment
func newCandle() *candle {

	return &candle{}
}

func (c *candle) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	// Yellow
	r := byte(rand.Intn(120) + 135)
	frame.Set(0, framebuffer.NewRgb(r, r, 0))
	// Red
	r = byte(rand.Intn(120) + 135)
	frame.Set(1, framebuffer.NewRgb(r, 0, 0))
	// Yellow
	r = byte(rand.Intn(120) + 135)
	frame.Set(2, framebuffer.NewRgb(r, r, 0))
}
