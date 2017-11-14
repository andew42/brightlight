package hue

// Represents a physical (or logical) hue luminary (bulb or collection of bulbs)
type Luminary struct {
	Type              string
	Name              string
	ModelId           string
	UniqueId          string
	ManufacturerName  string
	LuminaireUniqueId string
	SwVersion         string
}

// Add (or update exiting) luminary
func addOrUpdateLuminary(fs *fullState, luminary Luminary) {

	// Update existing light if it exists
	for k, v := range fs.Lights {
		if v.UniqueId == luminary.UniqueId {
			fs.Lights[k] = lightFromLuminary(luminary)
			return
		}
	}
	// Add new light
	fs.Lights[nextLightId(fs.Lights)] = lightFromLuminary(luminary)
}

func lightFromLuminary(luminary Luminary) *light {

	return &light{
		State: lightstate{
			Alert:          "none",
			Effect:         "none",
			ColorMode:      "xy",
			TransitionTime: 4,
			Reachable:      true,
		},
		LightType:         "LED Strip",
		Name:              luminary.Name,
		ModelId:           "LST001",
		UniqueId:          luminary.UniqueId,
		ManufacturerName:  luminary.ManufacturerName,
		LuminaireUniqueId: luminary.LuminaireUniqueId,
		SwVersion:         luminary.SwVersion,
	}
}

// Remove existing luminary
func removeLuminary(fs *fullState, luminary Luminary) {
	// TODO
}
