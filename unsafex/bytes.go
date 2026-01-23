package unsafex

import (
	"unsafe"
)

// UnsafeBytes converts a string to a byte slice without copying the underlying data.
// It creates a new byte slice header pointing to the string's data memory address.
// WARNING: If the string data is allocated in read-only memory, modifying the
// returned byte slice will cause a panic.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeString converts a byte slice to a string without copying the underlying data.
// It converts the byte slice header directly to a string header.
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
// This works because the memory layouts are similar - the byte slice header
// can be directly cast to a string header (ignoring the Cap field).
func UnsafeString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
