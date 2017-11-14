package hue

import "net/http"

type cmdContext struct {
	FullState              *fullState
	BrightlightCommandChan chan interface{}
	User                   string
	Method                 string
	Resource               []string
	ResourceUrl            string
	Body                   []byte
	ErrorAddress           string
	W                      http.ResponseWriter
}
