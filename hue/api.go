package hue

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Returns a handler for http requests for the hue bridge API emulation (/api/*)
func GetApiHandler(fsl *fullStateLocker, brightlightCommandChan chan interface{}) (func(w http.ResponseWriter, r *http.Request), error) {

	// Local Helper: Read request body as []byte
	readBody := func(r *http.Request) ([]byte, error) {

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithField("Error", err).Error("Hue API failed to read request body")
			return nil, err
		}
		return b, nil
	}

	// Local Helper: Finds the index of the count occurrence of c in s
	nthIndex := func(s string, c byte, count int) int {
		i := 0
		for ; count > 0 && i < len(s); i++ {
			if s[i] == c {
				count--
			}
		}
		if count > 0 {
			return -1
		}
		return i - 1
	}

	// The single threaded http request handler
	return func(w http.ResponseWriter, r *http.Request) {

		// Read request body as string
		b, err := readBody(r)
		if err != nil {
			return
		}

		// Log the request
		log.WithFields(log.Fields{
			"Method": r.Method,
			"Url":    r.URL,
			"Body":   string(b),
		}).Info("Hue API Request")

		// TODO
		//if r.Method == "DELETE" {
		//	return
		//}
		//time.Sleep(100 * time.Millisecond)

		// Split trimmed Url into components, first component is always api
		// otherwise we would not have been called. We immediately discard
		// the first api components [1:].
		trimmedUrl := strings.Trim(r.URL.String(), "/\\")
		cmd := strings.Split(trimmedUrl, "/")[1:]

		// If we get an error, return URL without /api prefix
		errorAddress := trimmedUrl[3:]

		// Get exclusive access to full state
		fs := fsl.Lock()
		defer fsl.Unlock()

		// Update times in config
		now := time.Now()
		fs.Config.Utc = now.UTC().Format("2006-01-02T15:04:05")
		fs.Config.LocalTime = now.Format("2006-01-02T15:04:05")

		// Create new user request (POST /api)
		// TODO simulate link button press with 30s timeout
		if len(cmd) == 0 {

			if r.Method == "POST" {
				if fs.Config.LinkButton {
					createNewUser(fs.Config, b, w)
					fs.Save()
				} else {
					reportError(w, newApiErrorLinkButtonNotPressed(errorAddress))
				}
			} else {
				reportError(w, newApiErrorMethodNotAvailable(errorAddress))
			}
			return
		}

		// Check user is a known user (allow (null) if link button is pressed)
		user := cmd[0]
		_, isWhiteListUser := fs.Config.WhiteList[user]
		isWhiteListUser = isWhiteListUser || (user == "nouser" && fs.Config.LinkButton)
		if !isWhiteListUser {
			reportError(w, newApiErrorUnauthorizedUser(errorAddress))
			return
		}

		// Discard user from the command collection
		cmd = cmd[1:]

		// GET /api/<username> is a special case which returns full state
		if len(cmd) == 0 {
			if r.Method == "GET" {
				respondWithJsonEncodedObject(w, &fs)
			} else {
				reportError(w, newApiErrorMethodNotAvailable(errorAddress))
			}
			return
		}

		// Build a context for command processing
		cmdContext := cmdContext{
			FullState:              fs,
			BrightlightCommandChan: brightlightCommandChan,
			User:                   user,
			Method:                 r.Method,
			Resource:               cmd[1:],
			// Remove /api/user-id prefix
			ResourceUrl:  trimmedUrl[nthIndex(trimmedUrl, '/', 2):],
			Body:         b,
			ErrorAddress: errorAddress,
			W:            w,
		}

		// Dispatch command based on resource type
		switch cmd[0] {
		case "lights":
			processLightsRequest(&cmdContext)

		case "groups":
			processGroupsRequest(&cmdContext)

		case "schedules":
			return // TODO

		case "scenes":
			processScenesRequest(&cmdContext)

		case "sensors":
			processSensorRequest(&cmdContext)

		case "rules":
			return // TODO

		case "config":
			processConfigRequest(&cmdContext)

		case "capabilities":
			return // TODO

		case "resourcelinks":
			return // TODO

		default:
			reportError(w, newApiErrorResourceNotAvailable(errorAddress))
			return
		}

		// Save state if we processed a possible state changing request
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
			fs.Save()
		}
	}, nil
}

// Returns API error as JSON TODO: What should the HTTP status be (it's currently 200 OK)
func reportError(w http.ResponseWriter, apierror apierror) {

	respondWithJsonEncodedObject(w, apierror)
}

// Sends supplied object encoded as JSON in response body
func respondWithJsonEncodedObject(w http.ResponseWriter, obj interface{}) {

	b, err := json.Marshal(obj)
	if err != nil {
		log.WithField("Error", err).Error("Hue API failed to marshal object")
		return
	}
	writeResponse(w, string(b))
}

// Write response body
func writeResponse(w http.ResponseWriter, responseBody string) {

	// These are the headers we get from a real hue bridge TODO VALIDATE ESPECIALLY Content-Length
	w.Header().Set("access-control-allow-credentials", "true")
	w.Header().Set("access-control-allow-headers", "Content-Type")
	w.Header().Set("access-control-allow-methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")
	w.Header().Set("access-control-allow-origin", "*")
	w.Header().Set("access-control-max-age", "3600")
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
	w.Header().Set("pragma", "no-cache")
	w.Header().Set("cache-control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("connection", "close")
	w.Header().Set("expires", "Mon, 1 Aug 2011 09:00:00 GMT")
	l, err := w.Write([]byte(responseBody))
	if err != nil || l != len(responseBody) {
		log.WithField("Error", err).Error("Hue API failed to write response")
	} else {
		log.WithField("Body", string(responseBody)).Info("Hue API Response")
	}
}
