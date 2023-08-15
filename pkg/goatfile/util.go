package goatfile

import (
	"bytes"
	"encoding/json"
	"path"
	"strings"
	"text/template"

	"github.com/studio-b12/goat/pkg/errs"
)

// ApplyTemplateBuf parses the given raw string as a template
// and applies the given values in params onto it returning
// the result as bytes buffer.
//
// If a key in the template is not present in the params,
// an error will be returned.
func ApplyTemplateBuf(raw string, params any) (*bytes.Buffer, error) {
	tmpl, err := template.New("").
		Funcs(builtinFuncsMap).
		Option("missingkey=error").
		Parse(raw)
	if err != nil {
		return nil, errs.WithPrefix("parsing template failed:", err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, params)
	if err != nil {
		return nil, errs.WithPrefix("executing template failed:", err)
	}

	return &out, err
}

// ApplyTemplate parses the given raw string as a template
// and applies the given values in params onto it returning
// the result as string.
//
// If a key in the template is not present in the params,
// an error will be returned.
func ApplyTemplate(raw string, params any) (string, error) {
	if m, ok := params.(map[string]any); ok {
		transformPrintable(m)
	}

	out, err := ApplyTemplateBuf(raw, params)
	if err != nil {
		return "", err
	}

	outStr := unescapeTemplateDelims(out.String())
	return outStr, nil
}

// ApplyTemplateToArray executes applyTemplate
// on all string instances in the given array
// or sub arrays.
func ApplyTemplateToArray(arr []any, params any) (err error) {
	for i, v := range arr {
		switch vt := v.(type) {
		case string:
			arr[i], err = ApplyTemplate(vt, params)
		case []any:
			err = ApplyTemplateToArray(vt, params)
		default:
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// ApplyTemplateToMap executes applyTemplate
// on all values in the given map.
func ApplyTemplateToMap(m map[string]any, params any) (err error) {
	for k, v := range m {
		switch vt := v.(type) {
		case ParameterValue:
			m[k], err = vt.ApplyTemplate(params)
		case string:
			m[k], err = ApplyTemplate(vt, params)
		case []any:
			err = ApplyTemplateToArray(vt, params)
		case map[string]any:
			err = ApplyTemplateToMap(vt, params)
		default:
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Extend takes a file path and adds the given extension
// to it if the path does not end with any file extension.
func Extend(v string, ext string) string {
	if path.Ext(v) == "" {
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

func mustMarshal(v any) string {
	r, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(r)
}

type printableList []any

func (t printableList) String() string {
	return mustMarshal(t)
}

type printableMap map[string]any

func (t printableMap) String() string {
	return mustMarshal(t)
}

func transformPrintable(m map[string]any) {
	for key, val := range m {
		switch tval := val.(type) {
		case []any:
			m[key] = printableList(tval)
		case map[string]any:
			m[key] = printableMap(tval)
			transformPrintable(tval)
		}
	}
}
