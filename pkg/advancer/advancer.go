package advancer

// Advancer allows to advance a waiter.
type Advancer interface {
	// Advance a waiter.
	Advance()
}

// Waiter allows to wait for an advancement.
type Waiter interface {
	// Wait for an advancement.
	Wait()
}
