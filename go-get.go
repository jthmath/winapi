// +build !windows

package winapi

import (
	"runtime"
)

func init() {
	str := `runtime.GOOS != "windows"`
	if runtime.GOOS != "windows" {
		panic(str)
	}
}
