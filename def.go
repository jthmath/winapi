//go:build windows

package winapi

const INVALID_HANDLE_VALUE HANDLE = HANDLE(^uintptr(0))

const NO_WND HWND = HWND(0)
