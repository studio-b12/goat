// Package goatfile provides functionalities to
// unmarshal and parse a Goatfile.
//
// Here you can find the Goatfile specification
// on which basis this parser in built on.
// https://github.com/studio-b12/goat/blob/main/docs/goatfile-spec.md
package goatfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	sectionNameSetup        = "setup"
	sectionNameSetupEach    = "setup-each"
	sectionNameTests        = "tests"
	sectionNameTeardown     = "teardown"
	sectionNameTeardownEach = "teardown-each"
)

const (
	optionNameQueryParams = "queryparams"
	optionNameHeader      = "header"
	optionNameHeaders     = "headers"
	optionNameBody        = "body"
	optionNameScript      = "script"
	optionNameOptions     = "options"
	abc                   = "asd"
)

// Goatfile holds all sections and
// their requests.
type Goatfile struct {
	Imports []string

	Setup        []Request
	SetupEach    []Request
	Tests        []Request
	Teardown     []Request
	TeardownEach []Request

	Path string
}

// Merge appends all requests in all sections of with
// to the current Goatfile.
func (t *Goatfile) Merge(with Goatfile) {
	t.Setup = append(t.Setup, with.Setup...)
	t.SetupEach = append(t.SetupEach, with.SetupEach...)
	t.Tests = append(t.Tests, with.Tests...)
	t.Teardown = append(t.Teardown, with.Teardown...)
	t.TeardownEach = append(t.TeardownEach, with.TeardownEach...)
}

// String returns the Goatfile as JSON encoded string.
func (t Goatfile) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}
	return string(b)
}

// Opts holds the specific request
// options.
type Opts struct {
	QueryParams map[string]any
	Options     map[string]any
}

// Request holds the specifications
// for a HTTP request with options
// and script commands.
type Request struct {
	Opts

	Method string
	URI    string
	Header http.Header
	Body   []byte
	Script string

	parsed bool
}

func newRequest() (r Request) {
	r.Header = http.Header{}
	return r
}

// ParseWithParams takes the given parameters
// and replaces placeholders within the request
// with values from the given params.
func (t *Request) ParseWithParams(params any) error {
	if t.parsed {
		return ErrTemplateAlreadyParsed
	}

	var err error

	t.URI, err = applyTemplate(t.URI, params)
	if err != nil {
		return err
	}

	for _, vals := range t.Header {
		for i, v := range vals {
			vals[i], err = applyTemplate(v, params)
			if err != nil {
				return err
			}
		}
	}

	bodyStr, err := applyTemplate(string(t.Body), params)
	if err != nil {
		return err
	}
	t.Body = []byte(bodyStr)

	t.Script, err = applyTemplate(t.Script, params)
	if err != nil {
		return err
	}

	applyTemplateToMap(t.QueryParams, params)
	applyTemplateToMap(t.Options, params)

	return nil
}

// ToHttpRequest returns a *http.Request built from the
// given Reuqest.
func (t Request) ToHttpRequest() (*http.Request, error) {
	uri, err := url.Parse(t.URI)
	if err != nil {
		return nil, fmt.Errorf("failed parsing URI: %s", err.Error())
	}

	query := uri.Query()

	for key, val := range t.Opts.QueryParams {
		if arr, ok := val.([]any); ok {
			for _, v := range arr {
				query.Add(key, toString(v))
			}
		} else {
			query.Add(key, toString(val))
		}
	}

	uri.RawQuery = query.Encode()

	var body io.Reader

	if len(t.Body) > 0 {
		body = bytes.NewBuffer(t.Body)
	}

	req, err := http.NewRequest(t.Method, uri.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header = t.Header

	return req, nil
}

func (t Request) String() string {
	return fmt.Sprintf("%s %s", t.Method, t.URI)
}

func toString(v any) string {
	return fmt.Sprintf("%v", v)
}