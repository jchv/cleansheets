package parser

import (
	"fmt"

	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

func (p *Parser) parseStatementItem() ast.Node {
	if n := p.parseStatement(); n != nil {
		return n
	}
	if n := p.parseDeclaration(); n != nil {
		return n
	}
	p.s.SyntaxError("expected declaration or statement")
	return nil
}

func (p *Parser) parseStatement() ast.Node {
	switch p.s.PeekAt(0).Type {
	case lexer.TokenPunctuatorOpenBrace:
		return p.parseBlock()
	case lexer.TokenKeywordVar:
		return p.parseVariableStatement()
	case lexer.TokenPunctuatorSemicolon:
		return p.parseEmptyExpression()
	case lexer.TokenKeywordIf:
		return p.parseIfStatement()
	case
		// Unary operators
		lexer.TokenPunctuatorIncrement, lexer.TokenPunctuatorDecrement, lexer.TokenKeywordDelete,
		lexer.TokenKeywordVoid, lexer.TokenKeywordTypeOf, lexer.TokenPunctuatorPlus,
		lexer.TokenPunctuatorMinus, lexer.TokenPunctuatorBitNot, lexer.TokenPunctuatorNot,
		// Primary Expression
		// Note the absence of: `{`, `function`, and `class`. `async` is
		// allowed if not followed by `function` with no newline. `let` is
		// allowed if not followed by `[`.
		// Additionally, we handle expressions starting with identifiers further down.
		lexer.TokenKeywordThis, lexer.TokenKeywordNull, lexer.TokenKeywordTrue,
		lexer.TokenKeywordFalse, lexer.TokenKeywordNew,
		lexer.TokenLiteralNumber, lexer.TokenLiteralString,
		lexer.TokenLiteralTemplate,
		lexer.TokenPunctuatorOpenBracket, lexer.TokenKeywordAsync, lexer.TokenKeywordLet,
		lexer.TokenPunctuatorOpenParen,
		// These will get relexed as a regexp, so they are valid to begin an expression.
		lexer.TokenPunctuatorDiv, lexer.TokenPunctuatorDivAssign:
		// Async function declaration (async [no line terminator] function)
		if p.s.PeekAt(0).Type == lexer.TokenKeywordAsync && p.s.PeekAt(1).Type == lexer.TokenKeywordFunction && !p.s.PeekAt(1).NewLine {
			return nil
		}
		if p.s.PeekAt(0).Type == lexer.TokenKeywordLet {
			if p.s.PeekAt(1).Type == lexer.TokenPunctuatorOpenBracket {
				// Array destructuring let (let [)
				return nil
			} else if p.s.PeekAt(1).Type == lexer.TokenPunctuatorOpenBrace {
				// Object destructuring let (let {)
				return nil
			} else if p.ctx.keywordToIdentifier(p.s.PeekAt(1), true).Type == lexer.TokenIdentifier {
				// Let with identifier (let ident)
				return nil
			}
		}
		return p.parseExpressionStatement()
	case lexer.TokenKeywordDo:
		return p.parseDoWhileStatement()
	case lexer.TokenKeywordWhile:
		return p.parseWhileStatement()
	case lexer.TokenKeywordFor:
		return p.parseForStatement()
	case lexer.TokenKeywordSwitch:
		return p.parseSwitchStatement()
	case lexer.TokenKeywordContinue:
		return p.parseContinueStatement()
	case lexer.TokenKeywordBreak:
		return p.parseBreakStatement()
	case lexer.TokenKeywordReturn:
		return p.parseReturnStatement()
	case lexer.TokenKeywordWith:
		return p.parseWithStatement()
	case lexer.TokenKeywordThrow:
		return p.parseThrowStatement()
	case lexer.TokenKeywordTry:
		return p.parseTryStatement()
	case lexer.TokenKeywordDebugger:
		return p.parseDebuggerStatement()
	case lexer.TokenIdentifier:
		fallthrough
	default:
		if p.ctx.keywordToIdentifier(p.s.PeekAt(0), false).Type == lexer.TokenIdentifier {
			if p.s.PeekAt(1).Type == lexer.TokenPunctuatorColon {
				return p.parseLabelledStatement()
			}
			return p.parseExpressionStatement()
		}
	}
	return nil
}

func (p *Parser) parseExpressionStatement() ast.Node {
	expr := p.parseExpression(exprOrderComma, 0)
	n := ast.ExpressionStatement{Expression: expr}
	n.SetStart(expr.Span().Start)
	n.SetEnd(expr.Span().End)
	p.expectSemicolon()
	return n
}

func (p *Parser) parseBlockOrShorthand() ast.Node {
	if p.s.PeekAt(0).Type == lexer.TokenPunctuatorOpenBrace {
		return p.parseBlock()
	} else {
		return p.parseExpression(exprOrderConditional, 0)
	}
}

func (p *Parser) parseBlock() ast.BlockStatement {
	n := ast.BlockStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenPunctuatorOpenBrace, "expected block opening brace `{`")

	// Early exit for empty block.
	if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseBrace {
		p.s.ScanExpect(lexer.TokenPunctuatorCloseBrace, "expected statement, declaration, or closing brace `}`")
		return n
	}

	ctx := p.ctx

	// Parse first statement so we can parse directives out of it.
	stmt := p.parseStatementItem()
	if expr, ok := stmt.(ast.ExpressionStatement); ok {
		if str, ok := expr.Expression.(ast.StringLiteral); ok {
			if str.Value == "use strict" {
				ctx.strictMode = true
				expr.Directive = "use strict"
			}
		}
		stmt = expr
	}
	n.Body = append(n.Body, stmt)

	for {
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorCloseBrace {
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBrace, "expected statement, declaration, or closing brace `}`")
			break
		}
		n.Body = append(n.Body, p.parseStatementItem())
	}

	p.ctx = ctx

	return n
}

