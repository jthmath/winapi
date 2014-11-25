package winapi

type HANDLE uintptr
type HWND uintptr
type HKEY uintptr
type HMENU uintptr
type HINSTANCE uintptr
type HMODULE HINSTANCE
type HICON uintptr
type HCURSOR uintptr
type HBRUSH uintptr

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
	x int32
	y int32
}
