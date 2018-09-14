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
	// {
	// 	 "key":3,
	//	 "name":"Sweet Shop",
	//	 "segments":[
	//	   {
	//	     "name":"All",
	//	     "z":1,
	//	     "animation":"Sweet Shop",
	//	     "params":[
	//	       {"key":60, "type":"range", "label":"Duration(frames)", "min":1, "max":100, "value":25},
	//	       {"key":61, "type":"range", "label":"Brightness", "min":20, "max":60, "value":50},
	//	       {"key":62, "type":"range", "label":"Min Saturation", "min":0, "max":99, "value":50}
	//	     ]
	//	   }
	//   ]
	// },

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithField("err", err.Error()).Error("RunAnimationsHandler bad body")
		http.Error(w, err.Error(), 400)
		return
	}

	// Un-marshal JSON
	var button animations.Button
	if err = json.Unmarshal(body, &button); err != nil {
		log.WithField("err", err.Error()).Error("RunAnimationsHandler bad body JSON")
		http.Error(w, err.Error(), 400)
		return
	}
	log.WithField("Decoded JSON", button).Info("RunAnimationsHandler called")

	// Perform the animation
	animations.RunAnimations(button.Segments)

	// Update button state (i.e. the button key for the animation we are running)
	updateButtonState(button.Key)

	// Return controller status
	allConnected := true
	for _, v := range controller.TeensyConnections() {
		allConnected = allConnected && v
	}
	d, _ := json.Marshal(allConnected)
	w.Write(d)
}
