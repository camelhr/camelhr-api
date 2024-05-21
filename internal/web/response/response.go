package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/log"
	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	ErrorText string `json:"error"`
}

// ErrorResponse writes an error response with the given status code.
// Appropriate error message will be sent in the response body.
// If the error is of type base.APIError with not-nil cause, the cause will be logged.
// General errors will not be logged.
func ErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	var message string

	var apiErr *base.APIError

	var notFoundErr *base.NotFoundError

	var validationErrors validator.ValidationErrors

	// send the error message for known errors
	if ok := errors.As(err, &apiErr); ok {
		message = apiErr.Error()

		cause := apiErr.Unwrap()
		if cause != nil {
			log.Error("%v", cause)
		}
	} else if ok := errors.As(err, &notFoundErr); ok {
		message = notFoundErr.Error()
	} else if ok := errors.As(err, &validationErrors); ok {
		message = buildValidationErrorMessage(validationErrors)
	} else {
		log.Error("%v", err)
	}

	JSON(w, statusCode, &errorResponse{ErrorText: message})
}

// JSON writes a JSON response with the given status code and value.
func JSON(w http.ResponseWriter, status int, v any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		log.Error("failed to encode response: %v", err)
		http.Error(w, "", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Error("failed to write response: %v", err)
	}
}

// OK writes an empty response with status code 200.
func OK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// Empty writes an empty response with the given status code.
func Empty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// Text writes a text response with the given status code and value.
func Text(w http.ResponseWriter, status int, v string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(v)); err != nil {
		log.Error("failed to write response: %v", err)
	}
}

// buildValidationErrorMessage extracts the validation error message from the given validation errors.
func buildValidationErrorMessage(errs validator.ValidationErrors) string {
	trans := base.ValidationTranslator()
	message := ""

	for _, fieldErr := range errs {
		if message != "" {
			message += ". "
		}

		message += fieldErr.Translate(trans)
	}

	return message
}
