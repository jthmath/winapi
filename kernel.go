package winapi

import (
	"errors"
	"syscall"
)

var dll_kernel *syscall.LazyDLL = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetLastError    *syscall.LazyProc = dll_kernel.NewProc("GetLastError")
	procExitProcess     *syscall.LazyProc = dll_kernel.NewProc("ExitProcess")
	procCreateFile      *syscall.LazyProc = dll_kernel.NewProc("CreateFileW")
	procReadFile        *syscall.LazyProc = dll_kernel.NewProc("ReadFile")
	procWriteFile       *syscall.LazyProc = dll_kernel.NewProc("WriteFile")
	procSetFilePointer  *syscall.LazyProc = dll_kernel.NewProc("SetFilePointer")
	procGetModuleHandle *syscall.LazyProc = dll_kernel.NewProc("GetModuleHandleW")
	procCloseHandle     *syscall.LazyProc = dll_kernel.NewProc("CloseHandle")
)

func GetLastError() uint32 {
	ret, _, _ := syscall.Syscall(procGetLastError.Addr(), 0, 0, 0, 0)
	return uint32(ret)
}

func ExitProcess(ExitCode uint32) {
	syscall.Syscall(procExitProcess.Addr(), 1, uintptr(ExitCode), 0, 0)
}

const (
	FILE_BEGIN   uint32 = 0
	FILE_CURRENT uint32 = 1
	FILE_END     uint32 = 2
)

func GetModuleHandle(ModuleName string) HMODULE {
	var r1 uintptr
	if ModuleName == "" {
		r1, _, _ = syscall.Syscall(procGetModuleHandle.Addr(), 1, 0, 0, 0)
	}
	return HMODULE(r1)
}

func CloseHandle(h HANDLE) (err error) {
	r1, _, e1 := syscall.Syscall(procCloseHandle.Addr(), 1, uintptr(h), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("CloseHandle -> 0")
		}
	}
	return
}

func SetFilePointer(hFile HANDLE, n int64) (NewPointer int64, err error) {
	return 1, nil
}
