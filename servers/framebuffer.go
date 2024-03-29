package servers

import (
	"encoding/json"
	"github.com/andew42/brightlight/framebuffer"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"strconv"
)

// Give each virtual frame buffer its own unique ID
var frameBufferListenerId = 0

// FrameBufferHandler Handle frame buffer web socket requests (web socket is closed when we return)
func FrameBufferHandler(ws *websocket.Conn) {

	// Not thread safe but good enough for debug output
	frameBufferListenerId++
	src, done := framebuffer.AddListener("Virtual Frame Buffer "+strconv.Itoa(frameBufferListenerId), false)
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
