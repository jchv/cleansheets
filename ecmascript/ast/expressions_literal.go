package ast

// NullLiteral is a node containing a null literal.
//
// For example:
//
//     null
//
// Would be represented as:
//
//     NullLiteral{}
type NullLiteral struct {
	BaseNode
}

// ESTree returns the corresponding ESTree representation for this node.
func (n NullLiteral) ESTree() interface{} {
	return struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
		Raw   string      `json:"raw"`
	}{
		Type:  "Literal",
		Value: nil,
		Raw:   "null",
	}
}

// BooleanLiteral is a node containing an ECMAScript boolean literal.
//
// For example:
//
//     true
//
// Would be represented as:
//
//     BooleanLiteral{
//         Value: true,
//         Raw: "true",
//     }
type BooleanLiteral struct {
	BaseNode
	Value bool
	Raw   string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n BooleanLiteral) ESTree() interface{} {
	return struct {
		Type  string `json:"type"`
		Value bool   `json:"value"`
		Raw   string `json:"raw"`
	}{
		Type:  "Literal",
		Value: n.Value,
		Raw:   n.Raw,
	}
}

// StringLiteral is a node containing an ECMAScript string literal.
//
// For example:
//
//     "test!"
//
// Would be represented as:
//
//     StringLiteral{
//         Value: "test!",
//         Raw: "\"test!\"",
//     }
type StringLiteral struct {
	BaseNode
	Value string
	Raw   string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n StringLiteral) ESTree() interface{} {
	return struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		Raw   string `json:"raw"`
	}{
		Type:  "Literal",
		Value: n.Value,
		Raw:   n.Raw,
	}
}

// NumberLiteral is a node containing an ECMAScript numeric literal.
//
// For example:
//
//     0.0
//
// Would be represented as:
//
//     NumberLiteral{
//         Value: 0,
//         Raw: "0.0",
//     }
type NumberLiteral struct {
	BaseNode
	Value float64
	Raw   string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n NumberLiteral) ESTree() interface{} {
	return struct {
		Type  string  `json:"type"`
		Value float64 `json:"value"`
		Raw   string  `json:"raw"`
	}{
		Type:  "Literal",
		Value: n.Value,
		Raw:   n.Raw,
	}
}

// RegExpLiteral is a node containing an ECMAScript regular expression literal.
//
// For example:
//
//     /a/g
//
// Would be represented as:
//
//     RegExpLiteral{
// 	       Pattern: "a",
//         Flags: "g",
//         Raw: "/a/g",
//     }
type RegExpLiteral struct {
	BaseNode
	Pattern string
	Flags   string
	Raw     string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n RegExpLiteral) ESTree() interface{} {
	return struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		Raw   string `json:"raw"`
		Regex struct {
			Pattern string `json:"pattern"`
			Flags   string `json:"flags"`
		} `json:"regex"`
	}{
		Type:  "Literal",
		Value: n.Raw,
		Raw:   n.Raw,
		Regex: struct {
			Pattern string `json:"pattern"`
			Flags   string `json:"flags"`
		}{
			Pattern: n.Pattern,
			Flags:   n.Flags,
		},
	}
}
