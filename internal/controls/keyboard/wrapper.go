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
	"BACKSPACE": C.VK_BACK,
	"DELETE":    C.VK_DELETE,
	"RETURN":    C.VK_RETURN,
	"TAB":       C.VK_TAB,
	"ESCAPE":    C.VK_ESCAPE,
	"UP":        C.VK_UP,
	"DOWN":      C.VK_DOWN,
	"RIGHT":     C.VK_RIGHT,
	"LEFT":      C.VK_LEFT,
	"HOME":      C.VK_HOME,
	"END":       C.VK_END,
	"PAGEUP":    C.VK_PRIOR,
	"PAGEDOWN":  C.VK_NEXT,

	"F1":  C.VK_F1,
	"F2":  C.VK_F2,
	"F3":  C.VK_F3,
	"F4":  C.VK_F4,
	"F5":  C.VK_F5,
	"F6":  C.VK_F6,
	"F7":  C.VK_F7,
	"F8":  C.VK_F8,
	"F9":  C.VK_F9,
	"F10": C.VK_F10,
	"F11": C.VK_F11,
	"F12": C.VK_F12,
	"F13": C.VK_F13,
	"F14": C.VK_F14,
	"F15": C.VK_F15,
	"F16": C.VK_F16,
	"F17": C.VK_F17,
	"F18": C.VK_F18,
	"F19": C.VK_F19,
	"F20": C.VK_F20,
	"F21": C.VK_F21,
	"F22": C.VK_F22,
	"F23": C.VK_F23,
	"F24": C.VK_F24,

	"META":        C.VK_LWIN,
	"LMETA":       C.VK_LWIN,
	"RMETA":       C.VK_RWIN,
	"ALT":         C.VK_MENU,
	"LALT":        C.VK_LMENU,
	"RALT":        C.VK_RMENU,
	"CONTROL":     C.VK_CONTROL,
	"LCONTROL":    C.VK_LCONTROL,
	"RCONTROL":    C.VK_RCONTROL,
	"SHIFT":       C.VK_SHIFT,
	"LSHIFT":      C.VK_LSHIFT,
	"RSHIFT":      C.VK_RSHIFT,
	"CAPSLOCK":    C.VK_CAPITAL,
	"SPACE":       C.VK_SPACE,
	"PRINTSCREEN": C.VK_SNAPSHOT,
	"INSERT":      C.VK_INSERT,
	"MENU":        C.VK_APPS,

	"NUMPAD_0":    C.VK_NUMPAD0,
	"NUMPAD_1":    C.VK_NUMPAD1,
	"NUMPAD_2":    C.VK_NUMPAD2,
	"NUMPAD_3":    C.VK_NUMPAD3,
	"NUMPAD_4":    C.VK_NUMPAD4,
	"NUMPAD_5":    C.VK_NUMPAD5,
	"NUMPAD_6":    C.VK_NUMPAD6,
	"NUMPAD_7":    C.VK_NUMPAD7,
	"NUMPAD_8":    C.VK_NUMPAD8,
	"NUMPAD_9":    C.VK_NUMPAD9,
	"NUMPAD_LOCK": C.VK_NUMLOCK,

	"NUMPAD_DECIMAL": C.VK_DECIMAL,
	"NUMPAD_PLUS":    C.VK_ADD,
	"NUMPAD_MINUS":   C.VK_SUBTRACT,
	"NUMPAD_MUL":     C.VK_MULTIPLY,
	"NUMPAD_DIV":     C.VK_DIVIDE,
	"NUMPAD_ENTER":   C.VK_RETURN,
	"NUMPAD_EQUAL":   C.VK_OEM_PLUS,

	"AUDIO_VOLUME_MUTE": C.VK_VOLUME_MUTE,
	"AUDIO_VOLUME_DOWN": C.VK_VOLUME_DOWN,
	"AUDIO_VOLUME_UP":   C.VK_VOLUME_UP,
	"AUDIO_PLAY":        C.VK_MEDIA_PLAY_PAUSE,
	"AUDIO_STOP":        C.VK_MEDIA_STOP,
	"AUDIO_PAUSE":       C.VK_MEDIA_PLAY_PAUSE,
	"AUDIO_PREV":        C.VK_MEDIA_PREV_TRACK,
	"AUDIO_NEXT":        C.VK_MEDIA_NEXT_TRACK,
}

// c implementation
type Keyboard_c struct{}

func (c *Keyboard_c) SendInputKeyChar(key rune, down bool) error {
	if code := int(C.sendInputKeyChar((C.char)(key), (C.bool)(down))); code != 0 {
		return pkgerrors.NewKeyboardInputError(code)
	}
	return nil
}

func (c *Keyboard_c) SendInputKeySpecialKey(k types.SpecialKeyboardKey, down bool) error {
	ckey, ok := specialKeyboardKey[k]
	if !ok {
		return pkgerrors.NewNotImplementedError("SendInputKeySpecialKey", string(k))
	}
	if code := int(C.sendInputKeyCode(ckey, (C.bool)(down))); code != 0 {
		return pkgerrors.NewKeyboardInputError(code)
	}
	return nil
}
