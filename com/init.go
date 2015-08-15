// +build windows

package com

import (
	"fmt"
)

func init() {
	if err := CoInitialize(nil); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
