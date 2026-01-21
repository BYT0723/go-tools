package unsafex

import (
	"unsafe"
)

// Bytes(string)
// 新建byte slice header, 将Data指向字符串的Data内存地址
// NOTE: 注意，如果string data内存非配在只读段内，修改[]byte内容会导致panic
func Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// String([]byte)
// 将byte slice header直接转换成string header
//
//	type SliceHeader struct {
//		Data uintptr
//		Len  int
//		Cap  int
//	}
//
//	type StringHeader struct {
//		Data uintptr
//		Len  int
//	}
//
// 利用内存布局相似的特性，直接将byte slice header转换成string header使用
func String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
