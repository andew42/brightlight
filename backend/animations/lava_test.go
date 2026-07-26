package animations

import (
	"math"
	"testing"

	"github.com/andew42/brightlight/framebuffer"
)

var lavaTestColour = framebuffer.NewRgb(255, 255, 255)

// Blob radii must stay within the min..max band, reach both ends of it
// and not all be the same size at once
func TestLavaBlobRadiiSweepSizeBand(t *testing.T) {
	const blobCount, sizeMin, sizeMax = 5, 8, 16
	l := newLava(lavaTestColour, lavaTestColour, 25, blobCount, sizeMin, sizeMax)

	base := 0.5 / float64(blobCount)
	wantMin, wantMax := base*sizeMin/lavaNaturalSize, base*sizeMax/lavaNaturalSize
	if math.Abs(l.radiusMin-wantMin) > 1e-9 || math.Abs(l.radiusMax-wantMax) > 1e-9 {
		t.Fatalf("size band %v..%v, want %v..%v", l.radiusMin, l.radiusMax, wantMin, wantMax)
	}

	lowest, highest, spread := math.Inf(1), math.Inf(-1), 0.0
	for f := uint(0); f < 20000; f++ {
		tSize := float64(f) * 0.04 * l.speed * lavaSizeRate
		frameLow, frameHigh := math.Inf(1), math.Inf(-1)
		for _, b := range l.blobs {
			r := l.blobRadius(b, tSize)
			if r < l.radiusMin-1e-9 || r > l.radiusMax+1e-9 {
				t.Fatalf("radius %v outside band at frame %d", r, f)
			}
			frameLow, frameHigh = math.Min(frameLow, r), math.Max(frameHigh, r)
		}
		lowest, highest = math.Min(lowest, frameLow), math.Max(highest, frameHigh)
		spread = math.Max(spread, frameHigh-frameLow)
	}

	if lowest > wantMin*1.02 || highest < wantMax*0.98 {
		t.Errorf("radii only covered %v..%v of the %v..%v band", lowest, highest, wantMin, wantMax)
	}
	if spread < (wantMax-wantMin)/2 {
		t.Errorf("blobs stay too close in size, widest spread %v", spread)
	}
}

// Size must change more slowly than the blobs drift across the segment
func TestLavaSizeChangesSlowerThanDrift(t *testing.T) {
	if lavaSizeRate >= 1 {
		t.Errorf("lavaSizeRate %v does not slow size changes", lavaSizeRate)
	}
}

// The sliders set the wrong way round should still give a usable band
func TestLavaSizeSlidersReversed(t *testing.T) {
	l := newLava(lavaTestColour, lavaTestColour, 10, 4, 20, 5)
	if l.radiusMin > l.radiusMax || l.radiusMin <= 0 {
		t.Errorf("band not normalised: %v..%v", l.radiusMin, l.radiusMax)
	}
}
