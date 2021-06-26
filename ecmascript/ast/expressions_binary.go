package ast

// BinaryOperator is an enumeration type for ECMAScript binary operators.
type BinaryOperator int

const (
	// BinaryExponentOp (**) is the operator for exponentiation.
	BinaryExponentOp BinaryOperator = iota

	// BinaryMultOp (*) is the operator for multiplication.
	BinaryMultOp

	// BinaryDivOp (/) is the operator for division.
	BinaryDivOp

	// BinaryModOp (%) is the operator for modulo.
	BinaryModOp

	// BinaryAddOp (+) is the operator for addition.
	BinaryAddOp

	// BinarySubOp (-) is the operator for subtraction.
	BinarySubOp

	// BinaryLShiftOp (<<) is the operator for a bitwise left shift.
	BinaryLShiftOp

	// BinaryRShiftOp (>>) is the operator for a bitwise right shift.
	BinaryRShiftOp

	// BinaryUnsignedRShiftOp (>>>) is the operator for an unsigned bitwise
	// right shift.
	BinaryUnsignedRShiftOp

	// BinaryLessThanOp (<) is the operator for a less-than relational
	// comparison.
	BinaryLessThanOp

	// BinaryGreaterThanOp (>) is the operator for a greater-than relational
	// comparison.
	BinaryGreaterThanOp

	// BinaryLessThanEqualOp (<=) is the operator for a less-than-or-equal-to
	// relational comparison.
	BinaryLessThanEqualOp

	// BinaryGreaterThanEqualOp (>=) is the operator for a greater-than-or-
	// equal-to relational comparison.
	BinaryGreaterThanEqualOp

	// BinaryInstanceOfOp (instanceof) is the operator for checking prototype
	// ancestry.
	BinaryInstanceOfOp

	// BinaryIn (in) is the operator for checking property existence.
	BinaryInOp

	// BinaryEqualOp (==) is the operator for checking value equality.
	BinaryEqualOp

	// BinaryNotEqualOp (!=) is the operator for checking value inequality.
	BinaryNotEqualOp

	// BinaryStrictEqualOp (===) is the operator for checking type and value
	// equality.
	BinaryStrictEqualOp

	// BinaryStrictNotEqualOp (!==) is the operator for checking type or value
	// inequality.
	BinaryStrictNotEqualOp

	// BinaryBitAndOp (&) is the operator for a bitwise AND operation.
	BinaryBitAndOp

	// BinaryBitXorOp (^) is the operator for a bitwise XOR operation.
	BinaryBitXorOp

	// BinaryBitOrOp (|) is the operator for a bitwise OR operation.
	BinaryBitOrOp

	// BinaryLogicalAndOp (&&) is the operator for a logical AND operation.
	BinaryLogicalAndOp

	// BinaryLogicalOrOp (||) is the operator for a logical OR operation.
	BinaryLogicalOrOp

	// BinaryCoalesceOp (??) is the operator for a null coalescing operation.
	BinaryCoalesceOp
)

// estreeBinaryOpMap maps from a BinaryOperator value to the corresponding
// EStree string.
var estreeBinaryOpMap = map[BinaryOperator]string{
	BinaryExponentOp:         "**",
	BinaryMultOp:             "*",
	BinaryDivOp:              "/",
	BinaryModOp:              "%",
	BinaryAddOp:              "+",
	BinarySubOp:              "-",
	BinaryLShiftOp:           "<<",
	BinaryRShiftOp:           ">>",
	BinaryUnsignedRShiftOp:   ">>>",
	BinaryLessThanOp:         "<",
	BinaryGreaterThanOp:      ">",
	BinaryLessThanEqualOp:    "<=",
	BinaryGreaterThanEqualOp: ">=",
	BinaryInstanceOfOp:       "instanceof",
	BinaryInOp:               "in",
	BinaryEqualOp:            "==",
	BinaryNotEqualOp:         "!=",
	BinaryStrictEqualOp:      "===",
	BinaryStrictNotEqualOp:   "!==",
	BinaryBitAndOp:           "&",
	BinaryBitXorOp:           "^",
	BinaryBitOrOp:            "|",
	BinaryLogicalAndOp:       "&&",
	BinaryLogicalOrOp:        "||",
	BinaryCoalesceOp:         "??",
}

// AssignmentOperator is an enumeration type for ECMAScript assignment
// operators.
type AssignmentOperator int

