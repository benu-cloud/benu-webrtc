package keyboard

import (
	"testing"

	"github.com/benu-cloud/benu-livestreaming-gst/internal/controls/keyboard"
	"github.com/benu-cloud/benu-livestreaming-gst/pkg/controls/types"
	"github.com/benu-cloud/benu-livestreaming-gst/pkg/pkgerrors"
	"github.com/stretchr/testify/assert"
)

type sendKeyCharTest struct {
	key      rune
	down     bool
	expected error
}

type sendKeySpecialKeyTest struct {
	key      types.SpecialKeyboardKey
	down     bool
	expected error
}

func TestSendKeyChar(t *testing.T) {
	sendKeyCharTests := []sendKeyCharTest{
		{'a', true, nil},
		{'b', false, nil},
		{'\n', true, nil},
	}
	var keyboard_impl Keyboard = &keyboard.Keyboard_c{}
	for _, test := range sendKeyCharTests {
		assert.Equal(t, keyboard_impl.SendInputKeyChar(test.key, test.down), test.expected)
	}
}

func TestSendKeySpecialKey(t *testing.T) {
	sendKeySpecialKeyTests := []sendKeySpecialKeyTest{
		{types.ESCAPE, true, nil},
		{types.CAPSLOCK, false, nil},
		{types.SpecialKeyboardKey("this_is_not_a_key"), true, &pkgerrors.NotImplementedError{Where: "SendInputKeySpecialKey", Feature: "this_is_not_a_key"}},
	}
	var keyboard_impl Keyboard = &keyboard.Keyboard_c{}
	for _, test := range sendKeySpecialKeyTests {
		assert.Equal(t, keyboard_impl.SendInputKeySpecialKey(test.key, test.down), test.expected)
	}
}
