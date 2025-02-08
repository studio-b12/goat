package util

import (
	"encoding/json"
	"fmt"
)

// SafeJsonMarshalIndent takes any object and
// decodes it into indented JSON, when possible.
//
// If an error is returned by the JSON marshaller,
// the error will be returned as string.
func SafeJsonMarshalIndent(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("<error: %s>", err.Error())
	}
	return string(b)
}
