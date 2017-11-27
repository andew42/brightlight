package config

import (
	"path"
	"io/ioutil"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

// A settings file is a map of named button column definitions
type SettingsFileDef map[string]ButtonColumnDef

// A button column definition is a slice of button definitions
type ButtonColumnDef []ButtonDef

// A button definition is a named slice of segments
type ButtonDef struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Segments []SegmentDef `json:"segments"`
}

type SegmentDef struct {
	Name      string `json:"segment"`
	Animation string `json:"animation"`
	Params    string `json:"params"`
}

type Preset struct {
	Name string
}

// Loads the user config json file containing preset (button) definitions set up by the user
func LoadUserPresets(configPath string) []Preset {

	presets := make([]Preset, 0)

	// Try loading user settings first
	fileContent, err := ioutil.ReadFile(path.Join(configPath, "/config/user.json"))
	if err != nil {
		log.WithField("error", err).Warn("Failed to open user.json")
		// No user settings, try loading the defaults
		if fileContent, err = ioutil.ReadFile(path.Join(configPath, "/config/default.json")); err != nil {
			log.WithField("error", err).Error("Failed to open default.json")
			// Return an empty preset list
			return presets
		}
	}

	// Parse the setting file
	var settings SettingsFileDef
	if err = json.Unmarshal(fileContent, &settings); err != nil {
		log.WithField("error", err).Error("Failed to unmarshal settings json")
		// Return an empty preset list
		return presets
	}

	// Return list of presets (button names)
	for _, col := range settings {
		for _, button := range col {
			presets = append(presets, Preset{Name: button.Name})
		}
	}
	return presets
}
