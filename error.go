//go:build windows

package winapi

import (
	"fmt"
	"syscall"
)

// 1. Windows Error Code

type WinErrorCode uint32

func (ec WinErrorCode) Error() string {
	var flags uint32 = FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_ARGUMENT_ARRAY | FORMAT_MESSAGE_IGNORE_INSERTS
	str, err := FormatMessage(flags, nil, uint32(ec), 0, nil)
	n := uint32(ec)
	if err == nil {
		return fmt.Sprintf("winapi error: %d(0x%08X) - ", n, n) + str
	} else {
		return fmt.Sprintf("winapi error: %d(0x%08X)", n, n)
	}
}

func MakeFromWinError(e syscall.Errno) error {
	return WinErrorCode(e)
}

// 2. HRESULT

type HRESULT int32

const (
	S_OK    HRESULT = 0
	S_FALSE HRESULT = 1

	E_NOTIMPL HRESULT = (0x80004001 & 0x7FFFFFFF) | (^0x7FFFFFFF)
)

func (hr HRESULT) Succeeded() bool {
	return hr >= 0
}

func (hr HRESULT) Failed() bool {
	return hr < 0
}

func (hr HRESULT) Error() string {
	var flags uint32 = FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_ARGUMENT_ARRAY | FORMAT_MESSAGE_IGNORE_INSERTS
	str, err := FormatMessage(flags, nil, uint32(int32(hr)), 0, nil)
	if err == nil {
		return fmt.Sprintf("error: HRESULT = %d(0x%08X) - ", int32(hr), uint32(hr)) + str
	} else {
		return fmt.Sprintf("error: HRESULT = %d(0x%08X)", int32(hr), uint32(hr))
	}
}
