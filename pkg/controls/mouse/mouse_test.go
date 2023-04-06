package mouse

import (
	"testing"

	"github.com/benu-cloud/benu-webrtc/internal/controls/mouse"
	"github.com/benu-cloud/benu-webrtc/pkg/controls/types"
	"github.com/stretchr/testify/assert"
)

type sentInputMoveTest struct {
	dx       int
	dy       int
	expected error
}

type sendInputKeyTest struct {
	key      types.MouseKey
	expected error
}

type sendInputScrollTest struct {
	direction types.MouseWheelDir
	magnitude int
	expected  error
}

func TestSendInputMove(t *testing.T) {
	sendInputMoveTests := []sentInputMoveTest{
		{10, 10, nil},
		{0, 0, nil},
		{-10, -10, nil},
	}
	var mouse_impl Mouse = &mouse.Mouse_c{}
	for _, test := range sendInputMoveTests {
		assert.Equal(t, mouse_impl.SendInputMove(test.dx, test.dy), test.expected)
	}
}

func TestSendInputKey(t *testing.T) {
	sendInputKeyTests := []sendInputKeyTest{
		{types.MMBDown, nil},
		{types.MMBUp, nil},
		{types.RMBUp, nil},
	}
	var mouse_impl Mouse = &mouse.Mouse_c{}
	for _, test := range sendInputKeyTests {
		assert.Equal(t, mouse_impl.SendInputKey(test.key), test.expected)
	}
}

func TestSendInputScroll(t *testing.T) {
	sendInputScrollTests := []sendInputScrollTest{
		{types.HWheel, 10, nil},
		{types.HWheel, -10, nil},
		{types.VWheel, 10, nil},
	}
	var mouse_impl Mouse = &mouse.Mouse_c{}
	for _, test := range sendInputScrollTests {
		assert.Equal(t, mouse_impl.SendInputScroll(test.direction, test.magnitude), test.expected)
	}
}