func (p *Parser) parseVariableStatement() ast.VariableDeclaration {
	n := p.parseVariableStatementNoSemicolon()
	p.expectSemicolon()
	p.setEnd(&n)
	return n
}

func (p *Parser) parseVariableStatementNoSemicolon() ast.VariableDeclaration {
	n := ast.VariableDeclaration{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordVar, "expected variable declaration")
	n.Declarations = p.parseVariableDeclarations()
	return n
}

func (p *Parser) parseVariableDeclarations() []ast.VariableDeclarator {
	v := []ast.VariableDeclarator{}
	for {
		v = append(v, p.parseVariableDeclaration())
		if p.s.PeekAt(0).Type != lexer.TokenPunctuatorComma {
			break
		}
		p.s.ScanExpect(lexer.TokenPunctuatorComma, "expected comma")
	}
	return v
}

func (p *Parser) parseVariableDeclaration() ast.VariableDeclarator {
	v := ast.VariableDeclarator{}

	t := p.ctx.keywordToIdentifier(p.s.PeekAt(0), false)
	switch t.Type {
	case lexer.TokenIdentifier:
		v.ID.Identifier = p.scanIdent("expected variable identifier")
	case lexer.TokenPunctuatorOpenBracket:
		v.ID.ArrayPattern = p.parseArrayBindingPattern()
	case lexer.TokenPunctuatorOpenBrace:
		v.ID.ObjectPattern = p.parseObjectBindingPattern()
	default:
		p.s.SyntaxError(fmt.Sprintf("unexpected token in variable declaration: %s", p.s.Scan().Source()))
	}

	if p.s.PeekAt(0).Type == lexer.TokenPunctuatorAssign {
		p.s.ScanExpect(lexer.TokenPunctuatorAssign, "expected `=`")
		v.Init = p.parseExpression(exprOrderAssign, 0)
	}

	return v
}

func (p *Parser) parseArrayBindingPattern() *ast.ArrayBindingPattern {
	p.s.ScanExpect(lexer.TokenPunctuatorOpenBracket, "expected array binding pattern")
	return p.parseArrayBindingPatternTail()
}

