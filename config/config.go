package config

import "time"

// Maximum led strip length (must match teensy firmware)
const MaxLedStripLen = 300
const StripsPerTeensy = 8

// Frame buffer update frequency used by animation and stats
// 20ms -> 50Hz, 25ms -> 40Hz *, 40ms -> 25Hz, 50ms -> 20Hz
const framePeriodMs = 25
const FramePeriodMs = framePeriodMs * time.Millisecond
const FrameFrequencyHz = 1 / (framePeriodMs / 1000.0)

// Statistics update period
const StatsUpdatePeriodMs = 1000 * time.Millisecond
