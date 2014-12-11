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

func UTF16FromMultiStrings(sl []string) ([]uint16, error) {
	n := len(sl)
	if n == 0 {
		return []uint16{0}, nil
	}
	//data := make([]uint16, 2048)
	for i := 0; i < n; i++ {
		if _, err := syscall.UTF16FromString(sl[i]); err != nil {
			return []uint16{}, err
		}
	}
	return []uint16{}, nil
}

// 拼接两个Unicode字符串
//
func Uint16SliceCat(a []uint16, b []uint16) []uint16 {
	var L int = 0
	s := make([]uint16, 0)
	L = len(a)
	for i := 0; i < L; i++ {
		s = append(s, a[i])
	}
	L = len(b)
	for i := 0; i < L; i++ {
		s = append(s, b[i])
	}
	return s
}
