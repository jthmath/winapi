// +build windows

package com

import (
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
	SetProgressState(hWnd winapi.HWND, tbpFlags int32) error
	SetProgressValue(hWnd winapi.HWND, completed uint64, total uint64) error
}

type _TaskbarList struct {
	_hWnd winapi.HWND
}

func (this *_TaskbarList) AddRef() uint32 {
	return 1
}

func (this *_TaskbarList) Release() uint32 {
	return 1
}

func (this *_TaskbarList) Init() error {
	return nil
}

func NewTaskbarList() TaskbarList {
	return nil
}
