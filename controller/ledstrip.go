package controller

import "math/rand"

// Maximum led strip length (must match teensy firmware)
const MaxLedStripLen = 300
const StripsPerTeensy = 8

type LedStrip struct {
	LeftToRight bool
	Leds        []Rgb
}

// leftToRight or clockwise
func NewLedStrip(leftToRight bool, len int) *LedStrip {

	var s LedStrip
	s.LeftToRight = leftToRight
	s.Leds = make([]Rgb, len)
	// TODO: REMOVE Initialise with 'random' values?
	for i := 0; i < len; i++ {
		s.Leds[i].Red = byte(i * rand.Intn(255))
		s.Leds[i].Green = byte(i * rand.Intn(255))
		s.Leds[i].Blue = byte(i * rand.Intn(255))
	}
	return &s
}
