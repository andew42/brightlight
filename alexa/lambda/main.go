// Lambda handler for the brightlight Alexa custom skill. Forwards the
// spoken button name to the brightlight /alexa/RunButton endpoint over
// HTTPS with a bearer token, trusting only the pinned brightlight
// self-signed certificate.
//
// Environment variables:
//
//	BRIGHTLIGHT_URL         e.g. https://myhouse.example.com:8443
//	BRIGHTLIGHT_ALEXA_TOKEN shared secret, same value as on the server
//	BRIGHTLIGHT_SERVER_CERT path to the server's PEM cert bundled in the
//	                        zip (default server-cert.pem)
//	ALEXA_SKILL_ID          if set, requests from other skills are refused
package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type alexaRequest struct {
	Session struct {
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
	} `json:"session"`
	Request struct {
		Type   string `json:"type"`
		Intent struct {
			Name  string `json:"name"`
			Slots map[string]struct {
				Value string `json:"value"`
			} `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
}

type alexaResponse struct {
	Version  string `json:"version"`
	Response struct {
		OutputSpeech struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"outputSpeech"`
		ShouldEndSession bool `json:"shouldEndSession"`
	} `json:"response"`
}

func speak(text string, endSession bool) alexaResponse {

	var r alexaResponse
	r.Version = "1.0"
	r.Response.OutputSpeech.Type = "PlainText"
	r.Response.OutputSpeech.Text = text
	r.Response.ShouldEndSession = endSession
	return r
}

func handle(ctx context.Context, request alexaRequest) (alexaResponse, error) {

	if skillId := os.Getenv("ALEXA_SKILL_ID"); skillId != "" &&
		request.Session.Application.ApplicationID != skillId {
		return alexaResponse{}, errors.New("request from unexpected skill id")
	}

	switch request.Request.Type {

	case "LaunchRequest":
		return speak("Which light setting would you like?", false), nil

	case "IntentRequest":
		switch request.Request.Intent.Name {

		case "RunButtonIntent":
			buttonName := request.Request.Intent.Slots["buttonName"].Value
			if buttonName == "" {
				return speak("Which light setting would you like?", false), nil
			}
			return runButton(buttonName), nil

		case "AMAZON.HelpIntent":
			return speak("Say the name of a light setting, for example rainbow.", false), nil

		default: // AMAZON.StopIntent, AMAZON.CancelIntent etc.
			return speak("Goodbye.", true), nil
		}

	default: // SessionEndedRequest must not include speech
		var r alexaResponse
		r.Version = "1.0"
		r.Response.ShouldEndSession = true
		return r, nil
	}
}

// Call the brightlight RunButton endpoint and turn the result into speech
func runButton(buttonName string) alexaResponse {

	client, err := newPinnedClient()
	if err != nil {
		return speak("The light server certificate is not configured.", true)
	}

	body, _ := json.Marshal(struct {
		Button string `json:"button"`
	}{buttonName})

	url := os.Getenv("BRIGHTLIGHT_URL") + "/alexa/RunButton"
	request, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+os.Getenv("BRIGHTLIGHT_ALEXA_TOKEN"))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return speak("I couldn't reach the light server.", true)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		var result struct{ Matched string }
		if json.NewDecoder(response.Body).Decode(&result) == nil && result.Matched != "" {
			return speak("OK, "+result.Matched+".", true)
		}
		return speak("OK.", true)
	case 404:
		return speak("I couldn't find a light setting called "+buttonName+".", true)
	default:
		return speak("The light server said no. Check its logs.", true)
	}
}

// HTTP client trusting only the brightlight self-signed certificate
func newPinnedClient() (*http.Client, error) {

	certPath := os.Getenv("BRIGHTLIGHT_SERVER_CERT")
	if certPath == "" {
		certPath = "server-cert.pem"
	}
	pem, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		return nil, errors.New("no certificate found in " + certPath)
	}
	return &http.Client{
		Timeout:   6 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}},
	}, nil
}

func main() {
	lambda.Start(handle)
}
