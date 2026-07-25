package controller

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"time"

	"log/slog"
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
	"windows": {"COM4", "COM3"},
	"linux":   {"/dev/ttyACM0", "/dev/ttyACM1"},
}

// Determine port name based on OS index and mapping table
func getPortName(portMappings map[string][]string, index int) string {

	portNames, ok := portMappings[runtime.GOOS]
	if !ok {
		slog.Warn("No port mappings for OS", "os", runtime.GOOS)
		return ""
	}

	if index < 0 || index >= len(portNames) {
		slog.Warn("No port mappings for index", "index", index)
		return ""
	}

	return portNames[index]
}

// Retry port open until it succeeds, or nil once shutdown is requested
func openUsbPortWithRetry(port string) *os.File {

	errorLogged := false
	for {
		select {
		case <-shutdownChan:
			return nil
		default:
		}
		f, err := os.OpenFile(port, os.O_RDWR, 0)
		if err == nil {
			slog.Info("openUsbPortWithRetry connected", "port", port)

			// Set raw mode on raspberry pi, if we don't set raw mode
			// xon/xoff character in the frame buffer cause problems
			if runtime.GOOS == "linux" {
				cmd := exec.Command("stty", "-F", port, "raw")
				if err := cmd.Run(); err != nil {
					slog.Error("openUsbPortWithRetry failed to set stty raw mode", "error", err.Error())
				}
			}

			return f
		}

		if !errorLogged {
			slog.Warn("openUsbPortWithRetry failed to open port", "error", err.Error())
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
