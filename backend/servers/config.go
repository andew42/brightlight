package servers

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path"
)

var configVersion = 0

// GetConfigHandler Handle HTTP requests to read and write config
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

			// Only support writing user.json (ui) and user-buttons.json (ui2)
			if r.URL.Path != "/config/user.json" && r.URL.Path != "/ui-config/user-buttons.json" {
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
				} else {
					// Let clients know the config has been updated
					configVersion++
					updateButtonPadVersion(configVersion)
				}
			}
		} else {
			log.WithField("Method", r.Method).Warn("Unknown config server method")
			http.Error(w, "Failed to write file", 405)
		}
	}
}
