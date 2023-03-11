package types

// types for mouse input
type (
	MouseKey      string
	MouseWheelDir bool
)

// all types (implementation independent)
const (
	LMBUp   MouseKey = "LMBUp"
	LMBDown MouseKey = "LMBDown"
	RMBUp   MouseKey = "RMBUp"
	RMBDown MouseKey = "RMBDown"
	MMBUp   MouseKey = "MMBUp"
	MMBDown MouseKey = "MMBDown"
	XMBUp   MouseKey = "XMBUp"
	XMBDown MouseKey = "XMBDown"

	VWheel MouseWheelDir = true
	HWheel MouseWheelDir = false
)
