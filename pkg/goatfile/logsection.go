package goatfile

type LogSection string

var _ (Action) = (*LogSection)(nil)

func (t LogSection) Type() ActionType {
	return ActionLogSection
}
