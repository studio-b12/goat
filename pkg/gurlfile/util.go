package gurlfile

import (
	"bytes"
	"fmt"
	"strings"
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

func removeComments(raw string) string {
	lines := strings.Split(raw, "\n")

	for i, line := range lines {
		cidx := strings.Index(line, "//")
		if cidx == -1 {
			continue
		}

		if cidx > 0 {
			if line[cidx-1] == ' ' {
				cidx -= 1
			} else {
				continue
			}
		}

		lines[i] = line[:cidx]
	}

	return strings.Join(lines, "\n")
}
