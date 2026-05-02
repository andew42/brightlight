package servers

import (
	"io"
	"log/slog"
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
			slog.Warn("configHandler path traversal attempt", "FullPath", fullPath)
			http.Error(w, "Forbidden", 403)
			return
		}

		if r.Method == "GET" {

			slog.Info("configHandler GET called", "FullPath", fullPath)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				slog.Warn("Failed to load config file", "Error", err.Error())
				http.Error(w, "Failed to load config file", 404)
			} else {
				w.Header().Set("Content-Type", "application/json")
				if _, err = w.Write(content); err != nil {
					slog.Warn("configHandler GET failed to write response", "Error", err.Error())
				}
			}
		} else if r.Method == "PUT" {

			slog.Info("configHandler PUT called", "FullPath", fullPath)

			// Only support writing user.json (ui) and user-buttons.json (ui2)
			if r.URL.Path != "/config/user.json" && r.URL.Path != "/ui-config/user-buttons.json" {
				slog.Warn("Unsupported config file name", "FileName", r.URL.Path)
				http.Error(w, "File name not allowed", 401)
				return
			}

			// Limit body to 10K regardless of whether Content-Length header is present
			r.Body = http.MaxBytesReader(w, r.Body, 10000)
			if content, err := io.ReadAll(r.Body); err != nil {
				slog.Warn("Failed to read PUT body content", "Error", err.Error())
				http.Error(w, "Failed to read PUT content", 400)
			} else {
				if err = os.WriteFile(fullPath, content, 0644); err != nil {
					slog.Error("Failed to write file", "Error", err.Error())
					http.Error(w, "Failed to write file", 507)
				} else {
					// Let clients know the config has been updated
					configVersion++
					updateButtonPadVersion(configVersion)
				}
			}
		} else {
			slog.Warn("Unknown config server method", "Method", r.Method)
			http.Error(w, "Method not allowed", 405)
		}
	}
}
