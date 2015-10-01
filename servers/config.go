package servers

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strings"
	"path"
	"io/ioutil"
)

type Config struct {
	ContentPath string
}

// Handle HTTP requests to read and write config
func (config Config) Handler(w http.ResponseWriter, r *http.Request) {

	// Colour value follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No config file specified", 400)
		log.Info("configHandler no config file specified")
		return
	}

	fullPath := path.Join(config.ContentPath, r.URL.Path)

	// TODO CHECK FILE NAME AND LENGTH

	if r.Method == "GET" {

		log.Info("configHandler with GET called " + fullPath)
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Error("Failed to load config file : " + err.Error())
			http.Error(w, "Failed to load config file", 404)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(content)
		}
	} else if r.Method == "PUT" {

		log.Info("configHandler with PUT called " + fullPath)
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Failed to read PUT content : " + err.Error())
			http.Error(w, "Failed to read PUT content", 400)
		} else {
			err = ioutil.WriteFile(fullPath, content, 0644)
			// TODO TEST ERROR
		}
	}
}
