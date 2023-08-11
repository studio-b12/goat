package executor

import (
	"fmt"

	"github.com/studio-b12/goat/pkg/clr"
	"github.com/zekrotja/rogu/log"
)

type Result struct {
	Setup    ResultSection
	Teardown ResultSection
	Tests    ResultSection
}

func (t *Result) Merge(other Result) {
	t.Setup.Merge(other.Setup)
	t.Teardown.Merge(other.Teardown)
	t.Tests.Merge(other.Tests)
}

func (t Result) All() int {
	return t.Setup.All() + t.Teardown.All() + t.Tests.All()
}

func (t Result) Failed() int {
	return t.Setup.Failed() + t.Teardown.Failed() + t.Tests.Failed()
}

func (t Result) Successfull() int {
	return t.Setup.Successfull() + t.Teardown.Successfull() + t.Tests.Successfull()
}

func (t Result) Log() {
	c := clr.ColorFGGreen
	if t.Failed() > 0 {
		c = clr.ColorFGRed
	}

	log.Info().
		Field("setup", fmt.Sprintf("%d/%d", t.Setup.Successfull(), t.Setup.Failed())).
		Field("tests", fmt.Sprintf("%d/%d", t.Tests.Successfull(), t.Tests.Failed())).
		Field("teardown", fmt.Sprintf("%d/%d", t.Teardown.Successfull(), t.Teardown.Failed())).
		Msg(clr.Print(clr.Format(
			fmt.Sprintf("Ran %d requests: %d succeeded and %d failed", t.All(), t.Successfull(), t.Failed()), c)))
}

func (t Result) Sum() (res ResultSection) {
	res.Merge(t.Setup)
	res.Merge(t.Tests)
	res.Merge(t.Teardown)

	return res
}

type ResultSection struct {
	failed int
	all    int
}

func (t *ResultSection) Merge(other ResultSection) {
	t.failed += other.failed
	t.all += other.all
}

func (t *ResultSection) Inc() {
	t.all++
}

func (t *ResultSection) IncFailed() {
	t.failed++
}

func (t ResultSection) All() int {
	return t.all
}

func (t ResultSection) Failed() int {
	return t.failed
}

func (t ResultSection) Successfull() int {
	return t.all - t.failed
}
