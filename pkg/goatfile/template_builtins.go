package goatfile

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var builtinFuncsMap = template.FuncMap{
	"base64":            builtin_base64,
	"base64Url":         builtin_base64Url,
	"base64Unpadded":    builtin_base64Unpadded,
	"base64UrlUnpadded": builtin_base64UrlUnpadded,
	"md5":               builtin_hasher(md5.New()),
	"sha1":              builtin_hasher(sha1.New()),
	"sha256":            builtin_hasher(sha256.New()),
	"sha512":            builtin_hasher(sha512.New()),
	"randomString":      builtin_randomString,
	"randomInt":         builtin_randomInt,
	"timestamp":         builtin_timestamp,
	"isset":             builtin_isset,
	"json":              builtin_json,
	"formatTimestamp":   builtin_formatTimestamp,
}

var dateFormats = map[string]string{
	"ANSIC":       time.ANSIC,
	"UNIXDATE":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339NANO": time.RFC3339Nano,
	"KITCHEN":     time.Kitchen,
	"STAMP":       time.Stamp,
	"STAMPMILLI":  time.StampMilli,
	"STAMPMICRO":  time.StampMicro,
	"STAMPNANO":   time.StampNano,
	"DATETIME":    time.DateTime,
	"DATEONLY":    time.DateOnly,
	"TIMEONLY":    time.TimeOnly,
}

func dateFormat(format string) string {
	if trans, ok := dateFormats[strings.ToUpper(format)]; ok {
		return trans
	}
	return format
}

func builtin_base64(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func builtin_base64Url(v string) string {
	return base64.URLEncoding.EncodeToString([]byte(v))
}

func builtin_base64Unpadded(v string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(v))
}

func builtin_base64UrlUnpadded(v string) string {
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
		return now.Format(dateFormat(formatOpt[0]))
	}

	return strconv.Itoa(int(now.Unix()))
}

func builtin_isset(m map[string]any, key string) bool {
	v, ok := m[key]
	return ok && v != nil
}

func builtin_json(v any, indent ...string) string {
	var (
		err error
		res []byte
	)

	if len(indent) > 0 {
		res, err = json.MarshalIndent(v, "", indent[0])
	} else {
		res, err = json.Marshal(v)
	}

	if err != nil {
		panic("failed encoding json: " + err.Error())
	}

	return string(res)
}

func builtin_formatTimestamp(v any, format ...string) string {
	var t time.Time

	switch tv := v.(type) {
	case time.Time:
		t = tv
	case string:
		if len(format) == 0 {
			panic("input format is missing")
		}
		var err error
		t, err = time.Parse(format[0], tv)
		if err != nil {
			panic("failed parsing input date: " + err.Error())
		}
		format = format[1:]
	}

	if len(format) > 0 {
		return t.Format(dateFormat(format[0]))
	}

	return strconv.Itoa(int(t.Unix()))
}
