// +build !windows

package com

func init() {
	panic(`runtime.GOOS != "windows"`)
}
