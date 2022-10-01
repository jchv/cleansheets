package parser

import (
	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

func (p *Parser) parseDeclaration() ast.Node {
	switch p.s.PeekAt(0).Type {
	case lexer.TokenKeywordFunction:
		return p.parseFunctionDeclaration()
	case lexer.TokenKeywordLet, lexer.TokenKeywordConst:
		return p.parseLexicalDeclaration()
	case lexer.TokenKeywordClass:
		return p.parseClassDeclaration()
	}
	return nil
}

func (p *Parser) parseFunctionDeclaration() ast.Node {
	s := p.s.Location()
	p.s.ScanExpect(lexer.TokenKeywordFunction, "expected function")
	// TODO: support eliding name when in `export default` context.
	name := p.scanIdent("expected identifier")
	// TODO: generator support
	p.s.ScanExpect(lexer.TokenPunctuatorOpenParen, "expected parameter list following function declaration")
	params := p.parseParametersTail()
	body := p.parseBlock()
	n := ast.FunctionDeclaration{
		ID:     name,
		Params: params,
		Body:   body,
	}
	n.SetStart(s)
	n.SetEnd(p.s.Location())
	return n
}

func (p *Parser) parseLexicalDeclaration() ast.VariableDeclaration {
	n := p.parseLexicalDeclarationNoSemicolon()
	p.expectSemicolon()
	p.setEnd(&n)
	return n
}

func (p *Parser) parseLexicalDeclarationNoSemicolon() ast.VariableDeclaration {
	n := ast.VariableDeclaration{}
	p.setStart(&n)
	defer p.setEnd(&n)

	switch p.s.Scan().Type {
	case lexer.TokenKeywordLet:
		n.Declarations = p.parseVariableDeclarations()
		n.Kind = ast.LetDeclaration
	case lexer.TokenKeywordConst:
		n.Declarations = p.parseVariableDeclarations()
		n.Kind = ast.ConstDeclaration
	default:
		p.s.SyntaxError("expected lexical declaration")
	}
	return n
}

func (p *Parser) parseClassDeclaration() ast.Node {
	n := ast.ClassDeclaration{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordClass, "expected class")
	n.ID = p.scanIdent("expected class name")

	if p.s.PeekAt(0).Type == lexer.TokenKeywordExtends {
		p.s.Scan()
		n.SuperClass = p.parseExpression(exprOrderMemberExpr, 0)
	}

	n.Body = p.parseClassBody()
	return n
}

func (p *Parser) parseClassBody() []ast.Node {
	p.s.ScanExpect(lexer.TokenPunctuatorOpenBrace, "expected '{'")

	n := []ast.Node{}

	for {
		peek := p.s.PeekAt(0)
		if peek.Type == lexer.TokenPunctuatorCloseBrace {
			p.s.Scan()
			break
		}

		// TODO: implement member variables...
		m := ast.MethodDefinition{}

		// Static specifier
		if peek.Type == lexer.TokenKeywordStatic {
			p.s.Scan()
			peek = p.s.PeekAt(0)
			m.Static = true
		}

		// Get/set specifier
		switch peek.Type {
		case lexer.TokenKeywordGet:
			p.s.Scan()
			m.Kind = ast.GetMethod

		case lexer.TokenKeywordSet:
			p.s.Scan()
			m.Kind = ast.SetMethod
		}

		// Identifier (possibly computed)
		t := p.s.Scan()
		switch t.Type {
		case lexer.TokenIdentifier:
			m.Key = ast.Identifier{Name: t.Literal}

		case lexer.TokenPunctuatorOpenBracket:
			m.Computed = true
			m.Key = p.parseExpression(exprOrderComma, 0)
			p.s.ScanExpect(lexer.TokenPunctuatorCloseBracket, "expected `]`")

		default:
			p.s.SyntaxError("expected method definition")
		}

		fn := ast.FunctionExpression{}
		fn.Params = p.parseParameters()
		fn.Body = p.parseBlock()
		fn.SetEnd(p.s.Location())
		m.Value = fn

		n = append(n, m)
	}

	return n
}
