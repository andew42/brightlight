package servers

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"github.com/andew42/brightlight/controller"
	"encoding/json"
	"github.com/andew42/brightlight/stats"
	"runtime/debug"
)

var gcStats debug.GCStats

// Handle stats web socket requests (web socket is closed when we return)
func GetStatsHandler(statistics *stats.Stats, fb *controller.FrameBuffer) (func(ws *websocket.Conn)) {

	return func(ws *websocket.Conn) {

		for {
			fb.Mutex.Lock()

			// Report stats for last second every second (50 frames)
			if statistics.FrameCount == stats.ResetFrame {
				// Update garbage collection count
				debug.ReadGCStats(&gcStats)
				statistics.GcCount = gcStats.NumGC
				// Render the stats as JSON (fails if the client has disappeared)
				rc, err := json.MarshalIndent(statistics, "", " ")
				if err == nil {
					_, err = ws.Write(rc)
				}
				if err != nil {
					fb.Mutex.Unlock()
					log.Info("statsSocketHandler" + err.Error())
					return
				}
			}
			// Wait for next frame buffer update
			fb.Cond.Wait()
			fb.Mutex.Unlock()
		}
	}
}
