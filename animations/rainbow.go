package animations

import (
	"time"
	"github.com/andew42/brightlight/controller"
)

type rainbow struct {
	seg             controller.Segment
	degreesPerFrame float32
	startDegree     float32
}

// period is the duration for a complete cycle of the rainbow
func newRainbow(seg controller.Segment, period time.Duration) *rainbow {
	var r rainbow
	r.seg = seg
	framesPerPeriod := float32(period) / float32(frameRate)
	r.degreesPerFrame = 360.0 / framesPerPeriod
	return &r
}

func (r *rainbow) animateNextFrame() {
	hue := r.startDegree
	hueIncrement := 360.0 / float32(r.seg.Len())
	for i := uint(0); i < r.seg.Len(); i++ {
		r.seg.Set(i, controller.NewRgbFromHsl(uint(hue), 100, 50))
		hue += hueIncrement
	}
	r.startDegree += r.degreesPerFrame
	if r.startDegree > 360 {
		r.startDegree -= 360
	}
}