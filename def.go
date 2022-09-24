//go:build windows

package winapi

const INVALID_HANDLE_VALUE HANDLE = HANDLE(^uintptr(0))
