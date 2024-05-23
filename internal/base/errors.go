package base

import "errors"

type APIError struct {
	httpStatusCode *int
	cause          error
	msg            string
}

// APIErrorOption is a function that modifies an property of an API error.
type APIErrorOption func(*APIError) *APIError

// Error returns the error message.
func (e *APIError) Error() string {
	return e.msg
}

// NewAPIError creates a new API error with the given message and options.
// Use it to send the error message in the response.
// Pass the options to customize the error.
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
// It defaults the message and cause to the given error.
// Use it to log and send the error message in the response.
// Pass the options to customize the error.
func WrapError(err error, options ...APIErrorOption) error {
	apiErr := &APIError{cause: err, msg: err.Error()}
	for _, option := range options {
		apiErr = option(apiErr)
	}

	return apiErr
}

// Unwrap unwraps the cause of the API error.
func (e *APIError) Unwrap() error {
	return e.cause
}

// HTTPStatusCode returns the http status code of the API error.
func (e *APIError) HTTPStatusCode() *int {
	return e.httpStatusCode
}

// ErrorCause sets the cause to the API error.
func ErrorCause(cause error) APIErrorOption {
	return func(apiErr *APIError) *APIError {
		apiErr.cause = cause
		return apiErr
	}
}

// ErrorHTTPStatus sets the http status code to the API error.
func ErrorHTTPStatus(httpStatusCode int) APIErrorOption {
	return func(apiErr *APIError) *APIError {
		apiErr.httpStatusCode = &httpStatusCode
		return apiErr
	}
}

// IsAPIError checks if the given error is an API error.
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
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

// InputValidationError is an error type that represents a validation error.
type InputValidationError struct {
	msg string
}

// Error returns the error message.
func (e *InputValidationError) Error() string {
	return e.msg
}

// NewInputValidationError creates a new validation error with the given message.
func NewInputValidationError(msg string) error {
	return &InputValidationError{msg: msg}
}

// IsInputValidationError checks if the given error is a validation error.
func IsInputValidationError(err error) bool {
	var validationErr *InputValidationError
	return errors.As(err, &validationErr)
}
