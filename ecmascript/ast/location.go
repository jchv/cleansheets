package ast

import (
	"fmt"
	"net/url"
)

// Location represents a single source location.
type Location struct {
	URI *url.URL

	Column, Row int
}

// Span represents a range from one location in source to another.
type Span struct {
	Start, End Location
}

// Span returns a span consisting of only this location.
func (l Location) Span() Span {
	return Span{l, l}
}

// String returns a string representing the source location.
func (l *Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.URI, l.Row, l.Column)
}

// String returns a string representation of the location span.
func (l *Span) String() string {
	a, b := l.Start, l.End
	if a.URI != b.URI {
		return fmt.Sprintf("%s:%d:%d-%s-%d-%d", a.URI, a.Row, a.Column, b.URI, b.Row, b.Column)
	}
	if a.Row != b.Row {
		return fmt.Sprintf("%s:%d:%d-%d-%d", a.URI, a.Row, a.Column, b.Row, b.Column)
	}
	if a.Column != b.Column {
		return fmt.Sprintf("%s:%d:%d-%d", a.URI, a.Row, a.Column, b.Column)
	}
	return fmt.Sprintf("%s:%d:%d", a.URI, a.Row, a.Column)
}
