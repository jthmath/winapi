//go:build windows

package winapi

type Resource interface {
	uint16 | string
}
