package servers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/animations"
	"net/http"
	"strconv"
	"strings"
)

// Handle HTTP requests to show strip lengths of room lights
func StripLenHandler(w http.ResponseWriter, r *http.Request) {

	// Strip index, length follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No parameters specified", 406)
		log.Info("stripLengthHandler called with no parameters")
		return
	}

	// index
	configStrings := strings.Split(r.URL.Path[extIndex+1:], ",")
	index, err := strconv.ParseInt(configStrings[0], 10, 32)
	if err != nil {
		index = -1
	}

	// length
	length, err := strconv.ParseInt(configStrings[1], 10, 32)
	if err != nil {
		length = -1
	}

	animations.AnimateStripLength(uint(index), uint(length))
	log.WithFields(log.Fields{"index": index, "length": length, "err": err}).Info("stripLengthHandler called")
}
