package pkgerrors

import "fmt"

// UnsupportedMessageTypeError indicates an unsupported message type has been received
type UnsupportedMessageTypeError struct {
	GotMessageType string
}

func (e *UnsupportedMessageTypeError) Error() string {
	if e.GotMessageType != "" {
		return fmt.Sprintf("UnsupportedMessageTypeError: message type '%s' is unsupported", e.GotMessageType)
	}
	return "UnsupportedMessageTypeError: unknown message type"
}

func NewUnsupportedMessageTypeError(gotMessageType string) error {
	return &UnsupportedMessageTypeError{
		GotMessageType: gotMessageType,
	}
}

// UnmarshalError indicates error in unmarshaling message
type UnmarshalError struct {
	Err error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("UnmarshalError: %v", e.Err)
}

func NewUnmarshalError(err error) error {
	return &UnmarshalError{
		Err: err,
	}
}

// UnmarshalError indicates error in marshaling message
type MarshalError struct {
	Err error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("MarshalError: %v", e.Err)
}

func NewMarshalError(err error) error {
	return &MarshalError{
		Err: err,
	}
}
