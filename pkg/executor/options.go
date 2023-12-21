package executor

import (
	"encoding/base64"
	"fmt"
	"time"
)

// AbortOptions wraps options that control the
// abort behavior of an execution batch.
type AbortOptions struct {
	NoAbort     bool
	AlwaysAbort bool
}

// AbortOptionsFromMap returns a new instance of
// AbortOptions extracted from the passed map.
func AbortOptionsFromMap(m map[string]any) AbortOptions {
	opt := AbortOptions{
		NoAbort:     false,
		AlwaysAbort: false,
	}

	if v, ok := m["noabort"].(bool); ok {
		opt.NoAbort = v
	}

	if v, ok := m["alwaysabort"].(bool); ok {
		opt.AlwaysAbort = v
	}

	return opt
}

// ExecOptions wraps options that control the
// execution of a request.
type ExecOptions struct {
	Condition bool
	Delay     time.Duration
}

// ExecOptionsFromMap returns a new instance of
// ExecOptions extracted from the passed map.
func ExecOptionsFromMap(m map[string]any) ExecOptions {
	opt := ExecOptions{
		Condition: true,
	}

	if v, ok := m["condition"].(bool); ok {
		opt.Condition = v
	}

	v, ok := m["delay"]
	if ok {
		switch vt := v.(type) {
		case int:
			opt.Delay = time.Duration(vt) * time.Millisecond
		case string:
			opt.Delay, _ = time.ParseDuration(vt)
		}
	}

	return opt
}

type AuthOptions struct {
	Type     string
	UserName string
	Password string
	Token    string
}

func AuthOptionsFromMap(m map[string]any) (opt AuthOptions, ok bool) {
	if m == nil {
		return opt, false
	}

	if v, ok := m["type"].(string); ok {
		opt.Type = v
	}

	if v, ok := m["username"].(string); ok {
		opt.UserName = v
	}

	if v, ok := m["password"].(string); ok {
		opt.Password = v
	}

	if v, ok := m["token"].(string); ok {
		opt.Token = v
	}

	return opt, true
}

func (t AuthOptions) HeaderValue() string {
	if t.UserName != "" && t.Password != "" {
		v := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", t.UserName, t.Password)))
		return fmt.Sprintf("basic %s", v)
	}

	if t.Type != "" {
		return fmt.Sprintf("%s %s", t.Type, t.Token)
	}

	return t.Token
}
