package goatfile

import (
	"bytes"
	"fmt"
	"github.com/studio-b12/goat/pkg/goatfile/ast"
	"io"
	"strconv"
	"strings"

	"github.com/studio-b12/goat/pkg/errs"
)

// Parser parses a Goatfile.
type Parser struct {
	fileDir string
	s       *scanner
	prevPos readerPos
	buf     struct {
		tok token  // last read token
		lit string // last read literal
		n   int    // buffer size
	}
}

// NewParser returns a new Parser scanning from
// the given Reader.
func NewParser(r io.Reader, fileDir string) *Parser {
	return &Parser{
		fileDir: fileDir,
		s:       newScanner(r),
	}
}

// Parse parses a Goatfile from the specified source.
func (t *Parser) Parse() (*ast.Goatfile, error) {
	var (
		err error
		gf  ast.Goatfile
	)

	defer func() {
		err = t.wrapErr(err)
	}()

	gf.Dir = t.fileDir

	for {
		tok, lit := t.scan()
		_ = lit

		switch tok {
		case tokCOMMENT, tokWS, tokLF:
			// TODO: set in ast
			continue

		case tokDELIMITER:
			continue

		case tokLOGSECTION:
			pos := t.astPos()
			sec := t.s.scanUntilLF()
			gf.Actions = append(gf.Actions, ast.LogSection{
				Pos:     pos,
				Content: strings.TrimSpace(sec),
			})
			continue

		case tokIDENT, tokSTRING:
			t.unscan()
			req, err := t.parseRequest()
			if err != nil {
				return nil, err
			}
			gf.Actions = append(gf.Actions, req)

		case tokEXECUTE:
			exec, err := t.parseExecute()
			if err != nil {
				return nil, err
			}
			gf.Actions = append(gf.Actions, exec)

		case tokUSE:
			use, err := t.parseUse()
			if err != nil {
				return nil, err
			}
			gf.Imports = append(gf.Imports, *use)

		case tokSECTION:
			sec, err := t.parseSection()
			if err != nil {
				return nil, err
			}
			gf.Sections = append(gf.Sections, sec)
		case tokEOF:
			return &gf, nil

		case tokBLOCKSTART:
			return nil, ErrBlockOutOfRequest

		default:
			return nil, ErrIllegalCharacter
		}
	}
}

func (t *Parser) scan() (tok token, lit string) {
	if t.buf.n != 0 {
		t.buf.n = 0
		return t.buf.tok, t.buf.lit
	}

	t.prevPos = t.s.readerPos

	t.buf.tok, t.buf.lit = t.s.scan()
	if t.buf.tok == tokCOMMENT {
		t.buf.tok = tokLF
		t.buf.lit = ""
	}

	return t.buf.tok, t.buf.lit
}

func (t *Parser) unscan() {
	t.buf.n = 1
}

func (t *Parser) scanSkipWS() (tok token, lit string) {
	tok, lit = t.scan()
	if tok == tokWS {
		return t.scan()
	}

	return tok, lit
}

func (t *Parser) wrapErr(err error) error {
	if err == nil {
		return nil
	}

	pErr := ParseError{}
	pErr.Inner = err
	pErr.Line = t.prevPos.line
	pErr.LinePos = t.prevPos.linepos

	return pErr
}

func (t *Parser) parseUse() (*ast.Import, error) {
	pos := t.astPos()

	tk, _ := t.scan()
	if tk != tokWS {
		return nil, ErrInvalidStringLiteral
	}

	tk, lit := t.s.scanString()
	if tk == tokILLEGAL {
		return nil, ErrInvalidStringLiteral
	}

	if lit == "" {
		return nil, ErrEmptyUsePath
	}

	return &ast.Import{
		Pos:  pos,
		Path: lit,
	}, nil
}

