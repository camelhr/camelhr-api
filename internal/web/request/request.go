package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// DecodeJSON decodes a JSON payload from the given reader into the given value
func DecodeJSON(r io.Reader, v any) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck
	return json.NewDecoder(r).Decode(v)
}

// GetURLParam returns the unescaped value of a URL parameter from the request
// This helps prevent dependency on a specific web framework
func URLParam(r *http.Request, param string) string {
	val := chi.URLParam(r, param)

	// https://github.com/go-chi/chi/issues/642
	// if RawPath is set, it means the URL was parsed from a raw URL and the value is already unescaped
	if r.URL.RawPath != "" {
		val, _ = url.PathUnescape(val)
	}
	return val
}

// GetURLParamID returns the value of a URL parameter from the request as an int64
func URLParamID(r *http.Request, param string) (int64, error) {
	id, err := strconv.ParseInt(URLParam(r, param), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' param from url", param)
	}
	if id <= 0 {
		return 0, fmt.Errorf("failed to parse '%s' param from url", param)
	}
	return id, nil
}
