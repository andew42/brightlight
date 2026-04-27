package stats

import (
	"github.com/andew42/brightlight/config"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"strconv"
	"time"
)

const (
	statFrameRenderTime = iota
	statFrameSyncJitter
	statSerialSendTime
	statSerialDroppedFrame
	statFrameRenderDroppedFrame
)

type statSample struct {
	StatType int    // const enum
	Src      string // Source of timing (e.g. Teensy 1 or rainbow
	Sample   time.Duration
}

func AddFrameRenderTimeSample(sample time.Duration) {
	statSampleChannel <- statSample{statFrameRenderTime, "Frame Render", sample}
}

// AddFrameSyncJitterSample Public calls to add new stats samples
func AddFrameSyncJitterSample(sample time.Duration) {
	statSampleChannel <- statSample{statFrameSyncJitter, "Frame Sync Jitter", sample}
}

func AddSerialSendTimeSample(src string, sample time.Duration) {
	statSampleChannel <- statSample{statSerialSendTime, src, sample}
}

func AddSerialDroppedFrame(src string) {
	statSampleChannel <- statSample{statSerialDroppedFrame, src, 0}
}

func AddFrameRenderDroppedFrame() {
	statSampleChannel <- statSample{statFrameRenderDroppedFrame, "", 0}
}

var statSampleChannel = make(chan statSample, 16)

// Parameter sent to internal add listener channel
type addListenerParams struct {
	name string
	src  chan *Stats
}

// AddListener is used to request a new statistics update every update period
// frame buffer changes. src is a channel down which new statistics are sent.
// The src channel is sent to the done channel when updates are no longer
// required
func AddListener(name string) (src chan *Stats, done chan<- chan *Stats) {
	newSrc := make(chan *Stats)
	addListener <- addListenerParams{name, newSrc}
	return newSrc, listenerDone
}

var addListener = make(chan addListenerParams)
var listenerDone = make(chan chan *Stats)

// StartDriver Stats driver waits for samples and publishes results to listeners
func StartDriver() {
	go func() {
		var gcStats debug.GCStats
		var listeners = make(map[chan *Stats]string)
		stats := NewStats()
		statsUpdate := time.Tick(config.StatsUpdatePeriodMs)
		for {
			select {
			// Time to send a stats update
			case <-statsUpdate:
				// Update the stats GC
				debug.ReadGCStats(&gcStats)
				stats.GcCount = gcStats.NumGC
				if len(gcStats.Pause) > 0 {
					pauseTime := float64(gcStats.Pause[0]) / float64(time.Millisecond)
					stats.GcPauseTime = strconv.FormatFloat(pauseTime, 'f', 2, 64)
				}
				// Send to all listeners, that are idle, the most recent frame buffer
				for k := range listeners {
					// Send the latest stats if listener has processed the last one
					select {
					case k <- stats:
					default:
					}
				}
				// New stats so we don't corrupt last one sent by reference
				stats = NewStats()

				// New stats sample from somewhere
			case s := <-statSampleChannel:
				{
					switch s.StatType {
					case statFrameRenderTime:
						stats.addFrameRenderTimeSample(s.Src, s.Sample)
					case statFrameSyncJitter:
						stats.addFrameSyncJitterSample(s.Src, s.Sample)
					case statSerialSendTime:
						stats.addSerialSendTimeSample(s.Src, s.Sample)
					case statSerialDroppedFrame:
						stats.addSerialDroppedFrame(s.Src)
					case statFrameRenderDroppedFrame:
						stats.addFrameRenderDroppedFrame()
					default:
						panic("Unknown stat type " + strconv.Itoa(s.StatType))
					}
				}

				// Process new listener requests
			case newListener := <-addListener:
				log.WithField("name", newListener.name).Info("Stats listener added")
				listeners[newListener.src] = newListener.name

				// Process remove listener request
			case listenerToRemove := <-listenerDone:
				log.WithField("name", listeners[listenerToRemove]).Info("Stats listener removed")
				delete(listeners, listenerToRemove)
			}
		}
	}()
}
