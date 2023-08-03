package goatfile

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"strconv"
	"text/template"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var builtinFuncsMap = template.FuncMap{
	"base64":       builtin_base64,
	"base64url":    builtin_base64Url,
	"md5":          builtin_hasher(md5.New()),
	"sha1":         builtin_hasher(sha1.New()),
	"sha256":       builtin_hasher(sha256.New()),
	"sha512":       builtin_hasher(sha512.New()),
	"randomString": builtin_randomString,
	"randomInt":    builtin_randomInt,
	"timestamp":    builtin_timestamp,
	"isset":        builtin_isset,
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

func builtin_randomString(lnOpt ...int) string {
	const defaultLen = 8
	const charSet = "abcdefhijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	ln := defaultLen
	if len(lnOpt) != 0 {
		ln = lnOpt[0]
	}

	buf := make([]byte, ln)
	for i := 0; i < ln; i++ {
		buf[i] = charSet[rng.Intn(len(charSet))]
	}

	return string(buf)
}

func builtin_randomInt(nOpt ...int) int {
	if len(nOpt) != 0 {
		return rng.Intn(nOpt[0])
	}

	return rng.Int()
}

func builtin_timestamp(formatOpt ...string) string {
	now := time.Now()

	if len(formatOpt) != 0 {
		return now.Format(formatOpt[0])
	}

	return strconv.Itoa(int(now.Unix()))
}

func builtin_isset(m map[string]any, key string) bool {
	v, ok := m[key]
	return ok && v != nil
}
