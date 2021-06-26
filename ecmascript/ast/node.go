package ast

import "reflect"

// BaseNode is a small struct that stores the source code span between two
// nodes and provides an embeddable base for Node interface implementations.
type BaseNode struct {
	span Span
}

func (b *BaseNode) clearSpan() {
	b.span = Span{}
}

// SetStart sets the start of the node's source code span.
func (b *BaseNode) SetStart(l Location) {
	b.span.Start = l
}

// SetEnd sets the end of the node's source code span.
func (b *BaseNode) SetEnd(l Location) {
	b.span.End = l
}

// Span returns the span of source code the node represents.
func (b BaseNode) Span() Span {
	return b.span
}

func (b BaseNode) isNode() {}

// Node is the interface type of an AST node.
type Node interface {
	// Span returns the span of source code the node represents.
	Span() Span

	// ESTree returns the corresponding ESTree representation for this node.
	// Because Node is an interface, beware that calling ESTree directly on a
	// nil Node value will cause a panic.
	ESTree() interface{}

	isNode()
}

func clearSpans(v reflect.Value) {
	// Drop pointer down to concrete level.
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			clearSpans(v.Index(i))
		}

	case reflect.Struct:
		if v.CanAddr() && v.Addr().CanInterface() {
			if b, ok := v.Addr().Interface().(*BaseNode); ok {
				b.clearSpan()
			}
		}
		for i, n := 0, v.NumField(); i < n; i++ {
			clearSpans(v.Field(i))
		}

	default:
		break
	}
}

// ClearSpans removes source code span data from the AST subtree.
func ClearSpans(n Node) {
	clearSpans(reflect.ValueOf(n))
}
