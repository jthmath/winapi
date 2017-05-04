package com

import (
	"syscall"
	"unsafe"

	"github.com/jthmath/winapi"
)

const (
	TBPF_NOPROGRESS    = 0x00000000
	TBPF_INDETERMINATE = 0x00000001
	TBPF_NORMAL        = 0x00000002
	TBPF_ERROR         = 0x00000004
	TBPF_PAUSED        = 0x00000008
)

type TaskbarList interface {
	Unknown
	Init() error
	SetProgressValue(hWnd winapi.HWND, completed uint64, total uint64) error
	SetProgressState(hWnd winapi.HWND, tbpFlags int32) error
}

type _TaskbarListVtbl struct {
	_IUnknownVtbl
	HrInit                uintptr
	AddTab                uintptr
	DeleteTab             uintptr
	ActivateTab           uintptr
	SetActiveAlt          uintptr
	MarkFullscreenWindow  uintptr
	SetProgressValue      uintptr
	SetProgressState      uintptr
	RegisterTab           uintptr
	UnregisterTab         uintptr
	SetTabOrder           uintptr
	SetTabActive          uintptr
	ThumbBarAddButtons    uintptr
	ThumbBarUpdateButtons uintptr
	ThumbBarSetImageList  uintptr
	SetOverlayIcon        uintptr
	SetThumbnailTooltip   uintptr
	SetThumbnailClip      uintptr
	SetTabProperties      uintptr
}

type _TaskbarList struct {
	pVtbl *_TaskbarListVtbl
}

func (this *_TaskbarList) AddRef() uint32 {
	r, _, _ := syscall.Syscall(this.pVtbl.AddRef, 1, uintptr(unsafe.Pointer(this)), 0, 0)
	return uint32(r)
}

func (this *_TaskbarList) Release() uint32 {
	r, _, _ := syscall.Syscall(this.pVtbl.Release, 1, uintptr(unsafe.Pointer(this)), 0, 0)
	return uint32(r)
}

func (this *_TaskbarList) Init() error {
	r, _, _ := syscall.Syscall(this.pVtbl.HrInit, 1, uintptr(unsafe.Pointer(this)), 0, 0)
	hr := winapi.HRESULT(r)
	if hr == 0 {
		return nil
	} else {
		return hr
	}
}

func (this *_TaskbarList) SetProgressValue(hWnd winapi.HWND,
	completed uint64, total uint64) error {
	r, _, _ := syscall.Syscall6(this.pVtbl.SetProgressValue, 4,
		uintptr(unsafe.Pointer(this)), uintptr(hWnd), uintptr(completed), uintptr(total),
		0, 0)
	hr := winapi.HRESULT(r)
	if hr == 0 {
		return nil
	} else {
		return hr
	}
}

func (this *_TaskbarList) SetProgressState(hWnd winapi.HWND, tbpFlags int32) error {
	r, _, _ := syscall.Syscall(this.pVtbl.SetProgressState, 3,
		uintptr(unsafe.Pointer(this)), uintptr(hWnd), uintptr(tbpFlags))
	hr := winapi.HRESULT(r)
	if hr == 0 {
		return nil
	} else {
		return hr
	}
}

var CLSID_TaskbarList = winapi.MustMakeGuid("56FDF344-FD6D-11d0-958A-006097C9A090")

var IID_ITaskbarList4 = winapi.GUID{
	Data1: 3292383128,
	Data2: 38353,
	Data3: 19434,
	Data4: [8]uint8{144, 48, 187, 153, 226, 152, 58, 26},
}

func NewTaskbarList() (itl TaskbarList, err error) {
	var pTaskbarList *_TaskbarList
	err = CoCreateInstance(&CLSID_TaskbarList, 0, CLSCTX_INPROC_SERVER, &IID_ITaskbarList4,
		uintptr(unsafe.Pointer(&pTaskbarList)))
	if err == nil {
		itl = pTaskbarList
	}
	return
}
