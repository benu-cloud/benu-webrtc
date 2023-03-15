package config

// video encoder
type VideoEncoder string

// Supported video encoders
// ! Must be compatible with encoders defined in C code
const (
	VP9    VideoEncoder = "VP9"
	H264   VideoEncoder = "H264"
	NVH264 VideoEncoder = "NVH264"
)

// Resolution
type Resolution struct {
	Height int
	Width  int
}

type PortNumber uint

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

// message broker settings
type MessageBrokerSettings struct {
	Host     string
	Port     PortNumber
	VHost    string
	Username string
	Password string
}
