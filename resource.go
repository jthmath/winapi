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
		// 制造一个悬挂指针
		// 不过这是安全的，因为 uint16 范围的地址肯定不会是一个有效的指针
		// WinAPI 函数不会误解
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
