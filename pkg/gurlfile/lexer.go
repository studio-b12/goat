package gurlfile

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type token int

const eof = rune(0)

const (
	// Special tokens
	ILLEGAL token = iota
	EOF
	WS
	LF

	// Literals
	IDENT

	// Control Characters
	COMMENT
	ESCAPE
	SECTION
	DELIMITER
	BLOCK_START
	BLOCK_END
	COLON
	COMMA
	ASSIGNMENT

	// Types
	STRING
	INTEGER
	FLOAT

	// Keywords
	USE
)

type scanner struct {
	r           *bufio.Reader
	linebuf     *bytes.Buffer
	lastline    *bytes.Buffer
	line        int
	lastlinepos int
	linepos     int
	pos         int
}

func newScanner(r io.Reader) *scanner {
	return &scanner{
		r:        bufio.NewReader(r),
		linebuf:  &bytes.Buffer{},
		lastline: &bytes.Buffer{},
	}
}

func (t *scanner) Scan() (tk token, lit string) {
	r := t.read()

	if isWhitespace(r) {
		t.unread()
		return t.scanWhitespace()
	}

	if isLetter(r) {
		t.unread()
		return t.scanIdent()
	}

	if isDigit(r) {
		t.unread()
		return t.scanNumber()
	}

	switch r {
	case '/':
		return t.scanComment()
	case '"', '\'':
		t.unread()
		return t.scanString()
	case '-':
		return t.scanDelimiter()
	case '#':
		return t.scanSection()

	case '[':
		return BLOCK_START, ""
	case ']':
		return BLOCK_END, ""
	case ':':
		return COLON, ""
	case ',':
		return COMMA, ""
	case '=':
		return ASSIGNMENT, ""
	case '\n':
		return LF, ""
	case eof:
		return EOF, ""
	}

	return ILLEGAL, string(r)
}

func (t *scanner) read() rune {
	t.pos++
	t.linepos++

	r, _, err := t.r.ReadRune()
	if err != nil {
		return eof
	}
	t.linebuf.WriteRune(r)

	if r == '\n' {
		t.line++
		t.lastlinepos = t.linepos - 1
		t.linepos = 0

		t.lastline.Reset()
		io.Copy(t.lastline, t.linebuf)
		t.linebuf.Reset()
	}

	return r
}

func (t *scanner) unread() {
	t.pos--
	if t.linepos == 0 {
		t.linepos = t.lastlinepos
		t.line--
	} else {
		t.linepos--
	}
	t.r.UnreadRune()
	if ln := t.linebuf.Len(); ln > 0 {
		t.linebuf.Truncate(t.linebuf.Len() - 1)
	}
}

func (t *scanner) readToLF() string {
	var b bytes.Buffer

	for {
		r := t.read()
		if r == eof || r == '\n' {
			break
		}
		b.WriteRune(r)
	}

	return strings.TrimSpace(b.String())
}

func (t *scanner) scanWhitespace() (tk token, lit string) {
	var b bytes.Buffer
	b.WriteRune(t.read())

	for {
		if r := t.read(); r == eof {
			break
		} else if !isWhitespace(r) {
			t.unread()
			break
		} else {
			b.WriteRune(r)
		}
	}

	return WS, b.String()
}

func (t *scanner) skipToLF() {
	for {
		r := t.read()
		if r == '\n' || r == eof {
			break
		}
	}
}

func (t *scanner) scanComment() (tk token, lit string) {
	if t.read() != '/' {
		return ILLEGAL, ""
	}

	t.skipToLF()

	return COMMENT, ""
}

func (t *scanner) scanDelimiter() (tk token, lit string) {
	for i := 0; i < 2; i++ {
		if t.read() != '-' {
			return ILLEGAL, ""
		}
	}

	t.skipToLF()

	return DELIMITER, ""
}

func (t *scanner) scanString() (tk token, lit string) {
	var b bytes.Buffer
	wrapper := rune(0)
	inString := false

	for {
		r := t.read()

		if r == eof || r == '\n' {
			if inString && wrapper != 0 {
				return ILLEGAL, ""
			}
			break
		}

		if inString {
			if isWhitespace(r) && wrapper == 0 {
				break
			}
			if r == wrapper {
				break
			}
			b.WriteRune(r)
		} else {
			if isWhitespace(r) {
				continue
			}
			if isStringWrapper(r) {
				wrapper = r
			} else {
				b.WriteRune(r)
			}
			inString = true
		}
	}

	return STRING, b.String()
}

func (t *scanner) scanUntilLF() string {
	var b bytes.Buffer

	for {
		r := t.read()
		if r == eof {
			t.unread()
			break
		}
		if r == '\n' {
			break
		}
		b.WriteRune(r)
	}

	return b.String()
}

func (t *scanner) scanSection() (tk token, lit string) {
	for i := 0; i < 2; i++ {
		if t.read() != '#' {
			return ILLEGAL, ""
		}
	}

	name := t.readToLF()

	return SECTION, name
}

func (t *scanner) scanIdent() (tk token, lit string) {
	var b bytes.Buffer
	b.WriteRune(t.read())

	for {
		if r := t.read(); r == eof {
			break
		} else if !isLetter(r) && !isDigit(r) && !isLiteralDelimiter(r) {
			t.unread()
			break
		} else {
			b.WriteRune(r)
		}
	}

	str := b.String()
	switch strings.ToLower(str) {
	case "use":
		return USE, ""
	}

	return IDENT, str
}

func (t *scanner) scanNumber() (tk token, lit string) {
	var b bytes.Buffer
	tk = INTEGER

	for {
		r := t.read()

		if r == '.' {
			tk = FLOAT
		} else if r == '_' {
			continue
		} else if !isDigit(r) {
			t.unread()
			break
		}

		b.WriteRune(r)
	}

	return tk, b.String()
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLiteralDelimiter(r rune) bool {
	return r == '_' || r == '-'
}

func isStringWrapper(r rune) bool {
	return r == '"' || r == '\''
}
