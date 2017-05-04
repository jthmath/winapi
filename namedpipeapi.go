package winapi

import (
	"errors"
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
	r1, _, e1 := syscall.Syscall(procConnectNamedPipe.Addr(), 2,
		uintptr(hNamedPipe), uintptr(unsafe.Pointer(po)), 0)
	if r1 == 0 {
		wec := WinErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("ConnectNamedPipe failed.")
		}
	}
	return
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
	r1, _, e1 := syscall.Syscall9(procCreateNamedPipe.Addr(), 8,
		uintptr(unsafe.Pointer(pName)),
		uintptr(dwOpenMode),
		uintptr(dwPipeMode),
		uintptr(nMaxInstances),
		uintptr(nOutBufferSize),
		uintptr(nInBufferSize),
		uintptr(nDefaultTimeOut),
		uintptr(unsafe.Pointer(pSecurityAttributes)),
		0)
	if h == INVALID_HANDLE_VALUE {
		wec := WinErrorCode(e1)
		if wec != 0 {
			err = wec
		} else {
			err = errors.New("CreateNamedPipe failed.")
		}
	} else {
		h = HANDLE(r1)
	}
	return
}
