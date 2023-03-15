package pkgerrors

import "fmt"

// BadCommanlineArgument indicates a badly formatted commandline argument
type BadCommanlineArgument struct {
	For            string
	Got            string
	ExpectedFormat string
}

func (e *BadCommanlineArgument) Error() string {
	return fmt.Sprintf("BadCommanlineArgument: for %s, got '%s'; however, the expected format is '%s'", e.For, e.Got, e.ExpectedFormat)
}

func NewBadCommanlineArgument(where string, got string, expectedFormat string) error {
	return &BadCommanlineArgument{
		For:            where,
		Got:            got,
		ExpectedFormat: expectedFormat,
	}
}
