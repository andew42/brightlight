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
		slog.Warn("BRIGHTLIGHT environment variable not set — static content will not be served (use Vite dev server)")
	} else {
		slog.Info("HTTP content base path", "contentBasePath", contentBasePath)
	}

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

	// ui-config lives alongside the backend binary; fall back to relative path for dev
	uiConfigDir := "ui-config"
	if len(contentBasePath) > 0 {
		uiConfigDir = contentBasePath + "/backend/ui-config"
	}
	http.HandleFunc("/api/ui-config/", servers.GetConfigHandler(uiConfigDir))
	http.HandleFunc("/api/RunAnimations/", servers.RunAnimationsHandler)
	http.HandleFunc("/api/StripLength/", servers.StripLenHandler)
	http.HandleFunc("/api/ButtonState", servers.ButtonStateHandler)
	http.HandleFunc("/api/FrameBuffer", servers.FrameBufferHandler)
	http.HandleFunc("/api/Stats", servers.StatsHandler)
	http.HandleFunc("/api/option/", servers.OptionHandler)

	// Set up static content serving (skipped when BRIGHTLIGHT is unset — Vite serves the frontend)
	if len(contentBasePath) > 0 {
		err := mime.AddExtensionType(".manifest", "text/cache-manifest")
		if err != nil {
			slog.Error("failed to add extension type")
			return
		}
		fs2 := http.FileServer(LoggedRedirectingDir{
			LoggedDir{http.Dir(contentBasePath + "/frontend/build")},
			[]string{"/buttons", "/virtual"}})
		http.Handle("/", fs2)
	}

	// Start web server — listen on all interfaces so both LAN IP and localhost work
	if localIP, err := config.GetLocalIP(); err == nil {
		slog.Info("serving on", "address", localIP+":8080")
	}
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("brightlight exited")
}
