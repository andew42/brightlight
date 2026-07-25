package controller

import (
	"os"
	"sync"
	"time"

	"log/slog"
)

// Closed by Shutdown to ask the driver go routines to quiesce and close
// their USB serial ports
var shutdownChan = make(chan struct{})

// Tracks driver go routines that have not yet finished shutting down
var shutdownWaitGroup sync.WaitGroup

// Shutdown asks the Teensy and relay drivers to stop writing, flush and
// close their USB serial ports, then waits for them (up to timeout).
// Exiting without this leaves the kernel to tear down a tty that may
// have a full output queue and USB transfers in flight; on a Raspberry
// Pi that has been seen to wedge the dwc_otg USB controller and freeze
// the whole machine (the Ethernet is on the same USB bus)
func Shutdown(timeout time.Duration) {

	close(shutdownChan)
	done := make(chan struct{})
	go func() {
		shutdownWaitGroup.Wait()
		close(done)
	}()
	select {
	case <-done:
		slog.Info("controller drivers shut down cleanly")
	case <-time.After(timeout):
		slog.Warn("controller driver shutdown timed out")
	}
}

// Discard any unsent output then close the port, so the tty is torn
// down with an empty queue
func closePortQuiesced(f *os.File) {

	flushOutputQueue(f.Fd())
	if err := f.Close(); err != nil {
		slog.Warn("closePortQuiesced failed to close port", "error", err.Error())
	}
}
