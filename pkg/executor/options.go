package executor

type AbortOptions struct {
	NoAbort     bool
	AlwaysAbort bool
}

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
