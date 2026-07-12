package servers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/andew42/brightlight/controller"
)

type cmd struct {
	Cmd   string
	Param string
}

// OptionHandler Handle HTTP requests to set a server option
func OptionHandler(w http.ResponseWriter, r *http.Request) {

	// JSON body of form
	// {"cmd": "outputMapping", "param": "Linear"},
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if err != nil {
		slog.Error("OptionHandler bad body", "err", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	// Un-marshal JSON into typed cmd
	var cmd cmd
	if err = json.Unmarshal(body, &cmd); err != nil {
		slog.Error("OptionHandler bad body JSON", "err", err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	slog.Info("OptionHandler called", "cmd", cmd.Cmd, "param", cmd.Param)

	// Perform the command
	switch cmd.Cmd {
	case "outputMapping":
		controller.SetOutputMapping(cmd.Param)
	default:
		slog.Warn("OptionHandler unknown command", "cmd", cmd.Cmd)
	}
}
