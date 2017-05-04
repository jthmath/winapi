package winapi

type HANDLE uintptr

type __HWND struct {
	unused int
}

type HWND uintptr

type HMENU uintptr
type HINSTANCE uintptr
type HMODULE uintptr

type HGDIOBJ uintptr
type HDC uintptr
type HICON uintptr
type HCURSOR uintptr
type HBRUSH uintptr
type HBITMAP uintptr

type OVERLAPPED struct {
	Internal     uintptr
	InternalHigh uintptr
	Offset       uint32
	OffsetHigh   uint32
	HEvent       HANDLE
}

type POINT struct {
	X int32
	Y int32
}

type ACCESS_MASK uint32

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}
