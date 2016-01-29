package framebuffer

type LedStrip struct {
	Clockwise bool
	Leds      []Rgb
}

func (src *LedStrip) CloneLedStrip() *LedStrip {

	clone := LedStrip{src.Clockwise, make([]Rgb, 0, len(src.Leds))}
	clone.Leds = append(clone.Leds, src.Leds...)
	return &clone
}

func NewLedStrip(clockwise bool, len int) *LedStrip {

	var s LedStrip
	s.Clockwise = clockwise
	s.Leds = make([]Rgb, len)
	return &s
}
