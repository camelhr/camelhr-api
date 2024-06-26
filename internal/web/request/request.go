package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/go-chi/chi/v5"
)

type requestContextKey int

const (
	CtxUserIDKey requestContextKey = iota
	CtxOrgIDKey
	CtxOrgSubdomainKey
)

var ErrInvalidPathParam = errors.New("invalid path parameter")

// decodeJSON decodes a JSON payload from the given reader into the given value.
// It also validates the fields using the validator.
func decodeJSON(r io.Reader, v any) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck // ignore the error as we are discarding the body

	if err := json.NewDecoder(r).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON payload: %w", err)
	}

	return nil
}

// validateRequestPayload validates the given value using the validator.
func validateRequestPayload(v any) error {
	// validate the request payload if it is a struct
	if reflect.TypeOf(v).Kind() == reflect.Ptr && reflect.TypeOf(v).Elem().Kind() == reflect.Struct {
		if err := base.Validator().Struct(v); err != nil {
			return fmt.Errorf("invalid request payload: %w", err)
		}
	}

	return nil
}

// DecodeAndValidateJSON decodes a JSON payload from the given reader into the given value.
// It also validates the fields using the validator.
// If the validation fails, it returns a validation error.
func DecodeAndValidateJSON(r io.Reader, v any) error {
	if err := decodeJSON(r, v); err != nil {
		return err
	}

	if err := validateRequestPayload(v); err != nil {
		return err
	}

	return nil
}

// GetURLParam returns the unescaped value of a URL parameter from the request.
// This helps prevent dependency on a specific web framework.
func URLParam(r *http.Request, param string) string {
	val := chi.URLParam(r, param)

	// https://github.com/go-chi/chi/issues/642
	// if RawPath is set, it means the URL was parsed from a raw URL and the value is already unescaped.
	if r.URL.RawPath != "" {
		val, _ = url.PathUnescape(val)
	}

	return val
}

// GetURLParamID returns the value of a ID URL parameter from the request as an int64.
func URLParamID(r *http.Request, param string) (int64, error) {
	id, err := strconv.ParseInt(URLParam(r, param), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' param from url, error: %w", param, err)
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid value %d for url param '%s': %w", id, param, ErrInvalidPathParam)
	}

	return id, nil
}
