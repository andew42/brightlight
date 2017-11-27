package servers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/hue"
	"github.com/andew42/brightlight/segments"
	"io/ioutil"
	"net/http"
	"path"
	"github.com/andew42/brightlight/config"
)

// Handle HTTP requests to read and write config
func GetConfigHandler(contentPath string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// Construct file system path to config
		fullPath := path.Join(contentPath, r.URL.Path)

		if r.Method == "GET" {

			log.WithField("FullPath", fullPath).Info("configHandler GET called")
			content, err := ioutil.ReadFile(fullPath)
			if err != nil {
				log.WithField("Error", err.Error()).Warn("Failed to load config file")
				http.Error(w, "Failed to load config file", 404)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(content)
			}
		} else if r.Method == "PUT" {

			log.WithField("FullPath", fullPath).Info("configHandler PUT called")

			// Only support user.json
			if r.URL.Path != "/config/user.json" {
				log.WithField("FileName", r.URL.Path).Warn("Unsupported config file name")
				http.Error(w, "File name not allowed", 401)
				return
			}

			// up to a size of 10K
			if r.ContentLength > 10000 {
				log.WithField("ContentLength", r.ContentLength).Warn("Config file content too large")
				http.Error(w, "Update content too large", 413)
				return
			}

			if content, err := ioutil.ReadAll(r.Body); err != nil {
				log.WithField("Error", err.Error()).Warn("Failed to read PUT body content")
				http.Error(w, "Failed to read PUT content", 400)
			} else {
				if err = ioutil.WriteFile(fullPath, content, 0644); err != nil {
					log.WithField("Error", err.Error()).Error("Failed to write file")
					http.Error(w, "Failed to write file", 507)
				}
			}
		} else {
			log.WithField("Method", r.Method).Warn("Unknown config server method")
			http.Error(w, "Failed to write file", 405)
		}
	}
}

// Push segment names and animations to hue bridge
func UpdateHueBridgeWithBrightlightConfig(contentPath string, u chan interface{}) {

	// Send the list of segment names
	for _, s := range segments.GetAllNamedSegmentNames() {
		u <- hue.SegmentUpdate{NewName: s}
	}

	// Send the list of user button names
	for _, p := range config.LoadUserPresets(contentPath) {
		u <- hue.PresetUpdate{NewName: p.Name}
	}
}
