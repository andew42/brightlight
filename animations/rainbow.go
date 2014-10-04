package animations

import "fmt"
import "github.com/andew42/brightlight/controller"

type rainbow struct {
	seg             controller.Segment
	degreesPerFrame float32
	startDegree     float32
}

// periodMs is the number of ms for a complete cycle of the rainbow
func newRainbow(seg controller.Segment, periodMs uint) *rainbow {
	var r rainbow
	r.seg = seg
	framesPerPeriod := float32(periodMs) / float32(frameRateMs)
	r.degreesPerFrame = 360.0 / framesPerPeriod
	fmt.Printf("%v %v %v", r.seg, r.degreesPerFrame, r.startDegree)
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
