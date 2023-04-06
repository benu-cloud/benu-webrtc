package pkgerrors

import "fmt"

// PublishError indicates an error publishing to rabbitmq
type PublishError struct {
	Err error
}

func (e *PublishError) Error() string {
	return fmt.Sprintf("PublishError: %v", e.Err)
}

func NewPublishError(err error) error {
	return &PublishError{
		Err: err,
	}
}

// ConsumeError indicates an error consuming from rabbitmq
type ConsumeError struct {
	Err error
}

func (e *ConsumeError) Error() string {
	return fmt.Sprintf("ConsumeError: %v", e.Err)
}

func NewConsumeError(err error) error {
	return &ConsumeError{
		Err: err,
	}
}
