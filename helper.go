package winapi

import (
	"syscall"
)

func SpecUTF16PtrFromString(s string) (*uint16, error) {
	if s == "" {
		return nil, nil
	} else {
		return syscall.UTF16PtrFromString(s)
	}
}

func UTF16PtrFromString(sl []string) (*uint16, error) {
	return nil, nil
}

func Uint32ToBinLittleEndian(n uint32) (r [4]byte) {
	r[0] = byte(n)
	r[1] = byte(n >> 8)
	r[2] = byte(n >> 16)
	r[3] = byte(n >> 24)
	return
}

func Uint64ToBinLittleEndian(n uint64) (r [8]byte) {
	r[0] = byte(n)
	r[1] = byte(n >> 8)
	r[2] = byte(n >> 16)
	r[3] = byte(n >> 24)
	r[4] = byte(n >> 32)
	r[5] = byte(n >> 40)
	r[6] = byte(n >> 48)
	r[7] = byte(n >> 56)
	return
}
