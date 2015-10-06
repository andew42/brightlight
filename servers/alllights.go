package servers

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/animations"
)

// Handle HTTP requests to set all lights to a specific colour
func AllLightsHandler(w http.ResponseWriter, r *http.Request) {

	// Colour value follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No colour specified", 406)
		log.Info("allLightsHandler no colour")
		return
	}

	colourValue, err := strconv.ParseInt(r.URL.Path[extIndex+1:], 16, 32)
	if err != nil {
		http.Error(w, "Invalid colour specified", 406)
		log.Info("allLightsHandler invalid colour")
		return
	}

	colourValueRgb := controller.NewRgbFromInt(int(colourValue))
	log.WithField("colourValueRgb", colourValueRgb).Info("allLightsHandler")
	animations.AnimateStaticColour(colourValueRgb)

	d, _ := json.Marshal(controller.IsDriverConnected())
	w.Write(d)
}
