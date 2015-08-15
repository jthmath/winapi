// +build windows

package winapi

import (
	"errors"
	"fmt"
)

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]uint8
}

func (this GUID) String() string {
	FormatString := "%08X-%04X-%04X-%04X-%02X%02X%02X%02X%02X%02X"
	a := this.Data4[0]
	b := this.Data4[1]
	u := uint16(a)<<8 | uint16(b)
	return fmt.Sprintf(FormatString, this.Data1, this.Data2, this.Data3, u,
		this.Data4[2], this.Data4[3], this.Data4[4],
		this.Data4[5], this.Data4[6], this.Data4[7])
}

/*
56FDF344-FD6D-11d0-958A-006097C9A090

var CLSID_TaskbarList = GUID{
	Data1: 1459483460,
	Data2: 64877,
	Data3: 4560,
	Data4: [8]uint8{149, 138, 0, 96, 151, 201, 160, 144},
}
*/

func MakeGuid(str string) (r GUID, err error) {
	var n int
	FormatString := "%08X-%04X-%04X-%04X-%02X%02X%02X%02X%02X%02X"
	var u uint16
	var guid GUID
	n, err = fmt.Sscanf(str, FormatString, &guid.Data1, &guid.Data2, &guid.Data3, &u,
		&guid.Data4[2], &guid.Data4[3], &guid.Data4[4],
		&guid.Data4[5], &guid.Data4[6], &guid.Data4[7])

	if err == nil {
		if n != 10 {
			err = errors.New("")
		} else {
			guid.Data4[0] = uint8(u >> 8)
			guid.Data4[1] = uint8(u)
			r = guid
		}
	}
	return
}

func MustMakeGuid(str string) GUID {
	guid, err := MakeGuid(str)
	if err == nil {
		return guid
	} else {
		panic(err)
	}
}
