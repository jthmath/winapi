//go:build windows

package winapi

import (
	"syscall"
)

var (
	dll_gdi = syscall.NewLazyDLL("gdi32.dll")

	procBitBlt             = dll_gdi.NewProc("BitBlt")
	procDeleteObject       = dll_gdi.NewProc("DeleteObject")
	procGetObject          = dll_gdi.NewProc("GetObject")
	procCreateCompatibleDC = dll_gdi.NewProc("CreateCompatibleDC")
	procSelectObject       = dll_gdi.NewProc("SelectObject")
	procDeleteDC           = dll_gdi.NewProc("DeleteDC")
)

var (
	dll_kernel = syscall.NewLazyDLL("kernel32.dll")

	procGetLastError     = dll_kernel.NewProc("GetLastError")
	procExitProcess      = dll_kernel.NewProc("ExitProcess")
	procCreateFile       = dll_kernel.NewProc("CreateFileW")
	procReadFile         = dll_kernel.NewProc("ReadFile")
	procWriteFile        = dll_kernel.NewProc("WriteFile")
	procSetFilePointer   = dll_kernel.NewProc("SetFilePointerEx")
	procGetModuleHandle  = dll_kernel.NewProc("GetModuleHandleW")
	procCloseHandle      = dll_kernel.NewProc("CloseHandle")
	procFormatMessage    = dll_kernel.NewProc("FormatMessageW")
	procCreateNamedPipe  = dll_kernel.NewProc("CreateNamedPipeW")
	procConnectNamedPipe = dll_kernel.NewProc("ConnectNamedPipe")
)

var (
	dll_user = syscall.NewLazyDLL("user32.dll")

	procDefWindowProc         = dll_user.NewProc("DefWindowProcW")
	procGetMessage            = dll_user.NewProc("GetMessageW")
	procRegisterClass         = dll_user.NewProc("RegisterClassExW")
	procMessageBox            = dll_user.NewProc("MessageBoxW")
	procCreateWindow          = dll_user.NewProc("CreateWindowExW")
	procShowWindow            = dll_user.NewProc("ShowWindow")
	procUpdateWindow          = dll_user.NewProc("UpdateWindow")
	procTranslateMessage      = dll_user.NewProc("TranslateMessage")
	procDispatchMessage       = dll_user.NewProc("DispatchMessageW")
	procPostQuitMessage       = dll_user.NewProc("PostQuitMessage")
	procDestroyWindow         = dll_user.NewProc("DestroyWindow")
	procLoadString            = dll_user.NewProc("LoadStringW")
	procLoadIcon              = dll_user.NewProc("LoadIconW")
	procLoadCursor            = dll_user.NewProc("LoadCursorW")
	procLoadBitmap            = dll_user.NewProc("LoadBitmapW")
	procLoadImage             = dll_user.NewProc("LoadImageW")
	procBeginPaint            = dll_user.NewProc("BeginPaint")
	procEndPaint              = dll_user.NewProc("EndPaint")
	procRegisterWindowMessage = dll_user.NewProc("RegisterWindowMessageW")

	// 菜单
	procAppendMenu      = dll_user.NewProc("AppendMenuW")
	procCreateMenu      = dll_user.NewProc("CreateMenu")
	procCreatePopupMenu = dll_user.NewProc("CreatePopupMenu")
	procDestroyMenu     = dll_user.NewProc("DestroyMenu")
)

var (
	dll_comdlg = syscall.NewLazyDLL("comdlg32.dll")

	procGetSaveFileName      = dll_comdlg.NewProc("GetSaveFileNameW")
	procGetOpenFileName      = dll_comdlg.NewProc("GetOpenFileNameW")
	procCommDlgExtendedError = dll_comdlg.NewProc("CommDlgExtendedError")
)
