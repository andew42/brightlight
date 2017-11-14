package servers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/hue"
)

// Handles light control commands from hue bridge (runs as go routine)
func HueAnimationHandler(commandChan chan interface{}, segmentNames []string) {

	// Start with all available segments turned off
	segmentAnimations := make([]animations.SegmentAction, 0)
	for _, s := range segmentNames {
		segmentAnimations = append(segmentAnimations, animations.SegmentAction{Segment: s, Animation: "Off", Params: "FFFFFF"})
	}

	// Loop processing a command per iteration
	for {
		command := <-commandChan
		log.WithField("command", command).Info("Command from hue bridge to brightlight")
		switch command.(type) {
		case hue.SegmentControl:
			processSegmentControl(segmentAnimations, command.(hue.SegmentControl))
		case hue.PresetControl:
			processPresetControl(segmentAnimations, command.(hue.PresetControl))
		default:
			log.Fatal("Invalid command type from hue bridge")
		}
		animations.RunAnimations(segmentAnimations)
	}
}

func processSegmentControl(segmentAnimations []animations.SegmentAction, command hue.SegmentControl) {

	i := findSegmentIndex(segmentAnimations, command.Name)
	if i == -1 {
		log.WithField("name", command.Name).Error("Invalid segment name")
		return
	}
	if command.On {
		segmentAnimations[i].Animation = "Static"
	} else {
		segmentAnimations[i].Animation = "Off"
	}
}

func processPresetControl(segmentAnimations []animations.SegmentAction, command hue.PresetControl) {
	// TODO
}

// Return the index to the found named animation in animations or -1 if not found
func findSegmentIndex(animations []animations.SegmentAction, name string) int {

	for i, a := range animations {
		if a.Segment == name {
			return i
		}
	}
	return -1
}
