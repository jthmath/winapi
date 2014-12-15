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
	procRegCreateKeyEx  *syscall.LazyProc = dll_Advapi.NewProc("RegCreateKeyExW")
	procRegOpenKeyEx    *syscall.LazyProc = dll_Advapi.NewProc("RegRegOpenKeyExW")
	procRegSetValueEx   *syscall.LazyProc = dll_Advapi.NewProc("RegSetValueExW")
	procRegCloseKey     *syscall.LazyProc = dll_Advapi.NewProc("RegCloseKey")
	procRegDeleteKeyEx  *syscall.LazyProc = dll_Advapi.NewProc("RegDeleteKeyExW")
	procRegQueryValueEx *syscall.LazyProc = dll_Advapi.NewProc("RegQueryValueExW")
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
	REG_UINT32                     uint32 = 4  // 32-bit number
	REG_LINK                       uint32 = 6  // Symbolic Link (unicode)
	REG_MULTI_SZ                   uint32 = 7  // Multiple Unicode strings
	REG_RESOURCE_LIST              uint32 = 8  // Resource list in the resource map
	REG_FULL_RESOURCE_DESCRIPTOR   uint32 = 9  // Resource list in the hardware description
	REG_RESOURCE_REQUIREMENTS_LIST uint32 = 10 //
	REG_UINT64                     uint32 = 11 // 64-bit number
)

const (
	KEY_WOW64_32KEY REGSAM = 0x0200
	KEY_WOW64_64KEY REGSAM = 0x0100
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
	case []byte:
		if Type != REG_BINARY {
			return TypeError
		} else {
			buf := Data.([]byte)
			pData = &buf[0]
			cbData = uint32(len(buf))
		}
	case uint32:
		if Type != REG_UINT32 {
			return TypeError
		} else {
			r := winapi.Uint32ToBinLittleEndian(Data.(uint32))
			pData = &r[0]
			cbData = 4
		}
	case uint64:
		if Type != REG_UINT64 {
			return TypeError
		} else {
			r := winapi.Uint64ToBinLittleEndian(Data.(uint64))
			pData = &r[0]
			cbData = 8
		}
	case string:
		if Type != REG_SZ {
			return TypeError
		} else {
			str := Data.(string)
			ustr, err := syscall.UTF16FromString(str)
			if err != nil {
				return err
			}
			pData = (*byte)(unsafe.Pointer(&ustr[0]))
			cbData = uint32(len(ustr)) * 2
		}
	case []string:
		if Type != REG_MULTI_SZ {
			return TypeError
		} else {
			str := Data.([]string)
			su := make([]uint16, 0)
			for i := 0; i < len(str); i++ {
				ustr, err := syscall.UTF16FromString(str[i])
				if err != nil {
					return err
				} else {
					su = append(su, ustr...)
				}
			}
			su = append(su, 0)
			pData = (*byte)(unsafe.Pointer(&su[0]))
			cbData = uint32(len(su)) * 2
		}
	default:
		return errors.New("SetValue不支持该类型")
	}
	return _SetValue(Key, pValueName, Reserved, Type, pData, cbData)
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
			err = errors.New("CloseKey failed.")
		}
	}
	return
}

func DeleteKey(Key HKEY, SubKey string, samDesired REGSAM, Reserved uint32) (err error) {
	pSubKey, err := syscall.UTF16PtrFromString(SubKey)
	if err != nil {
		return
	}
	r1, _, e1 := syscall.Syscall6(procRegSetValueEx.Addr(), 4,
		uintptr(Key),
		uintptr(unsafe.Pointer(pSubKey)),
		uintptr(samDesired),
		uintptr(Reserved),
		0, 0)
	n := int32(r1)
	if n != ERROR_SUCCESS {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = errors.New("DeleteKey failed.")
		}
	}
	return
}

func QueryValue(Key HKEY, ValueName string) (Type uint32, Data interface{}, err error) {
	pValueName, err := syscall.UTF16PtrFromString(ValueName)
	if err != nil {
		return
	}
	var dwType uint32
	var dwSize uint32
	err = _QueryValue(Key, pValueName, &dwType, nil, &dwSize)
	if err != nil {
		return
	}
	switch dwType {
	case REG_BINARY:
		sb := make([]byte, dwSize)
		err = _QueryValue(Key, pValueName, &dwType, &sb[0], &dwSize)
		if err != nil {
			return
		} else {
			Data = sb
		}
	case REG_SZ:
		if dwSize == 0 || dwSize%2 != 0 {
			err = errors.New("dwSize必须是正偶数")
			return
		}
		buf := make([]uint16, dwSize/2)
		err = _QueryValue(Key, pValueName, &dwType, (*byte)(unsafe.Pointer(&buf[0])), &dwSize)
		if err != nil {
			return
		} else {
			Data = syscall.UTF16ToString(buf)
		}
	case REG_MULTI_SZ:
		if dwSize == 0 || dwSize%2 != 0 {
			err = errors.New("dwSize必须是正偶数")
			return
		}
		buf := make([]uint16, dwSize/2)
		err = _QueryValue(Key, pValueName, &dwType, (*byte)(unsafe.Pointer(&buf[0])), &dwSize)
		if err != nil {
			return
		} else {
			ss, uerr := winapi.UTF16ToMultiString(buf)
			if uerr != nil {
				err = uerr
				return
			} else {
				Data = ss
			}
		}
	case REG_UINT32:
		var buf [4]byte
		err = _QueryValue(Key, pValueName, &dwType, &buf[0], &dwSize)
		if err != nil {
			return
		} else {
			Data = winapi.ByteArrayToUint32LittleEndian(buf)
		}
	case REG_UINT64:
		var buf [8]byte
		err = _QueryValue(Key, pValueName, &dwType, &buf[0], &dwSize)
		if err != nil {
			return
		} else {
			Data = winapi.ByteArrayToUint64LittleEndian(buf)
		}
	default:
		err = errors.New("不支持的类型")
		return
	}
	Type = dwType
	return
}

func _QueryValue(Key HKEY, ValueName *uint16, pType *uint32,
	pData *byte, pcbData *uint32) error {
	r1, _, _ := syscall.Syscall6(procRegQueryValueEx.Addr(), 6,
		uintptr(Key),
		uintptr(unsafe.Pointer(ValueName)),
		uintptr(0),
		uintptr(unsafe.Pointer(pType)),
		uintptr(unsafe.Pointer(pData)),
		uintptr(unsafe.Pointer(pcbData)))
	n := int32(r1)
	if n != ERROR_SUCCESS {
		return errors.New("_QueryValue failed.")
	} else {
		return nil
	}
}
