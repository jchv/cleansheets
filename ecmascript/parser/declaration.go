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
