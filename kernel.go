package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

var dll_kernel *syscall.LazyDLL = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetLastError    *syscall.LazyProc = dll_kernel.NewProc("GetLastError")
	procExitProcess     *syscall.LazyProc = dll_kernel.NewProc("ExitProcess")
	procCreateFile      *syscall.LazyProc = dll_kernel.NewProc("CreateFileW")
	procReadFile        *syscall.LazyProc = dll_kernel.NewProc("ReadFile")
	procWriteFile       *syscall.LazyProc = dll_kernel.NewProc("WriteFile")
	procSetFilePointer  *syscall.LazyProc = dll_kernel.NewProc("SetFilePointerEx")
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

func GetModuleHandle(ModuleName string) (h HINSTANCE, err error) {
	pStr, err := SpecUTF16PtrFromString(ModuleName)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.Syscall(procGetModuleHandle.Addr(), 1,
		uintptr(unsafe.Pointer(pStr)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
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

func SetFilePointer(hFile HANDLE,
	DistanceToMove int64, MoveMethod uint32) (NewPointer int64, err error) {
	var np int64
	r1, _, e1 := syscall.Syscall6(procSetFilePointer.Addr(), 4,
		uintptr(hFile),
		uintptr(DistanceToMove),
		uintptr(unsafe.Pointer(&np)),
		uintptr(MoveMethod),
		0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("SetFilePointer -> 0")
		}
	} else {
		NewPointer = np
	}
	return
}