func (p *Parser) parseArrayBindingPatternTail() *ast.ArrayBindingPattern {
	n := &ast.ArrayBindingPattern{}
	for {
		b := ast.BindingElement{}
		t := p.ctx.keywordToIdentifier(p.s.Scan(), false)
		switch t.Type {
		case lexer.TokenIdentifier:
			b.Value.Identifier = t.Literal

		case lexer.TokenPunctuatorComma:
			// Elision
			n.Elements = append(n.Elements, b)
			continue

		case lexer.TokenPunctuatorCloseBracket:
			return n

		case lexer.TokenPunctuatorOpenBracket:
			b.Value.ArrayPattern = p.parseArrayBindingPatternTail()

		case lexer.TokenPunctuatorOpenBrace:
			b.Value.ObjectPattern = p.parseObjectBindingPatternTail()

		case lexer.TokenPunctuatorEllipsis:
			t := p.ctx.keywordToIdentifier(p.s.PeekAt(0), false)
			switch t.Type {
			case lexer.TokenIdentifier:
				n.RestElement.Identifier = p.scanIdent("expected variable identifier")
			case lexer.TokenPunctuatorOpenBracket:
				n.RestElement.ArrayPattern = p.parseArrayBindingPattern()
			case lexer.TokenPunctuatorOpenBrace:
				n.RestElement.ObjectPattern = p.parseObjectBindingPattern()
			default:
				p.s.SyntaxError(fmt.Sprintf("unexpected token in rest pattern: %s", p.s.Scan().Source()))
			}
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBracket, "expected closing braket")
			return n

		default:
			p.s.SyntaxError(fmt.Sprintf("unexpected token in array binding pattern: %s", p.s.Scan().Source()))
		}

		// Default syntax
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorAssign, "expected default assignment `=`")
			b.Init = p.parseExpression(exprOrderAssign, 0)
		}

		n.Elements = append(n.Elements, b)

		t = p.s.Scan()
		switch t.Type {
		case lexer.TokenPunctuatorComma:
			continue

		case lexer.TokenPunctuatorCloseBracket:
			return n

		default:
			p.s.SyntaxError(fmt.Sprintf("expected `,` or `}`, but got: %s", t.Source()))
		}
	}
}

func (p *Parser) parseObjectBindingPattern() *ast.ObjectBindingPattern {
	p.s.ScanExpect(lexer.TokenPunctuatorOpenBrace, "expected object binding pattern")
	return p.parseObjectBindingPatternTail()
}

func (p *Parser) parseObjectBindingPatternTail() *ast.ObjectBindingPattern {
	n := &ast.ObjectBindingPattern{}
	for {
		b := ast.BindingProperty{}
		t := p.ctx.keywordToIdentifier(p.s.Scan(), false)
		switch t.Type {
		case lexer.TokenIdentifier:
			b.PropertyName = t.Literal

		case lexer.TokenPunctuatorEllipsis:
			n.RestElement = p.scanIdent("expected rest identifier")
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBrace, "expected closing brace")
			return n

		case lexer.TokenPunctuatorCloseBrace:
			return n

		default:
			p.s.SyntaxError(fmt.Sprintf("expected property name, `...`, or `}`, but got: %s", t.Source()))
		}

		// Binding syntax
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorColon {
			p.s.ScanExpect(lexer.TokenPunctuatorColon, "expected binding `:`")
			t = p.ctx.keywordToIdentifier(p.s.Scan(), false)
			switch t.Type {
			case lexer.TokenIdentifier:
				b.Value.Identifier = t.Literal

			case lexer.TokenPunctuatorOpenBracket:
				b.Value.ArrayPattern = p.parseArrayBindingPatternTail()

			case lexer.TokenPunctuatorOpenBrace:
				b.Value.ObjectPattern = p.parseObjectBindingPatternTail()

			default:
				p.s.SyntaxError(fmt.Sprintf("unexpected token in object binding pattern: %s", p.s.Scan().Source()))
			}
		}

		// Default syntax
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorAssign {
			p.s.ScanExpect(lexer.TokenPunctuatorAssign, "expected default assignment `=`")
			b.Init = p.parseExpression(exprOrderAssign, 0)
		}

		n.Properties = append(n.Properties, b)

		t = p.s.Scan()
		switch t.Type {
		case lexer.TokenPunctuatorComma:
			continue

		case lexer.TokenPunctuatorCloseBrace:
			return n

		default:
			p.s.SyntaxError(fmt.Sprintf("expected `,` or `}`, but got: %s", t.Source()))
		}
	}
}

func (p *Parser) parseEmptyExpression() ast.Node {
	n := ast.EmptyStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.expectSemicolon()
	return n
}

