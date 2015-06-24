// +build windows

package winapi

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	FORMAT_MESSAGE_IGNORE_INSERTS = 0x00000200
	FORMAT_MESSAGE_FROM_STRING    = 0x00000400
	FORMAT_MESSAGE_FROM_HMODULE   = 0x00000800
	FORMAT_MESSAGE_FROM_SYSTEM    = 0x00001000
	FORMAT_MESSAGE_ARGUMENT_ARRAY = 0x00002000
	FORMAT_MESSAGE_MAX_WIDTH_MASK = 0x000000FF
)

func FormatMessage(flags uint32, msgsrc interface{}, msgid uint32, langid uint32, args *byte) (string, error) {
	var b [300]uint16
	n, err := _FormatMessage(flags, msgsrc, msgid, langid, &b[0], 300, args)
	if err != nil {
		return "", err
	}
	for ; n > 0 && (b[n-1] == '\n' || b[n-1] == '\r'); n-- {
	}
	return string(utf16.Decode(b[:n])), nil
}

func _FormatMessage(flags uint32, msgsrc interface{}, msgid uint32, langid uint32, buf *uint16, nSize uint32, args *byte) (n uint32, err error) {
	r0, _, e1 := syscall.Syscall9(procFormatMessage.Addr(), 7,
		uintptr(flags), uintptr(0), uintptr(msgid), uintptr(langid),
		uintptr(unsafe.Pointer(buf)), uintptr(nSize),
		uintptr(unsafe.Pointer(args)), 0, 0)
	n = uint32(r0)
	if n == 0 {
		err = fmt.Errorf("winapi._FormatMessage error: %d", uint32(e1))
	}
	return
}

/*
typedef struct _SECURITY_ATTRIBUTES {
    DWORD  nLength;
    void   *pSecurityDescriptor;
    BOOL   bInheritHandle;
} SECURITY_ATTRIBUTES;
*/
type SECURITY_ATTRIBUTES struct {
	Length             uint32
	SecurityDescriptor uintptr
	InheritHandle      int32
}
