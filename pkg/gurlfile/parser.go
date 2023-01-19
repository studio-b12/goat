package gurlfile

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Parser struct {
	s   *scanner
	buf struct {
		tok token  // last read token
		lit string // last read literal
		n   int    // buffer size
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: newScanner(r)}
}

func (t *Parser) Parse() (gf Gurlfile, err error) {
	defer func() {
		err = t.wrapErr(err)
	}()

	for {
		tok, lit := t.scan()
		_ = lit

		switch tok {
		case COMMENT, WS, LF:
			continue

		case IDENT, STRING:
			t.unscan()
			err = t.parseRequest(&gf.Tests)

		case USE:
			err = t.parseUse(&gf)

		case SECTION:
			err = t.parseSection(lit, &gf)
		case EOF:
			return gf, nil

		case ILLEGAL:
			return Gurlfile{}, fmt.Errorf("illegal character")
		default:
			err = fmt.Errorf("unexpected error (%d '%s')", tok, lit)
		}

		if err != nil {
			return Gurlfile{}, err
		}
	}
}

func (t *Parser) scan() (tok token, lit string) {
	if t.buf.n != 0 {
		t.buf.n = 0
		return t.buf.tok, t.buf.lit
	}

	t.buf.tok, t.buf.lit = t.s.Scan()

	return t.buf.tok, t.buf.lit
}

func (t *Parser) scanSkipWS() (tok token, lit string) {
	tok, lit = t.scan()
	if tok == WS {
		return t.scan()
	}

	return tok, lit
}

func (t *Parser) unscan() {
	t.buf.n = 1
}

func (t *Parser) wrapErr(err error) error {
	if err == nil {
		return nil
	}

	pErr := ParseError{}
	pErr.Inner = err
	pErr.Line = t.s.line
	pErr.LinePos = t.s.linepos

	return pErr
}

func (t *Parser) parseUse(gf *Gurlfile) error {
	tk, lit := t.s.scanString()
	if tk == ILLEGAL {
		return fmt.Errorf("invalid string literal")
	}

	if lit == "" {
		return fmt.Errorf("empty use path")
	}

	gf.Imports = append(gf.Imports, lit)

	return nil
}

