package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/servers"
	"github.com/andew42/brightlight/stats"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"path"
	"runtime"
)

func main() {

	// Force logrus to use console colouring
	var forceColours = flag.Bool("logrusforcecolours", false, "force logrus to use console colouring")
	flag.Parse()
	if *forceColours {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}

	// What are we running on?
	log.WithFields(
		log.Fields{"gover": runtime.Version(), "goos": runtime.GOOS, "goarch": runtime.GOARCH}).Info("Brightlight")

	// Start drivers
	controller.StartDriver()
	renderer := make(chan *framebuffer.FrameBuffer)
	framebuffer.StartDriver(renderer)
	animations.StartDriver(renderer)

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
	http.HandleFunc("/StripLength/", servers.GetStripLenHandler())

	// Push frame buffer changes over a web socket
	http.Handle("/FrameBuffer", websocket.Handler(servers.GetFrameBufferHandler()))

	// Push stats info over a web socket
	var statistics = stats.NewStats()
	http.Handle("/Stats", websocket.Handler(servers.GetStatsHandler(&statistics)))

	// Start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error(err.Error())
	}

	log.Info("brightlight exited")
}
