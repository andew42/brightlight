package stats

import (
	"encoding/json"
	"github.com/andew42/brightlight/config"
	"sort"
	"strconv"
	"time"
)

// Statistics on animation frame times and serial send times
type Stats struct {
	FrameRenderTime StatsBlock
	FrameSyncJitter StatsBlock
	SerialSendTime  []StatsBlock
	FramePeriodMs   string
	FrameRateHz     string
	GcCount         int64
	GcPauseTime     string
}

type StatsBlock struct {
	Name          string
	SampleCount   int64
	DroppedFrames int64
	TotalTime     time.Duration
	MinTime       time.Duration
	MaxTime       time.Duration
	AverageTime   time.Duration
}

// ByAge implements sort.Interface for []StatsBlock based on name
type ByName []StatsBlock

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Marshal StatsBlock durations to 2dp we use an embedded struct with type alias technique
// http://choly.ca/post/go-json-marshalling/
func (sb *StatsBlock) MarshalJSON() ([]byte, error) {
	type Alias StatsBlock
	return json.Marshal(&struct {
		MinTime     string
		MaxTime     string
		AverageTime string
		*Alias
	}{
		MinTime:     strconv.FormatFloat(float64(sb.MinTime)/float64(time.Millisecond), 'f', 2, 64),
		MaxTime:     strconv.FormatFloat(float64(sb.MaxTime)/float64(time.Millisecond), 'f', 2, 64),
		AverageTime: strconv.FormatFloat(float64(sb.AverageTime)/float64(time.Millisecond), 'f', 2, 64),
		Alias:       (*Alias)(sb),
	})
}

func (sb *StatsBlock) resetStatsBlock() {
	sb.MinTime = 1<<63 - 1
}

func addSampleToBlock(block *StatsBlock, name string, sample time.Duration) {

	// Update the block
	block.Name = name
	block.SampleCount++
	block.TotalTime += sample
	if sample < block.MinTime {
		block.MinTime = sample
	}
	if sample > block.MaxTime {
		block.MaxTime = sample
	}

	block.AverageTime = time.Duration(block.TotalTime.Nanoseconds() / block.SampleCount)
}

func findBlockWithAddByName(blocks *[]StatsBlock, name string) *StatsBlock {

	// Return block if it exists
	for i, b := range *blocks {
		if b.Name == name {
			return &(*blocks)[i]
		}
	}
	// otherwise append a new block
	newStatsBlock := StatsBlock{}
	newStatsBlock.resetStatsBlock()
	newStatsBlock.Name = name
	*blocks = append(*blocks, newStatsBlock)
	// Sort the blocks by name so they don't bounce around in the UI
	sort.Sort(ByName(*blocks))
	return &(*blocks)[len(*blocks)-1]
}

func addSampleToBlockList(blocks *[]StatsBlock, name string, sample time.Duration) {

	// Find or add block
	block := findBlockWithAddByName(blocks, name)

	// Update the block
	block.SampleCount++
	block.TotalTime += sample
	if sample < block.MinTime {
		block.MinTime = sample
	}
	if sample > block.MaxTime {
		block.MaxTime = sample
	}

	block.AverageTime = time.Duration(block.TotalTime.Nanoseconds() / block.SampleCount)
}

// Constructor
func NewStats() *Stats {
	var s Stats
	s.reset()
	return &s
}

// Reset statistic
func (stats *Stats) reset() {

	stats.FrameRenderTime.resetStatsBlock()
	stats.FrameSyncJitter.resetStatsBlock()
	stats.FramePeriodMs = strconv.FormatFloat(float64(config.FramePeriodMs/time.Millisecond), 'f', 2, 64)
	stats.FrameRateHz = strconv.Itoa(config.FrameFrequencyHz)
}

// Adds a sample point for animation frame render
func (stats *Stats) addFrameRenderTimeSample(name string, sample time.Duration) {
	addSampleToBlock(&stats.FrameRenderTime, name, sample)
}

// Adds a sample point for frame sync timer jitter
func (stats *Stats) addFrameSyncJitterSample(name string, sample time.Duration) {
	addSampleToBlock(&stats.FrameSyncJitter, name, sample)
}

// Adds sample point for serial send
func (stats *Stats) addSerialSendTimeSample(name string, sample time.Duration) {
	addSampleToBlockList(&stats.SerialSendTime, name, sample)
}

// Adds a serial dropped frame
func (stats *Stats) addSerialDroppedFrame(name string) {
	block := findBlockWithAddByName(&stats.SerialSendTime, name)
	block.DroppedFrames++
}

// Adds a render dropped frame
func (stats *Stats) addFrameRenderDroppedFrame() {
	stats.FrameRenderTime.DroppedFrames++
}
