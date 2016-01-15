package servers

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/stats"
	"golang.org/x/net/websocket"
	"strconv"
)

// Give each stats listener its own unique ID
var statsListenerId = 0

// Handle stats web socket requests (web socket is closed when we return)
func GetStatsHandler() func(ws *websocket.Conn) {

	return func(ws *websocket.Conn) {
		// Not thread safe but good enough for debug output
		statsListenerId++
		src, done := stats.AddListener("Stats Listener " + strconv.Itoa(statsListenerId))
		for {
			select {
			// src sends us stats updates
			case stats := <-src:
				// Render the stats as JSON (fails if the client has disappeared)
				if rc, err := json.MarshalIndent(stats, "", " "); err == nil {
					_, err = ws.Write(rc)
				} else {
					log.Info("statsSocketHandler" + err.Error())
					// Un-subscribe before returning and closing connection
					done <- src
					return
				}
			}
		}
	}
}
