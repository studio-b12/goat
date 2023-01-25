package gurlfile

import (
	"fmt"
)

type ParameterValue string

func (t ParameterValue) Parse(params any) (any, error) {
	b, err := applyTemplateBuf(fmt.Sprintf("{{%s}}", t), params)
	if err != nil {
		return nil, err
	}

	return NewParser(b).parseValue()
}
