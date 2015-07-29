package main

import (
	"golang.org/x/net/websocket"
	"encoding/json"
	"fmt"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/stats"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
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
	fmt.Println("allLightsHandler Called")

	// Colour value follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No colour specified", 406)
		return
	}

	colourValue, err := strconv.ParseInt(r.URL.Path[extIndex+1:], 16, 32)
	if err != nil {
		http.Error(w, "Invalid colour specified", 406)
		return
	}

	animations.AnimateStaticColour(controller.NewRgbFromInt(int(colourValue)))

	d, _ := json.Marshal(controller.IsDriverConnected())
	w.Write(d)
}

// Handle HTTP requests to run a named animation
func animationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("animationHandler Called")

	// Animation name follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No animation name specified", 406)
		return
	}

	animationName := r.URL.Path[extIndex+1:]
	err := animations.Animate(animationName)
	if err != nil {
		http.Error(w, err.Error()+" "+animationName, 406)
		return
	}

	d, _ := json.Marshal(controller.IsDriverConnected())
	w.Write(d)
}

// Handle HTTP requests to update configuration
func configHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("configHandler Called")

	// Strip index, length follows request path
	extIndex := strings.LastIndex(r.URL.Path, `/`)
	if extIndex == -1 {
		http.Error(w, "No configuration specified", 406)
		return
	}

	// index and length
	config := strings.Split(r.URL.Path[extIndex+1:], ",")
	index, err := strconv.ParseInt(config[0], 10, 32)
	if err != nil || index < 0 || index > int64(len(fb.Strips)) {
		http.Error(w, "Invalid index specified", 406)
		return
	}
	length, err := strconv.ParseInt(config[1], 10, 32)
	if err != nil || length < 0 || length > controller.MaxLedStripLen {
		http.Error(w, "Invalid length specified", 406)
		return
	}

	// Turn off animations TODO

	// Light up strip for length
	testColour := controller.NewRgb(128, 128, 128)
	backgroundColour := controller.NewRgb(0, 0, 0)
	for s := 0; s < len(fb.Strips); s++ {
		for l := 0; l < controller.MaxLedStripLen && l < len(fb.Strips[s].Leds) ; l++ {
			if int64(s) == index && int64(l) < length {
				fb.Strips[s].Leds[l] = testColour
			} else {
				fb.Strips[s].Leds[l] = backgroundColour
			}
		}
	}
	fb.Flush()
}

// Handle frame buffer web socket requests (web socket is closed when we return)
func frameBufferSocketHandler(ws *websocket.Conn) {
	for {
		fb.Mutex.Lock()
		// Fails if the client has disappeared
		if err := sendFrameBufferToWebSocket(ws); err != nil {
			fb.Mutex.Unlock()
			log.Println(err)
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
				log.Println("Stats socket handler %v", err)
				return
			}
		}

		// Wait for next frame buffer update
		fb.Cond.Wait()
		fb.Mutex.Unlock()
	}
}

func main() {
	// What are we running on?
	fmt.Printf("environment: %v %v %v\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// Start serial and animation drivers
	controller.StartDriver(fb, &Statistics)
	animations.StartDriver(fb, &Statistics)

	// Figure out where the content directory is GOPATH may contain : separated paths
	goPath := strings.Split(os.Getenv("GOPATH"), ":")
	contentPath := path.Join(goPath[0], "src/github.com/andew42/brightlight/ui")

	// Set up web routes (first static content)
	fs := http.FileServer(http.Dir(contentPath))
	http.Handle("/", fs)
	// Requests to turn on all lights
	http.HandleFunc("/AllLights/", allLightsHandler)
	// Requests to run an animation
	http.HandleFunc("/Animation/", animationHandler)
	// Requests to show a configuration on the room lights
	http.HandleFunc("/Config/", configHandler)
	// Push frame buffer changes over a web socket
	http.Handle("/FrameBuffer", websocket.Handler(frameBufferSocketHandler))
	// Push stats info over a web socket
	http.Handle("/Stats", websocket.Handler(statsSocketHandler))
	// Start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("ListenAndServe: " + err.Error())
	}
	fmt.Print("main exiting")
}
