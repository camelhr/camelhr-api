package base

type APIError struct {
	cause error
	msg   string
}

// Error returns the error message.
func (e *APIError) Error() string {
	return e.msg
}

// NewAPIError creates a new API error with the given message.
// This error can be used to return a custom error message in the response.
func NewAPIError(msg string) error {
	return &APIError{msg: msg}
}

// WrapError wraps the given error into an API error.
func WrapError(err error) error {
	return &APIError{msg: err.Error(), cause: err}
}

// Unwrap unwraps the cause of the API error.
func (e *APIError) Unwrap() error {
	return e.cause
}

// NotFoundError is an error type that represents a not found error.
type NotFoundError struct {
	msg string
}

// Error returns the error message.
func (e *NotFoundError) Error() string {
	return e.msg
}

// NewNotFoundError creates a new not found error with the given message.
func NewNotFoundError(msg string) error {
	return &NotFoundError{msg: msg}
}
