package util

import (
	"encoding/json"
	"fmt"
)

// SafeJsonMarshalIndent takes any object and
// decodes it into indentated JSON, when possible.
//
// If an error is returned by the JSON marshaler,
// the error will be returned as string.
func SafeJsonMarshalIndent(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("<error: %s>", err.Error())
	}
	return string(b)
}
