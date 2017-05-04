package com

import (
	"syscall"
	"unsafe"
)

type Unknown interface {
	AddRef() uint32
	Release() uint32
}

type _IUnknownVtbl struct {
	_QueryInterface uintptr
	AddRef          uintptr
	Release         uintptr
}

type _IUnknown struct {
	pVtbl *_IUnknownVtbl
}

func (this *_IUnknown) AddRef() uint32 {
	r, _, _ := syscall.Syscall(this.pVtbl.AddRef, 1, uintptr(unsafe.Pointer(this)), 0, 0)
	return uint32(r)
}

func (this *_IUnknown) Release() uint32 {
	r, _, _ := syscall.Syscall(this.pVtbl.Release, 1, uintptr(unsafe.Pointer(this)), 0, 0)
	return uint32(r)
}
