package animations

import (
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	"time"
)

type rainbow struct {
	degreesPerFrame float32
	startDegree     float32
}

// period is the duration for a complete cycle of the rainbow
func newRainbow(period time.Duration) *rainbow {

	var r rainbow
	framesPerPeriod := float32(period) / float32(config.FramePeriodMs)
	r.degreesPerFrame = 360.0 / framesPerPeriod
	return &r
}

func (r *rainbow) clone() animator {
	clone := *r
	return &clone
}

func (r *rainbow) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	hue := r.startDegree
	hueIncrement := 360.0 / float32(frame.Len())
	for i := uint(0); i < frame.Len(); i++ {
		frame.Set(i, framebuffer.NewRgbFromHsl(uint(hue), 100, 50))
		hue += hueIncrement
	}
	r.startDegree += r.degreesPerFrame
	if r.startDegree > 360 {
		r.startDegree -= 360
	}
}
