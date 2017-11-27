package hue

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// HTTP Handler for /api/scenes
func processScenesRequest(c *cmdContext) {

	if c.Method == "GET" {
		if len(c.Resource) == 0 {
			respondWithJsonEncodedObject(c.W, &c.FullState.Scenes)
			return
		} else if len(c.Resource) == 1 {
			if s, ok := c.FullState.Scenes[c.Resource[0]]; ok {
				// Include light states in returned scene
				s1 := sceneWithLightstates(*s)
				if s1.LightStates == nil {
					s1.LightStates = make(map[string]lightstate)
				}
				respondWithJsonEncodedObject(c.W, s1)
				return
			}
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	} else if c.Method == "POST" {
		if len(c.Resource) == 0 {
			createScene(c)
			return
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	} else if c.Method == "PUT" {
		setScene(c)
		return
	} else if c.Method == "DELETE" {
		if len(c.Resource) == 0 {
			deleteScenes(c)
			return
		} else if len(c.Resource) == 1 {
			deleteScene(c, c.Resource[0])
			return
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}
	reportError(c.W, newApiErrorMethodNotAvailable(c.ErrorAddress))
}

type scene struct {
	Name        string   `json:"name"`
	Lights      []string `json:"lights"`
	Owner       string   `json:"owner"`
	Recycle     bool     `json:"recycle"`
	Locked      bool     `json:"locked"`
	AppData     AppData  `json:"appdata"`
	Picture     string   `json:"picture"`
	LastUpdated string   `json:"lastupdated"`
	Version     int      `json:"version"`
	// Don't send light states by default
	LightStates map[string]lightstate `json:"-"`
}

// Same as scene but with lightstates marshaled via json tag
type sceneWithLightstates struct {
	Name        string   `json:"name"`
	Lights      []string `json:"lights"`
	Owner       string   `json:"owner"`
	Recycle     bool     `json:"recycle"`
	Locked      bool     `json:"locked"`
	AppData     AppData  `json:"appdata"`
	Picture     string   `json:"picture"`
	LastUpdated string   `json:"lastupdated"`
	Version     int      `json:"version"`
	// Include light states
	LightStates map[string]lightstate `json:"lightstates"`
}

type AppData struct {
	Version int8   `json:"version"`
	Data    string `json:"data"`
}

type newSceneSuccessResponse struct {
	Id string `json:"id"`
}

// Adds a brightlight owned scene defined in brightlight button UI
func AddPresetScene(fs *fullState, name string) {

	// TODO: The idea here was to add the brightlight buttons as
	// preset scenes (I've just hardcoded All lights ID here) but
	// the iOS Hue app doesn't recognise the scene, I tried making
	// the Hue app the owner (rather than using brightlight) but
	// that didn't work. I think there is some 'secret' app data
	// required for it to be recognised.
	lights := make([]string, 1)
	lights[0] = "1"

	// Assumes scene doesn't exist (brightlight owned scenes are deleted at start up)
	s := scene{
		Name:        name,
		Lights:      lights,
		Owner:       "brightlight",
		Locked:      true,
		LastUpdated: time.Now().Format("2006-01-02T15:04:05"),
		Version:     2,
	}
	nextId := getNextSceneId(&fs.Scenes)
	fs.Scenes[nextId] = &s
}

// Body={"name":"Relax","lights":["1"],"recycle":false,"appdata":{"version":1,"data":"fskZZ_r01_d01"}
// Response=[{"success":{"id": "Abc123Def456Ghi"}}]
func createScene(c *cmdContext) {

	// Decode body as a scene
	var s scene
	err := json.Unmarshal(c.Body, &s)
	if err != nil {
		reportError(c.W, newApiErrorBodyContainsInvalidJson(c.ErrorAddress, err))
		return
	}

	// Check all lights exist TODO What does a real bridge return here?
	if !doAllLightsExist(&s.Lights, &c.FullState.Lights) {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Snapshot light states from lights
	s.LightStates = make(map[string]lightstate)
	for _, l := range s.Lights {
		s.LightStates[l] = c.FullState.Lights[l].State
	}

	// Set document version
	s.Version = 2

	// Set owner from command context
	s.Owner = c.User

	// Determine the next scene ID
	nextId := getNextSceneId(&c.FullState.Scenes)

	// Set the last update time
	s.LastUpdated = time.Now().Format("2006-01-02T15:04:05")

	// Add new scene
	c.FullState.Scenes[nextId] = &s

	// Return success
	SendSuccessResponse(c.W, newSceneSuccessResponse{
		nextId,
	})
}

func getNextSceneId(m *map[string]*scene) string {

	nextId := 0
	for k := range *m {

		i, err := strconv.Atoi(k)
		if err == nil {
			if i > nextId {
				nextId = i
			}
		}
	}
	return strconv.Itoa(nextId + 1)
}

// PUT /api/<username>/scenes/<id>/lightstates/<id> (modifying lightstates)
// PUT /api/<username>/scenes/<id> (modifying the scene name or lights list)
// https://developers.meethue.com/documentation/scenes-api#43_modify_scene
func setScene(c *cmdContext) {

	if !(len(c.Resource) == 1 || (len(c.Resource) == 3 && c.Resource[1] == "lightstates")) {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Extract the scene to modify
	s, ok := c.FullState.Scenes[c.Resource[0]]
	if !ok {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Perform modification
	if len(c.Resource) == 1 {
		// TODO (modifying the scene name or lights list)
		reportError(c.W, newApiErrorNotImplemented(c.ErrorAddress))
	} else {
		setSceneLightState(c, s)
	}
}

// URL = /api/<username>/scenes/<scene-id>/lightstates/<lightstate-id>
// Body = {"on":true,"bri":254,"ct":346}
// Response = [
//   {"success":{"address":"/scenes/scene-id/lights/light-id/state/on", "value":true}},
//   {"success":{"address":"/scenes/scene-id/lights/light-id/state/BRI", "value":254}},
//   {"success":{"address":"/scenes/scene-id/lights/light-id/state/ct", "value":346}}
// ]
func setSceneLightState(c *cmdContext, s *scene) {

	if s.LightStates == nil {
		// TODO REPORT ERROR
		return
	}

	// Check light state exists
	lightStateId := c.Resource[2]
	ls, ok := s.LightStates[lightStateId]
	if !ok {
		// TODO REPORT ERROR
		return
	}

	// Decode body as a lightstate
	var bodyLs lightstate
	if err := json.Unmarshal(c.Body, &bodyLs); err != nil {
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

	// Request resource is lightstates but response resource is lights
	resourceUrl := strings.Replace(c.ResourceUrl, "lightstates", "lights", 1)

	// Update light state fields from body
	if _, ok := bodyMap["on"]; ok {
		ls.On = bodyLs.On
		response = response.AppendSuccessResponse(resourceUrl+"/state/on", ls.On)
	}

	if _, ok := bodyMap["bri"]; ok {
		ls.Bri = bodyLs.Bri
		response = response.AppendSuccessResponse(resourceUrl+"/state/bri", ls.Bri)
	}

	if _, ok := bodyMap["hue"]; ok {
		ls.Hue = bodyLs.Hue
		response = response.AppendSuccessResponse(resourceUrl+"/state/hue", ls.Hue)
	}

	if _, ok := bodyMap["sat"]; ok {
		ls.Sat = bodyLs.Sat
		response = response.AppendSuccessResponse(resourceUrl+"/state/sat", ls.Sat)
	}

	if _, ok := bodyMap["xy"]; ok {
		ls.Xy = bodyLs.Xy
		response = response.AppendSuccessResponse(resourceUrl+"/state/xy", ls.Xy)
	}

	if _, ok := bodyMap["ct"]; ok {
		ls.Ct = bodyLs.Ct
		response = response.AppendSuccessResponse(resourceUrl+"/state/ct", ls.Ct)
	}

	if _, ok := bodyMap["alert"]; ok {
		ls.Alert = bodyLs.Alert
		response = response.AppendSuccessResponse(resourceUrl+"/state/alert", ls.Alert)
	}

	if _, ok := bodyMap["effect"]; ok {
		ls.Effect = bodyLs.Effect
		response = response.AppendSuccessResponse(resourceUrl+"/state/effect", ls.Effect)
	}

	if _, ok := bodyMap["transitiontime"]; ok {
		ls.TransitionTime = bodyLs.TransitionTime
		response = response.AppendSuccessResponse(resourceUrl+"/state/transitiontime", ls.TransitionTime)
	}

	if _, ok := bodyMap["colormode"]; ok {
		ls.ColorMode = bodyLs.ColorMode
		response = response.AppendSuccessResponse(resourceUrl+"/state/colormode", ls.ColorMode)
	}

	if _, ok := bodyMap["reachable"]; ok {
		ls.Reachable = bodyLs.Reachable
		response = response.AppendSuccessResponse(resourceUrl+"/state/reachable", ls.Reachable)
	}

	// Update scene last update time
	s.LastUpdated = time.Now().Format("2006-01-02T15:04:05")

	// Update lightstate in scene
	s.LightStates[lightStateId] = ls

	// Return response
	respondWithJsonEncodedObject(c.W, response)
}

// Delete all scenes
// Url = /api/<username>/scenes
// Response = [{"success": "/scenes deleted."}]
func deleteScenes(c *cmdContext) {

	// Remove the scenes and report success
	c.FullState.Scenes = make(map[string]*scene)
	SendSuccessResponse(c.W, c.ResourceUrl+" deleted.")
}

// Delete scene
// Url = /api/<username>/scenes/<id>
// Response = [{"success": "/scenes/1 deleted."}]
func deleteScene(c *cmdContext, sceneId string) {

	if _, ok := c.FullState.Scenes[sceneId]; !ok {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Remove the scene and report success
	delete(c.FullState.Scenes, sceneId)
	SendSuccessResponse(c.W, c.ResourceUrl+" deleted.")
}
