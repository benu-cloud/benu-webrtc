package mouse

/*
#cgo CFLAGS: -I${SRCDIR}/c
#cgo LDFLAGS: -L${SRCDIR}/c -lmouse
#include "mouse.h"
*/
import "C"

import (
	"fmt"

	"github.com/benu-cloud/benu-webrtc/pkg/controls/types"
	"github.com/benu-cloud/benu-webrtc/pkg/pkgerrors"
)

var mouseKey map[types.MouseKey]C.int = map[types.MouseKey]C.int{
	types.LMBUp:   C.MOUSEEVENTF_LEFTUP,
	types.LMBDown: C.MOUSEEVENTF_LEFTDOWN,
	types.RMBUp:   C.MOUSEEVENTF_RIGHTUP,
	types.RMBDown: C.MOUSEEVENTF_RIGHTDOWN,
	types.MMBUp:   C.MOUSEEVENTF_MIDDLEUP,
	types.MMBDown: C.MOUSEEVENTF_MIDDLEDOWN,
	types.XMBUp:   C.MOUSEEVENTF_XUP,
	types.XMBDown: C.MOUSEEVENTF_XDOWN,
}

// c implementation
type Mouse_c struct{}

var mouseWheelDir map[types.MouseWheelDir]C.int = map[types.MouseWheelDir]C.int{
	types.HWheel: C.MOUSEEVENTF_HWHEEL,
	types.VWheel: C.MOUSEEVENTF_WHEEL,
}

func (m *Mouse_c) SendInputMove(dx int, dy int) error {
	if code := C.sendInputMove(C.int(dx), C.int(dy)); code != C.ERROR_SUCCESS {
		return pkgerrors.NewMouseInputError(int(code))
	}
	return nil
}

func (m *Mouse_c) SendInputKey(key types.MouseKey) error {
	ckey, ok := mouseKey[key]
	if !ok {
		return pkgerrors.NewNotImplementedError("SendInputKey", string(key))
	}
	if code := C.sendInputKey(ckey); code != C.ERROR_SUCCESS {
		return pkgerrors.NewMouseInputError(int(code))
	}
	return nil
}

func (m *Mouse_c) SendInputScroll(direction types.MouseWheelDir, magnitude int) error {
	cdir, ok := mouseWheelDir[direction]
	if !ok {
		return pkgerrors.NewNotImplementedError("SendInputScroll", fmt.Sprintf("direction %v", direction))
	}
	if code := C.sendInputScroll(cdir, C.int(magnitude)); code != C.ERROR_SUCCESS {
		return pkgerrors.NewMouseInputError(int(code))
	}
	return nil
}
