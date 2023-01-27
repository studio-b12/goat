package set

// Set wraps a map where the value is struct{}.
type Set[T comparable] map[T]struct{}

// Add appends the given value to the set if
// it is not already existent in the Set.
//
// Returns true if the value has been added.
func (t Set[T]) Add(v T) bool {
	if t.Contains(v) {
		return false
	}

	t[v] = struct{}{}
	return true
}

// Contains returns true when the given value
// is contained in the Set.
func (t Set[T]) Contains(v T) bool {
	_, ok := t[v]
	return ok
}

// Remove deletes the given value from the Set
// if it is contained in the Set.
//
// Returns true if the value has been removed.
func (t Set[T]) Remove(v T) bool {
	if !t.Contains(v) {
		return false
	}

	delete(t, v)
	return true
}
