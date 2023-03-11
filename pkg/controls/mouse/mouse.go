package mouse

import "github.com/benu-cloud/benu-livestreaming-gst/pkg/controls/types"

type Mouse interface {
	SendInputMove(dx int, dy int) error
	SendInputKey(button types.MouseKey) error
	SendInputScroll(direction types.MouseWheelDir, magnitude int) error
}
