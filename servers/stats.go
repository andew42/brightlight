package servers

import (
	//"encoding/json"
	//log "github.com/Sirupsen/logrus"
	//"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/stats"
	"golang.org/x/net/websocket"
	"runtime/debug"
)

var gcStats debug.GCStats

// Handle stats web socket requests (web socket is closed when we return)
func GetStatsHandler(statistics *stats.Stats) func(ws *websocket.Conn) {

	return func(ws *websocket.Conn) {
		// TODO REWRITE STATISTICS
		//		for {
		//			fb.Mutex.Lock()
		//
		//			// Report stats for last second every second every second
		//			if statistics.FrameCount == config.FrameFrequencyHz {
		//				// Update garbage collection count
		//				debug.ReadGCStats(&gcStats)
		//				statistics.GcCount = gcStats.NumGC
		//				// Render the stats as JSON (fails if the client has disappeared)
		//				rc, err := json.MarshalIndent(statistics, "", " ")
		//				if err == nil {
		//					_, err = ws.Write(rc)
		//				}
		//				if err != nil {
		//					fb.Mutex.Unlock()
		//					log.Info("statsSocketHandler" + err.Error())
		//					return
		//				}
		//			}
		//			// Wait for next frame buffer update
		//			fb.Cond.Wait()
		//			fb.Mutex.Unlock()
		//		}
	}
}
