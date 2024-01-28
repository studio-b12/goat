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
		pos := t.astPos()
		tok, lit := t.scan()

		switch tok {
		case tokWS, tokLF:
			continue

		case tokCOMMENT:
			gf.Comments = append(gf.Comments, ast.Comment{
				Pos:     pos,
				Content: lit,
			})
			continue

		case tokDELIMITER:
			gf.Delimiters = append(gf.Delimiters, ast.Delimiter{Pos: pos, ExtraLen: len(lit)})
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
			req, comms, err := t.parseRequest()
			if err != nil {
				return nil, err
			}
			gf.Actions = append(gf.Actions, req)
			gf.Comments = append(gf.Comments, comms...)

		case tokEXECUTE:
			exec, comms, err := t.parseExecute()
			if err != nil {
				return nil, err
			}
			gf.Actions = append(gf.Actions, exec)
			gf.Comments = append(gf.Comments, comms...)

		case tokUSE:
			use, err := t.parseUse()
			if err != nil {
				return nil, err
			}
			gf.Imports = append(gf.Imports, *use)

		case tokSECTION:
			sec, comms, delims, err := t.parseSection()
			if err != nil {
				return nil, err
			}
			gf.Sections = append(gf.Sections, sec)
			gf.Comments = append(gf.Comments, comms...)
			gf.Delimiters = append(gf.Delimiters, delims...)
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

	imp := &ast.Import{
		Pos:  pos,
		Path: lit,
	}
	return imp, nil
}

func (t *Parser) parseExecute() (*ast.Execute, []ast.Comment, error) {
	var comments []ast.Comment

	tk, _ := t.scan()
	if tk != tokWS {
		return nil, nil, ErrInvalidStringLiteral
	}

	tk, lit := t.s.scanString()
	if tk == tokILLEGAL {
		return nil, nil, ErrInvalidStringLiteral
	}

	if lit == "" {
		return nil, nil, ErrEmptyCallPath
	}

	var (
		exec ast.Execute
		err  error
	)

	exec.Path = lit

	tok, _ := t.scanSkipWS()
	if tok != tokGROUPSTART {
		t.unscan()
		return &exec, nil, nil
	}

	end := tokGROUPEND
	exec.Parameters, comments, err = t.parseBlockEntries(&end)
	if err != nil {
		return nil, nil, err
	}
	t.scan() // re-scan closing group `)`

	tok, _ = t.scanSkipWS()
	if tok != tokRETURN {
		t.unscan()
		return &exec, comments, nil
	}

	tok, _ = t.scanSkipWS()
	if tok != tokGROUPSTART {
		return nil, nil, ErrMissingGroup
	}

	for {
		pos := t.astPos()
		tok, key := t.scanSkipWS()
		if tok == tokEOF {
			return nil, nil, ErrUnclosedGroup
		}
		if tok == tokLF {
			continue
		}
		if tok == tokGROUPEND {
			break
		}
		if tok != tokIDENT {
			return nil, nil, ErrIllegalCharacter
		}

		tok, _ = t.scanSkipWS()
		if tok != tokAS {
			return nil, nil, ErrIllegalCharacter
		}

		tok, val := t.scanSkipWS()
		if tok != tokIDENT {
			return nil, nil, ErrIllegalCharacter
		}

		exec.Returns.KVList = append(exec.Returns.KVList, ast.KV[string]{Pos: pos, Key: key, Value: val})
	}

	return &exec, comments, nil
}

func (t *Parser) parseSection() (sect ast.Section, comments []ast.Comment, delims []ast.Delimiter, err error) {
	sectionPos := t.astPos()

	name := SectionName(strings.ToLower(strings.TrimSpace(t.s.readToLF())))

	if name == SectionDefaults {
		pr, comms, err := t.parseDefaults()
		if err != nil {
			return nil, nil, nil, err
		}
		sect = ast.SectionDefaults{Pos: sectionPos, Request: *pr}
		comments = append(comments, comms...)
		return sect, comments, delims, nil
	}

	var actions []ast.Action

	for {
		pos := t.astPos()
		tok, lit := t.scan()

		if tok == tokLF || tok == tokWS {
			continue
		}

		if tok == tokDELIMITER {
			delims = append(delims, ast.Delimiter{Pos: pos, ExtraLen: len(lit)})
			continue
		}

		if tok == tokCOMMENT {
			comments = append(comments, ast.Comment{Pos: pos, Content: lit})
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
			exec, comms, err := t.parseExecute()
			if err != nil {
				return nil, nil, nil, err
			}
			actions = append(actions, exec)
			comments = append(comments, comms...)
			continue
		}

		t.unscan()
		req, comms, err := t.parseRequest()
		if err != nil {
			return nil, nil, nil, err
		}
		actions = append(actions, req)
		comments = append(comments, comms...)
	}

	switch name {
	case SectionSetup:
		sect = ast.SectionSetup{
			Pos:     sectionPos,
			Actions: actions,
		}
	case SectionTests:
		sect = ast.SectionTests{
			Pos:     sectionPos,
			Actions: actions,
		}
	case SectionTeardown:
		sect = ast.SectionTeardown{
			Pos:     sectionPos,
			Actions: actions,
		}
	default:
		return nil, nil, nil, ErrInvalidSection
	}

	return sect, comments, delims, nil
}

