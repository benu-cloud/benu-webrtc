package main

import (
	"fmt"
	"time"

	"github.com/benu-cloud/benu-livestreaming-gst/internal/controls/keyboard"
)

func main() {
	k := &keyboard.Keyboard_c{}
	time.Sleep(time.Second)
	fmt.Println("sending a")
	k.SendInputKeyChar('a', true)
	k.SendInputKeyChar('a', false)
}
