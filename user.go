package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	dll_user *syscall.LazyDLL = syscall.NewLazyDLL("user32.dll")

	procDefWindowProc    *syscall.LazyProc = dll_user.NewProc("DefWindowProcW")
	procGetMessage       *syscall.LazyProc = dll_user.NewProc("GetMessageW")
	procRegisterClass    *syscall.LazyProc = dll_user.NewProc("RegisterClassExW")
	procMessageBox       *syscall.LazyProc = dll_user.NewProc("MessageBoxW")
	procCreateWindow     *syscall.LazyProc = dll_user.NewProc("CreateWindowExW")
	procShowWindow       *syscall.LazyProc = dll_user.NewProc("ShowWindow")
	procUpdateWindow     *syscall.LazyProc = dll_user.NewProc("UpdateWindow")
	procTranslateMessage *syscall.LazyProc = dll_user.NewProc("TranslateMessage")
	procDispatchMessage  *syscall.LazyProc = dll_user.NewProc("DispatchMessageW")
	procPostQuitMessage  *syscall.LazyProc = dll_user.NewProc("PostQuitMessage")
	procDestroyWindow    *syscall.LazyProc = dll_user.NewProc("DestroyWindow")
)

const (
	WM_NULL    uint32 = 0x0000
	WM_CREATE  uint32 = 0x0001
	WM_DESTROY uint32 = 0x0002
	WM_MOVE    uint32 = 0x0003
	WM_SIZE    uint32 = 0x0005

	WM_ACTIVATE uint32 = 0x0006

	/*
	 * WM_ACTIVATE state values
	 */
	WA_INACTIVE    = 0
	WA_ACTIVE      = 1
	WA_CLICKACTIVE = 2

	WM_CLOSE uint32 = 0x0010
	WM_QUIT  uint32 = 0x0012

	WM_GETMINMAXINFO uint32 = 0x0024
)

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type WNDPROC func(HWND, uint32, uintptr, uintptr) uintptr

type WNDCLASS struct {
	Style         uint32
	PfnWndProc    WNDPROC
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	PszMenuName   string
	PszClassName  string
	HIconSmall    HICON
}

type _WNDCLASS struct {
	cbSize        uint32
	style         uint32
	pfnWndProcPtr uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     HINSTANCE
	hIcon         HICON
	hCursor       HCURSOR
	hbrBackground HBRUSH
	pszMenuName   *uint16
	pszClassName  *uint16
	hIconSmall    HICON
}

func newWndProc(proc WNDPROC) uintptr {
	return syscall.NewCallback(proc)
}

func RegisterClass(pWndClass *WNDCLASS) (atom uint16, err error) {
	if pWndClass == nil {
		return 0, error(syscall.EINVAL)
	}
	_pMenuName, err := syscall.UTF16PtrFromString(pWndClass.PszMenuName)
	if err != nil {
		return
	}
	_pClassName, err := syscall.UTF16PtrFromString(pWndClass.PszClassName)
	if err != nil {
		return
	}
	var wc _WNDCLASS
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	wc.style = pWndClass.Style
	wc.pfnWndProcPtr = newWndProc(pWndClass.PfnWndProc)
	wc.cbClsExtra = pWndClass.CbClsExtra
	wc.cbWndExtra = pWndClass.CbWndExtra
	wc.hInstance = pWndClass.HInstance
	wc.hIcon = pWndClass.HIcon
	wc.hCursor = pWndClass.HCursor
	wc.hbrBackground = pWndClass.HbrBackground
	wc.pszMenuName = _pMenuName
	wc.pszClassName = _pClassName
	wc.hIconSmall = pWndClass.HIconSmall
	r1, _, e1 := syscall.Syscall(procRegisterClass.Addr(), 1, uintptr(unsafe.Pointer(&wc)), 0, 0)
	n := uint16(r1)
	if n != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("RegisterClass failed.")
		}
	} else {
		atom = n
	}
	return
}

