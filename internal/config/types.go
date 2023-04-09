package config

type (
	// video encoder
	VideoEncoder int
	// port number
	PortNumber uint
)

// Supported video encoders
// ! Must be compatible with encoders defined in C code
const (
	VP9    VideoEncoder = 0
	H264   VideoEncoder = 1
	NVH264 VideoEncoder = 2
)

// Resolution
type Resolution struct {
	Height int
	Width  int
}

// stream settings
type StreamSettings struct {
	// video
	VideoResolution    Resolution
	VideoEncoder       VideoEncoder
	VideoBaseFramerate uint
	VideoBaseBitrate   uint
	VideoShowCursor    bool
	// audio
	AudioBaseBitrate       uint
	AudioBasePacketLossPct uint
}
