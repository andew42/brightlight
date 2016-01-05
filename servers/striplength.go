package servers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"net/http"
	"strconv"
	"strings"
)

func GetStripLenHandler(numberOfStrips int) func(http.ResponseWriter, *http.Request) {

	// Handle HTTP requests to show strip lengths of room lights
	return func(w http.ResponseWriter, r *http.Request) {

		// Strip index, length follows request path
		extIndex := strings.LastIndex(r.URL.Path, `/`)
		if extIndex == -1 {
			http.Error(w, "No parameters specified", 406)
			log.Info("stripLengthHandler called with no parameters")
			return
		}

		// index
		config := strings.Split(r.URL.Path[extIndex+1:], ",")
		index, err := strconv.ParseInt(config[0], 10, 32)
		if err != nil || index < 0 || index > int64(numberOfStrips) {
			index = -1
		}

		// length
		length, err := strconv.ParseInt(config[1], 10, 32)
		if err != nil || length < 0 || length > controller.MaxLedStripLen {
			length = -1
		}

		if err = animations.AnimateStripLength(uint(index), uint(length)); err != nil {
			http.Error(w, "Invalid index or length", 406)
		}

		log.WithFields(log.Fields{"index": index, "length": length, "err": err}).Info("stripLengthHandler called")
	}
}
