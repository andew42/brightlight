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
	"mime"
	"github.com/andew42/brightlight/hue"
	"github.com/andew42/brightlight/config"
)

// Wrap a Dir file system server object to log failures
type LoggedDir struct {
	http.Dir
}

func (d LoggedDir) Open(path string) (http.File, error) {

	f, err := d.Dir.Open(path)
	if err != nil {
		log.WithField("path", path).Info("No Static HTTP Content")
	}
	return f, err
}

// Main
func main() {

	// Force logrus to use console colouring
	var forceColours = flag.Bool("logrusforcecolours", false, "force logrus to use console colouring")
	flag.Parse()
	if *forceColours {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
	}

	// Report what are we running on
	log.WithFields(
		log.Fields{"gover": runtime.Version(), "goos": runtime.GOOS, "goarch": runtime.GOARCH}).Info("Brightlight")

	// Start drivers
	controller.StartTeensyDriver()
	controller.StartRelayDriver()
	renderer := make(chan *framebuffer.FrameBuffer)
	framebuffer.StartDriver(renderer)
	animations.StartDriver(renderer)
	stats.StartDriver()

	// Start Hue bridge emulator components
	addressAndPort, err := config.GetServerAddressAndPort()
	if err != nil {
		log.WithField("error", err).Error("Failed to get server address and port")
	} else {

		hue.StartUpnpResponder(addressAndPort)

		// Handlers for hue emulator discovery
		dh, err := hue.GetUpnpDescriptionHandler(addressAndPort)
		if err != nil {
			log.WithField("error", err).Fatal("Failed to create hue description.xml handler")
		} else {
			http.HandleFunc("/description.xml", dh)
			http.HandleFunc("/api/", hue.ApiHandler)
		}
	}

	// Figure out where the content directory is, GOPATH may contain : separated paths
	contentPath := path.Join(os.Getenv("GOPATH"), "src/github.com/andew42/brightlight/ui")
	log.WithField("contentPath", contentPath).Info("Serving content")

	// Set up web routes (default / is static content)
	mime.AddExtensionType(".manifest", "text/cache-manifest")
	fs := http.FileServer(LoggedDir{http.Dir(contentPath)})
	http.Handle("/", fs)

	// Config requires PUT (write) support
	http.HandleFunc("/config/", servers.GetConfigHandler(contentPath))

	// Requests to run zero or more animation (json payload)
	http.HandleFunc("/RunAnimations/", servers.RunAnimationsHandler)

	// Requests to show a strip length on the room lights
	http.HandleFunc("/StripLength/", servers.StripLenHandler)

	// Push frame buffer changes over a web socket
	http.Handle("/FrameBuffer", websocket.Handler(servers.FrameBufferHandler))

	// Push stats info over a web socket
	http.Handle("/Stats", websocket.Handler(servers.StatsHandler))

	// Request to set server options
	http.HandleFunc("/option/", servers.OptionHandler)

	// Start web server
	if err := http.ListenAndServe(config.GetServerPort(), nil); err != nil {
		log.Error(err.Error())
	}

	log.Info("brightlight exited")
}
