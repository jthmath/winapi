//go:build windows

package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

type WNDPROC func(HWND, uint32, uintptr, uintptr) uintptr

const (
	CS_VREDRAW         uint32 = 0x0001
	CS_HREDRAW         uint32 = 0x0002
	CS_DBLCLKS         uint32 = 0x0008
	CS_OWNDC           uint32 = 0x0020
	CS_CLASSDC         uint32 = 0x0040
	CS_PARENTDC        uint32 = 0x0080
	CS_NOCLOSE         uint32 = 0x0200
	CS_SAVEBITS        uint32 = 0x0800
	CS_BYTEALIGNCLIENT uint32 = 0x1000
	CS_BYTEALIGNWINDOW uint32 = 0x2000
	CS_GLOBALCLASS     uint32 = 0x4000
	CS_IME             uint32 = 0x00010000
	CS_DROPSHADOW      uint32 = 0x00020000
)

type WNDCLASS struct {
	Style         uint32
	PfnWndProc    WNDPROC
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	Menu          Resource
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

func RegisterClass(pWndClass *WNDCLASS) (cls uint16, err error) {
	if pWndClass == nil {
		err = errors.New("winapi: RegisterClass: pWndClass must not be nil")
		return
	}

	winStrClassName, err := syscall.UTF16PtrFromString(pWndClass.PszClassName)
	if err != nil {
		return
	}

	winStrMenuName, err := pWndClass.Menu.GetWinStr()
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
	wc.pszMenuName = winStrMenuName
	wc.pszClassName = winStrClassName
	wc.hIconSmall = pWndClass.HIconSmall

	r1, _, e1 := syscall.SyscallN(procRegisterClass.Addr(), 1, uintptr(unsafe.Pointer(&wc)))
	n := uint16(r1)
	if n == 0 {
		err = MakeFromWinError(e1)
		return
	}

	return n, nil
}

const (
	MB_OK                uint32 = 0x00000000
	MB_OKCANCEL          uint32 = 0x00000001
	MB_ABORTRETRYIGNORE  uint32 = 0x00000002
	MB_YESNOCANCEL       uint32 = 0x00000003
	MB_YESNO             uint32 = 0x00000004
	MB_RETRYCANCEL       uint32 = 0x00000005
	MB_CANCELTRYCONTINUE uint32 = 0x00000006
	MB_HELP              uint32 = 0x00004000

	MB_ICONERROR uint32 = 0x00000010
)

func MessageBox(hWnd HWND, text string, caption string, mbType uint32) (ret int32, err error) {
	pText, err := syscall.UTF16PtrFromString(text)
	if err != nil {
		return
	}
	pCaption, err := syscall.UTF16PtrFromString(caption)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.SyscallN(procMessageBox.Addr(),
		uintptr(hWnd),
		uintptr(unsafe.Pointer(pText)),
		uintptr(unsafe.Pointer(pCaption)),
		uintptr(mbType))
	n := int32(r1)
	if n == 0 {
		err = MakeFromWinError(e1)
		return
	}

	return n, nil
}

func DefWindowProc(hWnd HWND, message uint32, wParam uintptr, lParam uintptr) uintptr {
	ret, _, _ := syscall.SyscallN(procDefWindowProc.Addr(), uintptr(hWnd), uintptr(message), wParam, lParam)
	return ret
}

const CW_USEDEFAULT int32 = ^int32(0x7FFFFFFF) // 0x80000000

func CreateWindow(className string, windowName string, style uint32, exStyle uint32,
	x int32, y int32, width int32, height int32,
	wndParent HWND, menu HMENU, inst HINSTANCE, param uintptr) (hWnd HWND, err error) {
	pClassName, err := syscall.UTF16PtrFromString(className)
	if err != nil {
		return
	}
	pWindowName, err := syscall.UTF16PtrFromString(windowName)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.SyscallN(procCreateWindow.Addr(),
		uintptr(exStyle), uintptr(unsafe.Pointer(pClassName)), uintptr(unsafe.Pointer(pWindowName)), uintptr(style),
		uintptr(x), uintptr(y), uintptr(width), uintptr(height),
		uintptr(wndParent), uintptr(menu), uintptr(inst), uintptr(param))
	if r1 == 0 {
		err = MakeFromWinError(e1)
		return
	}

	return HWND(r1), nil
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

// ShowWindow
// Sets the specified window's show state.
//
// If the window was previously visible, the return value is nonzero.
// If the window was previously hidden, the return value is zero.
func ShowWindow(hWnd HWND, CmdShow int32) int32 {
	r1, _, _ := syscall.SyscallN(procShowWindow.Addr(), uintptr(hWnd), uintptr(CmdShow))
	return int32(r1)
}

func UpdateWindow(hWnd HWND) error {
	r1, _, _ := syscall.SyscallN(procUpdateWindow.Addr(), uintptr(hWnd))
	if r1 == 0 {
		return errors.New("winapi: UpdateWindow failed")
	} else {
		return nil
	}
}

func DestroyWindow(hWnd HWND) (err error) {
	r1, _, e1 := syscall.SyscallN(procDestroyWindow.Addr(), uintptr(hWnd))
	if n := int32(r1); n == 0 {
		return MakeFromWinError(e1)
	}
	return nil
}

func _LoadString(inst HINSTANCE, id uint16, buffer_ptr *uint16, buffer_max int32) (int32, error) {
	r1, _, e1 := syscall.SyscallN(procLoadString.Addr(),
		uintptr(inst), uintptr(id), uintptr(unsafe.Pointer(buffer_ptr)), uintptr(buffer_max))
	r := int32(r1)
	if r > 0 {
		return r, nil
	} else {
		return 0, MakeFromWinError(e1)
	}
}

func LoadString(hInstance HINSTANCE, id uint16) (string, error) {
	var err error
	var Len, Len1 int32
	var p *uint16 = nil
	Len, err = _LoadString(hInstance, id, (*uint16)(unsafe.Pointer(&p)), 0)

	if err == nil && Len > 0 {
		Buffer := make([]uint16, Len+1)
		Len1, err = _LoadString(hInstance, id, &Buffer[0], Len+1)
		if err == nil && Len == Len1 {
			Buffer[Len] = 0
			return syscall.UTF16ToString(Buffer), nil
		} else {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else {
		return "", errors.New("winapi: LoadString failed")
	}
}

func LoadBitmap[R ResourceConcept](hInst HINSTANCE, bigmapName R) (HBITMAP, error) {
	res := MakeResource(bigmapName)
	winStr, err := res.GetWinStr()
	if err != nil {
		return 0, err
	}

	r1, _, e1 := syscall.SyscallN(procLoadBitmap.Addr(), uintptr(hInst), uintptr(unsafe.Pointer(winStr)))
	if r1 != 0 {
		return HBITMAP(r1), nil
	} else {
		return 0, MakeFromWinError(e1)
	}
}

const (
	IDC_ARROW = 32512
)

func LoadCursor[R ResourceConcept](hinst HINSTANCE, cursorName R) (cursor HCURSOR, err error) {
	res := MakeResource(cursorName)
	winStr, err := res.GetWinStr()
	if err != nil {
		return 0, err
	}

	r1, _, e1 := syscall.SyscallN(procLoadCursor.Addr(), uintptr(hinst), uintptr(unsafe.Pointer(winStr)))
	if r1 != 0 {
		return HCURSOR(r1), nil
	} else {
		return 0, MakeFromWinError(e1)
	}
}

func LoadIcon[R ResourceConcept](hinst HINSTANCE, iconName R) (icon HICON, err error) {
	res := MakeResource(iconName)
	winStr, err := res.GetWinStr()
	if err != nil {
		return 0, err
	}

	r1, _, e1 := syscall.SyscallN(procLoadIcon.Addr(), uintptr(hinst), uintptr(unsafe.Pointer(winStr)))
	if r1 != 0 {
		return HICON(r1), nil
	} else {
		return 0, MakeFromWinError(e1)
	}
}

const (
	IMAGE_BITMAP = 0
	IMAGE_CURSOR = 2
	IMAGE_ICON   = 1
)

const (
	LR_CREATEDIBSECTION = 0x00002000
	LR_DEFAULTCOLOR     = 0x00000000
	LR_DEFAULTSIZE      = 0x00000040
	LR_LOADFROMFILE     = 0x00000010
	LR_LOADMAP3DCOLORS  = 0x00001000
	LR_LOADTRANSPARENT  = 0x00000020
	LR_MONOCHROME       = 0x00000001
	LR_SHARED           = 0x00008000
	LR_VGACOLOR         = 0x00000080
)

func LoadImage[R ResourceConcept](hinst HINSTANCE, name R, Type uint32,
	cxDesired int32, cyDesired int32, fLoad uint32) (hImage HANDLE, err error) {
	res := MakeResource(name)
	winStr, err := res.GetWinStr()
	if err != nil {
		return 0, err
	}

	r1, _, e1 := syscall.SyscallN(procLoadImage.Addr(),
		uintptr(hinst),
		uintptr(unsafe.Pointer(winStr)),
		uintptr(Type),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(fLoad))
	if r1 != 0 {
		return HANDLE(r1), nil
	} else {
		return 0, MakeFromWinError(e1)
	}
}

const (
	CTLCOLOR_MSGBOX HBRUSH = iota
	CTLCOLOR_EDIT
	CTLCOLOR_LISTBOX
	CTLCOLOR_BTN
	CTLCOLOR_DLG
	CTLCOLOR_SCROLLBAR
	CTLCOLOR_STATIC
	CTLCOLOR_MAX

	COLOR_SCROLLBAR           HBRUSH = 0
	COLOR_BACKGROUND          HBRUSH = 1
	COLOR_ACTIVECAPTION       HBRUSH = 2
	COLOR_INACTIVECAPTION     HBRUSH = 3
	COLOR_MENU                HBRUSH = 4
	COLOR_WINDOW              HBRUSH = 5
	COLOR_WINDOWFRAME         HBRUSH = 6
	COLOR_MENUTEXT            HBRUSH = 7
	COLOR_WINDOWTEXT          HBRUSH = 8
	COLOR_CAPTIONTEXT         HBRUSH = 9
	COLOR_ACTIVEBORDER        HBRUSH = 10
	COLOR_INACTIVEBORDER      HBRUSH = 11
	COLOR_APPWORKSPACE        HBRUSH = 12
	COLOR_HIGHLIGHT           HBRUSH = 13
	COLOR_HIGHLIGHTTEXT       HBRUSH = 14
	COLOR_BTNFACE             HBRUSH = 15
	COLOR_BTNSHADOW           HBRUSH = 16
	COLOR_GRAYTEXT            HBRUSH = 17
	COLOR_BTNTEXT             HBRUSH = 18
	COLOR_INACTIVECAPTIONTEXT HBRUSH = 19
	COLOR_BTNHIGHLIGHT        HBRUSH = 20

	COLOR_3DDKSHADOW HBRUSH = 21
	COLOR_3DLIGHT    HBRUSH = 22
	COLOR_INFOTEXT   HBRUSH = 23
	COLOR_INFOBK     HBRUSH = 24

	COLOR_HOTLIGHT                HBRUSH = 26
	COLOR_GRADIENTACTIVECAPTION   HBRUSH = 27
	COLOR_GRADIENTINACTIVECAPTION HBRUSH = 28
	COLOR_MENUHILIGHT             HBRUSH = 29
	COLOR_MENUBAR                 HBRUSH = 30

	COLOR_DESKTOP     = COLOR_BACKGROUND
	COLOR_3DFACE      = COLOR_BTNFACE
	COLOR_3DSHADOW    = COLOR_BTNSHADOW
	COLOR_3DHIGHLIGHT = COLOR_BTNHIGHLIGHT
	COLOR_3DHILIGHT   = COLOR_BTNHIGHLIGHT
	COLOR_BTNHILIGHT  = COLOR_BTNHIGHLIGHT
)

/*
 * Window Styles
 */
const (
	WS_OVERLAPPED   uint32 = 0x00000000
	WS_POPUP        uint32 = 0x80000000
	WS_CHILD        uint32 = 0x40000000
	WS_MINIMIZE     uint32 = 0x20000000
	WS_VISIBLE      uint32 = 0x10000000
	WS_DISABLED     uint32 = 0x08000000
	WS_CLIPSIBLINGS uint32 = 0x04000000
	WS_CLIPCHILDREN uint32 = 0x02000000
	WS_MAXIMIZE     uint32 = 0x01000000
	WS_CAPTION      uint32 = 0x00C00000 // WS_BORDER | WS_DLGFRAME
	WS_BORDER       uint32 = 0x00800000
	WS_DLGFRAME     uint32 = 0x00400000
	WS_VSCROLL      uint32 = 0x00200000
	WS_HSCROLL      uint32 = 0x00100000
	WS_SYSMENU      uint32 = 0x00080000
	WS_THICKFRAME   uint32 = 0x00040000
	WS_GROUP        uint32 = 0x00020000
	WS_TABSTOP      uint32 = 0x00010000

	WS_MINIMIZEBOX uint32 = 0x00020000
	WS_MAXIMIZEBOX uint32 = 0x00010000

	WS_TILED       uint32 = WS_OVERLAPPED
	WS_ICONIC      uint32 = WS_MINIMIZE
	WS_SIZEBOX     uint32 = WS_THICKFRAME
	WS_TILEDWINDOW uint32 = WS_OVERLAPPEDWINDOW

	/*
	 * Common Window Styles
	 */
	WS_OVERLAPPEDWINDOW uint32 = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU |
		WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX

	WS_POPUPWINDOW uint32 = WS_POPUP | WS_BORDER | WS_SYSMENU

	WS_CHILDWINDOW uint32 = WS_CHILD
)

/*
 * Extended Window Styles
 */
const (
	WS_EX_DLGMODALFRAME  uint32 = 0x00000001
	WS_EX_NOPARENTNOTIFY uint32 = 0x00000004
	WS_EX_TOPMOST        uint32 = 0x00000008
	WS_EX_ACCEPTFILES    uint32 = 0x00000010
	WS_EX_TRANSPARENT    uint32 = 0x00000020

	WS_EX_MDICHILD    uint32 = 0x00000040
	WS_EX_TOOLWINDOW  uint32 = 0x00000080
	WS_EX_WINDOWEDGE  uint32 = 0x00000100
	WS_EX_CLIENTEDGE  uint32 = 0x00000200
	WS_EX_CONTEXTHELP uint32 = 0x00000400

	WS_EX_RIGHT          uint32 = 0x00001000
	WS_EX_LEFT           uint32 = 0x00000000
	WS_EX_RTLREADING     uint32 = 0x00002000
	WS_EX_LTRREADING     uint32 = 0x00000000
	WS_EX_LEFTSCROLLBAR  uint32 = 0x00004000
	WS_EX_RIGHTSCROLLBAR uint32 = 0x00000000

	WS_EX_CONTROLPARENT uint32 = 0x00010000
	WS_EX_STATICEDGE    uint32 = 0x00020000
	WS_EX_APPWINDOW     uint32 = 0x00040000

	WS_EX_OVERLAPPEDWINDOW uint32 = WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE
	WS_EX_PALETTEWINDOW    uint32 = WS_EX_WINDOWEDGE | WS_EX_TOOLWINDOW | WS_EX_TOPMOST

	WS_EX_LAYERED uint32 = 0x00080000

	WS_EX_NOINHERITLAYOUT uint32 = 0x00100000 // Disable inheritence of mirroring by children

	WS_EX_NOREDIRECTIONBITMAP uint32 = 0x00200000

	WS_EX_LAYOUTRTL uint32 = 0x00400000 // Right to left mirroring

	WS_EX_COMPOSITED uint32 = 0x02000000
	WS_EX_NOACTIVATE uint32 = 0x08000000
)
