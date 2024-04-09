package request

import (
	"encoding/json"
	"io"
)

// DecodeJSON decodes a JSON payload from the given reader into the given value
func DecodeJSON(r io.Reader, v any) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck
	return json.NewDecoder(r).Decode(v)
}
