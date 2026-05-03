package servers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

// Track the currently active button and version of the button pad save
// Allows multiple UIs to refresh active button and button pad contents
// when it changes
type buttonState struct {
	ActiveButtonKey  int
	ButtonPadVersion int
}

var listenersMux sync.Mutex
var listeners []chan buttonState
var currentButtonState buttonState

// Called by run animations server to indicate current button pressed
func updateActiveButtonKey(key int) {
	updateCurrentButtonState(func() {
		currentButtonState.ActiveButtonKey = key
	})
}

// Called by config server to update button pad save version
func updateButtonPadVersion(ver int) {
	updateCurrentButtonState(func() {
		currentButtonState.ButtonPadVersion = ver
	})
}

func updateCurrentButtonState(f func()) {
	listenersMux.Lock()
	defer listenersMux.Unlock()
	f()
	for _, l := range listeners {
		l <- currentButtonState
	}
}

// Called when a web socket closes to remove its listener
func removeButtonListener(c chan buttonState) {
	listenersMux.Lock()
	defer listenersMux.Unlock()
	for i, l := range listeners {
		if l == c {
			// https://stackoverflow.com/questions/37334119
			listeners[len(listeners)-1], listeners[i] = listeners[i], listeners[len(listeners)-1]
			listeners = listeners[:len(listeners)-1]
			return
		}
	}
}

// Give each button state listener its own unique ID (for logging)
var buttonStateListenerId = 0

// ButtonStateHandler Handle button state web socket requests (web socket is closed
// when we return) We have one of these go routines per web socket request
func ButtonStateHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	listenersMux.Lock()

	// Create an id for this listener go routine
	buttonStateListenerId++
	listenerId := buttonStateListenerId
	slog.Info("adding button state listener", "id", listenerId)

	// Add our listener channel
	update := make(chan buttonState)
	listeners = append(listeners, update)

	// Copy the current button state
	bs := currentButtonState
	listenersMux.Unlock()

	// Send the current state immediately
	sendEvent := func(bs buttonState) error {
		slog.Info("sending button state", "id", listenerId, "state", bs)
		data, err := json.Marshal(bs)
		if err != nil {
			return err
		}
		if _, err = fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
			return err
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return nil
	}

	if err := sendEvent(bs); err != nil {
		removeButtonListener(update)
		return
	}

	for {
		select {
		case bs := <-update:
			if err := sendEvent(bs); err != nil {
				slog.Info("button state listener write error", "id", listenerId, "err", err)
				removeButtonListener(update)
				return
			}
		case <-r.Context().Done():
			slog.Info("closing button state listener", "id", listenerId)
			removeButtonListener(update)
			return
		}
	}
}
