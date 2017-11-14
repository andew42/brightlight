package hue

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
)

// Update sent from brightlight to Hue bridge
// These are brightlight segments which are mapped to virtual hue bulbs
type SegmentUpdate struct {
	OldName string // Set to delete or rename
	NewName string // Set to add or rename
}

// Update sent from brightlight to Hue bridge
// These are brightlight preset buttons which are mapped to virtual hue rooms or scenes
type PresetUpdate struct {
	OldName string // Set to delete or rename
	NewName string // Set to add or rename
}

// Commands sent from the Hue bridge to brightlight
// Turn a segment on or off
type SegmentControl struct {
	On   bool
	Name string
}

// Commands sent from the Hue bridge to brightlight
// Turn a preset button on or off
type PresetControl struct {
	On   bool
	Name string
}

// Loops processing updates from brightlight
func BrightlightUpdateHandler(fsl *fullStateLocker, brightlightUpdate chan interface{}) {

	// Delete any saved lights as they are rebuild each time
	fs := fsl.Lock()
	fs.Lights = make(map[string]*light)
	fsl.Unlock()

	// Mark each new luminary with an incremental ID
	hueId := 1

	// Wait for luminary updates used to add lights
	for {
		select {
		case u := <-brightlightUpdate:
			switch u.(type) {
			case SegmentUpdate:
				processSegmentUpdate(fsl, u.(SegmentUpdate), hueId)
				hueId++
			case PresetUpdate:
				// TODO
			default:
				log.WithField("update", u).Fatal("Unexpected update type")
			}
			//			processLuminaryUpdate(fsl, lu)
		}
	}
}

func processSegmentUpdate(fsl *fullStateLocker, su SegmentUpdate, hueId int) {

	// Only support add updates at present
	if su.NewName == "" || su.OldName != "" {
		log.WithField("segment update", su).Fatal("Unsupported update")
	}

	// TODO Just use All and Strip Three for debugging
	if su.NewName != "All" && su.NewName != "Strip Three" && su.NewName != "Strip Five"{
		return
	}

	fs := fsl.Lock()
	defer fsl.Unlock()
	addOrUpdateLuminary(fs, Luminary{
		Type:              "LED Strip",
		Name:              su.NewName,
		ModelId:           "LST001",
		UniqueId:          strconv.Itoa(hueId),
		ManufacturerName:  "LED",
		LuminaireUniqueId: strconv.Itoa(hueId) + ":4e:5b-0b", // TODO
		SwVersion:         "99999999",})
}