func (p *Parser) parseIfStatement() ast.Node {
	n := ast.IfStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordIf, "expected `if` statement")
	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(` after `if`")
	n.Test = p.parseExpression(exprOrderComma, 0)
	p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")
	n.Consequent = p.parseStatement()
	if p.s.PeekAt(0).Type == lexer.TokenKeywordElse {
		p.s.ScanExpect(lexer.TokenKeywordElse, "expected `else`")
		n.Alternate = p.parseStatement()
	}
	return n
}

func (p *Parser) parseDoWhileStatement() ast.Node {
	n := ast.DoWhileStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordDo, "expected `do` statement")
	n.Body = p.parseStatement()
	p.s.ScanExpect(lexer.TokenKeywordWhile, "expected `while` in do/while statement")
	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(` in `while` of do/while statement")
	n.Test = p.parseExpression(exprOrderComma, 0)
	p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)` in `while` of do/while statement")
	p.expectSemicolon()
	return n
}

func (p *Parser) parseWhileStatement() ast.Node {
	n := ast.WhileStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordWhile, "expected `while` statement")
	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(` in `while` of do/while statement")
	n.Test = p.parseExpression(exprOrderComma, 0)
	p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)` in `while` of do/while statement")
	n.Body = p.parseStatement()
	return n
}

func (p *Parser) parseForStatement() ast.Node {
	n := ast.ForStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordFor, "expected `for` statement")
	// TODO: async
	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(`")

	t := p.s.PeekAt(0)
	// TODO: let, const, more of/in cases, etc.
	if t.Type == lexer.TokenPunctuatorSemicolon {
		n.Init = nil
		p.expectSemicolon()
	} else {
		var v ast.Node
		if t.Type == lexer.TokenKeywordVar {
			v = p.parseVariableStatementNoSemicolon()
		} else {
			v = p.parseExpression(exprOrderComma, exprFlagDisallowIn)
		}
		// for in/of
		switch p.s.PeekAt(0).Type {
		case lexer.TokenKeywordIn:
			p.s.ScanExpect(lexer.TokenKeywordIn, "expected `in`")
			m := ast.ForInStatement{
				Left:  v,
				Right: p.parseExpression(exprOrderComma, 0),
			}
			m.SetStart(n.Span().Start)
			p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")
			m.Body = p.parseStatement()
			p.setEnd(&m)
			return m

		case lexer.TokenKeywordOf:
			p.s.ScanExpect(lexer.TokenKeywordOf, "expected `of`")
			m := ast.ForOfStatement{
				Left:  v,
				Right: p.parseExpression(exprOrderComma, 0),
			}
			m.SetStart(n.Span().Start)
			p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")
			m.Body = p.parseStatement()
			p.setEnd(&m)
			return m
		}
		n.Init = v
		p.expectSemicolon()
	}
	if p.s.PeekAt(0).Type != lexer.TokenPunctuatorSemicolon {
		n.Test = p.parseExpression(exprOrderComma, 0)
	}
	p.expectSemicolon()
	if p.s.PeekAt(0).Type != lexer.TokenPunctuatorCloseParen {
		n.Update = p.parseExpression(exprOrderComma, 0)
	}
	p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")
	n.Body = p.parseStatement()
	return n
}

