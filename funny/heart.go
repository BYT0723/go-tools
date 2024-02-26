package funny

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"
)

func Heart() error {
	size, err := getTermSize()
	if err != nil {
		return err
	}
	draw(int(size.Row), int(size.Col))
	return nil
}

type termSize struct {
	Row uint16
	Col uint16
}

func getTermSize() (size *termSize, err error) {
	win := termSize{0, 0}

	res, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&win)))
	if int(res) == -1 {
		return
	}
	return &win, nil
}

func draw(maxRow, maxCol int) {
	// 打印爱心形状
	for y := maxRow / 2; y >= -maxRow/2; y-- {
		for x := -maxCol / 2; x <= maxCol/2; x++ {
			if love(x, y) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
	}

	for row := maxRow / 2; row >= 0-maxRow/2; row-- {
		for col := 0 - maxCol/2; col <= maxCol; col++ {
			if love(col, row) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

// 爱心函数
func love(x, y int) bool {
	xf := float64(x) / 40.0
	yf := float64(y) / 20.0
	return math.Pow(math.Pow(xf, 2)+math.Pow(yf, 2)-1, 3)-math.Pow(xf, 2)*math.Pow(yf, 3) <= 0
}
