package controller

import "math/rand"

// Maximum led strip length (must match teensy firmware)
const MaxLedStripLen = 175

type LedStrip struct {
	LeftToRight bool
	Leds        []Rgb
}

func NewLedStrip(leftToRight bool, len int) *LedStrip {
	var s LedStrip
	s.LeftToRight = leftToRight
	s.Leds = make([]Rgb, len, MaxLedStripLen)
	// TODO: REMOVE Initialise with 'random' values?
	for i := 0; i < len; i++ {
		s.Leds[i].Red = byte(i * rand.Intn(255))
		s.Leds[i].Green = byte(i * rand.Intn(255))
		s.Leds[i].Blue = byte(i * rand.Intn(255))
	}
	return &s
}