func (t *Parser) parseRequest() (*ast.Request, []ast.Comment, error) {
	var (
		req      ast.Request
		comments []ast.Comment
	)

	// parse header

	req.Pos = t.astPos()

	tok, lit := t.scan()
	if tok != tokIDENT && tok != tokSTRING || lit == "" {
		return nil, nil, ErrInvalidRequestMethod
	}
	req.Head.Method = lit

	tok, _ = t.scan()
	if tok != tokWS && tok != tokLF {
		return nil, nil, ErrNoRequestURI
	}

	tok, lit = t.s.scanString()
	if tok != tokSTRING || lit == "" {
		return nil, nil, ErrNoRequestURI
	}
	req.Head.Url = lit

loop:
	for {
		pos := t.astPos()
		tok, lit = t.scan()

		switch tok {
		case tokCOMMENT:
			comments = append(comments, ast.Comment{Pos: pos, Content: lit})

		case tokBLOCKSTART:
			block, comms, err := t.parseBlock()
			if err != nil {
				return nil, nil, err
			}
			req.Blocks = append(req.Blocks, block)
			comments = append(comments, comms...)

		case tokWS, tokLF:
			continue loop
		case tokEOF, tokSECTION, tokLOGSECTION, tokDELIMITER:
			t.unscan()
			break loop

		default:
			return nil, nil, errs.WithSuffix(ErrInvalidToken, "(request)")
		}
	}

	return &req, comments, nil
}

func (t *Parser) parseDefaults() (*ast.PartialRequest, []ast.Comment, error) {
	var (
		req      ast.PartialRequest
		comments []ast.Comment
		err      error
	)

	req.Pos = t.astPos()

loop:
	for {
		pos := t.astPos()
		tok, lit := t.scan()

		switch tok {
		case tokCOMMENT:
			comments = append(comments, ast.Comment{Pos: pos, Content: lit})

		case tokBLOCKSTART:
			block, comms, err := t.parseBlock()
			if err != nil {
				return nil, nil, err
			}
			req.Blocks = append(req.Blocks, block)
			comments = append(comments, comms...)

		case tokWS, tokLF:
			continue loop
		case tokEOF, tokSECTION, tokLOGSECTION, tokDELIMITER:
			t.unscan()
			break loop

		default:
			err = errs.WithSuffix(ErrInvalidToken, "(request)")
		}

		if err != nil {
			return nil, nil, err
		}
	}

	return &req, comments, nil
}

func (t *Parser) parseBlock() (ast.RequestBlock, []ast.Comment, error) {
	var (
		blockHeader string
		comments    []ast.Comment
	)

	tok, lit := t.scanSkipWS()
	if tok != tokIDENT || lit == "" {
		return nil, nil, ErrInvalidBlockHeader
	}
	blockHeader = lit

	tok, _ = t.scan()
	if tok != tokBLOCKEND {
		return nil, nil, ErrInvalidBlockHeader
	}

	pos := t.astPos()
	tok, lit = t.scanSkipWS()
	if tok == tokCOMMENT {
		comments = append(comments, ast.Comment{Pos: pos, Content: lit})
	} else if tok != tokLF {
		return nil, nil, errs.WithSuffix(ErrInvalidToken, "(block)")
	}

	optName := optionName(strings.ToLower(blockHeader))

	switch optName {

	case optionNameQueryParams:
		data, comms, err := t.parseBlockEntries(nil)
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comms...)
		return ast.RequestQueryParams{KVList: data}, comments, nil

	case optionNameHeader:
		header, comms, err := t.parseHeaders()
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comms...)
		return ast.RequestHeader{HeaderEntries: header}, comments, nil

	case optionNameBody:
		raw, err := t.parseRaw()
		if err != nil {
			return nil, nil, err
		}
		return ast.RequestBody{DataContent: raw}, comments, nil

	case optionNamePreScript:
		raw, err := t.parseRaw()
		if err != nil {
			return nil, nil, err
		}
		return ast.RequestPreScript{DataContent: raw}, comments, nil

	case optionNameScript:
		raw, err := t.parseRaw()
		if err != nil {
			return nil, nil, err
		}
		return ast.RequestScript{DataContent: raw}, comments, nil

	case optionNameOptions:
		data, comms, err := t.parseBlockEntries(nil)
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comms...)
		return ast.RequestOptions{KVList: data}, comments, nil

	case optionNameAuth:
		data, comms, err := t.parseBlockEntries(nil)
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comms...)
		return ast.RequestAuth{KVList: data}, comments, nil

	default:
		return nil, nil, errs.WithSuffix(ErrInvalidBlockHeader,
			fmt.Sprintf("('%s')", blockHeader))
	}
}

