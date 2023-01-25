package executor

// AbortOptions wraps options that control the
// abort behavior of an execution batch.
type AbortOptions struct {
	NoAbort     bool
	AlwaysAbort bool
}

// AbortOptionsFromMap returns a new instance of
// AbortOptions extracted from the passed map.
func AbortOptionsFromMap(m map[string]any) AbortOptions {
	opt := AbortOptions{
		NoAbort:     false,
		AlwaysAbort: false,
	}

	if v, ok := m["noabort"].(bool); ok {
		opt.NoAbort = v
	}

	if v, ok := m["alwaysabort"].(bool); ok {
		opt.AlwaysAbort = v
	}

	return opt
}

// ExecOptions wraps options that control the
// execution of a request.
type ExecOptions struct {
	Condition bool
}

// ExecOptionsFromMap returns a new instance of
// ExecOptions extracted from the passed map.
func ExecOptionsFromMap(m map[string]any) ExecOptions {
	opt := ExecOptions{
		Condition: true,
	}

	if v, ok := m["condition"].(bool); ok {
		opt.Condition = v
	}

	return opt
}
