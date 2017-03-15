package controller

import (
	log "github.com/Sirupsen/logrus"
	"runtime"
	"time"
	"os/exec"
	"os"
	"errors"
)

// By OS a list of relay port names in index order
var relayPortMappings = map[string][]string{
	"darwin":  {"/dev/cu.usbserial"},
	"windows": {""},
	"linux":   {"/dev/ttyUSB0"},
}

// By OS a list of teensy port names in index order
var teensyPortMappings = map[string][]string{
	// Teensy 3.0 "/dev/cu.usbmodem103721"
	// Teensy 3.1 "/dev/cu.usbmodem103101"
	"darwin":  {"/dev/cu.usbmodem288181"},
	"windows": {"COM3", "COM34"},
	"linux":   {"/dev/ttyACM0", "/dev/ttyACM1"},
}

// Determine port name based on OS index and mapping table
func getPortName(portMappings map[string][]string, index int) string {

	portNames, ok := portMappings[runtime.GOOS]
	if !ok {
		log.WithField("os", runtime.GOOS).Warn("No port mappings for OS")
		return ""
	}

	if index < 0 || index >= len(portNames) {
		log.WithField("index", index).Warn("No port mappings for index")
		return ""
	}

	return portNames[index]
}

// Retry port open until it succeeds
func openUsbPortWithRetry(port string) *os.File {

	errorLogged := false
	for {
		f, err := os.OpenFile(port, os.O_RDWR, 0)
		if err == nil {
			log.WithField("port", port).Info("openUsbPortWithRetry connected")

			// Set raw mode on raspberry pi, if we don't set raw mode
			// xon/xoff character in the frame buffer cause problems
			if runtime.GOOS == "linux" {
				cmd := exec.Command("stty", "-F", port, "raw")
				if err := cmd.Run(); err != nil {
					log.WithField("error", err.Error()).Error("openUsbPortWithRetry failed to set stty raw mode")
				}
			}

			return f
		}

		if !errorLogged {
			log.WithField("error", err.Error()).Warn("openUsbPortWithRetry failed to open port")
			errorLogged = true
		}

		// Try again in a second
		time.Sleep(1000 * time.Millisecond)
	}
}

var readTimeoutError = errors.New("readUntilBufferFull timeout")

// Blocks reading until data is full or EOF or Error occurs
func readUntilBufferFull(f *os.File, data []byte, timeout time.Duration) error {

	doneTime := time.Now().Add(timeout)
	length, err := f.Read(data)
	for err == nil && length == 0 {
		if time.Now().After(doneTime) {
			return readTimeoutError
		}
		length, err = f.Read(data)
	}
	return err
}
