package keyboard

import "github.com/benu-cloud/benu-livestreaming-gst/pkg/controls/types"

type Keyboard interface {
	SendInputKeyChar(key rune, down bool) error
	SendInputKeySpecialKey(key types.SpecialKeyboardKey, down bool) error
}
