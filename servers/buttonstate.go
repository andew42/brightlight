package servers

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"sync"
)

var listenersMux sync.Mutex
var listeners []chan int

// Called by run animations server to indicate current button pressed
func updateButtonState(key int) {

	listenersMux.Lock()
	defer listenersMux.Unlock()
	for _, l := range listeners {
		l <- key
	}
	currentButtonKey = key
}

// Called by each web socket go routine to append a listener
func addButtonListener() chan int {

	listenersMux.Lock()
	defer listenersMux.Unlock()
	c := make(chan int)
	listeners = append(listeners, c)
	return c
}

// Called when a web socket closes to remove its listener
func removeButtonListener(c chan int) {

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

// Current button state so we can immediately tell new listeners
var currentButtonKey = 0

// Handle button state web socket requests (web socket is closed when
// we return) We have one of these go routines per web socket request
func ButtonStateHandler(ws *websocket.Conn) {

	// Not thread safe but good enough for debug output
	buttonStateListenerId++
	listenerId := buttonStateListenerId
	log.WithField("id", listenerId).Info("adding button state listener")

	// Add our listener
	src := addButtonListener()

	// Send the current state immediately
	if err := sendButtonStateToWebSocket(listenerId, currentButtonKey, ws, src); err != nil {
		return
	}

	// Watch for client close operations by setting up a read go routine, we
	// never expect anything from the client but the read fails on close
	// https://groups.google.com/forum/#!topic/golang-nuts/pXNSBx4wgAw
	close := make(chan int)
	go func() {
		websocket.Message.Receive(ws, nil)
		close <- 0
	}()

	for {
		select {
		case bs := <-src: // src sends us button state updates
			if err := sendButtonStateToWebSocket(listenerId, bs, ws, src); err != nil {
				return
			}
		case <-close: // close sends us read errors (i.e. socket closed by client)
			log.WithField("id", listenerId).Info("closing button state listener")
			removeButtonListener(src)
			return;
		}
	}
}

// Render button state key as JSON
func sendButtonStateToWebSocket(listenerId int, bs int, ws *websocket.Conn, c chan int) error {

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
