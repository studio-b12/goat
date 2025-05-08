package goatfile

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/studio-b12/goat/pkg/engine"
	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/goatfile/ast"
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

	parsed    bool
	preParsed bool
}

var _ Action = (*Request)(nil)

func newRequest() (t *Request) {
	t = new(Request)
	t.Header = http.Header{}
	t.Body = NoContent{}
	t.PreScript = NoContent{}
	t.Script = NoContent{}
	return t
}

func RequestFromAst(req *ast.Request, path string) (t *Request, err error) {
	if req == nil {
		return &Request{}, errors.New("request ast is nil")
	}

	t = newRequest()

	t.Path = path
	t.Method = req.Head.Method
	t.URI = req.Head.Url
	t.PosLine = req.Pos.Line + 1 // TODO: actually, this should start counting at 0 and the printer should add 1

	var additionalHeader http.Header

	for _, block := range req.Blocks {
		switch b := block.(type) {
		case ast.RequestHeader:
			t.Header = b.HeaderEntries.ToMultiMap()
		case ast.RequestOptions:
			t.Options = b.KVList.ToMap()
		case ast.RequestQueryParams:
			t.QueryParams = b.KVList.ToMap()
		case ast.RequestAuth:
			t.Auth = b.KVList.ToMap()
		case ast.RequestBody:
			t.Body, additionalHeader, err = DataFromAst(b.DataContent, path)
		case ast.RequestPreScript:
			t.PreScript, _, err = DataFromAst(b.DataContent, path)
		case ast.RequestScript:
			t.Script, _, err = DataFromAst(b.DataContent, path)
		case ast.FormData:
			t.Body, additionalHeader, err = DataFromAst(b, path)
		case ast.FormUrlEncoded:
			t.Body, additionalHeader, err = DataFromAst(b, path)
		default:
			err = fmt.Errorf("invalid request ast block type: %+v", block)
		}
	}

	for k, v := range additionalHeader {
		t.Header[k] = append(t.Header[k], v...)
	}

	if err != nil {
		return &Request{}, err
	}

	return t, nil
}

func PartialRequestFromAst(req ast.PartialRequest, path string) (t *Request, err error) {
	var fullReq ast.Request

	fullReq.Pos = req.Pos
	fullReq.Blocks = req.Blocks

	return RequestFromAst(&fullReq, path)
}

func (t *Request) Type() ActionType {
	return ActionRequest
}

// PreSubstituteWithParams takes the given parameters and replaces placeholders
// within specific parts of the request which shall be executed before the
// actual request is substituted (like PreScript).
func (t *Request) PreSubstituteWithParams(params any) error {
	if t.preParsed {
		return ErrTemplateAlreadyPreParsed
	}

	// Substitute PreScript

	preScriptStr, err := util.ReadReaderToString(t.PreScript.Reader())
	if err != nil {
		return errs.WithPrefix("reading preScript failed:", err)
	}

	preScriptStr, err = ApplyTemplate(preScriptStr, params)
	if err != nil {
		return err
	}
	t.PreScript = StringContent(preScriptStr)

	t.preParsed = true
	return nil
}

// SubstituteWithParams takes the given parameters
// and replaces placeholders within the request
// with values from the given params.
func (t *Request) SubstituteWithParams(params any) error {
	if t.parsed {
		return ErrTemplateAlreadyParsed
	}

	var err error

	// Substitute Options

	err = ApplyTemplateToMap(t.Options, params)
	if err != nil {
		return err
	}

	if v, ok := t.Options[conditionOptionName].(bool); ok && !v {
		return nil
	}

	// Substitute URI

	t.URI, err = ApplyTemplate(t.URI, params)
	if err != nil {
		return err
	}

	// Substitute QueryParams

	err = ApplyTemplateToMap(t.QueryParams, params)
	if err != nil {
		return err
	}

	// Substitute Auth

	err = ApplyTemplateToMap(t.Auth, params)
	if err != nil {
		return err
	}

	// Substitute Header

	for _, vals := range t.Header {
		for i, v := range vals {
			vals[i], err = ApplyTemplate(v, params)
			if err != nil {
				return err
			}
		}
	}

	// Substitute Body

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
	case FormData:
		err = ApplyTemplateToMap(body.fields, params)
		if err != nil {
			return err
		}
		t.Body = body
	case FormUrlEncoded:
		err = ApplyTemplateToMap(body.fields, params)
		if err != nil {
			return err
		}
		t.Body = body
	}

	// Substitute Script

	scriptStr, err := util.ReadReaderToString(t.Script.Reader())
	if err != nil {
		return errs.WithPrefix("reading script failed:", err)
	}

	scriptStr, err = ApplyTemplate(scriptStr, params)
	if err != nil {
		return err
	}
	t.Script = StringContent(scriptStr)

	return nil
}

// InsertRawDataIntoBody evaluates the raw bytes required for the request
func (t *Request) InsertRawDataIntoBody(state engine.State) error {
	body, ok := t.Body.(RawContent)
	if !ok {
		return nil
	}
	v, ok := state[body.varName]
	if !ok {
		return ErrVarNotFound
	}

	rv := util.UnwrapPointer(reflect.ValueOf(v))
	if rv.Kind() != reflect.Slice || rv.Type().Elem().Kind() != reflect.Uint8 {
		return errs.WithPrefix(fmt.Sprintf("$%v :", body.varName), ErrNotAByteArray)
	}

	body.value = rv.Bytes()
	t.Body = body
	return nil
}

// InsertRawDataIntoFormData evaluates the raw bytes required for the request
func (t *Request) InsertRawDataIntoFormData(state engine.State) error {
	body, ok := t.Body.(FormData)
	if !ok {
		return nil
	}
	for k, v := range body.fields {
		if vd, ok := v.(ast.RawDescriptor); ok {
			valeFromState, ok := state[vd.VarName]
			if !ok {
				return ErrVarNotFound
			}
			rv := util.UnwrapPointer(reflect.ValueOf(valeFromState))
			if rv.Kind() != reflect.Slice || rv.Type().Elem().Kind() != reflect.Uint8 {
				return errs.WithPrefix(fmt.Sprintf("$%v :", vd.VarName), ErrNotAByteArray)
			}
			vd.Data = rv.Bytes()
			body.fields[k] = vd
		}
	}

	t.Body = body
	return nil
}

// ToHttpRequest returns a *http.Request built from the
// given Reuqest.
func (t *Request) ToHttpRequest() (*http.Request, error) {
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

	if len(with.Auth) > 0 {
		t.Auth = mergeMaps(t.Auth, with.Auth)
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

func (t *Request) String() string {
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
