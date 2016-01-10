package framebuffer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/config"
	"strconv"
	"time"
)

// Frame buffer is a slice of strips A Mutex
// and Cond are used to broadcast changes
type FrameBuffer struct {
	Strips []LedStrip
}

// Create a frame buffer
func NewFrameBuffer() *FrameBuffer {

	var fb FrameBuffer

	// TODO: Make this more dynamic from config file?
	fb.Strips = make([]LedStrip, 0, config.StripsPerTeensy)
	// 0, 1 Unused strips
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 0))
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 0))

	// 2 Bed wall
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 168))

	// 3 Bed curtains
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 164))

	// 4 Bed ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 165))

	// 5 Dressing table wall
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 85))

	// 6 Dressing table ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 80))

	// 7 Dressing table curtain
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 162))

	// 8 Bathroom mirror wall
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 172))

	// 9 Bath ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 226))

	// 10 Bath+ wall
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 291))

	// 11 Bathroom mirror ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 162))

	// 12 Unused
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 0))

	// 13 Left of door ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(true, 88))

	// 14 Right of door ceiling
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 142))

	// 15 Right of door wall
	fb.Strips = append(fb.Strips, *NewLedStrip(false, 122))

	// Sanity check
	numberOfStrips := len(fb.Strips)
	if numberOfStrips <= 0 || numberOfStrips%config.StripsPerTeensy != 0 {
		log.WithField("StripsPerTeensy", strconv.Itoa(config.StripsPerTeensy)).Panic("framebuffer strips must be multiple of")
	}
	return &fb
}

// Parameter sent to internal add listener channel
type addListenerParams struct {
	name string
	src  chan *FrameBuffer
}

// Used internally to communicate an add listener request
var addListener = make(chan addListenerParams)

// Used by listeners to communicate a remove request
var listenerDone = make(chan chan *FrameBuffer)

// AddListener is used to request a new frame buffer update every time the
// frame buffer changes. src is a channel down which new frame buffers are
// sent. The src channel is sent to the done channel when updates are no
// longer required
func AddListener(name string) (src chan *FrameBuffer, done chan<- chan *FrameBuffer) {

	newSrc := make(chan *FrameBuffer)
	addListener <- addListenerParams{name, newSrc}
	return newSrc, listenerDone
}

// Acquire a frame buffer
func StartDriver(renderer chan *FrameBuffer) {

	go func() {
		// All the frame buffer listeners, can be added and removed dynamically to
		// support web page(s) with virtual frame buffer displays that come and go
		var listeners map[chan *FrameBuffer]string = make(map[chan *FrameBuffer]string)
		lastRenderedFrameBuffer := NewFrameBuffer()
		renderInProgress := false
		frameSync := time.Tick(config.FramePeriodMs)
		for {
			select {
			case <-frameSync:
				// Frame tick, first collect timer jitter
				//			started := time.Now()
				//			jitter := started.Sub(nextFrameTime)
				//			nextFrameTime = started.Add(config.FramePeriodMs)

				// Send to all listeners, that are idle, the most recent frame buffer
				for k, _ := range listeners {
					// Send the latest frame buffer if listener has processed the last one
					select {
					case k <- lastRenderedFrameBuffer:
					default:
					}
				}

				// Render next frame, if current one is done
				if !renderInProgress {
					renderInProgress = true
					renderer <- NewFrameBuffer()
				}

			// Process frame buffer render complete
			case lastRenderedFrameBuffer = <-renderer:
				renderInProgress = false

			// Process new listener requests
			case newListener := <-addListener:
				log.WithField("name", newListener.name).Info("Framebuffer listener added")
				listeners[newListener.src] = newListener.name

			// Process remove listener request
			case listenerToRemove := <-listenerDone:
				log.WithField("name", listeners[listenerToRemove]).Info("Framebuffer listener removed")
				delete(listeners, listenerToRemove)
			}
		}
	}()
}
