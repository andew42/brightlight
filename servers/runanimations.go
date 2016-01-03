package servers

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
)

// Handle HTTP requests to run zero or more animation specified in json payload
func RunAnimationsHandler(w http.ResponseWriter, r *http.Request) {

	// JSON body of form
	// [{"segmentId": "s1", "action": "static", "params": "6f16d4"},
	//  {"segmentId": "s2", "action": "static", "params": "6f16d4"}]}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithField("err", err.Error()).Error("RunAnimationsHandler bad body")
		http.Error(w, err.Error(), 400)
		return;
	}

	// Unmarshal JSON into typed slice
	var segments []animations.SegmentAction

	if err = json.Unmarshal(body, &segments); err != nil {
		log.WithField("err", err.Error()).Error("RunAnimationsHandler bad body JSON")
		http.Error(w, err.Error(), 400)
		return;
	}

	log.WithField("Decoded JSON", segments).Info("RunAnimationsHandler called")

	// Perform the animation
	animations.RunAnimations(segments)

	// Return controller status
	d, _ := json.Marshal(controller.IsDriverConnected())
	w.Write(d)
}
