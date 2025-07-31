//go:build windows

package winapi

import (
	"syscall"
	"time"
	"unsafe"
)

// Define the dwOpenMode values for CreateNamedPipe
const (
	PIPE_ACCESS_INBOUND  = 0x00000001
	PIPE_ACCESS_OUTBOUND = 0x00000002
	PIPE_ACCESS_DUPLEX   = 0x00000003
)

// Define the Named Pipe End flags for GetNamedPipeInfo
const (
	PIPE_CLIENT_END = 0x00000000
	PIPE_SERVER_END = 0x00000001
)

// Define the dwPipeMode values for CreateNamedPipe
const (
	PIPE_WAIT                  = 0x00000000
	PIPE_NOWAIT                = 0x00000001
	PIPE_READMODE_BYTE         = 0x00000000
	PIPE_READMODE_MESSAGE      = 0x00000002
	PIPE_TYPE_BYTE             = 0x00000000
	PIPE_TYPE_MESSAGE          = 0x00000004
	PIPE_ACCEPT_REMOTE_CLIENTS = 0x00000000
	PIPE_REJECT_REMOTE_CLIENTS = 0x00000008
)

// Define the well known values for CreateNamedPipe nMaxInstances
const PIPE_UNLIMITED_INSTANCES = 255

func ConnectNamedPipe(hNamedPipe HANDLE, po *OVERLAPPED) (err error) {
	r1, _, e1 := syscall.SyscallN(procConnectNamedPipe.Addr(), uintptr(hNamedPipe), uintptr(unsafe.Pointer(po)))
	if r1 == 0 {
		return MakeFromWinError(e1)
	}

	return nil
}

func CreateNamedPipe(
	name string,
	openMode uint32,
	pipeMode uint32,
	maxInstances uint32,
	outBufferSize uint32,
	inBufferSize uint32,
	defaultTimeOut time.Duration,
	sa *SECURITY_ATTRIBUTES) (h HANDLE, err error) {
	pName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return
	}

	dto := uint32(uint64(defaultTimeOut) / 1e6)

	h, err = _CreateNamedPipe(pName, openMode, pipeMode, maxInstances,
		outBufferSize, inBufferSize, dto, sa)
	return
}

func _CreateNamedPipe(pName *uint16, dwOpenMode uint32, dwPipeMode uint32,
	nMaxInstances uint32, nOutBufferSize uint32, nInBufferSize uint32,
	nDefaultTimeOut uint32, pSecurityAttributes *SECURITY_ATTRIBUTES) (h HANDLE, err error) {
	r1, _, e1 := syscall.SyscallN(procCreateNamedPipe.Addr(),
		uintptr(unsafe.Pointer(pName)),
		uintptr(dwOpenMode),
		uintptr(dwPipeMode),
		uintptr(nMaxInstances),
		uintptr(nOutBufferSize),
		uintptr(nInBufferSize),
		uintptr(nDefaultTimeOut),
		uintptr(unsafe.Pointer(pSecurityAttributes)))
	if h == INVALID_HANDLE_VALUE {
		err = MakeFromWinError(e1)
		return
	}

	h = HANDLE(r1)
	err = nil
	return
}
