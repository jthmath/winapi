package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

func GetLastError() WinErrorCode {
	ret, _, _ := syscall.Syscall(procGetLastError.Addr(), 0, 0, 0, 0)
	return WinErrorCode(uint32(ret))
}

func ExitProcess(ExitCode uint32) {
	syscall.Syscall(procExitProcess.Addr(), 1, uintptr(ExitCode), 0, 0)
}

func GetModuleHandle(ModuleName string) (h HINSTANCE, err error) {
	var a uintptr
	var pStr *uint16

	if ModuleName == "" {
		a = 0
	} else {
		pStr, err = syscall.UTF16PtrFromString(ModuleName)
		if err != nil {
			return
		} else {
			a = uintptr(unsafe.Pointer(pStr))
		}
	}

	r1, _, e1 := syscall.Syscall(procGetModuleHandle.Addr(), 1, a, 0, 0)
	if r1 == 0 {
		wec := WinErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("GetModuleHandle failed.")
		}
	} else {
		h = HINSTANCE(r1)
	}
	return
}

func CloseHandle(h HANDLE) (err error) {
	r1, _, e1 := syscall.Syscall(procCloseHandle.Addr(), 1, uintptr(h), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("CloseHandle failed.")
		}
	}
	return
}
