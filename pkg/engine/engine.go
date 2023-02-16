package engine

// Engine defines a service which can run scripts.
type Engine interface {
	// SetState sets the given state s
	// to the global state of the runtime.
	SetState(s State)

	// Set sets the given value by the given
	// name to the global context of the runtime.
	Set(name string, v any) error

	// Run executes the given script in the
	// runtime.
	Run(script string) error

	// State returns a map of all set
	// variables in the global state
	// which are not of the type 'function'.
	State() State
}
