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
	c := bufio.NewScanner(r)
	c.Split(ScanTokens)
	return &Scanner{
		c:   c,
		esf: Esf{},
	}
}

type Scanner struct {
	c   *bufio.Scanner
	buf *bytes.Buffer
	tok Tok
	esf Esf
	err error
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
	buf := s.esf.Bytes()
	return format.Source(buf)
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
	r := rune(data[0])
	if nonWord(r) {
		return &Tok{TOK_INVALID, data}
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
	TOK_INVALID
)

type Tok struct {
	t TokType
	b []byte
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
		e := s.esf.Event()
		e.Name = s.tok.b
		if string(e.Name) == "MetaEventInfo" {
			s.esf.meta = e
		}
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
		s.esf.Event().Attr().AddComment(s.tok.b)
		return scanType
	}
	if s.tok.t == TOK_LIST_END {
		s.esf.clearEvent()
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
		s.esf.Event().NewAttr().Type = []byte(t)
		return scanAttr
	}
	s.err = invalidToken
	return nil
}

func scanAttr(s *Scanner) scanFunc {
	if s.tok.t == TOK_WORD {
		s.esf.Event().Attr().Name = s.tok.b
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
	s.esf.Event().AddComment(s.tok.b)
}

// Build ESF
type Comment struct {
	Lines [][]byte
}

func (c *Comment) AddComment(b []byte) {
	var i int
	for ; i < len(b); i++ {
		if b[i] != '#' && !isSpace(rune(b[i])) {
			break
		}
	}
	c.Lines = append(c.Lines, b[i:])
}

type Esf struct {
	meta   *Event
	event  *Event
	Events []Event
}

func (e *Esf) Bytes() []byte {
	var buf bytes.Buffer
	var hasMeta bool

	if e.meta != nil {
		hasMeta = true
		buf.Write(e.meta.Bytes())
	}

	for _, event := range e.Events {
		event.meta = hasMeta
		buf.Write(event.Bytes())
	}

	return buf.Bytes()
}

func (e *Esf) Event() *Event {
	if e.event == nil {
		e.event = &Event{}
	}
	return e.event
}

func (e *Esf) clearEvent() {
	if e.event != nil && e.event != e.meta {
		e.Events = append(e.Events, *e.event)
	}
	e.event = nil
}

type Event struct {
	Comment
	Name  []byte
	Attrs []Attr
	attr  *Attr
	meta  bool
}

func (e *Event) Bytes() []byte {
	buf := &bytes.Buffer{}

	for _, c := range e.Comment.Lines {
		fmt.Fprintf(buf, "// %s\n", c)
	}
	fmt.Fprintf(buf, "type %s struct {\n", e.Name)
	if e.meta {
		fmt.Fprintf(buf, "  MetaEventInfo\n")
	}
	for _, attr := range e.Attrs {
		for _, c := range attr.Comment.Lines {
			fmt.Fprintf(buf, "// %s\n", c)
		}
		fmt.Fprintf(buf, "  %s %s\n", attr.Name, attr.Type)
	}
	fmt.Fprintf(buf, "}\n\n")

	return buf.Bytes()
}

func (e *Event) Attr() *Attr {
	if e.attr == nil {
		e.attr = &Attr{}
	}
	return e.attr
}

func (e *Event) NewAttr() *Attr {
	if e.attr != nil {
		e.Attrs = append(e.Attrs, *e.attr)
		e.attr = nil
	}
	return e.Attr()
}

type Attr struct {
	Comment
	Name, Type []byte
}
