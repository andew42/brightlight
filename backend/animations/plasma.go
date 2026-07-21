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
	t := float64(frameCount) * p.speed * 0.01

	for i := uint(0); i < frame.Len(); i++ {
		pos := float64(i) / float64(frame.Len())

		wave1 := math.Sin(pos*10.0 + t)
		wave2 := math.Sin(pos*8.0 - t*0.7)
		wave3 := math.Cos(pos*6.0 + t*0.5)
		wave4 := math.Sin(pos*12.0 - t*0.3)

		plasma := (wave1 + wave2 + wave3 + wave4) / 4.0
		normalized := (plasma + 1.0) / 2.0

		var hue uint
		var saturation uint

		switch p.palette {
		case 1: // Fire (red-orange-yellow)
			hue = uint(normalized * 60.0)
			saturation = 100
		case 2: // Ocean (blue-cyan-green)
			hue = uint(180.0 + normalized*60.0)
			saturation = 80
		case 3: // Forest (green shades)
			hue = uint(90.0 + normalized*60.0)
			saturation = 70
		default: // Rainbow
			hue = uint(normalized * 360.0)
			saturation = 100
		}

		frame.Set(i, framebuffer.NewRgbFromHsl(hue, saturation, uint(p.brightness)))
	}
}
