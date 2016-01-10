package framebuffer

type LedStrip struct {
	Clockwise bool
	Leds      []Rgb
}

func NewLedStrip(clockwise bool, len int) *LedStrip {

	var s LedStrip
	s.Clockwise = clockwise
	s.Leds = make([]Rgb, len)
	return &s
}
