package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"unicode"
	"unicode/utf8"
)

func NewScanner(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	s.Split(ScanTokens)
	return &Scanner{
		c:   s,
		buf: &bytes.Buffer{},
	}
}

type Scanner struct {
	c    *bufio.Scanner
	buf  *bytes.Buffer
	tok  Tok
	attr Tok
	err  error
}

func (s *Scanner) Scan() ([]byte, error) {
	f := scanEvent
	for s.scan() {
		// fmt.Printf("---> %s <---\n", s.c.Bytes())
		f = f(s)
	}
	if s.err != nil {
		return nil, s.err
	}
	return format.Source(s.buf.Bytes())
}

func (s *Scanner) scan() bool {
	if s.err == nil && s.c.Scan() {
		s.tok = *s.next()
		return true
	}
	return false
}

func (s *Scanner) next() *Tok {
	data := s.c.Bytes()

	if len(data) > 1 && data[0] == '#' {
		return &Tok{TOK_COMMENT, data}
	}
	if len(data) == 1 {
		switch data[0] {
		case '{':
			return &Tok{TOK_LIST_BEGIN, data}
		case '}':
			return &Tok{TOK_LIST_END, data}
		case ';':
			return &Tok{TOK_ATTR_END, data}
		}
	}
	switch string(data) {
	case
		"uint16",
		"int16",
		"uint32",
		"int32",
		"string",
		"ip_addr",
		"int64",
		"uint64",
		"boolean":
		return &Tok{TOK_TYPE, data}
	}
	return &Tok{TOK_WORD, data}
}

type TokType int

const (
	TOK_COMMENT TokType = iota
	TOK_LIST_BEGIN
	TOK_LIST_END
	TOK_ATTR_END
	TOK_TYPE
	TOK_WORD
)

type Tok struct {
	t TokType
	b []byte
}

func (t *Tok) Bytes() []byte {
	return t.b
}

var invalidToken = fmt.Errorf("invalid token")

// ScanTokens looks for:
// comments (//)
// {, ;, }
// words
func ScanTokens(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}
	var comment bool
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '#' && i == start {
			comment = true
			continue
		}
		if comment {
			if r == '\n' {
				return i + width, data[start:i], nil
			}
			continue
		}
		if nonWord(r) {
			if i == start {
				i += width
			}
			return i, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func nonWord(r rune) bool {
	return !(r == '_' || r == ':' ||
		('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		('0' <= r && r <= '9'))
}

var isSpace = unicode.IsSpace

type scanFunc func(s *Scanner) scanFunc

func scanEvent(s *Scanner) scanFunc {
	if s.tok.t == TOK_COMMENT {
		writeComment(s)
		return scanEvent
	}
	if s.tok.t == TOK_WORD {
		fmt.Fprintf(s.buf, "type %s struct { ", s.tok.b)
		return scanListBegin
	}
	s.err = invalidToken
	return nil
}

func scanListBegin(s *Scanner) scanFunc {
	if s.tok.t == TOK_COMMENT {
		writeComment(s)
		return scanListBegin
	}
	if s.tok.t == TOK_LIST_BEGIN {
		return scanType
	}
	s.err = invalidToken
	return nil
}

func scanType(s *Scanner) scanFunc {
	if s.tok.t == TOK_COMMENT {
		writeComment(s)
		return scanType
	}
	if s.tok.t == TOK_LIST_END {
		fmt.Fprintf(s.buf, "}\n")
		return scanEvent
	}
	if s.tok.t == TOK_TYPE {
		t := string(s.tok.b)
		switch t {
		case
			"uint16",
			"int16",
			"uint32",
			"int32",
			"string",
			"int64",
			"uint64":
		case "ip_addr":
			t = "net.IP"
		case "boolean":
			t = "bool"
		default:
			s.err = invalidToken
			return nil
		}
		s.tok.b = []byte(t)
		s.attr = s.tok
		return scanAttr
	}
	s.err = invalidToken
	return nil
}

func scanAttr(s *Scanner) scanFunc {
	if s.tok.t == TOK_WORD {
		fmt.Fprintf(s.buf, "  %s %s ", s.tok.b, s.attr.b)
		return scanAttrEnd
	}
	s.err = invalidToken
	return nil
}

func scanAttrEnd(s *Scanner) scanFunc {
	if s.tok.t == TOK_ATTR_END {
		return scanType
	}
	s.err = invalidToken
	return nil
}

func writeComment(s *Scanner) {
	fmt.Fprintf(s.buf, "// %s\n", s.tok.b)
}
