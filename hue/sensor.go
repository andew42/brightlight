package hue

import "time"

// HTTP Handler for /api/sensors
func processSensorRequest(c *cmdContext) {

	if c.Method == "GET" {
		if len(c.Resource) == 0 {
			respondWithJsonEncodedObject(c.W, &c.FullState.Sensors)
			return
		}
		if len(c.Resource) == 1 {

			if c.Resource[0] == "new" {
				// No new lights
				respondWithJsonEncodedObject(c.W, struct {
					LastScan string `json:"lastscan"`
				}{
					LastScan: time.Now().Format("2006-01-02T15:04:05"),
				})
				return
			}
			if s, ok := c.FullState.Sensors[c.Resource[0]]; ok {
				respondWithJsonEncodedObject(c.W, s)
				return
			}
		}
		reportError(c.W, newApiErrorMethodNotAvailable(c.ErrorAddress))
	}
}

type sensor struct {
}
