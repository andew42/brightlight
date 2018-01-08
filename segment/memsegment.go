package segment

import "github.com/andew42/brightlight/framebuffer"

type MemSegment struct {
	memory []framebuffer.Rgb
}

// An in memory segment (uninitialised)
func NewMemSegment(length uint) *MemSegment {

	return &MemSegment{memory: make([]framebuffer.Rgb, length)}
}

func (seg *MemSegment) Len() uint {

	return uint(len(seg.memory))
}

func (seg *MemSegment) Get(pos uint) framebuffer.Rgb {

	return seg.memory[pos]
}

func (seg *MemSegment) Set(pos uint, colour framebuffer.Rgb) {

	seg.memory[pos] = colour
}
