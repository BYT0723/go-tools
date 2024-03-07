package system

import (
	"syscall"
	"unsafe"
)

type TermSize struct {
	Row uint16
	Col uint16
}

func GetTermSize() (size *TermSize, err error) {
	win := TermSize{0, 0}

	res, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&win)))
	if int(res) == -1 {
		return
	}
	return &win, nil
}
