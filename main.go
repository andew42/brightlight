package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
	"encoding/json"
	"github.com/andew42/brightlight/servers"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/stats"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"flag"
)

// The controller's frame buffer
var fb = controller.NewFrameBuffer()
var Statistics stats.Stats = stats.NewStats()

// Render the frame buffer as JSON
func sendFrameBufferToWebSocket(ws *websocket.Conn) error {

	// Send back the frame buffer as JSON
	rc, err := json.MarshalIndent(fb, "", " ")
	if err != nil {
		return err
	}
	_, err = ws.Write(rc)
	return err
}

// Handle HTTP requests to set all lights to a specific colour
func allLightsHandler(w http.ResponseWriter, r *http.Request) {

	// Colour value follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No colour specified", 406)
		log.Info("allLightsHandler no colour")
		return
	}

	colourValue, err := strconv.ParseInt(r.URL.Path[extIndex+1:], 16, 32)
	if err != nil {
		http.Error(w, "Invalid colour specified", 406)
		log.Info("allLightsHandler invalid colour")
		return
	}

	colourValueRgb := controller.NewRgbFromInt(int(colourValue))
	log.WithField("colourValueRgb", colourValueRgb).Info("allLightsHandler")
	animations.AnimateStaticColour(colourValueRgb)

	d, _ := json.Marshal(controller.IsDriverConnected())
	w.Write(d)
}

// Handle HTTP requests to run a named animation
func animationHandler(w http.ResponseWriter, r *http.Request) {

	// Animation name follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No animation name specified", 406)
		log.Info("animationHandler no name")
		return
	}

	animationName := r.URL.Path[extIndex+1:]
	err := animations.Animate(animationName)
	if err != nil {
		http.Error(w, err.Error()+" "+animationName, 406)
		log.Info("animationHandler " + err.Error())
		return
	}

	log.WithField("animationName", animationName).Info("animationHandler")

	d, _ := json.Marshal(controller.IsDriverConnected())
	w.Write(d)
}

// Handle HTTP requests to show strip lengths of room lights
func stripLengthHandler(w http.ResponseWriter, r *http.Request) {

	// Strip index, length follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No parameters specified", 406)
		log.Info("stripLengthHandler called with no parameters")
		return
	}

	// index
	config := strings.Split(r.URL.Path[extIndex+1:], ",")
	index, err := strconv.ParseInt(config[0], 10, 32);
	if err != nil || index < 0 || index > int64(len(fb.Strips)) {
		index = -1
	}

	// length
	length, err := strconv.ParseInt(config[1], 10, 32)
	if err != nil || length < 0 || length > controller.MaxLedStripLen {
		length = -1
	}

	if err = animations.AnimateStripLength(uint(index), uint(length)); err != nil {
		http.Error(w, "Invalid index or length", 406)
	}

	log.WithFields(
	log.Fields{"index":index, "length":length, "err":err, }).Info("stripLengthHandler called")
}

// Handle frame buffer web socket requests (web socket is closed when we return)
func frameBufferSocketHandler(ws *websocket.Conn) {

	for {
		fb.Mutex.Lock()
		// Fails if the client has disappeared
		if err := sendFrameBufferToWebSocket(ws); err != nil {
			fb.Mutex.Unlock()
			log.Info("frameBufferSocketHandler " + err.Error())
			return
		}
		// Wait for next frame buffer update
		fb.Cond.Wait()
		fb.Mutex.Unlock()
	}
}

// Handle stats web socket requests (web socket is closed when we return)
func statsSocketHandler(ws *websocket.Conn) {

	for {
		fb.Mutex.Lock()

		// Report stats for last second every second (50 frames)
		if Statistics.FrameCount == stats.ResetFrame {
			// Render the stats as JSON (fails if the client has disappeared)
			rc, err := json.MarshalIndent(Statistics, "", " ")
			if err == nil {
				_, err = ws.Write(rc)
			}
			if err != nil {
				fb.Mutex.Unlock()
				log.Info("statsSocketHandler" + err.Error())
				return
			}
		}
		// Wait for next frame buffer update
		fb.Cond.Wait()
		fb.Mutex.Unlock()
	}
}

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
	controller.StartDriver(fb, &Statistics)
	animations.StartDriver(fb, &Statistics)

	// Figure out where the content directory is GOPATH may contain : separated paths
	goPath := strings.Split(os.Getenv("GOPATH"), ":")
	contentPath := path.Join(goPath[0], "src/github.com/andew42/brightlight/ui")

	// Set up web routes (first static content)
	fs := http.FileServer(http.Dir(contentPath))
	http.Handle("/", fs)

	// Config requires PUT (write) support
	configServer := servers.Config {contentPath}
	http.HandleFunc("/config/", configServer.Handler)

	// TODO: MOVE ALL HANDLERS INTO SERVERS PACKAGE LIKE CONFIG
	// Requests to turn on all lights
	http.HandleFunc("/AllLights/", allLightsHandler)

	// Requests to run an animation
	http.HandleFunc("/Animation/", animationHandler)

	// Requests to show a strip length on the room lights
	http.HandleFunc("/StripLength/", stripLengthHandler)

	// Push frame buffer changes over a web socket
	http.Handle("/FrameBuffer", websocket.Handler(frameBufferSocketHandler))

	// Push stats info over a web socket
	http.Handle("/Stats", websocket.Handler(statsSocketHandler))

	// Start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error(err.Error())
	}

	log.Info("brightlight exited")
}
