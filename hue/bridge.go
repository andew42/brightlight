package hue

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"path"
)

type NetworkInterfaceInfo struct {
	Ip      string // 192.168.0.2
	Mask    string // 255.255.255.0
	Gateway string // 192.168.0.1
	Port    string // :80
	Mac     string // b8:e8:59:2a:94:7a
}

func StartHueBridgeEmulator(nii NetworkInterfaceInfo, contentPath string,
	brightlightUpdateChan chan interface{}, brightlightCommandChan chan interface{}) error {

	// Handler for hue emulator discovery
	dh, err := GetUpnpDescriptionHandler(nii.Ip + nii.Port)
	if err != nil {
		log.WithField("Error", err).
			Error("Hue API failed to create description.xml handler")
		return err
	}

	// Retrieve full state of the bridge emulator (or default on first call)
	persistedHueStatePath := path.Join(contentPath, "/config/hue.json")
	fsl := NewFullStateLocker(persistedHueStatePath, nii)

	// Handler for hue emulator API
	apiHandler, err := GetApiHandler(fsl, brightlightCommandChan)
	if err != nil {
		log.WithField("Error", err).
			Error("Hue API handler couldn't be created")
		return err
	}

	// Start hue emulator components
	go BrightlightUpdateHandler(fsl, brightlightUpdateChan)
	StartUpnpResponder(nii.Ip + nii.Port)
	http.HandleFunc("/description.xml", dh)
	http.HandleFunc("/api/", apiHandler)
	log.Info("Hue API bridge emulation started")
	return nil
}