func (t *Parser) parseBlockEntries(exitToken *token) (ast.KVList[any], []ast.Comment, error) {
	m := make(ast.KVList[any], 0)
	var comments []ast.Comment

	for {
		pos := t.astPos()
		tok, lit := t.scanSkipWS()
		if tok == tokLF {
			continue
		}
		if tok == tokCOMMENT {
			comments = append(comments, ast.Comment{Pos: pos, Content: lit})
			continue
		}
		if tok == tokEOF ||
			(exitToken != nil && tok == *exitToken) ||
			(tok == tokDELIMITER || tok == tokBLOCKSTART || tok == tokSECTION) {

			t.unscan()
			break
		}

		if tok != tokIDENT {
			return nil, nil, ErrInvalidBlockEntryAssignment
		}

		key := lit

		tok, _ = t.scanSkipWS()
		if tok != tokASSIGNMENT {
			return nil, nil, ErrInvalidBlockEntryAssignment
		}

		val, comms, err := t.parseValue()
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comms...)

		m = append(m, ast.KV[any]{Pos: pos, Key: key, Value: val})
	}

	return m, comments, nil
}

func (t *Parser) parseHeaders() (header ast.HeaderEntries, comments []ast.Comment, err error) {

	for {
		pos := t.astPos()
		tok, lit := t.scanSkipWS()
		if tok == tokLF {
			continue
		}
		if tok == tokDELIMITER || tok == tokEOF || tok == tokBLOCKSTART || tok == tokSECTION {
			t.unscan()
			break
		}

		if tok == tokCOMMENT {
			comments = append(comments, ast.Comment{Pos: pos, Content: lit})
			continue
		}

		if tok != tokIDENT {
			return header, nil, ErrInvalidHeaderKey
		}
		key := lit

		tok, _ = t.scanSkipWS()
		if tok != tokCOLON {
			return header, nil, ErrInvalidHeaderSeparator
		}

		val := strings.TrimSpace(t.s.scanUntilLF())
		if val == "" {
			return header, nil, ErrNoHeaderValue
		}

		header.KVList = append(header.KVList, ast.KV[string]{Pos: pos, Key: key, Value: val})
	}

	return header, comments, nil
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
				lit := t.s.readToLF()
				t.s.unread()
				t.buf.tok = tokDELIMITER
				t.buf.lit = lit
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

func (t *Parser) parseValue() (any, []ast.Comment, error) {
	tok, lit := t.scanSkipWS()

	switch tok {
	case tokINTEGER:
		v, err := strconv.ParseInt(lit, 10, 64)
		return v, nil, err
	case tokFLOAT:
		v, err := strconv.ParseFloat(lit, 64)
		return v, nil, err
	case tokSTRING:
		return lit, nil, nil
	case tokIDENT:
		switch lit {
		case "true":
			return true, nil, nil
		case "false":
			return false, nil, nil
		default:
			return nil, nil, errs.WithSuffix(ErrInvalidLiteral, "(boolean expression expected)")
		}
	case tokBLOCKSTART:
		return t.parseArray()
	case tokPARAMETER:
		return ParameterValue(lit), nil, nil
	default:
		return nil, nil, errs.WithSuffix(ErrInvalidToken, "(value)")
	}
}

func (t *Parser) parseArray() ([]any, []ast.Comment, error) {
	var (
		arr      []any
		comments []ast.Comment
	)

loop:
	for {
		pos := t.astPos()
		tok, lit := t.scanSkipWS()

		switch tok {
		case tokBLOCKEND:
			break loop
		case tokCOMMA, tokLF:
			continue loop
		case tokCOMMENT:
			comments = append(comments, ast.Comment{Pos: pos, Content: lit})
			continue loop
		}

		t.unscan()

		val, comms, err := t.parseValue()
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comms...)
		arr = append(arr, val)
	}

	return arr, comments, nil
}

func (t *Parser) astPos() ast.Pos {
	return ast.Pos{
		Pos:     t.s.pos,
		Line:    t.s.line,
		LinePos: t.s.linepos,
	}
}
