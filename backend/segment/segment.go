package segment

import "github.com/andew42/brightlight/framebuffer"

// Segment A segment describes a logical LED strip
type Segment interface {
	Len() uint
	Get(pos uint) framebuffer.Rgb
	Set(pos uint, colour framebuffer.Rgb)
}
