package goatfile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/util"
)

const conditionOptionName = "condition"

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

	Path    string
	PosLine int

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

	err = ApplyTemplateToMap(t.Options, params)
	if err != nil {
		return err
	}

	if v, ok := t.Options[conditionOptionName].(bool); ok && !v {
		return nil
	}

	t.URI, err = ApplyTemplate(t.URI, params)
	if err != nil {
		return err
	}

	for _, vals := range t.Header {
		for i, v := range vals {
			vals[i], err = ApplyTemplate(v, params)
			if err != nil {
				return err
			}
		}
	}

	switch body := t.Body.(type) {
	case StringContent:
		bodyStr, err := ApplyTemplate(string(body), params)
		if err != nil {
			return err
		}
		t.Body = StringContent(bodyStr)
	case FileContent:
		body.filePath, err = ApplyTemplate(body.filePath, params)
		if err != nil {
			return err
		}
		t.Body = body
	}

	scriptStr, err := util.ReadReaderToString(t.Script.Reader())
	if err != nil {
		return errs.WithPrefix("reading script failed:", err)
	}

	scriptStr, err = ApplyTemplate(scriptStr, params)
	if err != nil {
		return err
	}
	t.Script = StringContent(scriptStr)

	err = ApplyTemplateToMap(t.QueryParams, params)
	if err != nil {
		return err
	}

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

func (t *Request) Merge(with *Request) {
	if t == nil || with == nil {
		return
	}

	if len(with.Header) > 0 {
		newHeaders := t.Header.Clone()
		for key, vals := range with.Header {
			for _, val := range vals {
				newHeaders.Add(key, val)
			}
		}
		t.Header = newHeaders
	}

	if len(with.QueryParams) > 0 {
		t.QueryParams = mergeMaps(t.QueryParams, with.QueryParams)
	}

	if len(with.Options) > 0 {
		t.Options = mergeMaps(t.Options, with.Options)
	}

	if IsNoContent(t.Body) && !IsNoContent(with.Body) {
		t.Body = with.Body
	}

	if IsNoContent(t.PreScript) && !IsNoContent(with.PreScript) {
		t.PreScript = with.PreScript
	}

	if IsNoContent(t.Script) && !IsNoContent(with.Script) {
		t.Script = with.Script
	}
}

func (t Request) String() string {
	return fmt.Sprintf("%s %s", t.Method, t.URI)
}

func toString(v any) string {
	return fmt.Sprintf("%v", v)
}

func mergeMaps[TK comparable, TV any](src, base map[TK]TV) map[TK]TV {
	new := map[TK]TV{}
	for key, val := range base {
		new[key] = val
	}
	for key, val := range src {
		new[key] = val
	}
	return new
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
