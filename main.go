package main

import (
	"golang.org/x/net/websocket"
	"encoding/json"
	"fmt"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"github.com/andew42/brightlight/stats"
	"io/ioutil"
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

// Read each file's content, in the content directory, into memory against its relative path
func cacheContentDirectory(path string, basePathLen int, cache map[string][]byte) error {
	info, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, i := range info {
		itemPath := path + "/" + i.Name()
		if i.IsDir() {
			err = cacheContentDirectory(itemPath, basePathLen, cache)
			if err != nil {
				return err
			}
		} else {
			cache[itemPath[basePathLen:]], err = ioutil.ReadFile(itemPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Handle HTTP requests for static content
func staticContentHandler(w http.ResponseWriter, r *http.Request) {
	// Crudely determine content type from path extension
	extIndex := strings.LastIndex(r.URL.Path, ".")
	if extIndex == -1 {
		http.Error(w, "Content file has no extension", 406)
		return
	}
	// Look up path in cache
	content, ok := contentCache[r.URL.Path]
	if !ok {
		http.Error(w, "File not found", 404)
		return
	}
	// Return content with correct type in header
	contentType := ""
	switch r.URL.Path[extIndex+1:] {
	case "js":
		contentType = "application/javascript"
		// TODO: Add HTML Type?
	default:
		contentType = "text/" + r.URL.Path[extIndex+1:]
	}
	w.Header().Set("Content-Type", contentType)
	fmt.Printf("%v CONTENT TYPE %v\n", r.URL.Path, contentType)
	w.Write(content)
}

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
		for l := 0; l < controller.MaxLedStripLen; l++ {
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

// HTTP content cache
var contentCache = make(map[string][]byte)

func main() {
	// What are we running on?
	fmt.Printf("environment: %v %v %v\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// Figure out where the content directory is GOPATH may contain : separated paths
	goPath := strings.Split(os.Getenv("GOPATH"), ":")
	contentPath := path.Join(goPath[0], "src/github.com/andew42/brightlight/content")

	// Create an in memory cache of the content directory
	err := cacheContentDirectory(contentPath, len(contentPath), contentCache)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: Print out cache summary
	for k, v := range contentCache {
		fmt.Println(k, len(v))
	}

	// Start serial and animation drivers
	controller.StartDriver(fb, &Statistics)
	animations.StartDriver(fb, &Statistics)

	// Start web handlers
	http.HandleFunc("/", staticContentHandler)
	http.HandleFunc("/AllLights/", allLightsHandler)
	http.HandleFunc("/Animation/", animationHandler)
	http.HandleFunc("/Config/", configHandler)
	http.Handle("/FrameBuffer", websocket.Handler(frameBufferSocketHandler))
	http.Handle("/Stats", websocket.Handler(statsSocketHandler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("ListenAndServe: " + err.Error())
	}

	fmt.Print("main exiting")
}
