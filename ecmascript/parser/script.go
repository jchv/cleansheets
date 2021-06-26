package parser

import (
	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

func (p *Parser) parseScript() ast.Node {
	m := ast.ScriptNode{}
	p.setStart(&m)
	defer p.setEnd(&m)

	for {
		if p.s.PeekAt(0).Type == lexer.TokenNone {
			break
		}
		m.Body = append(m.Body, p.parseStatementItem())
	}

	return m
}
