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
	// Size breathing frequency and phase. The random phase both puts each
	// blob at a different size to start with and keeps them from breathing
	// in step with one another
	w3, p3 float64
}

type lava struct {
	blobColour float64Rgb
	background float64Rgb
	speed      float64
	// Blob radius sweeps between these two, set by the size sliders
	radiusMin float64
	radiusMax float64
	blobs     []lavaBlob
}

// Rgb as floats to avoid per-pixel byte conversions while blending
type float64Rgb struct {
	r, g, b float64
}

func toFloat64Rgb(c framebuffer.Rgb) float64Rgb {
	return float64Rgb{float64(c.Red), float64(c.Green), float64(c.Blue)}
}

// All the sliders run 1..25. Speed maps exponentially onto a drift rate,
// 25 being the fastest, so the bottom of the slider is a barely
// perceptible creep. The size sliders scale the base radius linearly,
// with lavaNaturalSize giving the natural radius for the blob count.
const (
	lavaSliderSteps = 25
	lavaSlowestRate = 0.002
	lavaFastestRate = 1.0
	lavaNaturalSize = 12.0
	// Blobs grow and shrink at a fraction of their drift rate, so a blob
	// crosses the segment several times over one min-to-max-to-min cycle
	lavaSizeRate = 0.3
)

func clampSlider(v int) int {
	if v < 1 {
		return 1
	}
	if v > lavaSliderSteps {
		return lavaSliderSteps
	}
	return v
}

func lavaSpeedRate(speed int) float64 {
	speed = clampSlider(speed)
	return lavaSlowestRate * math.Pow(lavaFastestRate/lavaSlowestRate,
		float64(speed-1)/float64(lavaSliderSteps-1))
}

func newLava(blobColour framebuffer.Rgb, background framebuffer.Rgb, speed int, blobCount int, sizeMin int, sizeMax int) *lava {

	blobs := make([]lavaBlob, blobCount)
	for i := range blobs {
		blobs[i] = lavaBlob{
			w1: 0.6 + rand.Float64()*0.8,
			p1: rand.Float64() * 2 * math.Pi,
			w2: 1.3 + rand.Float64()*1.4,
			p2: rand.Float64() * 2 * math.Pi,
			w3: 0.6 + rand.Float64()*0.8,
			p3: rand.Float64() * 2 * math.Pi,
		}
	}

	// Tolerate the sliders being set the wrong way round
	sizeMin, sizeMax = clampSlider(sizeMin), clampSlider(sizeMax)
	if sizeMin > sizeMax {
		sizeMin, sizeMax = sizeMax, sizeMin
	}

	// Base blob radius as a fraction of the segment, smaller with more blobs
	base := 0.5 / float64(blobCount)

	return &lava{
		blobColour: toFloat64Rgb(blobColour),
		background: toFloat64Rgb(background),
		speed:      lavaSpeedRate(speed),
		radiusMin:  base * float64(sizeMin) / lavaNaturalSize,
		radiusMax:  base * float64(sizeMax) / lavaNaturalSize,
		blobs:      blobs,
	}
}

// Radius of a blob at size-time tSize. It sweeps the full min..max range,
// each blob at its own phase so they are all a different size at any one
// moment and none of them breathes in step with the others
func (l *lava) blobRadius(b lavaBlob, tSize float64) float64 {
	return l.radiusMin + (l.radiusMax-l.radiusMin)*(0.5+0.5*math.Sin(b.w3*tSize+b.p3))
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

	// Size changes more slowly than the blobs drift
	tSize := t * lavaSizeRate

	// Precompute each blob's position and radius for this frame
	positions := make([]float64, len(l.blobs))
	radii := make([]float64, len(l.blobs))
	for k, b := range l.blobs {
		// Slow primary drift plus a faster low amplitude wobble,
		// amplitudes sum to <0.5 so blobs stay within the segment
		positions[k] = 0.5 + 0.36*math.Sin(b.w1*t+b.p1) + 0.14*math.Sin(b.w2*t+b.p2)
		radii[k] = l.blobRadius(b, tSize)
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