func (t *Parser) parseExecute() (*ast.Execute, error) {
	tk, _ := t.scan()
	if tk != tokWS {
		return nil, ErrInvalidStringLiteral
	}

	tk, lit := t.s.scanString()
	if tk == tokILLEGAL {
		return nil, ErrInvalidStringLiteral
	}

	if lit == "" {
		return nil, ErrEmptyCallPath
	}

	var (
		exec ast.Execute
		err  error
	)

	exec.Path = lit

	tok, _ := t.scanSkipWS()
	if tok != tokGROUPSTART {
		t.unscan()
		return &exec, nil
	}

	end := tokGROUPEND
	exec.Parameters, err = t.parseBlockEntries(&end)
	if err != nil {
		return nil, err
	}
	t.scan() // re-scan closing group `)`

	tok, _ = t.scanSkipWS()
	if tok != tokRETURN {
		t.unscan()
		return &exec, nil
	}

	tok, _ = t.scanSkipWS()
	if tok != tokGROUPSTART {
		return nil, ErrMissingGroup
	}

	exec.Returns = make(ast.Assignments)
	for {
		tok, key := t.scanSkipWS()
		if tok == tokEOF {
			return nil, ErrUnclosedGroup
		}
		if tok == tokLF {
			continue
		}
		if tok == tokGROUPEND {
			break
		}
		if tok != tokIDENT {
			return nil, ErrIllegalCharacter
		}

		tok, _ = t.scanSkipWS()
		if tok != tokAS {
			return nil, ErrIllegalCharacter
		}

		tok, val := t.scanSkipWS()
		if tok != tokIDENT {
			return nil, ErrIllegalCharacter
		}

		exec.Returns[key] = val
	}

	return &exec, nil
}

func (t *Parser) parseSection() (ast.Section, error) {
	sectionPos := t.astPos()

	name := SectionName(strings.ToLower(strings.TrimSpace(t.s.readToLF())))

	if name == SectionDefaults {
		pr, err := t.parseDefaults()
		if err != nil {
			return nil, err
		}
		return ast.SectionDefaults{Pos: sectionPos, Request: *pr}, nil
	}

	var actions []ast.Action

	for {
		tok, _ := t.scan()
		if tok == tokLF || tok == tokWS || tok == tokDELIMITER {
			continue
		}

		if tok == tokEOF || tok == tokSECTION {
			t.unscan()
			break
		}

		if tok == tokLOGSECTION {
			pos := t.astPos()
			sec := t.s.scanUntilLF()
			actions = append(actions, ast.LogSection{
				Pos:     pos,
				Content: strings.TrimSpace(sec),
			})
			continue
		}

		if tok == tokEXECUTE {
			exec, err := t.parseExecute()
			if err != nil {
				return nil, err
			}
			actions = append(actions, exec)
			continue
		}

		t.unscan()
		req, err := t.parseRequest()
		if err != nil {
			return nil, err
		}
		actions = append(actions, req)
	}

	switch name {
	case SectionSetup:
		return ast.SectionSetup{
			Pos:     sectionPos,
			Actions: actions,
		}, nil
	case SectionTests:
		return ast.SectionTests{
			Pos:     sectionPos,
			Actions: actions,
		}, nil
	case SectionTeardown:
		return ast.SectionTeardown{
			Pos:     sectionPos,
			Actions: actions,
		}, nil
	default:
		return nil, ErrInvalidSection
	}
}

func (t *Parser) parseRequest() (*ast.Request, error) {
	var (
		req ast.Request
	)

	// parse header

	req.Pos = t.astPos()

	tok, lit := t.scan()
	if tok != tokIDENT && tok != tokSTRING || lit == "" {
		return nil, ErrInvalidRequestMethod
	}
	req.Head.Method = lit

	tok, _ = t.scan()
	if tok != tokWS && tok != tokLF {
		return nil, ErrNoRequestURI
	}

	tok, lit = t.s.scanString()
	if tok != tokSTRING || lit == "" {
		return nil, ErrNoRequestURI
	}
	req.Head.Url = lit

loop:
	for {
		tok, _ = t.scan()

		switch tok {
		case tokBLOCKSTART:
			block, err := t.parseBlock()
			if err != nil {
				return nil, err
			}
			req.Blocks = append(req.Blocks, block)

		case tokWS, tokLF:
			continue loop
		case tokEOF, tokSECTION, tokLOGSECTION:
			t.unscan()
			break loop
		case tokDELIMITER:
			break loop

		default:
			return nil, errs.WithSuffix(ErrInvalidToken, "(request)")
		}
	}

	return &req, nil
}

