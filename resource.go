//go:build windows

package winapi

import (
	"syscall"
	"unsafe"
)

type ResourceConcept interface {
	uint16 | string
}

type Resource struct {
	n uint16
	s string
}

func (r *Resource) GetWinStr() (*uint16, error) {
	if r.s == "" {
		// Creates a dangling pointer
		// This is safe, however, because an address in the uint16 range is definitely not a valid pointer
		// WinAPI won't misinterpret
		return (*uint16)(unsafe.Pointer(uintptr(r.n))), nil
	} else {
		return syscall.UTF16PtrFromString(r.s)
	}
}

func MakeResource[R ResourceConcept](r R) Resource {
	switch vv := any(r).(type) {
	case uint16:
		return Resource{
			n: vv,
			s: "",
		}
	case string:
		return Resource{
			n: 0,
			s: vv,
		}
	default:
		panic("big error")
	}
}
