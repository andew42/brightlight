package servers

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"github.com/andew42/brightlight/controller"
	"encoding/json"
)

// Handle frame buffer web socket requests (web socket is closed when we return)
func GetFrameBufferHandler(fb *controller.FrameBuffer) (func(ws *websocket.Conn)) {

	return func(ws *websocket.Conn) {

		for {
			fb.Mutex.Lock()
			// Fails if the client has disappeared
			if err := sendFrameBufferToWebSocket(fb, ws); err != nil {
				fb.Mutex.Unlock()
				log.Info("frameBufferSocketHandler " + err.Error())
				return
			}
			// Wait for next frame buffer update
			fb.Cond.Wait()
			fb.Mutex.Unlock()
		}
	}
}

// Render the frame buffer as JSON
func sendFrameBufferToWebSocket(fb *controller.FrameBuffer, ws *websocket.Conn) error {

	// Send back the frame buffer as JSON
	rc, err := json.MarshalIndent(fb, "", " ")
	if err != nil {
		return err
	}
	_, err = ws.Write(rc)
	return err
}
