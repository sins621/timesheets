package utils

import "unsafe"

func Bool2int(b bool) int {
	return int(*(*byte)(unsafe.Pointer(&b))) // xddddddd
}
