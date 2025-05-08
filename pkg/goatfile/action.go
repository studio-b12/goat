package goatfile

import (
	"fmt"

	"github.com/studio-b12/goat/pkg/goatfile/ast"
)

type ActionType int

const (
	ActionRequest = ActionType(iota + 1)
	ActionLogSection
	ActionExecute
)

func ActionFromAst(act ast.Action, path string) (Action, error) {
	switch a := act.(type) {
	case *ast.Request:
		return RequestFromAst(a, path)
	case ast.LogSection:
		return LogSection(a.Content), nil
	case *ast.Execute:
		return ExecuteFromAst(a, path)
	default:
		return nil, fmt.Errorf("invalid action ast type: %+v", act)
	}
}

// Action is used to determine the
// ActionType of an action definition
// used to cast the action to the specific
// Action implementation.
type Action interface {
	fmt.Stringer

	Type() ActionType
}
