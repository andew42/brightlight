package hue

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type newUserRequestBody struct {
	DeviceType string `json:"devicetype"`
}

type newUserResponseBody struct {
	UserName string `json:"username"`
}

// Create a new user
// https://developers.meethue.com/documentation/configuration-api
// Request: {"devicetype": "my_hue_app#iphone peter"}
// Response: [{"success":{"username": "83b7780291a6ceffbe0bd049104df"}}]
func createNewUser(c *config, body []byte, w http.ResponseWriter) {

	// Deserialize request body
	var r newUserRequestBody
	err := json.Unmarshal(body, &r)
	if err != nil {
		reportError(w, newApiErrorBodyContainsInvalidJson("", err))
		return
	}

	// Create new user id
	newUserId := createRandomUserId()
	now := time.Now().Format("2006-01-02T15:04:05")
	c.WhiteList[newUserId] = whitelistEntry{LastUseDate: now, CreateDate: now, Name: r.DeviceType}

	// Respond with new user id
	SendSuccessResponse(w, newUserResponseBody{
		newUserId,
	})
}

// A random 128 bit (32 chars) zero padded hex user ID
func createRandomUserId() string {
	a := rand.Uint64()
	b := rand.Uint64()
	return fmt.Sprintf("%016x%016x", a, b)
}
