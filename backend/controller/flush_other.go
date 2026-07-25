//go:build !linux

package controller

// Output queue flushing is only implemented for Linux (the deployment
// target); elsewhere closing the port without it is fine for development
func flushOutputQueue(fd uintptr) {
}
