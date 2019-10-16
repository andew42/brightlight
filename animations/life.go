package animations

import (
	"github.com/andew42/brightlight/segment"
	"github.com/andew42/brightlight/framebuffer"
	log "github.com/sirupsen/logrus"
)

type life struct {
	colour              framebuffer.Rgb
	framesPerGeneration uint
	rule                int
	autoRepeat          bool
	repeatGeneration	uint
	cachedGeneration    uint
	cachedState         []bool
	tempState           []bool
}

func newLife(colour framebuffer.Rgb, framesPerGeneration uint, rule int, autoRepeat bool) *life {

	return &life{colour: colour, framesPerGeneration: framesPerGeneration, rule: rule, autoRepeat: autoRepeat}
}

func (l *life) animateFrame(frameCount uint, frame segment.Segment) {

	// Require at least three leds
	if frame.Len() < 3 {
		return
	}

	// Fill cache with generation zero state if its empty
	if l.cachedState == nil {
		l.cachedState = getGenerationZeroState(frame.Len())
		l.tempState = make([]bool, frame.Len(), frame.Len())
	}

	// Don't support change of frame length
	if uint(len(l.cachedState)) != frame.Len() {
		log.Fatal("attempt to change frame length in life animation")
	}

	// Work out the required generation
	generation := frameCount / l.framesPerGeneration
	if l.repeatGeneration != 0 {
		generation = generation % l.repeatGeneration
	}

	// Get the state for that generation
	newState := l.getGenerationState(generation)

	updateFrameBuffer(frame, l.colour, newState)
}

func (l *life) getGenerationState(generation uint) []bool {

	if l.cachedGeneration == generation {
		return l.cachedState
	}

	if generation == 0 {
		l.cachedGeneration = 0;
		l.cachedState = getGenerationZeroState(uint(len(l.cachedState)))
		return l.cachedState
	}

	if l.cachedGeneration > generation {
		log.Fatal("attempt to regress generation in life animation")
	}

	// Move to the required generation
	for l.cachedGeneration < generation {
		getNextGenerationState(l.rule, l.cachedState, l.tempState)
		l.cachedState, l.tempState = l.tempState, l.cachedState
		l.cachedGeneration++
	}

	// If the state hasn't changed we may need to auto repeat
	if l.autoRepeat && areStatesEqual(l.cachedState, l.tempState) {
		l.cachedGeneration = 0;
		l.cachedState = getGenerationZeroState(uint(len(l.cachedState)))
		l.repeatGeneration = generation
	}

	return l.cachedState
}

func areStatesEqual(s1 []bool, s2 []bool) bool {

	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func getNextGenerationState(rule int, currentState []bool, nextState []bool) {

	// Iterate over the last frame creating next frame
	for i := 1; i < len(currentState)-1; i++ {
		// Calculate a value representing the three cells
		// that influence the current cell
		// https://www.wolframalpha.com/input/?i=rule+30
		var index uint
		if currentState[i-1] {
			index = 4
		}
		if currentState[i] {
			index += 2
		}
		if currentState[i+1] {
			index += 1
		}
		// Check bit in rule number corresponding to index
		nextState[i] = rule&(1<<index) != 0
	}
}

func getGenerationZeroState(length uint) []bool {

	// Add single true value in the middle
	// TODO support different starting conditions?
	s := make([]bool, length, length)
	s[length/2] = true
	return s
}

func updateFrameBuffer(frame segment.Segment, colour framebuffer.Rgb, src []bool) {

	for i := uint(0); i < frame.Len(); i++ {
		if src[i] {
			frame.Set(i, colour)
		}
	}
}
