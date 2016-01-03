package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"github.com/andew42/brightlight/servers"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/stats"
	"net/http"
	"os"
	"path"
	"runtime"
	"flag"
)

// The controller's frame buffer
var fb = controller.NewFrameBuffer()
var statistics = stats.NewStats()

func main() {

	// Force logrus to use console colouring
	var forceColours = flag.Bool("logrusforcecolours", false, "force logrus to use console colouring")
	flag.Parse()
	if *forceColours {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}

	// What are we running on?
	log.WithFields(
	log.Fields{"gover": runtime.Version(), "goos": runtime.GOOS, "goarch": runtime.GOARCH, }).Info("brightlight started")

	// Start serial and animation drivers
	controller.StartDriver(fb, &statistics)
	animations.StartDriver(fb, &statistics)

	// Figure out where the content directory is GOPATH may contain : separated paths
	contentPath := path.Join(os.Getenv("GOPATH"), "src/github.com/andew42/brightlight/ui")
	log.WithField("contentPath", contentPath).Info("Serving content")

	// Set up web routes (first static content)
	fs := http.FileServer(http.Dir(contentPath))
	http.Handle("/", fs)

	// Config requires PUT (write) support
	http.HandleFunc("/config/", servers.GetConfigHandler(contentPath))

	// Requests to run zero or more animation (json payload)
	http.HandleFunc("/RunAnimations/", servers.RunAnimationsHandler)

	// Requests to show a strip length on the room lights
	http.HandleFunc("/StripLength/", servers.GetStripLenHandler(len(fb.Strips)))

	// Push frame buffer changes over a web socket
	http.Handle("/FrameBuffer", websocket.Handler(servers.GetFrameBufferHandler(fb)))

	// Push stats info over a web socket
	http.Handle("/Stats", websocket.Handler(servers.GetStatsHandler(&statistics, fb)))

	// Start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error(err.Error())
	}

	log.Info("brightlight exited")
}
