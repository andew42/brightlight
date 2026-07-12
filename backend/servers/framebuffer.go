package servers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/andew42/brightlight/framebuffer"
)

// Give each virtual frame buffer its own unique ID (atomic as handlers run concurrently)
var frameBufferListenerId atomic.Int64

func FrameBufferHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	src, done := framebuffer.AddListener("Virtual Frame Buffer "+strconv.FormatInt(frameBufferListenerId.Add(1), 10), false)

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
