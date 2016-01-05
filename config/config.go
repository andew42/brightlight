package config

import "time"

// Frame buffer update frequency used by animation and stats
// 20ms -> 50Hz, 25ms -> 40Hz *, 40ms -> 25Hz, 50ms -> 20Hz
const framePeriodMs = 25
const FramePeriodMs = framePeriodMs * time.Millisecond
const FrameFrequencyHz = 1 / (framePeriodMs / 1000.0)
