package advancer

import "time"

// Ticker implements Waiter where the waiter
// will be automatically advanced after the
// given duration has passed.
type Ticker struct {
	*time.Ticker
}

// NewTicker returns a new Ticker with
// the given duration.
func NewTicker(d time.Duration) Ticker {
	return Ticker{Ticker: time.NewTicker(d)}
}

var _ Waiter = (*Ticker)(nil)

func (t Ticker) Wait() {
	<-t.C
}
