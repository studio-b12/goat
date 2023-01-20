package advancer

// None implements a shallow Advancer and Waiter
// where nothing will block on either advancing
// or waiting.
type None struct{}

var _ Advancer = (*None)(nil)
var _ Waiter = (*None)(nil)

func (t None) Advance() {}

func (t None) Wait() {}