func (t *Parser) parseDefaults() (*ast.PartialRequest, error) {
	var (
		req ast.PartialRequest
		err error
	)

	req.Pos = t.astPos()

loop:
	for {
		tok, _ := t.scan()

		switch tok {
		case tokBLOCKSTART:
			block, err := t.parseBlock()
			if err != nil {
				return nil, err
			}
			req.Blocks = append(req.Blocks, block)

		case tokWS, tokLF:
			continue loop
		case tokEOF, tokSECTION, tokLOGSECTION:
			t.unscan()
			break loop
		case tokDELIMITER:
			break loop

		default:
			err = errs.WithSuffix(ErrInvalidToken, "(request)")
		}

		if err != nil {
			return nil, err
		}
	}

	return &req, nil
}

func (t *Parser) parseBlock() (ast.RequestBlock, error) {
	var blockHeader string

	tok, lit := t.scanSkipWS()
	if tok != tokIDENT || lit == "" {
		return nil, ErrInvalidBlockHeader
	}
	blockHeader = lit

	tok, _ = t.scan()
	if tok != tokBLOCKEND {
		return nil, ErrInvalidBlockHeader
	}

	tok, _ = t.scanSkipWS()
	if tok != tokLF {
		return nil, errs.WithSuffix(ErrInvalidToken, "(block)")
	}

	optName := optionName(strings.ToLower(blockHeader))

	switch optName {

	case optionNameQueryParams:
		data, err := t.parseBlockEntries(nil)
		if err != nil {
			return nil, err
		}
		return ast.RequestQueryParams{KV: data}, nil

	case optionNameHeader:
		header, err := t.parseHeaders()
		if err != nil {
			return nil, err
		}
		return ast.RequestHeader{HeaderEntries: header}, nil

	case optionNameBody:
		raw, err := t.parseRaw()
		if err != nil {
			return nil, err
		}
		return ast.RequestBody{DataContent: raw}, nil

	case optionNamePreScript:
		raw, err := t.parseRaw()
		if err != nil {
			return nil, err
		}
		return ast.RequestPreScript{DataContent: raw}, nil

	case optionNameScript:
		raw, err := t.parseRaw()
		if err != nil {
			return nil, err
		}
		return ast.RequestScript{DataContent: raw}, nil

	case optionNameOptions:
		data, err := t.parseBlockEntries(nil)
		if err != nil {
			return nil, err
		}
		return ast.RequestOptions{KV: data}, nil

	case optionNameAuth:
		data, err := t.parseBlockEntries(nil)
		if err != nil {
			return nil, err
		}
		return ast.RequestAuth{KV: data}, nil

	default:
		return nil, errs.WithSuffix(ErrInvalidBlockHeader,
			fmt.Sprintf("('%s')", blockHeader))
	}
}

func (t *Parser) parseBlockEntries(exitToken *token) (ast.KV, error) {
	m := make(ast.KV)

	for {
		tok, lit := t.scanSkipWS()
		if tok == tokLF {
			continue
		}
		if tok == tokEOF ||
			(exitToken != nil && tok == *exitToken) ||
			(tok == tokDELIMITER || tok == tokBLOCKSTART || tok == tokSECTION) {

			t.unscan()
			break
		}

		if tok != tokIDENT {
			return nil, ErrInvalidBlockEntryAssignment
		}

		key := lit

		tok, _ = t.scanSkipWS()
		if tok != tokASSIGNMENT {
			return nil, ErrInvalidBlockEntryAssignment
		}

		val, err := t.parseValue()
		if err != nil {
			return nil, err
		}

		m[key] = val
	}

	return m, nil
}