func (p *Parser) parseSwitchStatement() ast.Node {
	n := ast.SwitchStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordSwitch, "expected `switch` statement")
	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(`")
	n.Discriminant = p.parseExpression(exprOrderComma, 0)
	p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")

	p.s.ScanExpect(lexer.TokenPunctuatorOpenBrace, "expected `{`")
	for {
		switch p.s.PeekAt(0).Type {
		case lexer.TokenKeywordCase:
			p.s.ScanExpect(lexer.TokenKeywordCase, "expected `case`")
			c := ast.SwitchCase{
				Test: p.parseExpression(exprOrderComma, 0),
			}
			p.s.ScanExpect(lexer.TokenPunctuatorColon, "expected `:`")
		caseStatements:
			for {
				switch p.s.PeekAt(0).Type {
				case lexer.TokenKeywordCase, lexer.TokenKeywordDefault, lexer.TokenPunctuatorCloseBrace:
					break caseStatements
				default:
					c.Consequent = append(c.Consequent, p.parseStatement())
				}
			}
			n.Cases = append(n.Cases, c)

		case lexer.TokenKeywordDefault:
			c := ast.SwitchCase{}
			p.s.ScanExpect(lexer.TokenKeywordDefault, "expected `default`")
			p.s.ScanExpect(lexer.TokenPunctuatorColon, "expected `:`")
		defaultStatements:
			for {
				switch p.s.PeekAt(0).Type {
				case lexer.TokenKeywordCase, lexer.TokenKeywordDefault, lexer.TokenPunctuatorCloseBrace:
					break defaultStatements
				default:
					c.Consequent = append(c.Consequent, p.parseStatement())
				}
			}
			n.Cases = append(n.Cases, c)

		case lexer.TokenPunctuatorCloseBrace:
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBrace, "expected `}`")
			return n
		}
	}
}

func (p *Parser) parseContinueStatement() ast.Node {
	n := ast.ContinueStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordContinue, "expected continue statement")
	t := p.ctx.keywordToIdentifier(p.s.PeekAt(0), false)
	if t.NewLine || t.Type != lexer.TokenIdentifier {
		p.expectSemicolon()
		return n
	}
	n.Label = p.scanIdent("expected identifier")

	p.expectSemicolon()
	return n
}

func (p *Parser) parseBreakStatement() ast.Node {
	n := ast.BreakStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordBreak, "expected break statement")
	t := p.ctx.keywordToIdentifier(p.s.PeekAt(0), false)
	if t.NewLine || t.Type != lexer.TokenIdentifier {
		p.expectSemicolon()
		return n
	}
	n.Label = p.scanIdent("expected identifier")

	p.expectSemicolon()
	return n
}

func (p *Parser) parseReturnStatement() ast.Node {
	n := ast.ReturnStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordReturn, "expected return statement")
	t := p.s.PeekAt(0)
	if t.NewLine || t.Type == lexer.TokenPunctuatorSemicolon || t.Type == lexer.TokenPunctuatorCloseBrace {
		p.expectSemicolon()
		return n
	}

	n.Argument = p.parseExpression(exprOrderComma, 0)
	p.expectSemicolon()
	return n
}

func (p *Parser) parseWithStatement() ast.Node {
	panic("unimplemented")
}

func (p *Parser) parseThrowStatement() ast.Node {
	n := ast.ThrowStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordThrow, "expected throw statement")
	if p.s.PeekAt(0).NewLine {
		p.s.SyntaxError("illegal newline after throw")
	}

	n.Argument = p.parseExpression(exprOrderComma, 0)
	p.expectSemicolon()
	return n
}

func (p *Parser) parseTryStatement() ast.Node {
	n := ast.TryStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordTry, "expected try statement")
	n.Block = p.parseBlock()
	if p.s.PeekAt(0).Type == lexer.TokenKeywordCatch {
		p.s.ScanExpect(lexer.TokenKeywordCatch, "expected catch statement")
		h := ast.CatchClause{}
		h.SetStart(p.s.Location())
		if p.s.PeekAt(0).Type == lexer.TokenPunctuatorOpenParen {
			p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected `(`")
			h.Param = p.parseCatchParameter()
			p.s.ScanExpect(lexer.TokenPunctuatorCloseParen, "expected `)`")
		}
		h.SetEnd(p.s.Location())
		h.Body = p.parseBlock()
		n.Handler = h
	}
	if p.s.PeekAt(0).Type == lexer.TokenKeywordFinally {
		p.s.ScanExpect(lexer.TokenKeywordFinally, "expected finally statement")
		n.Finalizer = p.parseBlock()
	}
	return n
}

func (p *Parser) parseCatchParameter() ast.BindingPattern {
	b := ast.BindingPattern{}
	t := p.ctx.keywordToIdentifier(p.s.Scan(), false)
	switch t.Type {
	case lexer.TokenIdentifier:
		b.Identifier = t.Literal
	case lexer.TokenPunctuatorOpenBracket:
		b.ArrayPattern = p.parseArrayBindingPatternTail()
	case lexer.TokenPunctuatorOpenBrace:
		b.ObjectPattern = p.parseObjectBindingPatternTail()
	}
	return b
}

func (p *Parser) parseDebuggerStatement() ast.Node {
	panic("unimplemented")
}

func (p *Parser) parseLabelledStatement() ast.Node {
	n := ast.LabeledStatement{}
	p.setStart(&n)
	defer p.setEnd(&n)

	n.Label = p.scanIdent("expected statement label")
	p.s.ScanExpect(lexer.TokenPunctuatorColon, "expected `:` after statement label")
	n.Body = p.parseStatement()
	return n
}