const (
	// AssignmentOp (=) is the assignment operator.
	AssignmentOp AssignmentOperator = iota

	// AssignmentMultOp (*=) is the compound assignment operator for
	// multiply-assign.
	AssignmentMultOp

	// AssignmentDivOp (/=) is the compound assignment operator for
	// divide-assign.
	AssignmentDivOp

	// AssignmentModOp (%=) is the compound assignment operator for
	// modulo-assign.
	AssignmentModOp

	// AssignmentAddOp (+=) is the compound assignment operator for
	// add-assign.
	AssignmentAddOp

	// AssignmentSubOp (-=) is the compound assignment operator for
	// subtract-assign.
	AssignmentSubOp

	// AssignmentLShiftOp (<<=) is the compound assignment operator for
	// bitwise-left-shift-assign.
	AssignmentLShiftOp

	// AssignmentRShiftOp (>>=) is the compound assignment operator for
	// bitwise-right-shift-assign.
	AssignmentRShiftOp

	// AssignmentUnsignedRShiftOp (>>>=) is the compound assignment operator
	// for bitwise-unsigned-right-shift-assign.
	AssignmentUnsignedRShiftOp

	// AssignmentBitAndOp (&=) is the compound assignment operator for
	// bitwise-and-assign.
	AssignmentBitAndOp

	// AssignmentBitXorOp (^=) is the compound assignment operator for
	// bitwise-xor-assign.
	AssignmentBitXorOp

	// AssignmentBitOrOp (|=) is the compound assignment operator for
	// bitwise-or-assign.
	AssignmentBitOrOp

	// AssignmentExponentOp (**=) is the compound assignment operator for
	// exponentiate-assign.
	AssignmentExponentOp

	// AssignmentLogicalAndOp (&&=) is the compound assignment operator for
	// logical-and-assign.
	AssignmentLogicalAndOp

	// AssignmentLogicalOr (||=) is the compound assignment operator for
	// logical-or-assign.
	AssignmentLogicalOr

	// AssignmentCoalesceOp (??=) is the compound assignment operator for
	// null-coalesce-assign.
	AssignmentCoalesceOp
)

// estreeAssignOpMap maps from a AssignmentOperator value to the corresponding
// EStree string.
var estreeAssignOpMap = map[AssignmentOperator]string{
	AssignmentOp:               "=",
	AssignmentMultOp:           "*=",
	AssignmentDivOp:            "/=",
	AssignmentModOp:            "%=",
	AssignmentAddOp:            "+=",
	AssignmentSubOp:            "-=",
	AssignmentLShiftOp:         "<<=",
	AssignmentRShiftOp:         ">>=",
	AssignmentUnsignedRShiftOp: ">>>=",
	AssignmentBitAndOp:         "&=",
	AssignmentBitXorOp:         "^=",
	AssignmentBitOrOp:          "|=",
	AssignmentExponentOp:       "**=",
	AssignmentLogicalAndOp:     "&&=",
	AssignmentLogicalOr:        "||=",
	AssignmentCoalesceOp:       "??=",
}

// BinaryExpression is a node for an ECMAScript binary expression statement.
//
// For example:
//
//     1 + 2
//
// Would be represented as:
//
//     BinaryExpression{
//         Operator: BinaryAddOp,
//         Left: NumberLiteral{Value: 1, ...},
//         Right: NumberLiteral{Value: 2, ...},
//     }
type BinaryExpression struct {
	BaseNode

	Operator BinaryOperator
	Left     Node
	Right    Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n BinaryExpression) ESTree() interface{} {
	nodeType := "BinaryExpression"
	if n.Operator == BinaryLogicalAndOp || n.Operator == BinaryLogicalOrOp {
		nodeType = "LogicalExpression"
	}

	return struct {
		Type     string      `json:"type"`
		Operator string      `json:"operator"`
		Left     interface{} `json:"left"`
		Right    interface{} `json:"right"`
	}{
		Type:     nodeType,
		Operator: estreeBinaryOpMap[n.Operator],
		Left:     estree(n.Left),
		Right:    estree(n.Right),
	}
}

// AssignmentExpression is the node for an ECMAScript assignment statement.
//
// For example:
//
//     i += 1
//
// Would be represented as:
//
//     AssignmentExpression{
//         Operator: AssignmentAddOp,
//         Left: Identifier{Name: "i"},
//         Right: NumberLiteral{Value: 2, ...},
//     }
type AssignmentExpression struct {
	BaseNode

	Operator AssignmentOperator
	Left     Node
	Right    Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n AssignmentExpression) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Operator string      `json:"operator"`
		Left     interface{} `json:"left"`
		Right    interface{} `json:"right"`
	}{
		Type:     "AssignmentExpression",
		Operator: estreeAssignOpMap[n.Operator],
		Left:     estree(n.Left),
		Right:    estree(n.Right),
	}
}
