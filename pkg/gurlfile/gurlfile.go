package gurlfile

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	sectionNameSetup        = "setup"
	sectionNameSetupEach    = "setup-each"
	sectionnameTests        = "tests"
	sectionNameTeardown     = "teardown"
	sectionNameTeardownEach = "teardown-each"
)

const (
	optionNameQueryParams = "queryparams"
)

// Gurlfile holds all sections and
// their requests.
type Gurlfile struct {
	Imports []string

	Setup        []Request
	SetupEach    []Request
	Tests        []Request
	Teardown     []Request
	TeardownEach []Request
}

// Merge appends all requests in all sections of with
// to the current Gurlfile.
func (t *Gurlfile) Merge(with Gurlfile) {
	t.Setup = append(t.Setup, with.Setup...)
	t.SetupEach = append(t.SetupEach, with.SetupEach...)
	t.Tests = append(t.Tests, with.Tests...)
	t.Teardown = append(t.Teardown, with.Teardown...)
	t.TeardownEach = append(t.TeardownEach, with.TeardownEach...)
}

// Options holds the specific request
// options.
type Options struct {
	QueryParams map[string]any
}

type context struct {
	raw     string
	section string
	index   int
}

func (t context) WrapErr(err error) error {
	var cErr ContextError

	cErr.context = t
	cErr.Inner = err

	return cErr
}

// Request holds the specifications
// for a HTTP request with options
// and script commands.
type Request struct {
	Options

	context

	Method string
	URI    string
	Header http.Header
	Body   []byte
	Script string
}

func newRequest() (r Request) {
	r.Header = http.Header{}
	return r
}

// ParseWithParams takes the given parameters
// and replaces placeholders within the request
// with values from the given params.
//
// Returns the new parsed request.
func (t Request) ParseWithParams(params any) (Request, error) {
	if t.raw == "" {
		return Request{}, ErrAlreadyParsed
	}

	return Request{}, nil
}

func (t Request) ToHttpRequest() (*http.Request, error) {
	uri, err := url.Parse(t.URI)
	if err != nil {
		return nil, fmt.Errorf("failed parsing URI: %s", err.Error())
	}

	for key, val := range t.Options.QueryParams {
		if arr, ok := val.([]any); ok {
			for _, v := range arr {
				uri.Query().Add(key, toString(v))
			}
		} else {
			uri.Query().Add(key, toString(val))
		}
	}

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
