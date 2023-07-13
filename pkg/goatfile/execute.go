package goatfile

var _ (Action) = (*Execute)(nil)

type Execute struct {
	File    string
	Params  map[string]any
	Returns map[string]string
}

func (t Execute) Type() ActionType {
	return ActionExecute
}
