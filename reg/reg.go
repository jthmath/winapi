package reg

import (
	"errors"
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
	procRegSetValueEx  *syscall.LazyProc = dll_Advapi.NewProc("RegSetValueExW")
	procRegCloseKey    *syscall.LazyProc = dll_Advapi.NewProc("RegCloseKey")
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

const (
	REG_NONE      uint32 = 0 // No value type
	REG_SZ        uint32 = 1 // Unicode nul terminated string
	REG_EXPAND_SZ uint32 = 2 // Unicode nul terminated string(with environment variable references)

	REG_BINARY                     uint32 = 3  // Free form binary
	REG_DWORD                      uint32 = 4  // 32-bit number
	REG_DWORD_LITTLE_ENDIAN        uint32 = 4  // 32-bit number (same as REG_DWORD)
	REG_DWORD_BIG_ENDIAN           uint32 = 5  // 32-bit number
	REG_LINK                       uint32 = 6  // Symbolic Link (unicode)
	REG_MULTI_SZ                   uint32 = 7  // Multiple Unicode strings
	REG_RESOURCE_LIST              uint32 = 8  // Resource list in the resource map
	REG_FULL_RESOURCE_DESCRIPTOR   uint32 = 9  // Resource list in the hardware description
	REG_RESOURCE_REQUIREMENTS_LIST uint32 = 10 //
	REG_QWORD                      uint32 = 11 // 64-bit number
	REG_QWORD_LITTLE_ENDIAN        uint32 = 11 // 64-bit number (same as REG_QWORD)
)

func CreateKey(hKey HKEY, SubKey string, Reserved uint32, Class string,
	Options uint32, samDesired REGSAM,
	sa *winapi.SecurityAttributes) (hkResult HKEY, Disposition uint32, err error) {
	pSubKey, err := syscall.UTF16PtrFromString(SubKey)
	if err != nil {
		return
	}
	pClass, err := winapi.SpecUTF16PtrFromString(Class)
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

func SetValue(Key HKEY, ValueName string, Reserved uint32,
	Type uint32, Data interface{}) error {
	pValueName, err := syscall.UTF16PtrFromString(ValueName)
	if err != nil {
		return err
	}
	TypeError := errors.New("函数 SetValue ，Type与实际的数据类型不匹配")
	var pData *byte
	var cbData uint32
	switch Data.(type) {
	case uint32:
		if Type != REG_DWORD {
			return TypeError
		} else {
			r := winapi.Uint32ToBinLittleEndian(Data.(uint32))
			pData = &r[0]
			cbData = 4
		}
	case uint64:
		if Type != REG_QWORD {
			return TypeError
		} else {
			r := winapi.Uint64ToBinLittleEndian(Data.(uint64))
			pData = &r[0]
			cbData = 8
		}
	case string:
		if Type != REG_SZ {
			return TypeError
		}
	case []string:
		if Type != REG_MULTI_SZ {
			return TypeError
		}
	default:
		return errors.New("SetValue不支持该类型")
	}
	err = _SetValue(Key, pValueName, Reserved, Type, pData, cbData)
	return err
}

func _SetValue(Key HKEY, ValueName *uint16, Reserved uint32,
	Type uint32, Data *byte, cbData uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procRegSetValueEx.Addr(), 6,
		uintptr(Key),
		uintptr(unsafe.Pointer(ValueName)),
		uintptr(Reserved),
		uintptr(Type),
		uintptr(unsafe.Pointer(Data)),
		uintptr(cbData))
	n := int32(r1)
	if n != ERROR_SUCCESS {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = fmt.Errorf("函数 RegSetValueEx 返回 %d", n)
		}
	}
	return
}

func CloseKey(Key HKEY) (err error) {
	r1, _, e1 := syscall.Syscall(procRegCloseKey.Addr(), 1,
		uintptr(unsafe.Pointer(Key)), 0, 0)
	if int32(r1) != ERROR_SUCCESS {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("failed to CloseKey")
		}
	}
	return
}
