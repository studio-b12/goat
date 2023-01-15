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

	t.rt.Set("assert", t.builtin_assert)
	t.rt.Set("debug", t.builtin_debug)
	t.rt.Set("info", t.builtin_info)
	t.rt.Set("warn", t.builtin_warn)
	t.rt.Set("error", t.builtin_error)
	t.rt.Set("fatal", t.builtin_fatal)

	return &t
}

func (t *Goja) SetState(s State) {
	for k, v := range s {
		t.Register(k, v)
	}
}

func (t *Goja) Register(name string, v any) error {
	return t.rt.Set(name, v)
}

func (t *Goja) Run(script string) error {
	_, err := t.rt.RunString(script)
	return err
}

func (t *Goja) State() State {
	values := make(State)
	for _, key := range t.rt.GlobalObject().Keys() {
		v := t.rt.Get(key)
		if v.ExportType().Kind() == reflect.Func {
			continue
		}
		values[key] = v.Export()
	}

	return values
}
