package pkgerrors

import "fmt"

// ! MUST MATCH ERROR CODES FROM C CODE
var streamErrorCodes = map[int]string{
	1: "ERROR_ENCODER_NOT_SUPPORTED",
	2: "ERROR_PIPELINE_ALREADY_CREATED",
	3: "ERROR_PIPELINE_PARSE_BAD_FORMAT",
	4: "ERROR_PIPELINE_BAD_STATE",
	5: "ERROR_PIPELINE_SET_STATE",
	6: "ERROR_LINKING_PEER",
	8: "ERROR_BAD_PEER_ID",
	7: "ERROR_PIPELINE_DOESNT_EXIST",
	9: "ERROR_BAD_SDP",
}

// CStreamError indicates an error from the stream C code
type CStreamError struct {
	ErrorCode    int
	ErrorMessage string
}

func (e *CStreamError) Error() string {
	if e.ErrorMessage == "" {
		errorStr, ok := streamErrorCodes[e.ErrorCode]
		if !ok {
			return fmt.Sprintf("CStreamError: got unknown error with code %d from stream", e.ErrorCode)
		}
		return fmt.Sprintf("CStreamError: got error %s with code %d from stream", errorStr, e.ErrorCode)
	} else {
		return fmt.Sprintf("CStreamError: got error '%s'", e.ErrorMessage)
	}
}

func NewCStreamError(errorCode int) error {
	return &CStreamError{
		ErrorCode: errorCode,
	}
}

func NewCStreamErrorWithMessage(errorMessage string) error {
	return &CStreamError{
		ErrorMessage: errorMessage,
	}
}

// StreamError indicates an error from the Go code
type StreamError struct {
	Err error
}

func (e *StreamError) Error() string {
	return fmt.Sprintf("StreamError: %v", e.Err)
}

func NewStreamError(err error) error {
	return &StreamError{
		Err: err,
	}
}
