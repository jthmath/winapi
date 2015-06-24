// +build windows

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

func CreateFile(FileName string, DesiredAccess uint32, ShareMode uint32,
	sa *SECURITY_ATTRIBUTES,
	CreationDisposition uint32, FlagsAndAttributes uint32,
	TemplateFile HANDLE) (HANDLE, error) {
	pFileName, err := syscall.UTF16PtrFromString(FileName)
	if err != nil {
		return 0, err
	}
	r1, _, e1 := syscall.Syscall9(procCreateFile.Addr(), 7,
		uintptr(unsafe.Pointer(pFileName)),
		uintptr(DesiredAccess),
		uintptr(ShareMode),
		uintptr(unsafe.Pointer(sa)),
		uintptr(CreationDisposition),
		uintptr(FlagsAndAttributes),
		uintptr(TemplateFile),
		0, 0)
	if r1 == 0 {
		wec := WinErrorCode(e1)
		if wec != 0 {
			return 0, wec
		} else {
			return 0, errors.New("CreateFile failed.")
		}
	} else {
		return HANDLE(r1), nil
	}
}

const (
	FILE_BEGIN   uint32 = 0
	FILE_CURRENT uint32 = 1
	FILE_END     uint32 = 2
)

func SetFilePointer(hFile HANDLE, DistanceToMove int64, MoveMethod uint32) (NewPointer int64, err error) {
	var np int64
	r1, _, e1 := syscall.Syscall6(procSetFilePointer.Addr(), 4,
		uintptr(hFile),
		uintptr(DistanceToMove),
		uintptr(unsafe.Pointer(&np)),
		uintptr(MoveMethod),
		0, 0)
	if r1 == 0 {
		wec := WinErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("SetFilePointer failed.")
		}
	} else {
		NewPointer = np
	}
	return
}
