// +build windows

package com

import (
	"syscall"
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
