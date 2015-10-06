package servers

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strings"
	"encoding/json"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/animations"
)

// Handle HTTP requests to run a named animation
func AnimationHandler(w http.ResponseWriter, r *http.Request) {

	// Animation name follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No animation name specified", 406)
		log.Info("animationHandler no name")
		return
	}

	animationName := r.URL.Path[extIndex+1:]
	err := animations.Animate(animationName)
	if err != nil {
		http.Error(w, err.Error()+" "+animationName, 406)
		log.Info("animationHandler " + err.Error())
		return
	}

	log.WithField("animationName", animationName).Info("animationHandler")

	d, _ := json.Marshal(controller.IsDriverConnected())
	w.Write(d)
}
