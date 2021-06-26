package parser

import (
	"fmt"

	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

func (p *Parser) parseModule() ast.Node {
	// Modules are always strict.
	p.ctx.strictMode = true

	m := ast.ModuleNode{}
	p.setStart(&m)
	defer p.setEnd(&m)

	for {
		if p.s.PeekAt(0).Type == lexer.TokenNone {
			break
		}
		m.Body = append(m.Body, p.parseModuleItem())
	}

	return m
}

func (p *Parser) parseModuleItem() ast.Node {
	switch p.s.PeekAt(0).Type {
	case lexer.TokenNone:
		return nil
	case lexer.TokenKeywordImport:
		return p.parseImportDecl()
	case lexer.TokenKeywordExport:
		return p.parseExportDecl()
	default:
		return p.parseStatementItem()
	}
}

func (p *Parser) parseImportDecl() ast.ImportDeclNode {
	n := ast.ImportDeclNode{}
	p.setStart(&n)
	defer p.setEnd(&n)

	p.s.ScanExpect(lexer.TokenKeywordImport, "expected `import` declaration")

	t := p.ctx.keywordToIdentifier(p.s.Scan(), false)
	switch t.Type {
	case lexer.TokenLiteralString:
		n.Module = t.StringConstant()
		p.expectSemicolon()
		return n

	case lexer.TokenIdentifier:
		n.DefaultBinding = &ast.ImportDefaultBinding{
			Identifier: t.Literal,
		}

		t = p.s.Scan()
		switch t.Type {
		case lexer.TokenPunctuatorComma:
			t = p.s.Scan()

		case lexer.TokenKeywordFrom:
			t = p.s.ScanExpect(lexer.TokenLiteralString, "expected module specifier after `from`")
			n.Module = t.StringConstant()
			p.expectSemicolon()
			return n

		default:
			p.s.SyntaxError(fmt.Sprintf("expected `,` or `from` after default import in import declaration, got %q", t.Source()))
		}
	}

	switch t.Type {
	case lexer.TokenPunctuatorMult:
		p.s.ScanExpect(lexer.TokenKeywordAs, "expected `as` after namespace binding operator `*`")
		n.NameSpace = &ast.NameSpaceImport{Identifier: p.scanIdent("expected namespace binding after `* as`")}

	case lexer.TokenPunctuatorOpenBrace:
		n.NamedImports = []ast.NamedImport{}

	importList:
		for {
			t = p.s.Scan()
			if t.Type == lexer.TokenPunctuatorCloseBrace {
				break importList
			}
			item := ast.NamedImport{
				Identifier: p.expectIdent(t, "expected import specifier in import list"),
			}
			t = p.s.Scan()
			switch t.Type {
			case lexer.TokenPunctuatorCloseBrace:
				n.NamedImports = append(n.NamedImports, item)
				break importList
			case lexer.TokenPunctuatorComma:
				n.NamedImports = append(n.NamedImports, item)
			case lexer.TokenKeywordAs:
				item.AsBinding = p.scanIdent("expected import binding after `as` in import list")
				t = p.s.Scan()
				switch t.Type {
				case lexer.TokenPunctuatorCloseBrace:
					n.NamedImports = append(n.NamedImports, item)
					break importList
				case lexer.TokenPunctuatorComma:
					n.NamedImports = append(n.NamedImports, item)
				}
			}
		}

	default:
		p.s.SyntaxError("expected namespace or named imports in import statement")
	}

	p.s.ScanExpect(lexer.TokenKeywordFrom, "expected `from` clause in import declaration")
	n.Module = p.s.ScanExpect(lexer.TokenLiteralString, "expected module specifier after `from`").StringConstant()

	p.expectSemicolon()

	return n
}

func (p *Parser) parseExportDecl() ast.Node {
	panic("unimplemented")
}
