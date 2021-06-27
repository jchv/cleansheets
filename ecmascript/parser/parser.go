package parser

import (
	"fmt"

	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/errs"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

// ParseMode specifies what mode to use when parsing the ECMAScript code.
type ParseMode int

const (
	// ScriptMode parses the ECMAScript code as a script.
	ScriptMode ParseMode = iota

	// ModuleMode parses the ECMAScript code as a module.
	ModuleMode

	// ExpressionMode parses the ECMAScript code as an expression.
	ExpressionMode
)

// ParseOptions are options that adjust how ECMAScript code should be parsed.
type ParseOptions struct {
	Mode ParseMode
}

// Parser parses ECMAScript code according to ECMA262.
type Parser struct {
	s   *Scanner
	ctx parseContext
}

// NewParser creates a new parser.
func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{s: NewScanner(l)}
}

// Parse parses ECMAScript code.
func (p *Parser) Parse(opt ParseOptions) (n ast.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case *errs.SyntaxError:
				err = t
			case *errs.EncodingError:
				err = t
			case *errs.ParserError:
				err = t
			default:
				panic(err)
			}
		}
	}()
	switch opt.Mode {
	case ScriptMode:
		return p.parseScript(), nil
	case ModuleMode:
		return p.parseModule(), nil
	case ExpressionMode:
		return p.parseExpression(exprOrderComma, 0), nil
	default:
		panic(fmt.Errorf("unexpected parse mode %d", opt.Mode))
	}
}

// scanIdent expects an identifier.
func (p *Parser) scanIdent(err string) string {
	return p.expectIdent(p.s.Scan(), err)
}

// forceScanIdent expects an identifier even when a reserved keyword is found.
func (p *Parser) forceScanIdent(err string) string {
	return p.forceIdent(p.s.Scan(), err)
}

// expectIdent expects an identifier.
func (p *Parser) expectIdent(t lexer.Token, err string) string {
	t = p.ctx.keywordToIdentifier(t, false)
	if t.Type != lexer.TokenIdentifier {
		p.s.SyntaxError(fmt.Sprintf("expected identifier, got %s: %s", t.Source(), err))
	}
	return t.Literal
}

// forceIdent forces conversion to identifier; for contexts where keywords can't appear.
func (p *Parser) forceIdent(t lexer.Token, err string) string {
	t = p.ctx.keywordToIdentifier(t, true)
	if t.Type != lexer.TokenIdentifier {
		p.s.SyntaxError(fmt.Sprintf("expected identifier, got %s: %s", t.Source(), err))
	}
	return t.Literal
}

// expectSemicolon expects either a semicolon, or an eligible newline for
// semicolon insertion.
func (p *Parser) expectSemicolon() {
	t := p.s.PeekAt(0)

	if t.Type != lexer.TokenPunctuatorSemicolon {
		// Part of the automatic semi-colon insertion algorithm.
		if t.NewLine || t.Type == lexer.TokenPunctuatorCloseBrace || t.Type == lexer.TokenNone {
			return
		}
	}

	p.s.ScanExpect(lexer.TokenPunctuatorSemicolon, "did you forget a semicolon?")
}

type spannedNode interface {
	ast.Node
	SetStart(ast.Location)
	SetEnd(ast.Location)
}

// setStart sets the start of a node to the current location.
func (p *Parser) setStart(s spannedNode) {
	s.SetStart(p.s.Location())
}

// setEnd sets the end of a node; ideal for use with defer.
func (p *Parser) setEnd(s spannedNode) {
	s.SetEnd(p.s.Location())
}
