package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/andew42/brightlight/animations"
	"github.com/andew42/brightlight/controller"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"runtime"
)

// The controller's frame buffer
var fb = controller.NewFrameBuffer()

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

	fb.SetColour(colourValue)
	fb.SignalChanged()
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
	}
}

// Handle web socket requests (web socket is closed when we return)
func frameBufferSocketHandler(ws *websocket.Conn) {
	processSocket(ws)
}

// Process a web socket request (each request in own go routine)
func processSocket(ws *websocket.Conn) {
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

// HTTP content cache
var contentCache = make(map[string][]byte)

func main() {
	// What are we running on?
	fmt.Printf("%v %v %v\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

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

	// Send frame buffer changes over to Teensy controller
	//go processTeensyUpdate()
	go controller.TeensyDriver(fb)

	http.HandleFunc("/", staticContentHandler)
	http.HandleFunc("/AllLights/", allLightsHandler)
	http.HandleFunc("/Animation/", animationHandler)
	http.Handle("/FrameBuffer", websocket.Handler(frameBufferSocketHandler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("ListenAndServe: " + err.Error())
	}

	fmt.Print("done")
}
