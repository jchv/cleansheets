package errs

import (
	"fmt"

	"github.com/jchv/cleansheets/ecmascript/ast"
)

// SyntaxError is emitted when the parser or lexer encounters invalid syntax.
type SyntaxError struct {
	Location ast.Location
	Err      error
}

// Unwrap returns the embedded error.
func (e *SyntaxError) Unwrap() error { return e.Err }

// Error implements the error interface.
func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s: syntax error: %s", &e.Location, e.Err)
}

// EncodingError is emitted when the scanner encounters an invalid sequence.
type EncodingError struct {
	Location ast.Location
	Err      error
}

// Unwrap returns the embedded error.
func (e *EncodingError) Unwrap() error { return e.Err }

// Error implements the error interface.
func (e *EncodingError) Error() string {
	return fmt.Sprintf("%s: encoding error: %s", &e.Location, e.Err)
}

// ParserError is returned when the parser encounters an error.
type ParserError struct {
	Location ast.Location
	Err      error
}

// Unwrap returns the embedded error.
func (e *ParserError) Unwrap() error { return e.Err }

// Error implements the error interface.
func (e *ParserError) Error() string {
	return fmt.Sprintf("%s: parser error: %s", &e.Location, e.Err)
}