func (t *Parser) parseSection(name string, gf *Gurlfile) error {
	var r *[]Request

	switch strings.ToLower(strings.TrimSpace(name)) {
	case "setup":
		r = &gf.Setup
	case "setup-each":
		r = &gf.SetupEach
	case "tests":
		r = &gf.Tests
	case "teardown":
		r = &gf.Teardown
	case "teardown-each":
		r = &gf.TeardownEach
	default:
		return fmt.Errorf("invalid section")
	}

	for {
		tok, _ := t.scan()
		if tok == LF || tok == WS {
			continue
		}

		if tok == EOF || tok == SECTION {
			t.unscan()
			break
		}

		t.unscan()
		err := t.parseRequest(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Parser) parseRequest(section *[]Request) (err error) {
	req := newRequest()

	// parse header

	tok, lit := t.scan()
	if tok != IDENT && tok != STRING || lit == "" {
		return fmt.Errorf("invalid request method")
	}
	req.Method = lit

	tok, lit = t.scan()
	if tok != WS && tok != LF {
		return fmt.Errorf("method must be followed by the request URI")
	}

	tok, lit = t.s.scanString()
	if tok != STRING || lit == "" {
		return fmt.Errorf("method must be followed by the request URI")
	}
	req.URI = lit

loop:
	for {
		tok, lit = t.scan()

		switch tok {
		case BLOCK_START:
			err = t.parseBlock(&req)

		case WS, LF:
			continue loop
		case EOF:
			t.unscan()
			break loop
		case DELIMITER:
			break loop

		default:
			err = fmt.Errorf("invalid token (request)")
		}

		if err != nil {
			return err
		}
	}

	*section = append(*section, req)
	return nil
}

func (t *Parser) parseBlock(req *Request) error {
	var blockHeader string

	tok, lit := t.scan()
	if tok != IDENT || lit == "" {
		return fmt.Errorf("invalid block header")
	}
	blockHeader = lit

	tok, _ = t.scan()
	if tok != BLOCK_END {
		return fmt.Errorf("invalid block header")
	}

	tok, _ = t.scan()
	if tok != LF {
		return fmt.Errorf("invalid token (block)")
	}

	switch strings.ToLower(blockHeader) {

	case "queryparams":
		data, err := t.parseBlockEntries()
		if err != nil {
			return err
		}
		req.QueryParams = data

	case "header", "headers":
		err := t.parseHeaders(req.Header)
		if err != nil {
			return err
		}

	case "body":
		raw, err := t.parseRaw()
		if err != nil {
			return err
		}
		req.Body = []byte(raw)

	case "script":
		raw, err := t.parseRaw()
		if err != nil {
			return err
		}
		req.Script = raw

	default:
		return fmt.Errorf("invalid block header '%s'", blockHeader)
	}

	return nil
}

func (t *Parser) parseBlockEntries() (map[string]any, error) {
	m := map[string]any{}

	for {
		tok, lit := t.scanSkipWS()
		if tok == LF {
			continue
		}
		if tok == DELIMITER || tok == EOF || tok == BLOCK_START || tok == SECTION {
			t.unscan()
			break
		}

		if tok != IDENT {
			return nil, fmt.Errorf("block entry must start with an assignment")
		}

		key := lit

		tok, lit = t.scanSkipWS()
		if tok != ASSIGNMENT {
			return nil, fmt.Errorf("block entry key must be followed by an assignment")
		}

		val, err := t.parseValue()
		if err != nil {
			return nil, err
		}

		m[key] = val
	}

	return m, nil
}

func (t *Parser) parseHeaders(header http.Header) error {
	for {
		tok, lit := t.scanSkipWS()
		if tok == LF {
			continue
		}
		if tok == DELIMITER || tok == EOF || tok == BLOCK_START || tok == SECTION {
			t.unscan()
			break
		}

		if tok != IDENT {
			return fmt.Errorf("header values must start with a key")
		}
		key := lit

		tok, _ = t.scanSkipWS()
		if tok != COLON {
			return fmt.Errorf("header key and value must be separated by a colon (:)")
		}

		val := strings.TrimSpace(t.s.scanUntilLF())
		header.Add(key, val)
	}

	return nil
}

func (t *Parser) parseRaw() (string, error) {
	var out bytes.Buffer

	inEscape := false

	for {
		r := t.s.read()

		if !inEscape {
			if out.Len() > 3 && string(out.Bytes()[out.Len()-4:]) == "\n---" {
				t.buf.tok = DELIMITER
				t.buf.lit = ""
				t.unscan()
				out.Truncate(out.Len() - 4)
				break
			}
			if out.Len() > 0 && out.Bytes()[out.Len()-1] == '[' {
				t.buf.tok = BLOCK_START
				t.buf.lit = ""
				t.s.unread()
				t.unscan()
				out.Truncate(out.Len() - 1)
				break
			}
		}

		out.WriteRune(r)

		if out.Len() == 4 && out.String() == "```\n" ||
			out.Len() > 3 && string(out.Bytes()[out.Len()-4:]) == "\n```" {
			if inEscape {
				inEscape = false
			} else {
				inEscape = true
			}
			if out.Len() == 3 {
				out.Reset()
			} else {
				out.Truncate(out.Len() - 4)
			}
			continue
		}

	}

	return out.String(), nil
}

func (t *Parser) parseValue() (any, error) {
	tok, lit := t.scanSkipWS()

	switch tok {
	case INTEGER:
		return strconv.ParseInt(lit, 10, 64)
	case FLOAT:
		return strconv.ParseFloat(lit, 64)
	case STRING:
		return lit, nil
	case BLOCK_START:
		return t.parseArray()
	}

	return nil, fmt.Errorf("invalid token (value)")
}

func (t *Parser) parseArray() ([]any, error) {
	var arr []any

loop:
	for {
		val, err := t.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)

		tok, _ := t.scanSkipWS()
		switch tok {
		case BLOCK_END:
			break loop
		case COMMA:
			continue loop
		default:
			return nil, fmt.Errorf("invalid token (value array)")
		}
	}

	return arr, nil
}
