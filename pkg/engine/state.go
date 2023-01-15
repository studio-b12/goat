package engine

// State holds a key-value definition of globally
// availabe variables in a runtime.
type State map[string]any

// Merge applies the entries from with to the
// current state. Already set keys will be
// overwritten.
func (t State) Merge(with State) {
	for k, v := range with {
		t[k] = v
	}
}
