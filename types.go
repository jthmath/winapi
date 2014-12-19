package winapi

type HANDLE uintptr
type HWND uintptr
type HMENU uintptr
type HINSTANCE uintptr
type HMODULE HINSTANCE

type HGDIOBJ uintptr
type HDC uintptr
type HICON uintptr
type HCURSOR uintptr
type HBRUSH uintptr
type HBITMAP uintptr

type SecurityAttributes struct {
	Length             uint32
	SecurityDescriptor uintptr
	InheritHandle      int32
}

type Overlapped struct {
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
