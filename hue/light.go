package hue

import (
	"encoding/json"
	"strconv"
	"time"
)

// HTTP Handler for /api/lights
func processLightsRequest(c *cmdContext) {

	if c.Method == "GET" {
		if len(c.Resource) == 0 {
			// Return lights
			respondWithJsonEncodedObject(c.W, &c.FullState.Lights)
			return
		}
		if len(c.Resource) == 1 && c.Resource[0] == "new" {
			// No new lights
			respondWithJsonEncodedObject(c.W, struct {
				LastScan string `json:"lastscan"`
			}{
				LastScan: time.Now().Format("2006-01-02T15:04:05"),
			})
			return
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}
	if c.Method == "PUT" {
		setLight(c)
		return
	}
	reportError(c.W, newApiErrorMethodNotAvailable(c.ErrorAddress))
}

type light struct {
	State             lightstate `json:"state"`
	LightType         string     `json:"type"`
	Name              string     `json:"name"`
	ModelId           string     `json:"modelid"`
	UniqueId          string     `json:"uniqueid"`
	ManufacturerName  string     `json:"manufacturername"`
	LuminaireUniqueId string     `json:"luminaireuniqueid"`
	SwVersion         string     `json:"swversion"`
}

type lightstate struct {
	On             bool       `json:"on"`
	Bri            uint8      `json:"bri"`
	Hue            uint16     `json:"hue"`
	Sat            uint8      `json:"sat"`
	Xy             [2]float32 `json:"xy"`
	Ct             uint16     `json:"ct"`
	Alert          string     `json:"alert"`
	Effect         string     `json:"effect"`
	TransitionTime uint16     `json:"transitiontime"`
	ColorMode      string     `json:"colormode"`
	Reachable      bool       `json:"reachable"`
}

func nextLightId(lights map[string]*light) string {

	nextId := 0
	for k, _ := range lights {
		if i, err := strconv.Atoi(k); err == nil {
			if i > nextId {
				nextId = i
			}
		}
	}
	return strconv.Itoa(nextId + 1)
}

// Check light ids all exist
func doAllLightsExist(test *[]string, lights *map[string]*light) bool {

	// Check lights exist, just log if they don't
	for _, l := range *test {
		_, ok := (*lights)[l]
		if !ok {
			return false
		}
	}
	return true
}

// PUT /api/<username>/lights/<id>/state (modifying light attribute)
// PUT /api/<username>/lights/<id> (modifying the light name)
func setLight(c *cmdContext) {

	if !(len(c.Resource) == 1 || (len(c.Resource) == 2 && c.Resource[1] == "state")) {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Extract the light to modify
	l, ok := c.FullState.Lights[c.Resource[0]]
	if !ok {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Perform modification
	if len(c.Resource) == 1 {
		setLightAttributes(c, l)
	} else {
		setLightState(c, l)
	}
}

// Rename
// Url=/api/4d65822107fcfd5278629a0f5f3f164f/lights/1
// {"name":"Bedroom Light"}
// [{"success":{"/lights/1/name":"Bedroom Light"}}]
func setLightAttributes(c *cmdContext, l *light) {

	// TODO
	reportError(c.W, newApiErrorNotImplemented(c.ErrorAddress))
}

// Url=/api/4d65822107fcfd5278629a0f5f3f164f/lights/1/state
// Body={"on":true, "bri_inc":1}
//[
//{"success":{"/lights/1/state/on":true}},
//{"success":{"/lights/1/state/bri":200}}
//]
func setLightState(c *cmdContext, l *light) {

	// Decode body as a light state
	var newLs lightstate
	if err := json.Unmarshal(c.Body, &newLs); err != nil {
		reportError(c.W, newApiErrorBodyContainsInvalidJson(c.ErrorAddress, err))
		return
	}

	// Decode the body as a map so we can determine what was set
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(c.Body, &bodyMap); err != nil {
		reportError(c.W, newApiErrorBodyContainsInvalidJson(c.ErrorAddress, err))
		return
	}

	response := apiResponseList{}

	// Update light state fields from body
	if _, ok := bodyMap["on"]; ok {
		l.State.On = newLs.On
		response = response.AppendSuccessResponse(c.ResourceUrl+"/on", l.State.On)
		// Action command
		c.BrightlightCommandChan <- SegmentControl{On: l.State.On, Name: l.Name}
	}

	if _, ok := bodyMap["bri"]; ok {
		l.State.Bri = newLs.Bri
		response = response.AppendSuccessResponse(c.ResourceUrl+"/bri", l.State.Bri)
		// TODO: Action command
	}

	// TODO handle all other properties

	// TODO handle the *_inc properties

	// TODO Action light state change

	// Return response
	respondWithJsonEncodedObject(c.W, response)
}
