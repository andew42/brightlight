package servers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/andew42/brightlight/controller"
	"io/ioutil"
	"net/http"
)

type cmd struct {
	Cmd   string
	Param string
}

// Handle HTTP requests to set a server option
func OptionHandler(w http.ResponseWriter, r *http.Request) {

	// JSON body of form
	// {"cmd": "outputMapping", "param": "Linear"},
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithField("err", err.Error()).Error("OptionHandler bad body")
		http.Error(w, err.Error(), 400)
		return
	}

	// Un-marshal JSON into typed cmd
	var cmd cmd
	if err = json.Unmarshal(body, &cmd); err != nil {
		log.WithField("err", err.Error()).Error("OptionHandler bad body JSON")
		http.Error(w, err.Error(), 400)
		return
	}

	log.WithFields(map[string]interface{}{"cmd": cmd.Cmd, "param": cmd.Param}).Info("OptionHandler called")

	// Perform the command
	switch cmd.Cmd {
	case "outputMapping":
		controller.SetOutputMapping(cmd.Param)
	default:
		log.WithField("cmd", cmd.Cmd).Warn("OptionHandler unknown command")
	}
}
