package parser

import (
	"fmt"

	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

type exprOrder int

const (
	exprOrderComma exprOrder = iota
	exprOrderAssign
	exprOrderConditional
	exprOrderLogicalOr
	exprOrderLogicalAnd
	exprOrderBitwiseOr
	exprOrderBitwiseXor
	exprOrderBitwiseAnd
	exprOrderEqualityExpr
	exprOrderRelationalExpr
	exprOrderShiftExpr
	exprOrderAdditiveExpr
	exprOrderMultiplicativeExpr
	exprOrderExponentExpr
	exprOrderUnaryExpr
	exprOrderLHSExpr
	exprOrderCallExpr
	exprOrderMemberExpr
	exprOrderPrimaryExpr
)

type exprFlags int

const (
	exprFlagDisallowIn exprFlags = 1 << iota
)

// parseExpression parses an expression up to a certain level of operator
// precedence.
//
// For example, if you pass exprOrderPrimaryExpr, the lowest precedence, this
// function will only parse primary expressions; it will return an error if
// it is unable to. However, if you parse exprOrderLHSExpr, it will continue to
// parse operators and subexpressions until it either reaches an LHS operator
// or a token that it does not know how to parse on an expression boundary.
//
// Flags mainly control context-specific behavior, such as allowing the 'in'
// operator. Note that flags may or may not propagate to sub-expressions,
// depending on exactly what kind of sub-expression it is.
func (p *Parser) parseExpression(order exprOrder, flags exprFlags) ast.Node {
	var n ast.Node
	s := p.s.Location()
	t := p.ctx.keywordToIdentifier(p.s.Scan(), false)

	invalidprimary := func() {
		p.s.SyntaxError(fmt.Sprintf("unexpected token `%s`, expected primary expression", t.Source()))
	}

	wrap := func(n spannedNode, precedence exprOrder) ast.Node {
		if order > precedence {
			invalidprimary()
		}
		n.SetStart(s)
		n.SetEnd(p.s.Location())
		return n
	}

	wrapbinary := func(op ast.BinaryOperator, next exprOrder) ast.Node {
		m := ast.BinaryExpression{Operator: op}
		m.Left = n
		m.Right = p.parseExpression(next, flags)
		m.SetStart(s)
		m.SetEnd(p.s.Location())
		return m
	}

	wrapassign := func(op ast.AssignmentOperator, next exprOrder) ast.Node {
		m := ast.AssignmentExpression{Operator: op}
		m.Left = n
		m.Right = p.parseExpression(next, flags)
		m.SetStart(s)
		m.SetEnd(p.s.Location())
		return m
	}

	// Can't be Div/DivAssign here, relex as a regex. NOTE: if we are peeked
	// ahead at this point, this will fail.
	re := lexer.ReToken{}
	if t.Type == lexer.TokenPunctuatorDiv || t.Type == lexer.TokenPunctuatorDivAssign {
		re = p.s.ReScan()
		t = re.Token
	}

	switch t.Type {
	// Unary operators
	case lexer.TokenPunctuatorIncrement:
		// TODO: should add order for update operator?
		n = wrap(&ast.UpdateExpression{Operator: ast.UpdatePreIncrementOp, Argument: p.parseExpression(exprOrderLHSExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenPunctuatorDecrement:
		// TODO: should add order for update operator?
		n = wrap(&ast.UpdateExpression{Operator: ast.UpdatePreDecrementOp, Argument: p.parseExpression(exprOrderLHSExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenKeywordDelete:
		n = wrap(&ast.UnaryExpression{Operator: ast.UnaryDeleteOp, Argument: p.parseExpression(exprOrderUnaryExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenKeywordVoid:
		n = wrap(&ast.UnaryExpression{Operator: ast.UnaryVoidOp, Argument: p.parseExpression(exprOrderUnaryExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenKeywordTypeOf:
		n = wrap(&ast.UnaryExpression{Operator: ast.UnaryTypeOfOp, Argument: p.parseExpression(exprOrderUnaryExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenPunctuatorPlus:
		n = wrap(&ast.UnaryExpression{Operator: ast.UnaryPlusOp, Argument: p.parseExpression(exprOrderUnaryExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenPunctuatorMinus:
		n = wrap(&ast.UnaryExpression{Operator: ast.UnaryMinusOp, Argument: p.parseExpression(exprOrderUnaryExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenPunctuatorBitNot:
		n = wrap(&ast.UnaryExpression{Operator: ast.UnaryBitNotOp, Argument: p.parseExpression(exprOrderUnaryExpr, flags)}, exprOrderUnaryExpr)
	case lexer.TokenPunctuatorNot:
		n = wrap(&ast.UnaryExpression{Operator: ast.UnaryNotOp, Argument: p.parseExpression(exprOrderUnaryExpr, flags)}, exprOrderUnaryExpr)

	// Primary Expression
	case lexer.TokenKeywordThis:
		n = ast.ThisExpression{}
	case lexer.TokenIdentifier:
		n = ast.Identifier{Name: t.Literal}
	case lexer.TokenKeywordNull:
		n = ast.NullLiteral{}
	case lexer.TokenKeywordTrue:
		n = ast.BooleanLiteral{Value: true, Raw: t.Literal}
	case lexer.TokenKeywordFalse:
		n = ast.BooleanLiteral{Value: false, Raw: t.Literal}
	case lexer.TokenLiteralNumber:
		n = ast.NumberLiteral{Value: t.NumberConstant(), Raw: t.Literal}
	case lexer.TokenLiteralString:
		n = ast.StringLiteral{Value: t.StringConstant(), Raw: t.Literal}
	case lexer.TokenPunctuatorOpenBracket:
		n = p.parseArrayTail(s)
	case lexer.TokenPunctuatorOpenBrace:
		n = p.parseObjectTail(s)
	case lexer.TokenKeywordFunction:
		t = p.ctx.keywordToIdentifier(p.s.Scan(), false)
		name := ""
		if t.Type == lexer.TokenIdentifier {
			name = t.Literal
			t = p.s.Scan()
		}
		// TODO: generator support
		if t.Type != lexer.TokenPunctuatorOpenParen {
			p.s.SyntaxError("expected parameter list following function expression head")
		}
		params := p.parseParametersTail()
		body := p.parseBlock()
		m := ast.FunctionExpression{
			ID:     name,
			Params: params,
			Body:   body,
		}
		m.SetStart(s)
		m.SetEnd(p.s.Location())
		n = m
	case lexer.TokenKeywordNew:
		ctor := p.parseExpression(exprOrderMemberExpr, flags)
		m := ast.NewExpression{
			Callee: ctor,
		}
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorOpenParen {
			m.Arguments = p.parseArguments()
		}
		m.SetStart(s)
		m.SetEnd(p.s.Location())
		n = m
	case lexer.TokenKeywordAsync:
		panic("unimplemented: async function/generator expression")
	case lexer.TokenKeywordClass:
		panic("unimplemented: class expression")
	case lexer.TokenLiteralRegExp:
		m := ast.RegExpLiteral{
			Raw:     t.Literal,
			Pattern: re.Pattern,
			Flags:   re.Flags,
		}
		m.SetStart(s)
		m.SetEnd(p.s.Location())
		n = m
	case lexer.TokenLiteralTemplate:
		panic("unimplemented: template literal")
	case lexer.TokenPunctuatorOpenParen:
		// TODO: Need to support arrow functions! This is tricky because we
		// need to accept any production that would be a valid expression OR
		// argument list, then interpret it after we can peek to see if there
		// is a => after it or not.
		m := ast.ParenthesizedExpression{Expression: p.parseExpression(exprOrderComma, 0)}
		p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)` operator")
		m.SetStart(s)
		m.SetEnd(p.s.Location())
		n = m
	default:
		invalidprimary()
	}

	if order >= exprOrderPrimaryExpr {
		return n
	}

	for {
		// exprOrderLHSExpr
		t = p.s.PeekAt(0)
		if t.Type == lexer.TokenPunctuatorDot {
			p.s.ScanExpect(lexer.TokenPunctuatorDot, "expected `.` operator")
			m := ast.MemberExpression{
				Object:   n,
				Computed: false,
				Property: ast.Identifier{
					Name: p.forceScanIdent("expected property name after `.` operator"),
				},
			}
			m.SetStart(s)
			m.SetEnd(p.s.Location())
			n = m
			continue
		} else if t.Type == lexer.TokenPunctuatorOpenBracket {
			p.s.ScanExpect(lexer.TokenPunctuatorOpenBracket, "expected `[` operator")
			m := ast.MemberExpression{
				Object:   n,
				Computed: true,
				Property: p.parseExpression(exprOrderAssign, 0),
			}
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBracket, "expected `]` operator")
			m.SetStart(s)
			m.SetEnd(p.s.Location())
			n = m
			continue
		}
		if order >= exprOrderMemberExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorOpenParen {
			m := ast.CallExpression{
				Callee:    n,
				Arguments: p.parseArguments(),
			}
			m.SetStart(s)
			m.SetEnd(p.s.Location())
			n = m
			continue
		}
		if order >= exprOrderCallExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorOptionalChain {
			p.s.ScanExpect(lexer.TokenPunctuatorDot, "expected `?.` operator")
			if p.s.PeekAt(0).Type == lexer.TokenPunctuatorOpenBracket {
				p.s.ScanExpect(lexer.TokenPunctuatorOpenBracket, "expected `[` operator")
				m := ast.MemberExpression{
					Object:   n,
					Computed: true,
					Property: p.parseExpression(exprOrderAssign, 0),
					Optional: true,
				}
				p.s.ScanExpect(lexer.TokenPunctuatorCloseBracket, "expected `]` operator")
				m.SetStart(s)
				m.SetEnd(p.s.Location())
				n = m
			} else if p.s.PeekAt(0).Type == lexer.TokenPunctuatorOpenParen {
				m := ast.CallExpression{
					Callee:    n,
					Optional:  true,
					Arguments: p.parseArguments(),
				}
				m.SetStart(s)
				m.SetEnd(p.s.Location())
				n = m
			} else {
				m := ast.MemberExpression{
					Object:   n,
					Computed: false,
					Property: ast.Identifier{
						Name: p.forceScanIdent("expected property name after `.` operator"),
					},
					Optional: true,
				}
				m.SetStart(s)
				m.SetEnd(p.s.Location())
				n = m
			}
			continue
		}
		if order >= exprOrderLHSExpr {
			break
		}

		// TODO: should add order for update?
		if t.Type == lexer.TokenPunctuatorIncrement {
			p.s.ScanExpect(lexer.TokenPunctuatorIncrement, "expected `++` operator")
			n = wrap(&ast.UpdateExpression{Operator: ast.UpdatePostIncrementOp, Argument: n}, exprOrderUnaryExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorDecrement {
			p.s.ScanExpect(lexer.TokenPunctuatorDecrement, "expected `--` operator")
			n = wrap(&ast.UpdateExpression{Operator: ast.UpdatePostDecrementOp, Argument: n}, exprOrderUnaryExpr)
			continue
		}
		if order >= exprOrderUnaryExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorExponent {
			p.s.ScanExpect(lexer.TokenPunctuatorExponent, "expected `**` operator")
			n = wrapbinary(ast.BinaryExponentOp, exprOrderUnaryExpr)
			continue
		}
		if order >= exprOrderExponentExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorMult {
			p.s.ScanExpect(lexer.TokenPunctuatorMult, "expected `*` operator")
			n = wrapbinary(ast.BinaryMultOp, exprOrderExponentExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorDiv {
			p.s.ScanExpect(lexer.TokenPunctuatorDiv, "expected `/` operator")
			n = wrapbinary(ast.BinaryDivOp, exprOrderExponentExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorMod {
			p.s.ScanExpect(lexer.TokenPunctuatorMod, "expected `%` operator")
			n = wrapbinary(ast.BinaryModOp, exprOrderExponentExpr)
			continue
		}
		if order >= exprOrderMultiplicativeExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorPlus {
			p.s.ScanExpect(lexer.TokenPunctuatorPlus, "expected `+` operator")
			n = wrapbinary(ast.BinaryAddOp, exprOrderMultiplicativeExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorMinus {
			p.s.ScanExpect(lexer.TokenPunctuatorMinus, "expected `-` operator")
			n = wrapbinary(ast.BinarySubOp, exprOrderMultiplicativeExpr)
			continue
		}
		if order >= exprOrderAdditiveExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorLShift {
			p.s.ScanExpect(lexer.TokenPunctuatorLShift, "expected `<<` operator")
			n = wrapbinary(ast.BinaryLShiftOp, exprOrderAdditiveExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorRShift {
			p.s.ScanExpect(lexer.TokenPunctuatorRShift, "expected `>>` operator")
			n = wrapbinary(ast.BinaryRShiftOp, exprOrderAdditiveExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorUnsignedRShift {
			p.s.ScanExpect(lexer.TokenPunctuatorUnsignedRShift, "expected `>>>` operator")
			n = wrapbinary(ast.BinaryUnsignedRShiftOp, exprOrderAdditiveExpr)
			continue
		}
		if order >= exprOrderShiftExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorLessThan {
			p.s.ScanExpect(lexer.TokenPunctuatorLessThan, "expected `<` operator")
			n = wrapbinary(ast.BinaryLessThanOp, exprOrderShiftExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorGreaterThan {
			p.s.ScanExpect(lexer.TokenPunctuatorGreaterThan, "expected `>` operator")
			n = wrapbinary(ast.BinaryGreaterThanOp, exprOrderShiftExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorLessThanEqual {
			p.s.ScanExpect(lexer.TokenPunctuatorLessThanEqual, "expected `<=` operator")
			n = wrapbinary(ast.BinaryLessThanEqualOp, exprOrderShiftExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorGreaterThanEqual {
			p.s.ScanExpect(lexer.TokenPunctuatorGreaterThanEqual, "expected `>=` operator")
			n = wrapbinary(ast.BinaryGreaterThanEqualOp, exprOrderShiftExpr)
			continue
		} else if t.Type == lexer.TokenKeywordInstanceOf {
			p.s.ScanExpect(lexer.TokenKeywordInstanceOf, "expected `instanceof` operator")
			n = wrapbinary(ast.BinaryInstanceOfOp, exprOrderShiftExpr)
			continue
		} else if flags&exprFlagDisallowIn == 0 && t.Type == lexer.TokenKeywordIn {
			p.s.ScanExpect(lexer.TokenKeywordIn, "expected `in` operator")
			n = wrapbinary(ast.BinaryInOp, exprOrderShiftExpr)
			continue
		}
		if order >= exprOrderRelationalExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorEqual {
			p.s.ScanExpect(lexer.TokenPunctuatorEqual, "expected `==` operator")
			n = wrapbinary(ast.BinaryEqualOp, exprOrderRelationalExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorNotEqual {
			p.s.ScanExpect(lexer.TokenPunctuatorNotEqual, "expected `!=` operator")
			n = wrapbinary(ast.BinaryNotEqualOp, exprOrderRelationalExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorStrictEqual {
			p.s.ScanExpect(lexer.TokenPunctuatorStrictEqual, "expected `===` operator")
			n = wrapbinary(ast.BinaryStrictEqualOp, exprOrderRelationalExpr)
			continue
		} else if t.Type == lexer.TokenPunctuatorStrictNotEqual {
			p.s.ScanExpect(lexer.TokenPunctuatorStrictNotEqual, "expected `!==` operator")
			n = wrapbinary(ast.BinaryStrictNotEqualOp, exprOrderRelationalExpr)
			continue
		}
		if order >= exprOrderEqualityExpr {
			break
		}

		if t.Type == lexer.TokenPunctuatorBitAnd {
			p.s.ScanExpect(lexer.TokenPunctuatorBitAnd, "expected `&` operator")
			n = wrapbinary(ast.BinaryBitAndOp, exprOrderEqualityExpr)
			continue
		}
		if order >= exprOrderBitwiseAnd {
			break
		}

		if t.Type == lexer.TokenPunctuatorBitXor {
			p.s.ScanExpect(lexer.TokenPunctuatorBitXor, "expected `^` operator")
			n = wrapbinary(ast.BinaryBitXorOp, exprOrderBitwiseAnd)
			continue
		}
		if order >= exprOrderBitwiseXor {
			break
		}

		if t.Type == lexer.TokenPunctuatorBitOr {
			p.s.ScanExpect(lexer.TokenPunctuatorBitOr, "expected `|` operator")
			n = wrapbinary(ast.BinaryBitXorOp, exprOrderBitwiseXor)
			continue
		}
		if order >= exprOrderBitwiseOr {
			break
		}

		if t.Type == lexer.TokenPunctuatorLogicalAnd {
			p.s.ScanExpect(lexer.TokenPunctuatorLogicalAnd, "expected `&&` operator")
			n = wrapbinary(ast.BinaryLogicalAndOp, exprOrderBitwiseOr)
			continue
		}
		if order >= exprOrderLogicalAnd {
			break
		}

		if t.Type == lexer.TokenPunctuatorLogicalOr {
			p.s.ScanExpect(lexer.TokenPunctuatorLogicalOr, "expected `||` operator")
			n = wrapbinary(ast.BinaryLogicalOrOp, exprOrderLogicalAnd)
			continue
		} else if t.Type == lexer.TokenPunctuatorNullCoalesce {
			p.s.ScanExpect(lexer.TokenPunctuatorNullCoalesce, "expected `??` operator")
			n = wrapbinary(ast.BinaryCoalesceOp, exprOrderLogicalAnd)
			continue
		}
		if order >= exprOrderLogicalOr {
			break
		}

		if t.Type == lexer.TokenPunctuatorQuestionMark {
			p.s.ScanExpect(lexer.TokenPunctuatorQuestionMark, "expected `?` operator in conditional expression")
			a := p.parseExpression(exprOrderAssign, 0)
			p.s.ScanExpect(lexer.TokenPunctuatorColon, "expected `:` operator in conditional expression")
			b := p.parseExpression(exprOrderAssign, 0)
			m := ast.ConditionalExpression{
				Test:       n,
				Consequent: a,
				Alternate:  b,
			}
			m.SetStart(s)
			m.SetEnd(p.s.Location())
			n = m
			continue
		}
		if order >= exprOrderConditional {
			break
		}

		if t.Type == lexer.TokenPunctuatorAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorAssign, "expected `=` operator")
			n = wrapassign(ast.AssignmentOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorMultAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorMultAssign, "expected `*=` operator")
			n = wrapassign(ast.AssignmentMultOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorDivAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorDivAssign, "expected `/=` operator")
			n = wrapassign(ast.AssignmentDivOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorModAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorModAssign, "expected `%=` operator")
			n = wrapassign(ast.AssignmentModOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorPlusAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorPlusAssign, "expected `+=` operator")
			n = wrapassign(ast.AssignmentAddOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorMinusAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorMinusAssign, "expected `-=` operator")
			n = wrapassign(ast.AssignmentSubOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorLShiftAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorLShiftAssign, "expected `<<=` operator")
			n = wrapassign(ast.AssignmentLShiftOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorRShiftAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorRShiftAssign, "expected `>>=` operator")
			n = wrapassign(ast.AssignmentRShiftOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorUnsignedRShiftAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorUnsignedRShiftAssign, "expected `>>>=` operator")
			n = wrapassign(ast.AssignmentUnsignedRShiftOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorBitAndAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorBitAndAssign, "expected `&=` operator")
			n = wrapassign(ast.AssignmentBitAndOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorBitXorAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorBitXorAssign, "expected `^=` operator")
			n = wrapassign(ast.AssignmentBitXorOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorBitOrAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorBitOrAssign, "expected `|=` operator")
			n = wrapassign(ast.AssignmentBitOrOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorExponentAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorExponentAssign, "expected `**=` operator")
			n = wrapassign(ast.AssignmentExponentOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorLogicalAndAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorLogicalAndAssign, "expected `&&=` operator")
			n = wrapassign(ast.AssignmentLogicalAndOp, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorLogicalOrAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorLogicalOrAssign, "expected `||=` operator")
			n = wrapassign(ast.AssignmentLogicalOr, exprOrderAssign)
			continue
		} else if t.Type == lexer.TokenPunctuatorNullCoalesceAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorNullCoalesceAssign, "expected `??=` operator")
			n = wrapassign(ast.AssignmentCoalesceOp, exprOrderAssign)
			continue
		}
		if order >= exprOrderAssign {
			break
		}
		if t.Type == lexer.TokenPunctuatorComma {
			p.s.ScanExpect(lexer.TokenPunctuatorComma, "expected `,` operator")
			if seq, ok := n.(ast.SequenceExpression); ok {
				seq.Expressions = append(seq.Expressions, p.parseExpression(exprOrderAssign, flags))
				n = seq
			} else {
				seq := ast.SequenceExpression{Expressions: []ast.Node{n}}
				seq.SetStart(s)
				seq.SetEnd(p.s.Location())
				seq.Expressions = append(seq.Expressions, p.parseExpression(exprOrderAssign, flags))
				n = seq
			}
			continue
		}
		if order >= exprOrderComma {
			break
		}

		// Matched nothing; end of expression.
		break
	}

	return n
}

// Parses an array assuming a `[` was already consumed.
func (p *Parser) parseArrayTail(start ast.Location) ast.Node {
	n := ast.ArrayExpression{}
	n.SetStart(start)
	defer p.setEnd(&n)

	for {
		for p.s.PeekAt(0).Type == lexer.TokenPunctuatorComma {
			n.Elements = append(n.Elements, nil)
			p.s.ScanExpect(lexer.TokenPunctuatorComma, "expected `,`")
		}
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseBracket {
			break
		}
		n.Elements = append(n.Elements, p.parseExpression(exprOrderAssign, 0))
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorComma {
			p.s.ScanExpect(lexer.TokenPunctuatorComma, "expected `,`")
		}
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseBracket {
			break
		}
	}

	p.s.ScanExpect(lexer.TokenPunctuatorCloseBracket, "expected `]`")
	return n
}

// Parses an object assuming a `{` was already consumed.
func (p *Parser) parseObjectTail(start ast.Location) ast.Node {
	n := ast.ObjectExpression{}
	n.SetStart(start)
	defer p.setEnd(&n)

	atEndOfPropertyKey := func() bool {
		// Colon ends the property key when not using shorthand, otherwise
		// comma or close brace could also end the property key. Finally, when
		// using method shorthand, an open paren can also end the key.
		t := p.s.PeekAt(0).Type
		return t == lexer.TokenPunctuatorColon ||
			t == lexer.TokenPunctuatorComma ||
			t == lexer.TokenPunctuatorCloseBrace ||
			t == lexer.TokenPunctuatorOpenParen
	}

	for {
		// On first iteration: ends empty object. On other iterations: ends
		// object after trailing comma.
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseBrace {
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBrace, "expected `}`")
			return n
		}

		// Keeps track of specifiers that are specified for the method
		// shorthand.
		async := false
		generator := false

		prop := ast.Property{Kind: ast.InitProperty}

		// Until we get to the identifier, keep track of the position of the
		// last token. We need this to know the identifier span.
		pos := p.s.Location()

		// Handle specifiers before keyword.
		t := p.s.Scan()

		// We need to special case if we have started on a computed key because
		// an arbitrary number of tokens will be the computed expression.
		startedOnComputedKey := t.Type == lexer.TokenPunctuatorOpenBracket

		// If we did not start on a computed key, and the last token retrieved
		// did not put us on a token that ends the property, we can look for
		// specifiers.
		if !startedOnComputedKey && !atEndOfPropertyKey() {
			switch t.Type {
			case lexer.TokenKeywordGet:
				prop.Kind = ast.GetProperty

			case lexer.TokenKeywordSet:
				prop.Kind = ast.SetProperty

			case lexer.TokenKeywordAsync:
				async = true

				// Async generator (async *)
				if p.s.PeekAt(0).Type == lexer.TokenPunctuatorMult {
					generator = true

					// Don't need to update position yet; it'll get taken care
					// of below when the next token is read.
					t = p.s.Scan()
				}

			case lexer.TokenPunctuatorMult:
				generator = true

			default:
				// We don't know what is wrong here.
				// TODO: better error message heuristics here?
				p.s.SyntaxError("invalid property syntax")
			}

			pos = p.s.Location()
			t = p.s.Scan()
		}

		// Next, handle identifier...
		t = p.ctx.keywordToIdentifier(t, true)
		switch t.Type {
		case lexer.TokenIdentifier:
			// Normal identifier.
			id := ast.Identifier{Name: t.Literal}
			id.SetStart(pos)
			id.SetEnd(p.s.Location())
			prop.Key = id

		case lexer.TokenLiteralString:
			// String literal.
			id := ast.StringLiteral{Value: t.StringConstant(), Raw: t.Literal}
			id.SetStart(pos)
			id.SetEnd(p.s.Location())
			prop.Key = id

		case lexer.TokenLiteralNumber:
			// Number literal.
			id := ast.NumberLiteral{Value: t.NumberConstant(), Raw: t.Literal}
			id.SetStart(pos)
			id.SetEnd(p.s.Location())
			prop.Key = id

		case lexer.TokenPunctuatorOpenBracket:
			// Computed identifier.
			prop.Computed = true
			prop.Key = p.parseExpression(exprOrderComma, 0)
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBracket, "expected `]`")

		default:
			p.s.SyntaxError("expected property name")
		}

		peek := p.s.PeekAt(0)

		switch {
		case prop.Kind == ast.GetProperty || prop.Kind == ast.SetProperty:
			// Getter/setter
			fn := ast.FunctionExpression{}
			fn.Params = p.parseParameters()
			fn.Body = p.parseBlock()
			fn.SetEnd(p.s.Location())
			prop.Value = fn

		case peek.Type == lexer.TokenPunctuatorColon:
			// Normal init property
			if async || generator {
				p.s.SyntaxError("expected method")
			}

			p.s.ScanExpect(lexer.TokenPunctuatorColon, "expected `:`")
			prop.Value = p.parseExpression(exprOrderAssign, 0)

		case peek.Type == lexer.TokenPunctuatorOpenParen:
			// Method short-hand property
			ctx := p.ctx
			p.ctx.async = async
			p.ctx.generator = generator

			fn := ast.FunctionExpression{
				Async:     async,
				Generator: generator,
			}

			fn.SetStart(p.s.Location())
			fn.Params = p.parseParameters()
			fn.Body = p.parseBlock()
			fn.SetEnd(p.s.Location())

			prop.Value = fn
			prop.Method = true

			p.ctx = ctx

		case peek.Type == lexer.TokenPunctuatorComma ||
			peek.Type == lexer.TokenPunctuatorCloseBrace:
			// Shorthand syntax. We don't need to do anything, but we should
			// disallow this from happening with a computed property.
			if prop.Computed {
				p.s.SyntaxError("shorthand not allowed for computed property")
			}

			// We also should not allow this when async/generator is specified.
			if async || generator {
				p.s.SyntaxError("expected method")
			}

		default:
			p.s.SyntaxError("expected `,` or `}`")
		}

		n.Properties = append(n.Properties, prop)

		// Object ends after a property.
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseBrace {
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBrace, "expected `}`")
			return n
		}

		// Comma before next property, or before ending after a trailing comma.
		p.s.ScanExpect(lexer.TokenPunctuatorComma, "expected `,` or `}`")
	}
}

// Parses arguments.
func (p *Parser) parseArguments() []ast.Node {
	n := []ast.Node{}

	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(`")
	if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseParen {
		p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")
		return n
	}
	for {
		// TODO: spread elements
		n = append(n, p.parseExpression(exprOrderAssign, 0))
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorComma {
			p.s.ScanExpect(lexer.TokenPunctuatorComma, "expected `,`")
		}
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseParen {
			p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")
			return n
		}
	}
}

// Parses parameters.
func (p *Parser) parseParameters() ast.FormalParameters {
	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(`")
	return p.parseParametersTail()
}

func (p *Parser) parseParametersTail() ast.FormalParameters {
	n := ast.FormalParameters{}

	for {
		b := ast.BindingElement{}
		t := p.ctx.keywordToIdentifier(p.s.Scan(), false)
		switch t.Type {
		case lexer.TokenIdentifier:
			b.Value.Identifier = t.Literal

		case lexer.TokenPunctuatorCloseParen:
			return n

		case lexer.TokenPunctuatorOpenBracket:
			b.Value.ArrayPattern = p.parseArrayBindingPattern()

		case lexer.TokenPunctuatorOpenBrace:
			b.Value.ObjectPattern = p.parseObjectBindingPattern()

		case lexer.TokenPunctuatorEllipsis:
			n.RestParameter = p.scanIdent("expected identifier for rest parameter")
			p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected closing paren")
			return n

		default:
			p.s.SyntaxError(fmt.Sprintf("unexpected token in formal parameter list: %s", p.s.Scan().Source()))
		}

		// Default syntax
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorAssign, "expected default assignment `=`")
			b.Init = p.parseExpression(exprOrderAssign, 0)
		}

		n.Parameters = append(n.Parameters, b)

		t = p.s.Scan()
		switch t.Type {
		case lexer.TokenPunctuatorComma:
			continue

		case lexer.TokenPunctuatorCloseParen:
			return n

		default:
			p.s.SyntaxError(fmt.Sprintf("expected `,` or `)`, but got: %s", t.Source()))
		}
	}
}
