// +build windows

package winapi

//
// AccessSystemAcl access type
//
const ACCESS_SYSTEM_SECURITY = 0x01000000

//
// MaximumAllowed access type
//
const MAXIMUM_ALLOWED = 0x02000000

//
//  These are the generic rights.
//
const (
	GENERIC_READ    = 0x80000000
	GENERIC_WRITE   = 0x40000000
	GENERIC_EXECUTE = 0x20000000
	GENERIC_ALL     = 0x10000000
)

const (
	FILE_SHARE_READ   = 0x00000001
	FILE_SHARE_WRITE  = 0x00000002
	FILE_SHARE_DELETE = 0x00000004
)
