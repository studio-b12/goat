package goatfile

import (
	"fmt"
)

type LogSection string

var _ Action = (*LogSection)(nil)

func (t LogSection) Type() ActionType {
	return ActionLogSection
}

func (t LogSection) String() string {
	return fmt.Sprintf("LogSection('%s')", string(t))
}
