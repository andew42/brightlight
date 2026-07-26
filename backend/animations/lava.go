package animations

import (
	"math"
	"math/rand"

	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
)

// Lava effect: distinct blobs of one colour drift slowly through a
// background colour, merging and separating like wax in a lava lamp.
// Each blob's position and size are driven by summed sine waves with
// randomised frequencies and phases. The blobs contribute to a
// metaball-style field which is soft-thresholded so blob interiors are
// solid colour with rounded edges that blend into the background.

type lavaBlob struct {
	// Drift frequencies and phases (primary slow drift + faster wobble)
	w1, p1 float64
	w2, p2 float64
	// Size breathing frequency and phase
	w3, p3 float64
	// Relative size multiplier
	size float64
}

type lava struct {
	blobColour float64Rgb
	background float64Rgb
	speed      float64
	radius     float64
	blobs      []lavaBlob
}

// Rgb as floats to avoid per-pixel byte conversions while blending
type float64Rgb struct {
	r, g, b float64
}

func toFloat64Rgb(c framebuffer.Rgb) float64Rgb {
	return float64Rgb{float64(c.Red), float64(c.Green), float64(c.Blue)}
}

// Speed 1..10 maps exponentially onto a drift rate. 10 is the fastest and
// matches what used to be speed 5; 1 is ten times slower than the old
// slowest setting so the lamp can be made to creep.
const (
	lavaSlowestRate = 0.02
	lavaFastestRate = 1.0
)

func lavaSpeedRate(speed int) float64 {
	if speed < 1 {
		speed = 1
	} else if speed > 10 {
		speed = 10
	}
	return lavaSlowestRate * math.Pow(lavaFastestRate/lavaSlowestRate, float64(speed-1)/9)
}

func newLava(blobColour framebuffer.Rgb, background framebuffer.Rgb, speed int, blobCount int, blobSize int) *lava {

	blobs := make([]lavaBlob, blobCount)
	for i := range blobs {
		blobs[i] = lavaBlob{
			w1:   0.6 + rand.Float64()*0.8,
			p1:   rand.Float64() * 2 * math.Pi,
			w2:   1.3 + rand.Float64()*1.4,
			p2:   rand.Float64() * 2 * math.Pi,
			w3:   0.9 + rand.Float64()*1.1,
			p3:   rand.Float64() * 2 * math.Pi,
			size: 0.7 + rand.Float64()*0.6,
		}
	}

	return &lava{
		blobColour: toFloat64Rgb(blobColour),
		background: toFloat64Rgb(background),
		speed:      lavaSpeedRate(speed),
		// Base blob radius as a fraction of the segment, smaller with more
		// blobs, scaled by blob size 1..10 (5 is the natural size)
		radius: 0.5 / float64(blobCount) * float64(blobSize) / 5.0,
		blobs:  blobs,
	}
}

func smoothstep(edge0 float64, edge1 float64, x float64) float64 {
	t := (x - edge0) / (edge1 - edge0)
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return t * t * (3 - 2*t)
}

// animation interface
func (l *lava) animateFrame(frameCount uint, frame segment.Segment) {

	if frame.Len() == 0 {
		return
	}

	// Frame period is 40ms so t advances 0.04 * speed per real second
	t := float64(frameCount) * 0.04 * l.speed

	// Precompute each blob's position and radius for this frame
	positions := make([]float64, len(l.blobs))
	radii := make([]float64, len(l.blobs))
	for k, b := range l.blobs {
		// Slow primary drift plus a faster low amplitude wobble,
		// amplitudes sum to <0.5 so blobs stay within the segment
		positions[k] = 0.5 + 0.36*math.Sin(b.w1*t+b.p1) + 0.14*math.Sin(b.w2*t+b.p2)
		// Radius breathes around its base size
		radii[k] = l.radius * b.size * (1 + 0.25*math.Sin(b.w3*t+b.p3))
	}

	for i := uint(0); i < frame.Len(); i++ {
		pos := float64(i) / float64(frame.Len())

		// Metaball field: gaussian contribution from each blob so
		// nearby blobs merge smoothly into a single larger blob
		field := 0.0
		for k := range l.blobs {
			d := (pos - positions[k]) / radii[k]
			field += math.Exp(-d * d * 3.0)
		}

		// Soft threshold: solid blob interior, rounded blended edge
		a := smoothstep(0.15, 0.75, field)

		frame.Set(i, framebuffer.NewRgb(
			byte(l.background.r+(l.blobColour.r-l.background.r)*a),
			byte(l.background.g+(l.blobColour.g-l.background.g)*a),
			byte(l.background.b+(l.blobColour.b-l.background.b)*a)))
	}
}
