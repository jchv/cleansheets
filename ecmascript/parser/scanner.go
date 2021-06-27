package parser

import (
	"errors"
	"fmt"

	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/errs"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

// Scanner provides lookahead for scanning tokens.
type Scanner struct {
	l *lexer.Lexer

	last []lexer.Token
	loc  []ast.Location
}

// NewScanner creates a new scanner.
func NewScanner(l *lexer.Lexer) *Scanner {
	return &Scanner{l: l}
}

// Location returns the current source code location.
func (s *Scanner) Location() ast.Location {
	if len(s.loc) > 0 {
		return s.loc[0]
	}
	return s.l.Location()
}

// PeekAt peeks into the future of the lexer. Calling this function will lex
// up to i tokens into the future.
func (s *Scanner) PeekAt(i int) lexer.Token {
	for len(s.last) <= i {
		s.loc = append(s.loc, s.Location())
		s.last = append(s.last, s.l.Lex())
	}
	return s.last[i]
}

// PeekLen returns how far we are peeked into the future.
func (s *Scanner) PeekLen() int {
	return len(s.last)
}

// Scan returns the next lexical token.
func (s *Scanner) Scan() lexer.Token {
	if len(s.last) > 0 {
		t := s.last[0]
		s.last = s.last[1:]
		s.loc = s.loc[1:]
		return t
	}
	return s.l.Lex()
}

// ReScan relexes the last token as a regular expression. Panics if we are
// currently peeked into the future, since ReScan changes the future.
func (s *Scanner) ReScan() lexer.ReToken {
	if len(s.last) > 0 {
		panic("internal error")
	}
	return s.l.ReLex()
}

// ScanExpect scans and panics if the token is not of the expected type.
func (s *Scanner) ScanExpect(typ lexer.TokenType, err string) lexer.Token {
	t := s.Scan()
	if t.Type != typ {
		if t.Type == lexer.TokenNone {
			s.SyntaxError(fmt.Sprintf("expected %s, got eof: %s", typ, err))
		} else {
			s.SyntaxError(fmt.Sprintf("expected %s, got %q: %s", typ, t.Source(), err))
		}
	}
	return t
}

// SyntaxError panics with a syntax error with the given string.
func (s *Scanner) SyntaxError(err string) {
	panic(&errs.SyntaxError{
		Location: s.Location(),
		Err:      errors.New(err),
	})
}
