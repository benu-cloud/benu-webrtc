package pkgerrors

import "fmt"

// KeyboardInputError indicates a proble sending keyboard input to OS
type KeyboardInputError struct {
	// Refer to https://learn.microsoft.com/en-us/windows/win32/debug/system-error-codes
	Code int
}

func (e *KeyboardInputError) Error() string {
	return fmt.Sprintf("KeyboardInputError: Code %d", e.Code)
}

func NewKeyboardInputError(code int) error {
	return &KeyboardInputError{
		Code: code,
	}
}

// MouseInputError indicates that the sending of mouse input to OS failed
type MouseInputError struct {
	// Refer to https://learn.microsoft.com/en-us/windows/win32/debug/system-error-codes
	Code int
}

func (e *MouseInputError) Error() string {
	return fmt.Sprintf("MouseInputError: Code %d", e.Code)
}

func NewMouseInputError(code int) error {
	return &MouseInputError{
		Code: code,
	}
}
