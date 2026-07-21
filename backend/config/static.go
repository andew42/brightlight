package config

import (
	"os"
	"strings"
	"time"
)

// Site names the hardware installation this server drives. Selected at
// runtime with the BRIGHTLIGHT_SITE environment variable ("titania" or
// "bedroom"); anything else falls back to titania.
var Site = func() string {
	if strings.EqualFold(os.Getenv("BRIGHTLIGHT_SITE"), "bedroom") {
		return "bedroom"
	}
	return "titania"
}()

// Titania (or bedroom)
var Titania = Site == "titania"

// MaxLedStripLen Maximum led strip length (must match Teensy firmware)
const MaxLedStripLen = 300
const StripsPerTeensy = 8

// Frame buffer update frequency used by animation and stats
// 20ms -> 50Hz, 25ms -> 40Hz, 40ms -> 25Hz, 50ms -> 20Hz
const framePeriodMs = 40
const FramePeriodMs = framePeriodMs * time.Millisecond
const FrameFrequencyHz = 1 / (framePeriodMs / 1000.0)

// StatsUpdatePeriodMs Statistics update period
const StatsUpdatePeriodMs = 1000 * time.Millisecond