// (NOTE): for MessageBox function
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
			err = errors.New("MessageBox failed.")
		}
	} else {
		ret = n
	}
	return
}

func DefWindowProc(hWnd HWND, message uint32, wParam uintptr, lParam uintptr) uintptr {
	ret, _, _ := syscall.Syscall6(procDefWindowProc.Addr(), 4, uintptr(hWnd), uintptr(message), wParam, lParam, 0, 0)
	return ret
}

const CW_USEDEFAULT int32 = ^int32(0x7FFFFFFF) // 0x80000000

func CreateWindow(ClassName string, WindowName string, Style uint32, ExStyle uint32,
	X int32, Y int32, Width int32, Height int32,
	WndParent HWND, Menu HMENU, inst HINSTANCE, Param uintptr) (hWnd HWND, err error) {
	pClassName, err := syscall.UTF16PtrFromString(ClassName)
	if err != nil {
		return
	}
	pWindowName, err := syscall.UTF16PtrFromString(WindowName)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.Syscall12(procCreateWindow.Addr(), 12,
		uintptr(ExStyle), uintptr(unsafe.Pointer(pClassName)), uintptr(unsafe.Pointer(pWindowName)), uintptr(Style),
		uintptr(X), uintptr(Y), uintptr(Width), uintptr(Height),
		uintptr(WndParent), uintptr(Menu), uintptr(inst), uintptr(Param))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("CreateWindow failed.")
		}
	} else {
		hWnd = HWND(r1)
	}
	return
}

const (
	SW_HIDE            int32 = 0
	SW_SHOWNORMAL      int32 = 1
	SW_NORMAL          int32 = 1
	SW_SHOWMINIMIZED   int32 = 2
	SW_SHOWMAXIMIZED   int32 = 3
	SW_MAXIMIZE        int32 = 3
	SW_SHOWNOACTIVATE  int32 = 4
	SW_SHOW            int32 = 5
	SW_MINIMIZE        int32 = 6
	SW_SHOWMINNOACTIVE int32 = 7
	SW_SHOWNA          int32 = 8
	SW_RESTORE         int32 = 9
	SW_SHOWDEFAULT     int32 = 10
	SW_FORCEMINIMIZE   int32 = 11
	SW_MAX             int32 = 11
)

func ShowWindow(hWnd HWND, CmdShow int32) error {
	r1, _, _ := syscall.Syscall(procShowWindow.Addr(), 2, uintptr(hWnd), uintptr(CmdShow), 0)
	if r1 == 0 {
		return errors.New("CreateWindow failed.")
	} else {
		return nil
	}
}

func UpdateWindow(hWnd HWND) error {
	r1, _, _ := syscall.Syscall(procUpdateWindow.Addr(), 1, uintptr(hWnd), 0, 0)
	if r1 == 0 {
		return errors.New("UpdateWindow failed.")
	} else {
		return nil
	}
}

func TranslateMessage(pMsg *MSG) error {
	r1, _, _ := syscall.Syscall(procTranslateMessage.Addr(), 1, uintptr(unsafe.Pointer(pMsg)), 0, 0)
	if r1 == 0 {
		return errors.New("TranslateMessage failed.")
	} else {
		return nil
	}
}

func DispatchMessage(pMsg *MSG) uintptr {
	r1, _, _ := syscall.Syscall(procDispatchMessage.Addr(), 1, uintptr(unsafe.Pointer(pMsg)), 0, 0)
	return r1
}

func PostQuitMessage(ExitCode int32) {
	syscall.Syscall(procPostQuitMessage.Addr(), 1, uintptr(ExitCode), 0, 0)
}

func DestroyWindow(hWnd HWND) (err error) {
	r1, _, e1 := syscall.Syscall(procDestroyWindow.Addr(), 1, uintptr(hWnd), 0, 0)
	if n := int32(r1); n == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("DestroyWindow failed.")
		}
	}
	return
}
