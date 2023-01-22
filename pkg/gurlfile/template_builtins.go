package gurlfile

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"text/template"
)

var builtinFuncsMap = template.FuncMap{
	"base64":    builtin_base64,
	"base64url": builtin_base64Url,
	"md5":       builtin_hasher(md5.New()),
	"sha1":      builtin_hasher(sha1.New()),
	"sha256":    builtin_hasher(sha256.New()),
	"sha512":    builtin_hasher(sha512.New()),
}

func builtin_base64(v string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(v))
}

func builtin_base64Url(v string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(v))
}

func builtin_hasher(hsh hash.Hash) func(string) string {
	return func(s string) string {
		io.WriteString(hsh, s)
		defer hsh.Reset()
		return fmt.Sprintf("%x", hsh.Sum(nil))
	}
}
