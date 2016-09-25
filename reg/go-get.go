// +build !windows

package reg

func init() {
	panic(`runtime.GOOS != "windows"`)
}
