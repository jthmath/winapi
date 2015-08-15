// +build windows

package com

import (
	"syscall"
	"unsafe"

	"github.com/jthmath/winapi"
)

var (
	dll_ole32 = syscall.NewLazyDLL("ole32.dll")

	procCoInitialize     = dll_ole32.NewProc("CoInitialize")
	procCoCreateInstance = dll_ole32.NewProc("CoCreateInstance")
)

/*
var (
	procCoInitialize, _       = modole32.FindProc("CoInitialize")
	procCoInitializeEx, _     = modole32.FindProc("CoInitializeEx")
	procCoUninitialize, _     = modole32.FindProc("CoUninitialize")
	procCoCreateInstance, _   = modole32.FindProc("CoCreateInstance")
	procCoTaskMemFree, _      = modole32.FindProc("CoTaskMemFree")
	procCLSIDFromProgID, _    = modole32.FindProc("CLSIDFromProgID")
	procCLSIDFromString, _    = modole32.FindProc("CLSIDFromString")
	procStringFromCLSID, _    = modole32.FindProc("StringFromCLSID")
	procStringFromIID, _      = modole32.FindProc("StringFromIID")
	procIIDFromString, _      = modole32.FindProc("IIDFromString")
	procGetUserDefaultLCID, _ = modkernel32.FindProc("GetUserDefaultLCID")
	procCopyMemory, _         = modkernel32.FindProc("RtlMoveMemory")
	procVariantInit, _        = modoleaut32.FindProc("VariantInit")
	procVariantClear, _       = modoleaut32.FindProc("VariantClear")
	procSysAllocString, _     = modoleaut32.FindProc("SysAllocString")
	procSysAllocStringLen, _  = modoleaut32.FindProc("SysAllocStringLen")
	procSysFreeString, _      = modoleaut32.FindProc("SysFreeString")
	procSysStringLen, _       = modoleaut32.FindProc("SysStringLen")
	procCreateDispTypeInfo, _ = modoleaut32.FindProc("CreateDispTypeInfo")
	procCreateStdDispatch, _  = modoleaut32.FindProc("CreateStdDispatch")
	procGetActiveObject, _    = modoleaut32.FindProc("GetActiveObject")
)
*/

const (
	CLSCTX_INPROC_SERVER          = 0x1
	CLSCTX_INPROC_HANDLER         = 0x2
	CLSCTX_LOCAL_SERVER           = 0x4
	CLSCTX_INPROC_SERVER16        = 0x8
	CLSCTX_REMOTE_SERVER          = 0x10
	CLSCTX_INPROC_HANDLER16       = 0x20
	CLSCTX_RESERVED1              = 0x40
	CLSCTX_RESERVED2              = 0x80
	CLSCTX_RESERVED3              = 0x100
	CLSCTX_RESERVED4              = 0x200
	CLSCTX_NO_CODE_DOWNLOAD       = 0x400
	CLSCTX_RESERVED5              = 0x800
	CLSCTX_NO_CUSTOM_MARSHAL      = 0x1000
	CLSCTX_ENABLE_CODE_DOWNLOAD   = 0x2000
	CLSCTX_NO_FAILURE_LOG         = 0x4000
	CLSCTX_DISABLE_AAA            = 0x8000
	CLSCTX_ENABLE_AAA             = 0x10000
	CLSCTX_FROM_DEFAULT_CONTEXT   = 0x20000
	CLSCTX_ACTIVATE_32_BIT_SERVER = 0x40000
	CLSCTX_ACTIVATE_64_BIT_SERVER = 0x80000
	CLSCTX_ENABLE_CLOAKING        = 0x100000
	CLSCTX_APPCONTAINER           = 0x400000
	CLSCTX_ACTIVATE_AAA_AS_IU     = 0x800000
	CLSCTX_PS_DLL                 = 0x80000000
)

func CoCreateInstance(rclsid *winapi.GUID,
	u uintptr,
	ClsContext uint32,
	riid *winapi.GUID,
	ppv uintptr) error {
	r, _, _ := syscall.Syscall6(procCoCreateInstance.Addr(), 5,
		uintptr(unsafe.Pointer(rclsid)),
		uintptr(u),
		uintptr(ClsContext),
		uintptr(unsafe.Pointer(riid)),
		ppv,
		0)
	hr := winapi.HRESULT(r)
	if hr == 0 {
		return nil
	} else {
		return hr
	}
}
