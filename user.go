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
	procLoadString       *syscall.LazyProc = dll_user.NewProc("LoadStringW")
	procLoadIcon         *syscall.LazyProc = dll_user.NewProc("LoadIconW")
	procLoadImage        *syscall.LazyProc = dll_user.NewProc("LoadImageW")
)

const (
	WM_NULL    uint32 = 0x0000
	WM_CREATE  uint32 = 0x0001
	WM_DESTROY uint32 = 0x0002
	WM_MOVE    uint32 = 0x0003
	WM_SIZE    uint32 = 0x0005

	WM_ACTIVATE uint32 = 0x0006

	WM_PAINT uint32 = 0x000F

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
	if n == 0 {
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

	MB_ICONERROR uint32 = 0x00000010
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

func ErrorBox(err error) error {
	var e error
	if err == nil {
		_, e = MessageBox(0, "<nil>", "error", MB_OK)
	} else {
		_, e = MessageBox(0, err.Error(), "error", MB_OK|MB_ICONERROR)
	}
	return e
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

// 返回值：如果窗口事先是可见的，返回true
//       如果窗口事先是隐藏的，返回false
func ShowWindow(hWnd HWND, CmdShow int32) bool {
	r1, _, _ := syscall.Syscall(procShowWindow.Addr(), 2, uintptr(hWnd), uintptr(CmdShow), 0)
	return r1 != 0
}

func UpdateWindow(hWnd HWND) error {
	r1, _, _ := syscall.Syscall(procUpdateWindow.Addr(), 1, uintptr(hWnd), 0, 0)
	if r1 == 0 {
		return errors.New("UpdateWindow failed.") // 该函数没有对应的GetLastError值
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

/*
int WINAPI LoadString(
  _In_opt_  HINSTANCE hInstance,
  _In_      UINT uID,
  _Out_     LPTSTR lpBuffer,
  _In_      int nBufferMax
);
*/

func LoadString(hInstance HINSTANCE, id uint32, BufferMax int) (str string, err error) {
	CommonErrorString := "winapi.LoadString: "
	if BufferMax <= 0 {
		err = errors.New(CommonErrorString + "BufferMax <= 0.")
		return
	}

	buf := make([]uint16, BufferMax)
	r1, _, e1 := syscall.Syscall6(procLoadString.Addr(), 4,
		uintptr(hInstance),
		uintptr(id),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(BufferMax),
		0, 0)
	n := int32(r1)
	if n <= 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New(CommonErrorString + "syscall failed.")
		}
	} else {
		str = syscall.UTF16ToString(buf)
	}

	return
}

func LoadIconById(hinst HINSTANCE, id uint16) (icon HICON, err error) {
	r1, _, e1 := syscall.Syscall(procLoadIcon.Addr(), 2,
		uintptr(hinst), MakeIntResource(id), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("LoadIconById failed.")
		}
	} else {
		icon = HICON(r1)
	}
	return
}

func LoadIconByName(hinst HINSTANCE, name string) (icon HICON, err error) {
	pName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return
	}

	r1, _, e1 := syscall.Syscall(procLoadIcon.Addr(), 2,
		uintptr(hinst), uintptr(unsafe.Pointer(pName)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("LoadIconByName failed.")
		}
	} else {
		icon = HICON(r1)
	}
	return
}

const ( // LoadImage函数的uType参数
	IMAGE_BITMAP = 0
	IMAGE_CURSOR = 2
	IMAGE_ICON   = 1
)

const ( // LoadImage函数的fuLoad参数
	LR_CREATEDIBSECTION uint32 = 0x00002000
	LR_DEFAULTCOLOR     uint32 = 0x00000000
	LR_DEFAULTSIZE      uint32 = 0x00000040
	LR_LOADFROMFILE     uint32 = 0x00000010
	LR_LOADMAP3DCOLORS  uint32 = 0x00001000
	LR_LOADTRANSPARENT  uint32 = 0x00000020
	LR_MONOCHROME       uint32 = 0x00000001
	LR_SHARED           uint32 = 0x00008000
	LR_VGACOLOR         uint32 = 0x00000080
)

/*
HANDLE WINAPI LoadImage(
  _In_opt_  HINSTANCE hinst,
  _In_      LPCTSTR lpszName,
  _In_      UINT uType,
  _In_      int cxDesired,
  _In_      int cyDesired,
  _In_      UINT fuLoad
);
*/

func LoadImageById(hinst HINSTANCE, id uint16, Type uint32,
	cxDesired int32, cyDesired int32, fLoad uint32) (hImage HANDLE, err error) {
	r1, _, e1 := syscall.Syscall6(procLoadImage.Addr(), 6,
		uintptr(hinst),
		MakeIntResource(id),
		uintptr(Type),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(fLoad))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("LoadImageById failed.")
		}
	} else {
		hImage = HANDLE(r1)
	}
	return
}

func LoadImageByName(hinst HINSTANCE, name string, Type uint32,
	cxDesired int32, cyDesired int32, fLoad uint32) (hImage HANDLE, err error) {
	pName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.Syscall6(procLoadImage.Addr(), 6,
		uintptr(hinst),
		uintptr(unsafe.Pointer(pName)),
		uintptr(Type),
		uintptr(cxDesired),
		uintptr(cyDesired),
		uintptr(fLoad))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("LoadImageByName failed.")
		}
	} else {
		hImage = HANDLE(r1)
	}
	return
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

	COLOR_HOTLIGHT                HBRUSH = 26 // 上一个是24，所以这里不能直接用iota
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
	WS_CAPTION      uint32 = 0x00C00000 /* WS_BORDER | WS_DLGFRAME  */
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

func GetMessage(pMsg *MSG, hWnd HWND, wMsgFilterMin uint32, wMsgFilterMax uint32) int32 {
	r1, _, _ := syscall.Syscall6(procGetMessage.Addr(), 4,
		uintptr(unsafe.Pointer(pMsg)),
		uintptr(hWnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax),
		0, 0)
	return int32(r1)
}
