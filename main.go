package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/servers"
	"github.com/andew42/brightlight/stats"
	"golang.org/x/net/websocket"
	"mime"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
)

// Wrap a Dir file system server object to log failures
type LoggedDir struct {
	http.Dir
}

// Log requests for non existent content
func (d LoggedDir) Open(path string) (http.File, error) {

	f, err := d.Dir.Open(path)
	if err != nil {
		log.WithField("path", path).
			Info("requested static HTTP content not found")
	}
	return f, err
}

// Wrap a Dir file system server object to redirect
type LoggedRedirectingDir struct {
	LoggedDir
	Prefix []string
}

// Redirect paths in Prefix slice by remove the entry prefix
// e.g. for entry /buttons path /buttons/index.html -> /index.html
func (d LoggedRedirectingDir) Open(path string) (http.File, error) {

	for _, p := range d.Prefix {
		if strings.HasPrefix(path, p) {
			path = strings.TrimPrefix(path, p)
			log.WithField("prefix", p).Info("redirecting")
			break
		}
	}
	return d.LoggedDir.Open(path)
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
		log.Fields{"gover": runtime.Version(), "goos": runtime.GOOS, "goarch": runtime.GOARCH}).
		Info("brightlight environment")

	// Figure out where the content directory is, GOPATH may contain : separated paths
	// contentPath1 supports the original ractive ui
	contentBasePath := path.Join(os.Getenv("GOPATH"), "src/github.com/andew42/brightlight")
	log.WithFields(
		log.Fields{"contentBasePath": contentBasePath}).
		Info("HTTP content base path")

	// Start drivers
	controller.StartTeensyDriver()
	controller.StartRelayDriver()
	renderer := make(chan *framebuffer.FrameBuffer)
	framebuffer.StartDriver(renderer)
	animations.StartDriver(renderer)
	stats.StartDriver()

	// Set up static content serving
	mime.AddExtensionType(".manifest", "text/cache-manifest")
	// Serve original ui on /ui
	fs1 := http.FileServer(LoggedDir{http.Dir(contentBasePath)})
	http.Handle("/ui/", fs1)
	// Serve new react ui on /
	fs2 := http.FileServer(LoggedRedirectingDir{
		LoggedDir{http.Dir(contentBasePath + "/ui2/build")},
		[]string{"/buttons", "/virtual"}})
	//	[]string{}})
	http.Handle("/", fs2)

	// Config requires PUT (write) support
	http.HandleFunc("/config/", servers.GetConfigHandler(contentBasePath+"/ui"))
	http.HandleFunc("/ui-config/", servers.GetConfigHandler(contentBasePath+"/ui2/build"))

	// Requests to run zero or more animation (json payload)
	http.HandleFunc("/RunAnimations/", servers.RunAnimationsHandler)

	// Requests to show a strip length on the room lights
	http.HandleFunc("/StripLength/", servers.StripLenHandler)

	// Push button state changes over a web socket to keep UIs in sync
	http.Handle("/ButtonState", websocket.Handler(servers.ButtonStateHandler))

	// Push frame buffer changes over a web socket for virtual framebuffer debugging
	http.Handle("/FrameBuffer", websocket.Handler(servers.FrameBufferHandler))

	// Push stats info over a web socket
	http.Handle("/Stats", websocket.Handler(servers.StatsHandler))

	// Request to set server options
	http.HandleFunc("/option/", servers.OptionHandler)

	// Start web server
	ipAndPort, err := config.GetLocalIP()
	if err != nil {
		log.WithField("err", err).
			Fatal("Failed to find an IP address on which to serve content")
	}
	ipAndPort += ":8080"
	log.WithField("address", ipAndPort).
		Info("serving old UI at /ui/html and new UI at /")
	if err := http.ListenAndServe(ipAndPort, nil); err != nil {
		log.Error(err.Error())
	}

	log.Info("brightlight exited")
}