func (t *Parser) parseHeaders() (ast.HeaderEntries, error) {
	header := make(ast.HeaderEntries)

	for {
		tok, lit := t.scanSkipWS()
		if tok == tokLF {
			continue
		}
		if tok == tokDELIMITER || tok == tokEOF || tok == tokBLOCKSTART || tok == tokSECTION {
			t.unscan()
			break
		}

		if tok != tokIDENT {
			return nil, ErrInvalidHeaderKey
		}
		key := lit

		tok, _ = t.scanSkipWS()
		if tok != tokCOLON {
			return nil, ErrInvalidHeaderSeparator
		}

		val := strings.TrimSpace(t.s.scanUntilLF())
		if val == "" {
			return nil, ErrNoHeaderValue
		}

		header[key] = append(header[key], val)
	}

	return header, nil
}

func (t *Parser) parseRaw() (ast.DataContent, error) {
	var out bytes.Buffer

	inEscape := false

	r := t.s.read()
	if r == '@' {
		tk, file := t.s.scanString()
		if tk != tokSTRING {
			return nil, ErrInvalidFileDescriptor
		}
		return ast.FileDescriptor{Path: file}, nil
	}

	t.s.unread()

	for {
		if !inEscape {
			if out.Len() > 3 && string(out.Bytes()[out.Len()-4:]) == "\n---" {
				for {
					if t.s.read() != '-' {
						break
					}
				}
				t.s.unread()
				t.buf.tok = tokDELIMITER
				t.buf.lit = ""
				t.unscan()
				out.Truncate(out.Len() - 4)
				break
			}
			if out.Len() > 1 && string(out.Bytes()[out.Len()-2:]) == "\n[" {
				t.buf.tok = tokBLOCKSTART
				t.buf.lit = ""
				t.unscan()
				out.Truncate(out.Len() - 2)
				break
			}
			if out.Len() > 3 && string(out.Bytes()[out.Len()-4:]) == "\n###" {
				if t.s.read() != '#' {
					t.buf.tok = tokSECTION
					t.buf.lit = ""
				} else {
					if r = t.s.read(); r != '#' {
						return nil, ErrInvalidLogSection
					}
					// t.s.unread()
					t.buf.tok = tokLOGSECTION
					t.buf.lit = ""
				}

				t.unscan()
				out.Truncate(out.Len() - 4)
				break
			}
		}

		r := t.s.read()

		if r == eof {
			if inEscape {
				return ast.NoContent{}, ErrOpenEscapeBlock
			}
			t.s.unread()
			break
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

	outStr := out.String()
	if outStr == "" {
		return ast.NoContent{}, nil
	}

	return ast.TextBlock{Content: outStr}, nil
}

func (t *Parser) parseValue() (any, error) {
	tok, lit := t.scanSkipWS()

	switch tok {
	case tokINTEGER:
		return strconv.ParseInt(lit, 10, 64)
	case tokFLOAT:
		return strconv.ParseFloat(lit, 64)
	case tokSTRING:
		return lit, nil
	case tokIDENT:
		switch lit {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, errs.WithSuffix(ErrInvalidLiteral, "(boolean expression expected)")
		}
	case tokBLOCKSTART:
		return t.parseArray()
	case tokPARAMETER:
		return ParameterValue(lit), nil
	}

	return nil, errs.WithSuffix(ErrInvalidToken, "(value)")
}

func (t *Parser) parseArray() ([]any, error) {
	var arr []any

loop:
	for {
		tok, _ := t.scanSkipWS()
		switch tok {
		case tokBLOCKEND:
			break loop
		case tokCOMMA, tokLF:
			continue loop
		}

		t.unscan()

		val, err := t.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
	}

	return arr, nil
}

func (t *Parser) astPos() ast.Pos {
	return ast.Pos{
		Pos:     t.prevPos.pos,
		Line:    t.prevPos.line,
		LinePos: t.prevPos.linepos,
	}
}
