//go:build windows

package winapi

import (
	"syscall"
)

var (
	dllGdi = syscall.NewLazyDLL("gdi32.dll")

	procBitBlt             = dllGdi.NewProc("BitBlt")
	procDeleteObject       = dllGdi.NewProc("DeleteObject")
	procGetObject          = dllGdi.NewProc("GetObject")
	procCreateCompatibleDC = dllGdi.NewProc("CreateCompatibleDC")
	procSelectObject       = dllGdi.NewProc("SelectObject")
	procDeleteDC           = dllGdi.NewProc("DeleteDC")
)

var (
	dllKernel = syscall.NewLazyDLL("kernel32.dll")

	procGetLastError     = dllKernel.NewProc("GetLastError")
	procExitProcess      = dllKernel.NewProc("ExitProcess")
	procCreateFile       = dllKernel.NewProc("CreateFileW")
	procReadFile         = dllKernel.NewProc("ReadFile")
	procWriteFile        = dllKernel.NewProc("WriteFile")
	procSetFilePointer   = dllKernel.NewProc("SetFilePointerEx")
	procGetModuleHandle  = dllKernel.NewProc("GetModuleHandleW")
	procCloseHandle      = dllKernel.NewProc("CloseHandle")
	procFormatMessage    = dllKernel.NewProc("FormatMessageW")
	procCreateNamedPipe  = dllKernel.NewProc("CreateNamedPipeW")
	procConnectNamedPipe = dllKernel.NewProc("ConnectNamedPipe")
)

var (
	dllUser = syscall.NewLazyDLL("user32.dll")

	procDefWindowProc         = dllUser.NewProc("DefWindowProcW")
	procGetMessage            = dllUser.NewProc("GetMessageW")
	procRegisterClass         = dllUser.NewProc("RegisterClassExW")
	procMessageBox            = dllUser.NewProc("MessageBoxW")
	procCreateWindow          = dllUser.NewProc("CreateWindowExW")
	procShowWindow            = dllUser.NewProc("ShowWindow")
	procUpdateWindow          = dllUser.NewProc("UpdateWindow")
	procTranslateMessage      = dllUser.NewProc("TranslateMessage")
	procDispatchMessage       = dllUser.NewProc("DispatchMessageW")
	procPostQuitMessage       = dllUser.NewProc("PostQuitMessage")
	procDestroyWindow         = dllUser.NewProc("DestroyWindow")
	procLoadString            = dllUser.NewProc("LoadStringW")
	procLoadIcon              = dllUser.NewProc("LoadIconW")
	procLoadCursor            = dllUser.NewProc("LoadCursorW")
	procLoadBitmap            = dllUser.NewProc("LoadBitmapW")
	procLoadImage             = dllUser.NewProc("LoadImageW")
	procBeginPaint            = dllUser.NewProc("BeginPaint")
	procEndPaint              = dllUser.NewProc("EndPaint")
	procRegisterWindowMessage = dllUser.NewProc("RegisterWindowMessageW")

	// menu
	procAppendMenu      = dllUser.NewProc("AppendMenuW")
	procCreateMenu      = dllUser.NewProc("CreateMenu")
	procCreatePopupMenu = dllUser.NewProc("CreatePopupMenu")
	procDestroyMenu     = dllUser.NewProc("DestroyMenu")
)
