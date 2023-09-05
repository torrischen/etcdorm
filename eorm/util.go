package eorm

import "unsafe"

func btos(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func stob(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		int
	}{s, len(s)}))
}
