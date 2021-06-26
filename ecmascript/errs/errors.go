package errs

import (
	"fmt"

	"github.com/jchv/cleansheets/ecmascript/ast"
)

type SyntaxError struct {
	Location ast.Location
	Err      error
}

func (e *SyntaxError) Unwrap() error { return e.Err }

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s: syntax error: %s", &e.Location, e.Err)
}

type EncodingError struct {
	Location ast.Location
	Err      error
}

func (e *EncodingError) Unwrap() error { return e.Err }

func (e *EncodingError) Error() string {
	return fmt.Sprintf("%s: encoding error: %s", &e.Location, e.Err)
}

type ParserError struct {
	Location ast.Location
	Err      error
}

func (e *ParserError) Unwrap() error { return e.Err }

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s: parser error: %s", &e.Location, e.Err)
}
