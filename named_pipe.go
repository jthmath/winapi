package winapi

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
