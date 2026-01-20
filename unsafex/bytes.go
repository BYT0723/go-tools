package unsafex

import (
	"unsafe"
)

// UnsafeBytes()
// 新建byte slice header, 将Data指向字符串的Data内存地址
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString()
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
func UnsafeString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
