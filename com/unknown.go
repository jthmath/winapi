// +build windows

package com

// golang版的Unknown没有QueryInterface，因为不需要
type Unknown interface {
	AddRef() uint32
	Release() uint32
}
