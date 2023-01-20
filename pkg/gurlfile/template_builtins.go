package gurlfile

import (
	"encoding/base64"
	"text/template"
)

var builtinFuncsMap = template.FuncMap{
	"base64":    builtin_base64,
	"base64Url": builtin_base64Url,
}

func builtin_base64(v string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(v))
}

func builtin_base64Url(v string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(v))
}
