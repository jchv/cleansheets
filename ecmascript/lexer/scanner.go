package lexer

import (
	"errors"
	"io"
	"net/url"

	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/errs"
)

// EOFRune is returned when the scanner hits an EOF error.
const EOFRune = rune(-1)

// Scanner provides additional logic on top of a RuneScanner.
type Scanner struct {
	r io.RuneScanner

	uri      *url.URL
	col, row int

	eof bool
}

// NewScanner creates a new scanner for the given RuneScanner and URL.
func NewScanner(r io.RuneScanner, uri *url.URL) *Scanner {
	return &Scanner{
		r:   r,
		uri: uri,
		col: 1,
		row: 1,
	}
}

// Location returns the current source code location.
func (s *Scanner) Location() ast.Location {
	column := s.col

	if column < 0 {
		column = 1
	}

	return ast.Location{
		URI:    s.uri,
		Column: column,
		Row:    s.row,
	}
}

// Read reads a rune and returns it. On EOF, EOFRune is returned.
func (s *Scanner) Read() rune {
	r, _, err := s.r.ReadRune()

	if errors.Is(err, io.EOF) {
		s.eof = true
		return EOFRune
	}

	if err != nil {
		panic(&errs.EncodingError{
			Location: s.Location(),
			Err:      err,
		})
	}

	// Increment source location. On newline, we set col to -col. This allows
	// us to know when we're unreading a line terminator (because col will be
	// negative) and what to restore it to without needing additional state.
	if _, ok := lineterms[r]; ok {
		s.row++
		if s.col > 0 {
			// Last read was not a newline
			s.col = -s.col
		} else if s.col < 0 {
			// Last read was a newline- treat it as having been column 1.
			s.col = -1
		}
	} else {
		if s.col < 0 {
			s.col = 1
		}
		s.col++
	}

	return r
}

// Unread unreads a rune. If we are at EOF, this will not call the underlying
// RuneReader, so it is safe to unread at EOF.
func (s *Scanner) Unread() {
	if !s.eof {
		err := s.r.UnreadRune()

		if err != nil {
			panic(&errs.ParserError{
				Location: s.Location(),
				Err:      err,
			})
		}
	}

	// If negative: we just read a line terminal rune. Invert col and
	// decrement row.
	// If positive: we read any other rune. Just decrement col.
	if s.col < 0 {
		s.col = -s.col
		s.row--
	} else {
		s.col--
	}
}
