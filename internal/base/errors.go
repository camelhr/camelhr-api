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

// NotFoundError is an error type that represents a not found error.
type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string {
	return e.msg
}

// NewNotFoundError creates a new not found error with the given message.
func NewNotFoundError(msg string) error {
	return &NotFoundError{msg: msg}
}
