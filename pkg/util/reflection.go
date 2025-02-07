package util

import (
	"reflect"
)

// UnwrapPointer takes a reflect.Value and calls Elem() on it until the
// value is not a pointer anymore. The resulting unwrapped value is
// returned.
func UnwrapPointer(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
