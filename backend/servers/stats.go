package servers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/andew42/brightlight/stats"
)

// Give each stats listener its own unique ID
var statsListenerId = 0

func StatsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	statsListenerId++
	src, done := stats.AddListener("Stats Listener " + strconv.Itoa(statsListenerId))

	for {
		select {
		case statsUpdate := <-src:
			data, err := json.Marshal(statsUpdate)
			if err != nil {
				slog.Info("statsHandler marshal error", "err", err)
				done <- src
				return
			}
			if _, err = fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				slog.Info("statsHandler write error", "err", err)
				done <- src
				return
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			done <- src
			return
		}
	}
}
