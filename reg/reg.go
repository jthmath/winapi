package reg

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/jthmath/winapi"
)

var (
	dll_Advapi *syscall.LazyDLL = syscall.NewLazyDLL("Advapi32.dll")
)

var (
	procRegCreateKeyEx *syscall.LazyProc = dll_Advapi.NewProc("RegCreateKeyExW")
	procRegOpenKeyEx   *syscall.LazyProc = dll_Advapi.NewProc("RegRegOpenKeyExW")
)

type HKEY uintptr

type REGSAM winapi.ACCESS_MASK

const ERROR_SUCCESS int32 = 0

const (
	HKEY_CLASSES_ROOT   HKEY = 0x80000000
	HKEY_CURRENT_USER   HKEY = 0x80000001
	HKEY_LOCAL_MACHINE  HKEY = 0x80000002
	HKEY_USERS          HKEY = 0x80000003
	HKEY_CURRENT_CONFIG HKEY = 0x80000005
)

const (
	REG_OPTION_BACKUP_RESTORE uint32 = 0x00000004
	REG_OPTION_CREATE_LINK    uint32 = 0x00000002
	REG_OPTION_NON_VOLATILE   uint32 = 0x00000000
	REG_OPTION_VOLATILE       uint32 = 0x00000001
)

const (
	KEY_ALL_ACCESS REGSAM = 0xF003F
)

const (
	REG_CREATED_NEW_KEY     uint32 = 0x00000001
	REG_OPENED_EXISTING_KEY uint32 = 0x00000002
)

func CreateKey(hKey HKEY, SubKey string, Reserved uint32, Class string,
	Options uint32, samDesired REGSAM,
	sa *winapi.SecurityAttributes) (hkResult HKEY, Disposition uint32, err error) {
	pSubKey, err := syscall.UTF16PtrFromString(SubKey)
	if err != nil {
		return
	}
	pClass, err := syscall.UTF16PtrFromString(Class)
	if err != nil {
		return
	}
	var hKeyResult HKEY
	var disposition uint32
	r1, _, e1 := syscall.Syscall9(procRegCreateKeyEx.Addr(), 9,
		uintptr(hKey),
		uintptr(unsafe.Pointer(pSubKey)),
		uintptr(Reserved),
		uintptr(unsafe.Pointer(pClass)),
		uintptr(Options),
		uintptr(samDesired),
		uintptr(unsafe.Pointer(sa)),
		uintptr(unsafe.Pointer(&hKeyResult)),
		uintptr(unsafe.Pointer(&disposition)))
	n := int32(r1)
	if n != ERROR_SUCCESS {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = fmt.Errorf("函数 RegCreateKeyEx 返回 %d", n)
		}
	} else {
		hkResult = hKeyResult
		Disposition = disposition
	}

	return
}
