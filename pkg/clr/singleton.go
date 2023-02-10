package clr

var singleton = &Printer{}

// SetEnable sets the enable state on the
// global printer instance.
func SetEnable(v bool) {
	singleton.SetEnable(v)
}

// Print takes one or more values and prints
// then to a string.
//
// When a formatWrapper is passed, the format
// will be applied depending on the global
// Printers enabled state.
func Print(v ...any) string {
	return singleton.Print(v...)
}
