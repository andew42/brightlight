package servers

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"io/ioutil"
	"net/http"
)

// Handle HTTP requests to run zero or more animation specified in json payload
func RunAnimationsHandler(w http.ResponseWriter, r *http.Request) {

	// JSON body of form
	// [
	//	 {"name":"Bedroom", "animation":"Sweet Shop", "params":[{"key":60, "type":"speed", "value":50}]},
	//	 {"name":"Bathroom", "animation":"Rainbow", "params":[{"key":40, "type":"speed", "value":50}]}
	// ]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithField("err", err.Error()).Error("RunAnimationsHandler bad body")
		http.Error(w, err.Error(), 400)
		return
	}

	// Un-marshal JSON into typed slice
	var segments []animations.SegmentAction

	if err = json.Unmarshal(body, &segments); err != nil {
		log.WithField("err", err.Error()).Error("RunAnimationsHandler bad body JSON")
		http.Error(w, err.Error(), 400)
		return
	}

	log.WithField("Decoded JSON", segments).Info("RunAnimationsHandler called")

	// Perform the animation
	animations.RunAnimations(segments)

	// Return controller status
	allConnected := true
	for _, v := range controller.TeensyConnections() {
		allConnected = allConnected && v
	}
	d, _ := json.Marshal(allConnected)
	w.Write(d)
}
