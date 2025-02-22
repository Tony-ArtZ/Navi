package utils

import (
	"os"
	"syscall"
	"unsafe"
)

func EnableRawMode() {
	var termios syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCGETS, uintptr(unsafe.Pointer(&termios)))
	termios.Lflag &^= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCSETS, uintptr(unsafe.Pointer(&termios)))
}

func DisableRawMode() {
	var termios syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCGETS, uintptr(unsafe.Pointer(&termios)))
	termios.Lflag |= syscall.ICANON | syscall.ECHO
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TCSETS, uintptr(unsafe.Pointer(&termios)))
}

func ReadKey() (rune, error) {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return rune(buf[0]), nil
}
