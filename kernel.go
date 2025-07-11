//go:build windows

package winapi

import (
	"syscall"
	"unsafe"
)

func GetLastError() uint32 {
	ret, _, _ := syscall.SyscallN(procGetLastError.Addr())
	return uint32(ret)
}

func ExitProcess(ExitCode uint32) {
	syscall.SyscallN(procExitProcess.Addr(), uintptr(ExitCode))
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

	r1, _, e1 := syscall.SyscallN(procGetModuleHandle.Addr(), a)
	if r1 == 0 {
		err = MakeFromWinError(e1)
		return
	}

	h = HINSTANCE(r1)
	return
}

func CloseHandle(h HANDLE) error {
	r1, _, e1 := syscall.SyscallN(procCloseHandle.Addr(), uintptr(h))
	if r1 == 0 {
		return MakeFromWinError(e1)
	}

	return nil
}
