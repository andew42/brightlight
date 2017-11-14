package hue

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

// HTTP Handler for /api/groups
func processGroupsRequest(c *cmdContext) {

	if c.Method == "GET" {
		if len(c.Resource) == 0 {
			respondWithJsonEncodedObject(c.W, &c.FullState.Groups)
			return
		} else if len(c.Resource) == 1 {
			if g, ok := c.FullState.Groups[c.Resource[0]]; ok {
				respondWithJsonEncodedObject(c.W, g)
				return
			}
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	} else if c.Method == "POST" {
		if len(c.Resource) == 0 {
			createGroup(c)
			return
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	} else if c.Method == "PUT" {
		setGroup(c)
		return
	} else if c.Method == "DELETE" {
		if len(c.Resource) == 0 {
			deleteGroups(c)
			return
		} else if len(c.Resource) == 1 {
			deleteGroup(c, c.Resource[0])
			return
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}
	reportError(c.W, newApiErrorMethodNotAvailable(c.ErrorAddress))
}

type group struct {
	Action   lightstate `json:"action,omitempty"`
	Lights   []string   `json:"lights"`
	Name     string     `json:"name"`
	Type     string     `json:"type,omitempty"`
	ModelId  string     `json:"modelid,omitempty"`
	UniqueId string     `json:"uniqueid,omitempty"`
	Class    string     `json:"class,omitempty"`
}

type newGroupSuccessResponse struct {
	Id string `json:"id"`
}

// https://developers.meethue.com/documentation/groups-api#22_create_group
// Body={"name":"Groupie","lights":["1"],"type":"Room","class":"Bedroom"}
// Response [{"success":{"id":"1"}}]
func createGroup(c *cmdContext) {

	var g group
	err := json.Unmarshal(c.Body, &g)
	if err != nil {
		log.WithField("Error", err).Error("Failed to decode group")
		return
	}

	// Check all lights exist TODO What does a real bridge return here?
	if !doAllLightsExist(&g.Lights, &c.FullState.Lights) {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Set action to lightstate of first light
	if len(g.Lights) > 0 {
		if l0, ok := c.FullState.Lights[g.Lights[0]]; ok {
			g.Action = l0.State
		}
	}

	// Add the group to full state under next available id
	nextId := getNextGroupId(&c.FullState.Groups)
	c.FullState.Groups[nextId] = &g

	// Return success
	SendSuccessResponse(c.W, newGroupSuccessResponse{
		nextId,
	})
}

func getNextGroupId(groups *map[string]*group) string {

	nextId := 0
	for k := range *groups {

		i, err := strconv.Atoi(k)
		if err == nil {
			if i > nextId {
				nextId = i
			}
		}
	}
	return strconv.Itoa(nextId + 1)
}

// PUT /api/<username>/groups/<id>/action (modifying group attribute)
// PUT /api/<username>/groups/<id> (modifying the group name or lights list)
func setGroup(c *cmdContext) {

	if !(len(c.Resource) == 1 || (len(c.Resource) == 2 && c.Resource[1] == "action")) {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Extract the group to modify
	g, ok := c.FullState.Groups[c.Resource[0]]
	if !ok {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Perform modification
	if len(c.Resource) == 1 {
		setGroupAttributes(c, g)
	} else {
		setGroupState(c, g)
	}
}

// Url=/api/4d65822107fcfd5278629a0f5f3f164f/groups/1
// Body={"name":"Beddy","lights":["1"],"class":"Bedroom"}
//[
//{"success":{"/groups/1/lights":["1"]}},
//{"success":{"/groups/1/name":"Bedroom"}}
//]
func setGroupAttributes(c *cmdContext, g *group) {

	// Decode body as a group
	var bodyGroup group
	if err := json.Unmarshal(c.Body, &bodyGroup); err != nil {
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

	// Update group fields from body
	if _, ok := bodyMap["name"]; ok {
		g.Name = bodyGroup.Name
		response = response.AppendSuccessResponse(c.ResourceUrl+"/name", g.Name)
	}

	if _, ok := bodyMap["lights"]; ok {
		g.Lights = bodyGroup.Lights
		response = response.AppendSuccessResponse(c.ResourceUrl+"/lights", g.Lights)
	}

	if _, ok := bodyMap["class"]; ok {
		g.Class = bodyGroup.Class
		response = response.AppendSuccessResponse(c.ResourceUrl+"/name", g.Class)
	}

	// Return response
	respondWithJsonEncodedObject(c.W, response)
}

// Url=/api/4d65822107fcfd5278629a0f5f3f164f/groups/1/action
// Body={"on":true}
// [
// {"success":{"/groups/1/action/on":, "value": true}}
// ]
func setGroupState(c *cmdContext, g *group) {

	// Decode the body as a map
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(c.Body, &bodyMap); err != nil {
		reportError(c.W, newApiErrorBodyContainsInvalidJson(c.ErrorAddress, err))
		return
	}

	response := apiResponseList{}

	// Do we have an on field?
	if on, onOk := bodyMap["on"]; onOk {
		for _, lightId := range g.Lights {
			lightName := c.FullState.Lights[lightId].Name
			c.BrightlightCommandChan <- SegmentControl{On: on.(bool), Name: lightName}
		}
		response = response.AppendSuccessResponse(c.ResourceUrl+"/on", on)
	}

	// Return response
	respondWithJsonEncodedObject(c.W, response)
}

// Delete all groups
// Url = /api/<username>/groups
// Response = [{"success": "/groups deleted."}]
func deleteGroups(c *cmdContext) {

	// Remove the groups and report success
	c.FullState.Groups = make(map[string]*group)
	SendSuccessResponse(c.W, c.ResourceUrl+" deleted.")
}

// Delete group
// Url = /api/<username>/groups/<id>
// Response = [{"success": "/groups/1 deleted."}]
func deleteGroup(c *cmdContext, groupId string) {

	if _, ok := c.FullState.Groups[groupId]; !ok {
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}

	// Remove the group and report success
	delete(c.FullState.Groups, groupId)
	SendSuccessResponse(c.W, c.ResourceUrl+" deleted.")
}
