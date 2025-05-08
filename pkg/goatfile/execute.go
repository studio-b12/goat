package goatfile

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/studio-b12/goat/pkg/goatfile/ast"
)

var _ Action = (*Execute)(nil)

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
	t.Params = a.Parameters.ToMap()
	t.Returns = a.Returns.ToMap()

	return t, nil
}

func (t Execute) Type() ActionType {
	return ActionExecute
}

func (t Execute) String() string {
	var sb strings.Builder

	sb.WriteString("Execute(")
	sb.WriteString(path.Join(t.Path, t.File))
	for k, v := range t.Params {
		_, _ = fmt.Fprintf(&sb, " %s=%v", k, v)
	}
	if len(t.Returns) > 0 {
		sb.WriteString(" ->")
		for k, v := range t.Returns {
			_, _ = fmt.Fprintf(&sb, " %s>%s", k, v)
		}
	}
	sb.WriteRune(')')

	return sb.String()
}
