package hue

import "net/http"

// Response = [
//   {"success":{"address":"/scenes/scene-id/lights/light-id/state/on", "value":true}},
//   {"success":{"address":"/scenes/scene-id/lights/light-id/state/BRI", "value":254}},
//   {"success":{"address":"/scenes/scene-id/lights/light-id/state/ct", "value":346}}
// ]

type apiResponseList []interface{}

type apiSuccessResponseRow struct {
	Success interface{} `json:"success"`
}

type apiResponseRowContent struct {
	Address string      `json:"address"`
	Value   interface{} `json:"value"`
}

// Append a success response row to the response list
func (l apiResponseList) AppendSuccessResponse(address string, value interface{}) apiResponseList {

	r := apiSuccessResponseRow{
		Success: apiResponseRowContent{
			Address: address,
			Value:   value,
		},
	}
	return append(l, r)
}

// Sends a single success line as a response
// e.g. [{"success":{"username": "83b7780291a6ceffbe0bd049104df"}}]
func SendSuccessResponse(w http.ResponseWriter, responseObject interface{}) {

	// Respond with new user id
	var rl apiResponseList
	sr := apiSuccessResponseRow{
		responseObject,
	}
	rl = append(rl, sr)
	respondWithJsonEncodedObject(w, rl)
}
