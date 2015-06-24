// +build windows

package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

const _MAX_READ = 1 << 30

const (
	CREATE_NEW        = 1
	CREATE_ALWAYS     = 2
	OPEN_EXISTING     = 3
	OPEN_ALWAYS       = 4
	TRUNCATE_EXISTING = 5
)

func _ReadFile(hFile HANDLE,
	buf *byte, NumberOfBytesToRead uint32,
	pNumberOfBytesRead *uint32,
	pOverlapped *OVERLAPPED) error {
	r1, _, e1 := syscall.Syscall6(procReadFile.Addr(), 5,
		uintptr(hFile), uintptr(unsafe.Pointer(buf)), uintptr(NumberOfBytesToRead),
		uintptr(unsafe.Pointer(pNumberOfBytesRead)), uintptr(unsafe.Pointer(pOverlapped)), 0)
	if r1 == 0 {
		wec := WinErrorCode(e1)
		if wec != 0 {
			return wec
		} else {
			return errors.New("winapi: ReadFile failed.")
		}
	} else {
		return nil
	}
}

func ReadFile(hFile HANDLE, buf []byte, pOverlapped *OVERLAPPED) (int, error) {
	if buf == nil {
		return 0, errors.New("ReadFile: 必须提供有效的缓冲区")
	}
	var Len int = len(buf)
	if Len <= 0 || Len > _MAX_READ {
		return 0, errors.New("ReadFile: 缓冲区长度必须大于零且不超过_MAX_READ")
	}
	var NumberOfBytesRead uint32 = 0
	err := _ReadFile(hFile, &buf[0], uint32(Len), &NumberOfBytesRead, pOverlapped)
	if err != nil {
		return 0, err
	} else {
		return int(NumberOfBytesRead), nil
	}
}
