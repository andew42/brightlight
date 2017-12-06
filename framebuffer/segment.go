package framebuffer

// A segment describes a logical LED strip
type Segment interface {
	Len() uint
	Get(pos uint) Rgb
	Set(pos uint, colour Rgb)
}
