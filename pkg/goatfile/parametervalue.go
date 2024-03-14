package goatfile

import (
	"fmt"
)

// ParameterValue holds a go template value as string
// allowing to apply parameters onto it and parsing
// the result.
type ParameterValue string

// ApplyTemplate applies the passed params onto the teplate
// value and parses the result using a new instance of Parser
// as sub-parser.
func (t ParameterValue) ApplyTemplate(params any) (any, error) {
	b, err := ApplyTemplateBuf(fmt.Sprintf("{{%s}}", t), params)
	if err != nil {
		return nil, err
	}

	v, _, err := NewParser(b, "").parseValue()
	return v, err
}
