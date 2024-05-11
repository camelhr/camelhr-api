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
// If the error is an APIError, the error message will be used in the response.
// Otherwise, the error message will be empty.
func ErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	var message string

	var apiErr *base.APIError

	var validationErr validator.ValidationErrors

	// send the error message if error is an APIError or a validation error
	if ok := errors.As(err, &apiErr); ok {
		message = apiErr.Error()
	} else if ok := errors.As(err, &validationErr); ok {
		message = validationErr.Error()
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
