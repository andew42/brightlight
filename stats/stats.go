package stats

import "time"

// Statistics reset after every 50 frames
const ResetFrame = 50

// Statistics on animation frame times and serial send times
type Stats struct {
	// Animation
	FrameCount int64
	TotalFrameTime time.Duration
	MinFrameTime time.Duration
	MaxFrameTime time.Duration
	AverageFrameTime time.Duration
	TotalJitter time.Duration
	MinJitter time.Duration
	MaxJitter time.Duration
	AverageJitter time.Duration
	// Serial
	SendCount int64
	TotalSendTime time.Duration
	MinSendTime time.Duration
	MaxSendTime time.Duration
	AverageSendTime time.Duration
}

// Constructor
func NewStats() Stats {
	var s Stats
	s.Reset()
	return s
}

// Reset statistic
func (stats *Stats) Reset() {
	stats.FrameCount = 0
	stats.TotalFrameTime = 0
	stats.MinFrameTime = 1<<63 - 1
	stats.MaxFrameTime = 0
	stats.AverageFrameTime = 0
	stats.TotalJitter = 0
	stats.MinJitter = 1<<63 - 1
	stats.MaxJitter = 0
	stats.AverageJitter = 0
	stats.SendCount = 0
	stats.TotalSendTime = 0
	stats.MinSendTime = 1<<63 - 1
	stats.MaxSendTime = 0
	stats.AverageSendTime = 0
}

// Adds a sample point for animation frame
func (stats *Stats) AddAnimation(frameTime time.Duration, jitter time.Duration) {

	// Reset statistics every 50 frames, socket listeners should have update after last frame
	if stats.FrameCount == ResetFrame {
		stats.Reset()
	}
	stats.FrameCount++

	// Duration statistics
	stats.TotalFrameTime += frameTime
	if frameTime < stats.MinFrameTime {
		stats.MinFrameTime = frameTime
	}
	if frameTime > stats.MaxFrameTime {
		stats.MaxFrameTime = frameTime
	}
	stats.AverageFrameTime = time.Duration(stats.TotalFrameTime.Nanoseconds() / stats.FrameCount)

	// Jitter statistics
	stats.TotalJitter += jitter
	if jitter < stats.MinJitter {
		stats.MinJitter = jitter
	}
	if jitter > stats.MaxJitter {
		stats.MaxJitter = jitter
	}
	stats.AverageJitter = time.Duration(stats.TotalJitter.Nanoseconds() / stats.FrameCount)
}

// Adds sample point for serial send
func (stats *Stats) AddSerial(sendTime time.Duration) {

	stats.SendCount++

	// Serial statistics
	stats.TotalSendTime += sendTime
	if sendTime < stats.MinSendTime {
		stats.MinSendTime = sendTime
	}
	if sendTime > stats.MaxSendTime {
		stats.MaxSendTime = sendTime
	}
	stats.AverageSendTime = time.Duration(stats.TotalSendTime.Nanoseconds() / stats.SendCount)
}
