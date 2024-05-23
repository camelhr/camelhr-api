package base

import "errors"

type APIError struct {
	cause error
	msg   string
}

// APIErrorOption is a function that modifies an property of an API error.
type APIErrorOption func(*APIError) *APIError

// Error returns the error message.
func (e *APIError) Error() string {
	return e.msg
}

// NewAPIError creates a new API error with the given message and options.
// Use it to send the error message in the response.
// Use the options to add more information to the error.
// If cause is added using base.ErrorCause option, the cause will be logged.
//
// Example: Create a new API error with a just a message.
//
//	err := base.NewAPIError("invalid email")
//	response.ErrorResponse(w, http.StatusBadRequest, err)
//
// Here, the error message will be sent in the response body. The error will not be logged.
//
// Example: Create a new API error with a message and a cause.
//
//	err := base.NewAPIError("invalid email", base.ErrorCause(err))
//	response.ErrorResponse(w, http.StatusBadRequest, err)
//
// Here, the error message will be sent in the response body and the cause will be logged.
func NewAPIError(msg string, options ...APIErrorOption) error {
	apiErr := &APIError{msg: msg}
	for _, option := range options {
		apiErr = option(apiErr)
	}

	return apiErr
}

// WrapError wraps the given error.
// It sets message and cause to the given error.
// Use it to log and send the error message in the response.
func WrapError(err error) error {
	return &APIError{cause: err, msg: err.Error()}
}

// Unwrap unwraps the cause of the API error.
func (e *APIError) Unwrap() error {
	return e.cause
}

// ErrorCause adds a cause to the API error.
func ErrorCause(cause error) APIErrorOption {
	return func(apiErr *APIError) *APIError {
		apiErr.cause = cause
		return apiErr
	}
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

// IsNotFoundError checks if the given error is a not found error.
func IsNotFoundError(err error) bool {
	var notFoundErr *NotFoundError
	return errors.As(err, &notFoundErr)
}
