package animations

import (
	"math"

	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/segment"
)

type plasma struct {
	speed      float64
	brightness int
	palette    int // 0=rainbow, 1=fire, 2=ocean, 3=forest
}

func newPlasma(speed float64, brightness int, palette int) *plasma {
	return &plasma{
		speed:      speed,
		brightness: brightness,
		palette:    palette,
	}
}

func (p *plasma) animateFrame(frameCount uint, frame segment.Segment) {
	// Time-based animation - convert frame count to time offset
	t := float64(frameCount) * p.speed * 0.01

	for i := uint(0); i < frame.Len(); i++ {
		// Normalized position (0.0 to 1.0)
		pos := float64(i) / float64(frame.Len())

		// Create plasma effect using multiple sine waves with different frequencies
		// This creates the organic, flowing lava lamp effect
		wave1 := math.Sin(pos*10.0 + t)
		wave2 := math.Sin(pos*8.0 - t*0.7)
		wave3 := math.Cos(pos*6.0 + t*0.5)
		wave4 := math.Sin(pos*12.0 - t*0.3)

		// Combine waves to create complex plasma pattern
		// Each wave contributes to the overall value
		plasma := (wave1 + wave2 + wave3 + wave4) / 4.0

		// Normalize to 0.0 - 1.0 range
		normalized := (plasma + 1.0) / 2.0

		// Convert to hue based on selected palette
		var hue uint
		var saturation uint

		switch p.palette {
		case 1: // Fire palette (red-orange-yellow)
			hue = uint(normalized*60.0) // 0-60 degrees (red to yellow)
			saturation = 100
		case 2: // Ocean palette (blue-cyan-green)
			hue = uint(180.0 + normalized*60.0) // 180-240 degrees (cyan to blue)
			saturation = 80
		case 3: // Forest palette (green shades)
			hue = uint(90.0 + normalized*60.0) // 90-150 degrees (green range)
			saturation = 70
		default: // Rainbow palette (full spectrum)
			hue = uint(normalized * 360.0)
			saturation = 100
		}

		// Create color with calculated hue and configured brightness
		colour := framebuffer.NewRgbFromHsl(hue, saturation, uint(p.brightness))
		frame.Set(i, colour)
	}
}
