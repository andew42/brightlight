package hue

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

// Respond to the hue bridge API requests
func ApiHandler(w http.ResponseWriter, r *http.Request) {

	// Dump API request
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		log.WithFields(log.Fields{
			"Body": string(b), "Header": r.Header, "Url": r.URL}).Info("Hue Api Request")
	} else {
		log.WithField("Error", err).Error("Failed to read Hue Api request body")
	}

	// TODO
	//w.Header().Set("Content-Type", "application/xml")
	//l, err := w.Write(setupUrlContent)
	//if err != nil || l != len(setupUrlContent) {
	//	log.WithField("error", err).Warn("Hue Bridge Emulator failed to serve http request")
	//}
}
