package ast

// UpdateOperator is an enumeration type for ECMAScript update operators.
type UpdateOperator int

const (
	// UpdatePreIncrementOp is the operator for pre-increment, i.e. ++i
	UpdatePreIncrementOp UpdateOperator = iota

	// UpdatePreDecrementOp is the operator for pre-decrement, i.e. --i
	UpdatePreDecrementOp

	// UpdatePostIncrementOp is the operator for post-increment, i.e. i++
	UpdatePostIncrementOp

	// UpdatePostDecrementOp is the operator for post-decrement, i.e. i--
	UpdatePostDecrementOp
)

// estreeUpdateOpMap maps from a UpdateOperator value to the corresponding
// EStree string.
var estreeUpdateOpMap = map[UpdateOperator]string{
	UpdatePreIncrementOp:  "++",
	UpdatePreDecrementOp:  "--",
	UpdatePostIncrementOp: "++",
	UpdatePostDecrementOp: "--",
}

// estreeUpdateOpPrefixMap maps from a UpdateOperator value to the value of the
// `prefix` field of the ESTree node.
var estreeUpdateOpPrefixMap = map[UpdateOperator]bool{
	UpdatePreIncrementOp:  true,
	UpdatePreDecrementOp:  true,
	UpdatePostIncrementOp: false,
	UpdatePostDecrementOp: false,
}

// UnaryOperator is an enumeration type for ECMAScript unary operators.
type UnaryOperator int

const (
	// UnaryDeleteOp (delete) is the operator for deleting properties.
	UnaryDeleteOp UnaryOperator = iota

	// UnaryVoidOp (void) is the operator for discarding an expression value.
	UnaryVoidOp

	// UnaryTypeOfOp (typeof) is the operator for getting the primitive type
	// of a value.
	UnaryTypeOfOp

	// UnaryPlusOp (+) is an operator that converts a value to a numeric value.
	UnaryPlusOp

	// UnaryMinusOp (-) is the negation operator.
	UnaryMinusOp

	// UnaryBitNotOp (~) is the bitwise not operator.
	UnaryBitNotOp

	// UnaryNotOp (!) is the logical not operator.
	UnaryNotOp
)

// estreeUnaryOpMap maps from a UnaryOperator value to the corresponding EStree
// string.
var estreeUnaryOpMap = map[UnaryOperator]string{
	UnaryDeleteOp: "delete",
	UnaryVoidOp:   "void",
	UnaryTypeOfOp: "typeof",
	UnaryPlusOp:   "+",
	UnaryMinusOp:  "-",
	UnaryBitNotOp: "~",
	UnaryNotOp:    "!",
}

// estreeUnaryOpPrefixMap maps from a UnaryOperator value to the value of the
// `prefix` field of the ESTree node.
var estreeUnaryOpPrefixMap = map[UnaryOperator]bool{
	UnaryDeleteOp: true,
	UnaryVoidOp:   true,
	UnaryTypeOfOp: true,
	UnaryPlusOp:   true,
	UnaryMinusOp:  true,
	UnaryBitNotOp: true,
	UnaryNotOp:    true,
}

// UpdateExpression is the node for an ECMAScript update expression statement.
type UpdateExpression struct {
	BaseNode

	Operator UpdateOperator
	Argument Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n UpdateExpression) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Operator string      `json:"operator"`
		Argument interface{} `json:"argument"`
		Prefix   bool        `json:"prefix"`
	}{
		Type:     "UpdateExpression",
		Operator: estreeUpdateOpMap[n.Operator],
		Argument: estree(n.Argument),
		Prefix:   estreeUpdateOpPrefixMap[n.Operator],
	}
}

// UnaryExpression is the AST node for an unary expression statement.
type UnaryExpression struct {
	BaseNode

	Operator UnaryOperator
	Argument Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n UnaryExpression) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Operator string      `json:"operator"`
		Argument interface{} `json:"argument"`
		Prefix   bool        `json:"prefix"`
	}{
		Type:     "UnaryExpression",
		Operator: estreeUnaryOpMap[n.Operator],
		Argument: estree(n.Argument),
		Prefix:   estreeUnaryOpPrefixMap[n.Operator],
	}
}
