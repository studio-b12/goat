package gurlfile

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

// applyTemplateBuf parses the given raw string as a template
// and applies the given values in params onto it returning
// the result as bytes buffer.
//
// If a key in the template is not present in the params,
// an error will be returned.
func applyTemplateBuf(raw string, params any) (*bytes.Buffer, error) {
	tmpl, err := template.New("").
		Funcs(builtinFuncsMap).
		Option("missingkey=error").
		Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parsing template failed: %s", err.Error())
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, params)
	if err != nil {
		return nil, fmt.Errorf("executing template failed: %s", err.Error())
	}

	return &out, err
}

// applyTemplate parses the given raw string as a template
// and applies the given values in params onto it returning
// the result as string.
//
// If a key in the template is not present in the params,
// an error will be returned.
func applyTemplate(raw string, params any) (string, error) {
	out, err := applyTemplateBuf(raw, params)
	if err != nil {
		return "", err
	}

	outStr := unescapeTemplateDelims(out.String())
	return outStr, nil
}

// applyTemplateToArray executes applyTemplate
// on all string instances in the given array
// or sub arrays.
func applyTemplateToArray(arr []any, params any) (err error) {
	for i, v := range arr {
		switch vt := v.(type) {
		case string:
			arr[i], err = applyTemplate(vt, params)
		case []any:
			err = applyTemplateToArray(vt, params)
		default:
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// applyTemplateToMap executes applyTemplate
// on all values in the given map.
func applyTemplateToMap(m map[string]any, params any) (err error) {
	for k, v := range m {
		switch vt := v.(type) {
		case ParameterValue:
			m[k], err = vt.ApplyTemplate(params)
		case string:
			m[k], err = applyTemplate(vt, params)
		case []any:
			err = applyTemplateToArray(vt, params)
		case map[string]any:
			err = applyTemplateToMap(vt, params)
		default:
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// extend takes a file path and adds the given extension
// to it if the path does not end with any file extension.
func extend(v string, ext string) string {
	if filepath.Ext(v) == "" {
		return v + "." + ext
	}

	return v
}

// crlf2lf converts all CRLF line endings in the given
// string to LF line endings and returns the result.
func crlf2lf(v string) string {
	return strings.ReplaceAll(v, "\r\n", "\n")
}

// unescapeTemplateDelims unescapes escaped
// template delimiter characters.
// For example, "\{\{.foo\}\}" becomes "{{.foo}}".
func unescapeTemplateDelims(v string) string {
	v = strings.ReplaceAll(v, "\\{", "{")
	v = strings.ReplaceAll(v, "\\}", "}")
	return v
}
