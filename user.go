package winapi

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	dll_user *syscall.LazyDLL = syscall.NewLazyDLL("user32.dll")

	procCreateWindow *syscall.LazyProc = dll_user.NewProc("CreateWindowExW")
	procMessageBox   *syscall.LazyProc = dll_user.NewProc("MessageBoxW")
)

const (
	MB_OK                uint32 = 0x00000000
	MB_OKCANCEL          uint32 = 0x00000001
	MB_ABORTRETRYIGNORE  uint32 = 0x00000002
	MB_YESNOCANCEL       uint32 = 0x00000003
	MB_YESNO             uint32 = 0x00000004
	MB_RETRYCANCEL       uint32 = 0x00000005
	MB_CANCELTRYCONTINUE uint32 = 0x00000006
	MB_HELP              uint32 = 0x00004000
)

func MessageBox(hWnd HWND, Text string, Caption string, Type uint32) (ret int32, err error) {
	pText, err := syscall.UTF16PtrFromString(Text)
	if err != nil {
		return
	}
	pCaption, err := syscall.UTF16PtrFromString(Caption)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.Syscall6(procMessageBox.Addr(), 4,
		uintptr(hWnd),
		uintptr(unsafe.Pointer(pText)),
		uintptr(unsafe.Pointer(pCaption)),
		uintptr(Type),
		0, 0)
	n := int32(r1)
	if n == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = fmt.Errorf("MessageBox -> %d", n)
		}
	} else {
		ret = n
	}
	return
}
