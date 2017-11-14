package hue

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"os"
)

type fullState struct {
	savePath      string               `json:"-"`
	Lights        map[string]*light    `json:"lights"`
	Groups        map[string]*group    `json:"groups"`
	Config        *config              `json:"config"`
	Schedules     map[string]*schedule `json:"schedules"`
	Scenes        map[string]*scene    `json:"scenes"`
	Sensors       map[string]*sensor   `json:"sensors"`
	Rules         map[string]*string   `json:"rules"`
	ResourceLinks map[string]*string   `json:"resourcelinks"`
}

func newFullState(savePath string) *fullState {

	return &fullState{
		savePath:      savePath,
		Lights:        make(map[string]*light),
		Groups:        make(map[string]*group),
		Config:        newConfig(),
		Schedules:     make(map[string]*schedule),
		Scenes:        make(map[string]*scene),
		Sensors:       make(map[string]*sensor),
		Rules:         make(map[string]*string),
		ResourceLinks: make(map[string]*string),
	}
}

// Retrieve bridge state from disc if possible otherwise return a new default one
func getFullState(persistedHueStatePath string, nii NetworkInterfaceInfo) *fullState {

	var err error
	var f *os.File
	if f, err = os.Open(persistedHueStatePath); err == nil {
		defer f.Close()
		// Here we opened the file, attempt to deserialize
		var fs fullState
		decoder := json.NewDecoder(f)
		if err := decoder.Decode(&fs); err == nil {
			fs.savePath = persistedHueStatePath
			fs.Config.setNetworkInfo(nii)
			return &fs
		}
	}
	logrus.WithField("error", err).
		Warn("hue API failed to load state, using default state")
	fs := newFullState(persistedHueStatePath)
	fs.Config.setNetworkInfo(nii)
	return fs
}

// Persist hue state to disc
func (fs *fullState) Save() error {

	var err error
	if f, err := os.Create(fs.savePath); err == nil {
		defer f.Close()
		// Here we created the file, attempt to serialize
		encoder := json.NewEncoder(f)
		err = encoder.Encode(fs)
	}

	if err != nil {
		logrus.WithField("error", err).Error("Hue API failed to save state")
	}

	return err
}
