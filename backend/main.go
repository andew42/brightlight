package main

import (
	"log/slog"
	"mime"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/config"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/framebuffer"
	"github.com/andew42/brightlight/servers"
	"github.com/andew42/brightlight/stats"
	"golang.org/x/net/websocket"
)

// LoggedDir Wrap a Dir file system server object to log failures
type LoggedDir struct {
	http.Dir
}

// Open Log requests for non existent content
func (d LoggedDir) Open(path string) (http.File, error) {

	f, err := d.Dir.Open(path)
	if err != nil {
		slog.Info("requested static HTTP content not found", "path", path)
	}
	return f, err
}

// LoggedRedirectingDir Wrap a Dir file system server object to redirect
type LoggedRedirectingDir struct {
	LoggedDir
	Prefix []string
}

// Open Redirect paths in Prefix slice by remove the entry prefix
// e.g. for entry /buttons path /buttons/index.html -> /index.html
func (d LoggedRedirectingDir) Open(path string) (http.File, error) {

	for _, p := range d.Prefix {
		if strings.HasPrefix(path, p) {
			path = strings.TrimPrefix(path, p)
			slog.Info("redirecting", "prefix", p)
			break
		}
	}
	return d.LoggedDir.Open(path)
}

// Main
func main() {

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Report what are we running on
	slog.Info("environment", "gover", runtime.Version(), "goos", runtime.GOOS, "goarch", runtime.GOARCH)

	// Figure out where the content directory is by loading BRIGHTLIGHT
	contentBasePath := os.Getenv("BRIGHTLIGHT")
	if len(contentBasePath) == 0 {
		slog.Error("BRIGHTLIGHT environment variable not set (web root)")
		os.Exit(1)
	}
	// contentBasePath := "C:/Users/Andrew/GolandProjects/brightlight"
	slog.Info("HTTP content base path", "contentBasePath", contentBasePath)

	// Start drivers
	controller.StartTeensyDriver()
	controller.StartRelayDriver()
	renderer := make(chan *framebuffer.FrameBuffer)
	framebuffer.StartDriver(renderer)
	animations.StartDriver(renderer)
	stats.StartDriver()

	// Dump out the list of named segments with lengths
	//fb := framebuffer.NewFrameBuffer()
	//for _, s := range segment.GetAllNamedSegmentNames() {
	//	seg, _ := segment.GetNamedSegment(s)
	//	rs := seg.GetSegment(fb)
	//	fmt.Printf("%s %d\n", s, rs.Len())
	//}

	// Set up static content serving
	mime.AddExtensionType(".manifest", "text/cache-manifest")
	// Serve react frontend on /
	fs2 := http.FileServer(LoggedRedirectingDir{
		LoggedDir{http.Dir(contentBasePath + "/frontend/build")},
		[]string{"/buttons", "/virtual"}})
	http.Handle("/", fs2)

	// ui-config requires PUT (write) support for saving button config
	http.HandleFunc("/ui-config/", servers.GetConfigHandler(contentBasePath+"/frontend/build"))

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
		slog.Error("Failed to find an IP address on which to serve content", "err", err)
		os.Exit(1)
	}
	ipAndPort += ":8080"
	slog.Info("serving frontend at /", "address", ipAndPort)
	if err := http.ListenAndServe(ipAndPort, nil); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("brightlight exited")
}
