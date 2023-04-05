package goatfile

type ActionType int

const (
	ActionRequest = ActionType(iota + 1)
	ActionLogSection
)

// Action is used to determine the
// ActionType of an action definition
// used to cast the action to the specific
// Action implementation.
type Action interface {
	Type() ActionType
}
