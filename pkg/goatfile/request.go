package goatfile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/util"
)

// Request holds the specifications
// for a HTTP request with options
// and script commands.
type Request struct {
	Opts

	Method    string
	URI       string
	Header    http.Header
	Body      Data
	PreScript Data
	Script    Data

	parsed bool
}

var _ Action = (*Request)(nil)

func newRequest() (r Request) {
	r.Header = http.Header{}
	r.Body = NoContent{}
	r.PreScript = NoContent{}
	r.Script = NoContent{}
	return r
}

func (t Request) Type() ActionType {
	return ActionRequest
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

	if strData, ok := t.Body.(StringContent); ok {
		bodyStr, err := applyTemplate(string(strData), params)
		if err != nil {
			return err
		}
		t.Body = StringContent(bodyStr)
	}

	scriptStr, err := util.ReadReaderToString(t.Script.Reader())
	if err != nil {
		return errs.WithPrefix("reading script failed:", err)
	}

	scriptStr, err = applyTemplate(scriptStr, params)
	if err != nil {
		return err
	}
	t.Script = StringContent(scriptStr)

	applyTemplateToMap(t.QueryParams, params)
	applyTemplateToMap(t.Options, params)

	return nil
}

// ToHttpRequest returns a *http.Request built from the
// given Reuqest.
func (t Request) ToHttpRequest() (*http.Request, error) {
	uri, err := url.Parse(t.URI)
	if err != nil {
		return nil, errs.WithPrefix("failed parsing URI:", err)
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

	bodyReader, err := t.Body.Reader()
	if err != nil {
		return nil, errs.WithPrefix("failed reading body data:", err)
	}

	if bodyReader != nil {
		body = bodyReader
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

type requestParseChecker struct {
	*Request

	set map[optionName]struct{}
}

func wrapIntoRequestParseChecker(req *Request) *requestParseChecker {
	return &requestParseChecker{
		Request: req,
		set:     make(map[optionName]struct{}),
	}
}

func (t *requestParseChecker) Check(opt optionName) error {
	if _, ok := t.set[opt]; ok {
		return errs.WithPrefix(fmt.Sprintf("[%s]:", opt), ErrSectionDefinedMultiple)
	}
	t.set[opt] = struct{}{}
	return nil
}
