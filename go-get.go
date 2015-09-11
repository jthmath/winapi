// +build !windows

package winapi

import (
	"fmt"
	"runtime"
)

func init() {
	str := `runtime.GOOS != "windows"`
	if runtime.GOOS != "windows" {
		fmt.Println(str)
		panic(str)
	}
}
