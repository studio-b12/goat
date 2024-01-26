package goatfile

import (
	"errors"
	"github.com/studio-b12/goat/pkg/goatfile/ast"
)

var _ (Action) = (*Execute)(nil)

type Execute struct {
	File    string
	Params  map[string]any
	Returns map[string]string

	Path string
}

func ExecuteFromAst(a *ast.Execute, path string) (t Execute, err error) {
	if a == nil {
		return Execute{}, errors.New("execute ast is nil")
	}

	t.File = a.Path
	t.Path = path
	t.Params = a.Parameters
	t.Returns = a.Returns

	return t, nil
}

func (t Execute) Type() ActionType {
	return ActionExecute
}
