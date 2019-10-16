package servers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
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

// Handle button state web socket requests (web socket is closed when
// we return) We have one of these go routines per web socket request
func ButtonStateHandler(ws *websocket.Conn) {

	listenersMux.Lock()

	// Create an id for this listener go routine
	buttonStateListenerId++
	listenerId := buttonStateListenerId
	log.WithField("id", listenerId).Info("adding button state listener")

	// Add our listener channel
	update := make(chan buttonState)
	listeners = append(listeners, update)

	// Copy the current button state
	buttonState := currentButtonState

	listenersMux.Unlock()

	// Send the current state immediately
	if err := sendButtonStateToWebSocket(listenerId, buttonState, ws, update); err != nil {
		return
	}

	// Watch for client closeSocket operations by setting up a read go routine, we
	// never expect anything from the client but the read fails on closeSocket
	// https://groups.google.com/forum/#!topic/golang-nuts/pXNSBx4wgAw
	closeSocket := make(chan int)
	go func() {
		websocket.Message.Receive(ws, nil)
		closeSocket <- 0
	}()

	for {
		select {
		case bs := <-update: // update sends us button state updates
			if err := sendButtonStateToWebSocket(listenerId, bs, ws, update); err != nil {
				return
			}
		case <-closeSocket: // closeSocket sends us read errors (i.e. socket closed by client)
			log.WithField("id", listenerId).Info("closing button state listener")
			removeButtonListener(update)
			return
		}
	}
}

// Render button state key as JSON
func sendButtonStateToWebSocket(listenerId int, bs buttonState, ws *websocket.Conn, c chan buttonState) error {

	log.WithFields(log.Fields{"id": listenerId, "state": bs}).Info("sending button state")

	// Send back the frame buffer as JSON
	rc, err := json.MarshalIndent(bs, "", " ")
	if err == nil {
		_, err = ws.Write(rc)
	}

	if err != nil {
		log.Info("buttonStateSocketHandler " + err.Error())
		// Un-subscribe before returning and closing connection
		removeButtonListener(c)
	}

	return err
}
