package keyboard

/*
#cgo CFLAGS: -I${SRCDIR}/c
#cgo LDFLAGS: -L${SRCDIR}/c -lkeyboard
#include "keyboard.h"
*/
import "C"

import (
	"github.com/benu-cloud/benu-livestreaming-gst/pkg/controls/types"
	"github.com/benu-cloud/benu-livestreaming-gst/pkg/pkgerrors"
)

var specialKeyboardKey map[types.SpecialKeyboardKey]C.int = map[types.SpecialKeyboardKey]C.int{
	types.BACKSPACE: C.VK_BACK,
	types.DELETE:    C.VK_DELETE,
	types.RETURN:    C.VK_RETURN,
	types.TAB:       C.VK_TAB,
	types.ESCAPE:    C.VK_ESCAPE,
	types.UP:        C.VK_UP,
	types.DOWN:      C.VK_DOWN,
	types.RIGHT:     C.VK_RIGHT,
	types.LEFT:      C.VK_LEFT,
	types.HOME:      C.VK_HOME,
	types.END:       C.VK_END,
	types.PAGEUP:    C.VK_PRIOR,
	types.PAGEDOWN:  C.VK_NEXT,

	types.F1:  C.VK_F1,
	types.F2:  C.VK_F2,
	types.F3:  C.VK_F3,
	types.F4:  C.VK_F4,
	types.F5:  C.VK_F5,
	types.F6:  C.VK_F6,
	types.F7:  C.VK_F7,
	types.F8:  C.VK_F8,
	types.F9:  C.VK_F9,
	types.F10: C.VK_F10,
	types.F11: C.VK_F11,
	types.F12: C.VK_F12,
	types.F13: C.VK_F13,
	types.F14: C.VK_F14,
	types.F15: C.VK_F15,
	types.F16: C.VK_F16,
	types.F17: C.VK_F17,
	types.F18: C.VK_F18,
	types.F19: C.VK_F19,
	types.F20: C.VK_F20,
	types.F21: C.VK_F21,
	types.F22: C.VK_F22,
	types.F23: C.VK_F23,
	types.F24: C.VK_F24,

	types.META:        C.VK_LWIN,
	types.LMETA:       C.VK_LWIN,
	types.RMETA:       C.VK_RWIN,
	types.ALT:         C.VK_MENU,
	types.LALT:        C.VK_LMENU,
	types.RALT:        C.VK_RMENU,
	types.CONTROL:     C.VK_CONTROL,
	types.LCONTROL:    C.VK_LCONTROL,
	types.RCONTROL:    C.VK_RCONTROL,
	types.SHIFT:       C.VK_SHIFT,
	types.LSHIFT:      C.VK_LSHIFT,
	types.RSHIFT:      C.VK_RSHIFT,
	types.CAPSLOCK:    C.VK_CAPITAL,
	types.SPACE:       C.VK_SPACE,
	types.PRINTSCREEN: C.VK_SNAPSHOT,
	types.INSERT:      C.VK_INSERT,
	types.MENU:        C.VK_APPS,

	types.NUMPAD_0:    C.VK_NUMPAD0,
	types.NUMPAD_1:    C.VK_NUMPAD1,
	types.NUMPAD_2:    C.VK_NUMPAD2,
	types.NUMPAD_3:    C.VK_NUMPAD3,
	types.NUMPAD_4:    C.VK_NUMPAD4,
	types.NUMPAD_5:    C.VK_NUMPAD5,
	types.NUMPAD_6:    C.VK_NUMPAD6,
	types.NUMPAD_7:    C.VK_NUMPAD7,
	types.NUMPAD_8:    C.VK_NUMPAD8,
	types.NUMPAD_9:    C.VK_NUMPAD9,
	types.NUMPAD_LOCK: C.VK_NUMLOCK,

	types.NUMPAD_DECIMAL: C.VK_DECIMAL,
	types.NUMPAD_PLUS:    C.VK_ADD,
	types.NUMPAD_MINUS:   C.VK_SUBTRACT,
	types.NUMPAD_MUL:     C.VK_MULTIPLY,
	types.NUMPAD_DIV:     C.VK_DIVIDE,
	types.NUMPAD_ENTER:   C.VK_RETURN,
	types.NUMPAD_EQUAL:   C.VK_OEM_PLUS,

	types.AUDIO_VOLUME_MUTE: C.VK_VOLUME_MUTE,
	types.AUDIO_VOLUME_DOWN: C.VK_VOLUME_DOWN,
	types.AUDIO_VOLUME_UP:   C.VK_VOLUME_UP,
	types.AUDIO_PLAY:        C.VK_MEDIA_PLAY_PAUSE,
	types.AUDIO_STOP:        C.VK_MEDIA_STOP,
	types.AUDIO_PAUSE:       C.VK_MEDIA_PLAY_PAUSE,
	types.AUDIO_PREV:        C.VK_MEDIA_PREV_TRACK,
	types.AUDIO_NEXT:        C.VK_MEDIA_NEXT_TRACK,
}

// c implementation
type Keyboard_c struct{}

func (c *Keyboard_c) SendInputKeyChar(key rune, down bool) error {
	if code := C.sendInputKeyChar((C.char)(key), (C.bool)(down)); code != C.ERROR_SUCCESS {
		return pkgerrors.NewKeyboardInputError(int(code))
	}
	return nil
}

func (c *Keyboard_c) SendInputKeySpecialKey(k types.SpecialKeyboardKey, down bool) error {
	ckey, ok := specialKeyboardKey[k]
	if !ok {
		return pkgerrors.NewNotImplementedError("SendInputKeySpecialKey", string(k))
	}
	if code := C.sendInputKeyCode(ckey, (C.bool)(down)); code != C.ERROR_SUCCESS {
		return pkgerrors.NewKeyboardInputError(int(code))
	}
	return nil
}
