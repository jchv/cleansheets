package parser

import "github.com/jchv/cleansheets/ecmascript/lexer"

// parseContext provides context-specific parsing state and utilities.
type parseContext struct {
	strictMode bool
	async      bool
	generator  bool
}

// keywordToIdentifier converts a keyword to an identifier, if permissible in
// the context.
func (ctx *parseContext) keywordToIdentifier(token lexer.Token, force bool) lexer.Token {
	reservation, ok := reservedWords[token.Type]
	if !ok {
		return token
	}

	if !force {
		switch reservation {
		case reservedAlways:
			return token
		case reservedAsync:
			if ctx.async {
				return token
			}
		case reservedGenerator:
			if ctx.generator {
				return token
			}
		case reservedStrict:
			if ctx.strictMode {
				return token
			}
		default:
			break
		}
	}

	return lexer.Token{Type: lexer.TokenIdentifier, Literal: token.Literal}
}
