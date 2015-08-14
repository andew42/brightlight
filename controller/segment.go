package controller

// A segment describes a logical LED strip
type Segment interface {

	Len() uint
	Set(pos uint, colour Rgb)
}
