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

// ErrorResponse writes an error response for the given error.
// Appropriate http status code will be used based on the error type.
// Error message will be sent in the http response for applicable errors.
// If the error is of type base.APIError with not-nil cause, the cause will be logged.
// If the error is of type base.APIError with not-nil http status code, it will be used. Otherwise, default is 500.
// Appropriate errors will be logged.
func ErrorResponse(w http.ResponseWriter, err error) {
	var apiErr *base.APIError

	var notFoundErr *base.NotFoundError

	var inputValidationErr *base.InputValidationError

	var validationErrors validator.ValidationErrors

	// handle api error
	if ok := errors.As(err, &apiErr); ok {
		statusCode, message := processAPIError(apiErr)
		JSON(w, statusCode, &errorResponse{ErrorText: message})

		return
	}

	// handle not found error
	if ok := errors.As(err, &notFoundErr); ok {
		statusCode, message := processNotFoundError(notFoundErr)
		JSON(w, statusCode, &errorResponse{ErrorText: message})

		return
	}

	// handle input validation error
	if ok := errors.As(err, &inputValidationErr); ok {
		statusCode, message := processInputValidationError(inputValidationErr)
		JSON(w, statusCode, &errorResponse{ErrorText: message})

		return
	}

	// handle validation errors
	if ok := errors.As(err, &validationErrors); ok {
		statusCode, message := processValidationErrors(validationErrors)
		JSON(w, statusCode, &errorResponse{ErrorText: message})

		return
	}

	// log and send generic error
	log.Error("%v", err)
	JSON(w, http.StatusInternalServerError, &errorResponse{ErrorText: ""})
}

// processAPIError processes the API error and returns the status code and message.
func processAPIError(apiErr *base.APIError) (int, string) {
	statusCode := http.StatusInternalServerError
	message := apiErr.Error()

	cause := apiErr.Unwrap()
	if cause != nil {
		log.Error("%v", cause)
	}

	if httpStatusCode := apiErr.HTTPStatusCode(); httpStatusCode != nil {
		statusCode = *httpStatusCode
	}

	return statusCode, message
}

// processNotFoundError processes the not found error and returns the status code and message.
func processNotFoundError(notFoundErr *base.NotFoundError) (int, string) {
	return http.StatusNotFound, notFoundErr.Error()
}

// processInputValidationError processes the input validation error and returns the status code and message.
func processInputValidationError(inputValidationErr *base.InputValidationError) (int, string) {
	return http.StatusBadRequest, inputValidationErr.Error()
}

// processValidationErrors processes the validation errors and returns the status code and message.
func processValidationErrors(validationErrors validator.ValidationErrors) (int, string) {
	trans := base.ValidationTranslator()
	message := ""

	for _, fieldErr := range validationErrors {
		if message != "" {
			message += ". "
		}

		message += fieldErr.Translate(trans)
	}

	return http.StatusBadRequest, message
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

// SetCookie sets a cookie in the response.
func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,                  // time in seconds until the cookie expires
		HttpOnly: true,                    // prevent javascript access
		Secure:   true,                    // only send over https
		SameSite: http.SameSiteStrictMode, // do not send on cross-site requests. prevent csrf
	})
}

// RemoveCookie removes a cookie from the response.
func RemoveCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // delete the cookie
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// Redirect redirects the request to the given url with status code 302.
func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
}
