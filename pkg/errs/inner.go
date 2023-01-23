package errs

// InnerError wraps an inner error and
// implements the Unwrap() function to
// unwrap the inner error.
type InnerError struct {
	Inner error
}

func (t InnerError) Unwrap() error {
	return t.Inner
}
