//go:build windows

package winapi

import (
	"testing"
)

func TestMessageBox(t *testing.T) {
	MessageBox(0, "fuck", "shit", MB_OK)
}
