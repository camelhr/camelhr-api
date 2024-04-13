package base

type APIError struct {
	msg string
}

func (e *APIError) Error() string {
	return e.msg
}

// NewAPIError creates a new API error with the given message.
// This error can be used to return a custom error message in the response.
func NewAPIError(msg string) error {
	return &APIError{msg: msg}
}
