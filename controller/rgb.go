package controller

import "strconv"

// Single LED colour
type Rgb struct {

	Red   byte
	Green byte
	Blue  byte
}

// Rgb constructors
func NewRgb(r byte, g byte, b byte) Rgb {

	return Rgb{r, g, b}
}

func NewRgbFromInt(colour int) Rgb {

	return Rgb{byte(colour >> 16), byte(colour >> 8), byte(colour)}
}

// Convert led RGB (3 bytes) to JSON colour (int)
func (led Rgb) MarshalJSON() ([]byte, error) {

	rc := make([]byte, 0, 16)
	rc = strconv.AppendInt(rc, int64(led.Red)<<16+int64(led.Green)<<8+int64(led.Blue), 10)
	return rc, nil
}

// Convert HSL (Hue, Saturation, Lightness) to RGB (Red, Green, Blue)
//
//   hue:        0 to 359 - position on the color wheel, 0=red, 60=orange,
//                            120=yellow, 180=green, 240=blue, 300=violet
//
//   saturation: 0 to 100 - how bright or dull the color, 100=full, 0=gray
//
//   lightness:  0 to 100 - how light the color is, 100=white, 50=color, 0=black
//
func NewRgbFromHsl(hue uint, saturation uint, lightness uint) Rgb {

	if hue > 359 {
		hue = hue%360
	}
	if saturation > 100 {
		saturation = 100
	}
	if lightness > 100 {
		lightness = 100
	}

	// algorithm from: http://www.easyrgb.com/index.php?X=MATH&H=19#text19
	var x Rgb
	if saturation == 0 {
		x.Red = byte(lightness * 255 / 100)
		x.Green = x.Red
		x.Blue = x.Green
	} else {
		var v2 uint
		if lightness < 50 {
			v2 = lightness*(100+saturation)
		} else {
			v2 = ((lightness+saturation)*100)-(saturation*lightness)
		}
		v1 := lightness * 200 - v2

		// Red
		var h uint
		if hue < 240 {
			h = hue+120
		} else {
			h = hue-240
		}
		x.Red = byte(h2rgb(v1, v2, h) * 255 / 600000)

		// Green
		x.Green = byte(h2rgb(v1, v2, hue) * 255 / 600000)

		// Blue
		if hue >= 120 {
			h = hue-120
		} else {
			h = hue+240
		}
		x.Blue = byte(h2rgb(v1, v2, h) * 255 / 600000)
	}
	return x
}

func h2rgb(v1 uint, v2 uint, hue uint) uint {

	if hue < 60 {
		return v1 * 60 + (v2 - v1) * hue
	}
	if hue < 180 {
		return v2 * 60
	}
	if hue < 240 {
		return v1 * 60 + (v2 - v1) * (240 - hue)
	}
	return v1 * 60
}
