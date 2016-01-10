package servers

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/framebuffer"
	"golang.org/x/net/websocket"
	"strconv"
)

// Give each virtual frame buffer its own unique ID
var listenerId = 0

// Handle frame buffer web socket requests (web socket is closed when we return)
func GetFrameBufferHandler() func(ws *websocket.Conn) {

	return func(ws *websocket.Conn) {
		// Not thread safe but good enough for debug output
		listenerId++
		src, done := framebuffer.AddListener("Virtual Frame Buffer " + strconv.Itoa(listenerId))
		for {
			select {
			// src sends us frame buffer updates
			case fb := <-src:
				// Fails if the client has disappeared
				if err := sendFrameBufferToWebSocket(fb, ws); err != nil {
					log.Info("frameBufferSocketHandler " + err.Error())
					// Un-subscribe before returning and closing connection
					done <- src
					return
				}
			}
		}
	}
}

// Render the frame buffer as JSON
func sendFrameBufferToWebSocket(fb *framebuffer.FrameBuffer, ws *websocket.Conn) error {

	// Send back the frame buffer as JSON
	rc, err := json.MarshalIndent(fb, "", " ")
	if err != nil {
		return err
	}
	_, err = ws.Write(rc)
	return err
}
