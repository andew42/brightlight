package servers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/andew42/brightlight/framebuffer"
)

// Give each virtual frame buffer its own unique ID
var frameBufferListenerId = 0

func FrameBufferHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	frameBufferListenerId++
	src, done := framebuffer.AddListener("Virtual Frame Buffer "+strconv.Itoa(frameBufferListenerId), false)

	for {
		select {
		case fb := <-src:
			data, err := json.Marshal(fb)
			if err != nil {
				slog.Info("frameBufferHandler marshal error", "err", err)
				done <- src
				return
			}
			if _, err = fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				slog.Info("frameBufferHandler write error", "err", err)
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
