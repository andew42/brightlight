package animations

import (
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/framebuffer"
	"time"
	"github.com/andew42/brightlight/segment"
)

type rainbow struct {
	framesPerCycle  uint
	degreesPerFrame float32
}

// period is the duration for a complete cycle of the rainbow
func newRainbow(period time.Duration) *rainbow {

	if period < config.FramePeriodMs {
		period = config.FramePeriodMs
	}

	var r rainbow
	r.framesPerCycle = uint(float32(period) / float32(config.FramePeriodMs))
	r.degreesPerFrame = 360.0 / float32(r.framesPerCycle)
	return &r
}

func (r *rainbow) animateFrame(frameCount uint, frame segment.Segment) {

	// Work out the phase (starting angle) for this frame
	phase := float32(frameCount%r.framesPerCycle) * r.degreesPerFrame
	phaseIncrementPerPixel := 360.0 / float32(frame.Len())
	for i := uint(0); i < frame.Len(); i++ {
		frame.Set(i, framebuffer.NewRgbFromHsl(uint(phase), 100, 50))
		phase += phaseIncrementPerPixel
	}
}
