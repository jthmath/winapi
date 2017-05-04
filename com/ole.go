package com

import (
	"syscall"

	"github.com/jthmath/winapi"
)

func CoInitialize(a interface{}) error {
	r1, _, _ := syscall.Syscall(procCoInitialize.Addr(), 1, 0, 0, 0)
	hr := winapi.HRESULT(r1)
	if r1 == 0 {
		return nil
	} else {
		return hr
	}
}
