package engine

import (
	"reflect"

	"github.com/dop251/goja"
)

// Goja is the Engine implementation using
// ECMAScript 5.
type Goja struct {
	rt *goja.Runtime
}

var _ Engine = (*Goja)(nil)

// NewGoja initializes the Goja engine runtime
// and sets builtin functions to the global scope.
func NewGoja() Engine {
	var t Goja

	t.rt = goja.New()

	t.Set("assert", t.builtin_assert)
	t.Set("assert_eq", t.builtin_assert_eq)
	t.Set("debug", t.builtin_debug)
	t.Set("debugf", t.builtin_debugf)
	t.Set("info", t.builtin_info)
	t.Set("infof", t.builtin_infof)
	t.Set("warn", t.builtin_warn)
	t.Set("warnf", t.builtin_warnf)
	t.Set("error", t.builtin_error)
	t.Set("errorf", t.builtin_errorf)
	t.Set("fatal", t.builtin_fatal)
	t.Set("fatalf", t.builtin_fatalf)
	t.Set("print", t.builtin_print)
	t.Set("printf", t.builtin_printf)
	t.Set("println", t.builtin_println)
	t.Set("jq", t.builtin_jq)

	return &t
}

func (t *Goja) SetState(s State) {
	for k, v := range s {
		t.Set(k, v)
	}
}

func (t *Goja) Set(name string, v any) error {
	return t.rt.Set(name, v)
}

func (t *Goja) Run(script string) error {
	_, err := t.rt.RunString(script)
	if gojaException, ok := err.(*goja.Exception); ok {
		// Extract Goja Exceptions into a new exception
		// wrapper so that we can handle how error
		// messages are printed.
		var ex Exception
		ex.Inner = gojaException
		val := gojaException.Value()
		if val != nil {
			ex.Msg = val.String()
		}
		return ex
	}
	return err
}

func (t *Goja) State() State {
	values := make(State)
	for _, key := range t.rt.GlobalObject().Keys() {
		v := t.rt.Get(key)
		typ := v.ExportType()
		// Don't extract <null> values or function
		// type instances.
		if typ == nil || typ.Kind() == reflect.Func {
			continue
		}
		values[key] = v.Export()
	}

	return values
}
