package gurlfile

import (
	"bytes"
	"fmt"
	"text/template"
)

func applyTemplate(raw string, params any) (string, error) {
	tmpl, err := template.New("").Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parsing template failed: %s", err.Error())
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, params)
	if err != nil {
		return "", fmt.Errorf("executing template failed: %s", err.Error())
	}

	return out.String(), nil
}
