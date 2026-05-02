package servers

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var configVersion = 0

// GetConfigHandler Handle HTTP requests to read and write config
func GetConfigHandler(contentPath string) func(http.ResponseWriter, *http.Request) {

	// Compute the canonical base path once for traversal checks
	cleanBase := filepath.Clean(contentPath)

	return func(w http.ResponseWriter, r *http.Request) {

		// Construct file system path to config and guard against path traversal
		fullPath := filepath.Join(contentPath, filepath.FromSlash(r.URL.Path))
		if !strings.HasPrefix(filepath.Clean(fullPath)+string(filepath.Separator), cleanBase+string(filepath.Separator)) {
			log.WithField("FullPath", fullPath).Warn("configHandler path traversal attempt")
			http.Error(w, "Forbidden", 403)
			return
		}

		if r.Method == "GET" {

			log.WithField("FullPath", fullPath).Info("configHandler GET called")
			content, err := os.ReadFile(fullPath)
			if err != nil {
				log.WithField("Error", err.Error()).Warn("Failed to load config file")
				http.Error(w, "Failed to load config file", 404)
			} else {
				w.Header().Set("Content-Type", "application/json")
				if _, err = w.Write(content); err != nil {
					log.WithField("Error", err.Error()).Warn("configHandler GET failed to write response")
				}
			}
		} else if r.Method == "PUT" {

			log.WithField("FullPath", fullPath).Info("configHandler PUT called")

			// Only support writing user.json (ui) and user-buttons.json (ui2)
			if r.URL.Path != "/config/user.json" && r.URL.Path != "/ui-config/user-buttons.json" {
				log.WithField("FileName", r.URL.Path).Warn("Unsupported config file name")
				http.Error(w, "File name not allowed", 401)
				return
			}

			// Limit body to 10K regardless of whether Content-Length header is present
			r.Body = http.MaxBytesReader(w, r.Body, 10000)
			if content, err := io.ReadAll(r.Body); err != nil {
				log.WithField("Error", err.Error()).Warn("Failed to read PUT body content")
				http.Error(w, "Failed to read PUT content", 400)
			} else {
				if err = os.WriteFile(fullPath, content, 0644); err != nil {
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
			http.Error(w, "Method not allowed", 405)
		}
	}
}
