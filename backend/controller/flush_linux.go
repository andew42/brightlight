//go:build linux

package controller

import "syscall"

// tcflush(fd, TCOFLUSH): discard output that has been written to the tty
// but not yet transmitted. Values from asm-generic/ioctls.h and termios.h
func flushOutputQueue(fd uintptr) {

	const TCFLSH = 0x540B
	const TCOFLUSH = 1
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, fd, TCFLSH, TCOFLUSH)
}
