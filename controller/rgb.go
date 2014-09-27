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
	var x Rgb
	x.Red = r
	x.Green = g
	x.Blue = b
	return x
}

func NewRgbFromInt(colour int) Rgb {
	var x Rgb
	x.Red = byte(colour >> 16)
	x.Green = byte(colour >> 8)
	x.Blue = byte(colour)
	return x
}

// Convert led RGB (3 bytes) to JSON colour (int)
func (led Rgb) MarshalJSON() ([]byte, error) {
	rc := make([]byte, 0, 16)
	rc = strconv.AppendInt(rc, int64(led.Red)<<16+int64(led.Green)<<8+int64(led.Blue), 10)
	return rc, nil
}
