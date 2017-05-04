/*
错误相关，包括 Windows Error Code 和 HRESULT
*/

package winapi

import (
	"fmt"
)

// 1. Windows Error Code

type WinErrorCode uint32

func (this WinErrorCode) Error() string {
	var flags uint32 = FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_ARGUMENT_ARRAY | FORMAT_MESSAGE_IGNORE_INSERTS
	str, err := FormatMessage(flags, nil, uint32(this), 0, nil)
	n := uint32(this)
	if err == nil {
		return fmt.Sprintf("winapi error: %d(0x%08X) - ", n, n) + str
	} else {
		return fmt.Sprintf("winapi error: %d(0x%08X)", n, n)
	}
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
