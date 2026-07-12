package servers

import (
	"net/http"
	"strconv"
	"strings"

	"log/slog"

	"github.com/andew42/brightlight/animations"
)

// StripLenHandler Handle HTTP requests to show strip lengths of room lights
func StripLenHandler(w http.ResponseWriter, r *http.Request) {

	// Strip index, length follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No parameters specified", 406)
		slog.Info("stripLengthHandler called with no parameters")
		return
	}

	configStrings := strings.Split(r.URL.Path[extIndex+1:], ",")
	if len(configStrings) != 2 {
		http.Error(w, "Expected index,length parameters", 406)
		slog.Info("stripLengthHandler called with wrong parameter count")
		return
	}

	// index
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
	slog.Info("stripLengthHandler called", "index", index, "length", length, "err", err)
}
