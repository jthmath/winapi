//go:build windows

package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	WM_NULL    uint32 = 0x0000
	WM_CREATE  uint32 = 0x0001
	WM_DESTROY uint32 = 0x0002
	WM_MOVE    uint32 = 0x0003
	WM_SIZE    uint32 = 0x0005

	WM_ACTIVATE uint32 = 0x0006

	WM_PAINT uint32 = 0x000F

	/*
	 * WM_ACTIVATE state values
	 */
	WA_INACTIVE    = 0
	WA_ACTIVE      = 1
	WA_CLICKACTIVE = 2

	WM_CLOSE uint32 = 0x0010
	WM_QUIT  uint32 = 0x0012

	WM_GETMINMAXINFO uint32 = 0x0024

	WM_COMMAND uint32 = 0x0111
)

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func GetMessage(pMsg *MSG, hWnd HWND, wMsgFilterMin uint32, wMsgFilterMax uint32) int32 {
	r1, _, _ := syscall.SyscallN(procGetMessage.Addr(),
		uintptr(unsafe.Pointer(pMsg)),
		uintptr(hWnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax))
	return int32(r1)
}

func TranslateMessage(pMsg *MSG) error {
	r1, _, _ := syscall.SyscallN(procTranslateMessage.Addr(), uintptr(unsafe.Pointer(pMsg)))
	if r1 == 0 {
		return errors.New("winapi: TranslateMessage failed")
	} else {
		return nil
	}
}

func DispatchMessage(pMsg *MSG) uintptr {
	r1, _, _ := syscall.SyscallN(procDispatchMessage.Addr(), uintptr(unsafe.Pointer(pMsg)))
	return r1
}

func PostQuitMessage(ExitCode int32) {
	syscall.SyscallN(procPostQuitMessage.Addr(), uintptr(ExitCode))
}

func RegisterWindowMessage(str string) (message uint32, err error) {
	p, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		return
	}

	r1, _, e1 := syscall.SyscallN(procRegisterWindowMessage.Addr(), uintptr(unsafe.Pointer(p)))
	if r1 == 0 {
		err = MakeFromWinError(e1)
		return
	}

	message = uint32(r1)
	return
}
