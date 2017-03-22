package animations

import (
	"math/rand"
	"github.com/andew42/brightlight/framebuffer"
)

// import "github.com/pborges/huejack"

// Twinkle inspired by:
// https://github.com/daterdots/LEDs/blob/master/TwinkleSparkle2015/TwinkleSparkle2015.ino

// How quickly to cool down a pixel
const cooling = 5

// Percentage of pixies to heat up each frame
const heatPercentage = 5

type twinkle struct {
	pixies []byte
}

// ctor
func newTwinkle() *twinkle {

	var t twinkle
	return &t
}

func (t *twinkle) clone() animator {
	// TODO: DEEP COPY
	panic("Not Implemented")
}

func randInt31Between(min int, max int) int {
	x := max - min
	return rand.Intn(x) + min
}

// 8 bit saturation math subtraction
func qsub8(left uint8, right uint8) uint8 {
	if right >= left {
		return 0
	}
	return left - right
}

// animation interface
func (t *twinkle) animateNextFrame(frameCount int, frame framebuffer.Segment) {

	// Create the pixies if not already created
	if t.pixies == nil {
		t.pixies = make([]byte, frame.Len())
	}

	// "Cool" down every location on the strip a little
	for i := 0; i < len(t.pixies); i++ {
		t.pixies[i] = qsub8(t.pixies[i], cooling)
	}

	// Adding heat every frame is too much so do it every 2 frames
	if (frameCount % 2) == 0 {

		// Add some heat to some random pixies
		numPixiesToHeat := (len(t.pixies) * heatPercentage) / 100
		if numPixiesToHeat < 1 {
			numPixiesToHeat = 1
		}
		for j := 0; j < numPixiesToHeat; j++ {
			index := rand.Intn(len(t.pixies) - 1)
			heat := uint8(randInt31Between(50, 255))
			if heat > t.pixies[index] {
				t.pixies[index] = heat
			}
		}
	}

	// Map the heat to a white colour
	for i := uint(0); i < frame.Len(); i++ {
		frame.Set(i, framebuffer.NewRgb(t.pixies[i], t.pixies[i], t.pixies[i]))
	}
}
