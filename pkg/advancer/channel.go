package advancer

// Channel implements Advancer and Waiter for a
// go channel. Advancing will block the current
// go routine until the channel has been awaited.
// Wait will block the current goroutine until
// the waiter has been advanced.
type Channel chan struct{}

var _ Advancer = (*Channel)(nil)
var _ Waiter = (*Channel)(nil)

func (t Channel) Advance() {
	t <- struct{}{}
}

func (t Channel) Wait() {
	<-t
}
