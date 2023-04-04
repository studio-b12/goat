package goatfile

type ActionType int

const (
	ActionRequest = ActionType(iota + 1)
	ActionLogSection
)

type Action interface {
	Type() ActionType
}
