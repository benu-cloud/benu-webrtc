package pkgerrors

import "fmt"

// NotImplementedError indicates missing implementation for a given input
type NotImplementedError struct {
	Where   string
	Feature string
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("NotImplementedError: %s not implemented in %s", e.Feature, e.Where)
}

func NewNotImplementedError(where string, feature string) error {
	return &NotImplementedError{
		Where:   where,
		Feature: feature,
	}
}
