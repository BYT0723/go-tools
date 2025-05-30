//go:build windows

package system

import "golang.org/x/sys/windows"

func GetOEMCP() uint32 {
	ret, _, _ := windows.NewLazyDLL("kernel32.dll").NewProc("GetOEMCP").Call()
	return uint32(ret)
}
