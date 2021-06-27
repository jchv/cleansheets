package lexer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/errs"
)

// This lexer does most of the dirty work of turning text into meaningful
// lexical tokens to be parsed. It requires some additional state to be passed
// from the parser to resolve unfortunate ambiguities in ECMA262's grammar.
// Note that keywords, even ones which may be identifiers, will always be lexed
// as keywords; the parser should automatically treat these keywords as
// identifiers as appropriate.

// Lexer lexes ECMAScript code according to ECMA262, 2022 edition section 12.
type Lexer struct {
	s         *Scanner
	lastToken Token
	newLine   bool
}

// Location returns the current source location of the lexer.
func (l *Lexer) Location() ast.Location {
	return l.s.Location()
}

// NewLexer creates a new lexer.
func NewLexer(s *Scanner) *Lexer {
	return &Lexer{s: s}
}

// Lex returns the next token by scanning the input stream.
func (l *Lexer) Lex() Token {
	t := l.consumeNextToken()
	if l.newLine {
		t.NewLine = true
		l.newLine = false
	}
	l.lastToken = t
	return t
}

// ReLex relexes the last token as a regular expression.
func (l *Lexer) ReLex() ReToken {
	t := l.consumeRegex(l.lastToken)
	l.lastToken = t.Token
	return t
}

// consumeRegex lexes a regex, using the passed token as initial state.
func (l *Lexer) consumeRegex(t Token) ReToken {
	lit := &strings.Builder{} // Literal - includes all runes
	pat := &strings.Builder{} // Pattern - includes runes in pattern part
	flg := &strings.Builder{} // Flag - includes runes in flag part

	// Take the passed token and treat it as the start of the pattern.
	lit.WriteString(t.Source())
	pat.WriteString(t.Source()[1:])

patternLoop:
	for {
		r := l.s.Read()
		lit.WriteRune(r)

		switch r {
		case '/':
			// End of pattern
			break patternLoop

		case '[':
			// Consume character class. It is necessary to do this because / is
			// allowed in a character class.
			pat.WriteRune(r)
			for {
				r := l.s.Read()
				lit.WriteRune(r)
				pat.WriteRune(r)

				if r == '\\' {
					r := l.s.Read()
					lit.WriteRune(r)
					pat.WriteRune(r)
				} else if r == ']' {
					break
				} else if r == EOFRune {
					panic(&errs.SyntaxError{
						Location: l.s.Location(),
						Err:      errors.New("unexpected EOF"),
					})
				}
			}

		case '\\':
			// Escape sequence.
			r = l.s.Read()
			lit.WriteRune(r)

			if r == '/' || r == '\\' {
				pat.WriteRune(r)
			} else {
				pat.WriteRune('\\')
				pat.WriteRune(r)
			}

		case EOFRune:
			panic(&errs.SyntaxError{
				Location: l.s.Location(),
				Err:      errors.New("unexpected EOF"),
			})

		default:
			pat.WriteRune(r)
		}
	}

	// Flag loop
	for {
		r := l.s.Read()
		if !isIdentifierContinue(r) {
			l.s.Unread()
			break
		}
		flg.WriteRune(r)
		lit.WriteRune(r)
	}

	return ReToken{
		Token: Token{
			Type:    TokenLiteralRegExp,
			Literal: lit.String(),
		},
		Pattern: pat.String(),
		Flags:   flg.String(),
	}
}

// Consumes a multi-line comment, eating until after the next */.
func (l *Lexer) consumeMultiLineComment() {
	var r rune
	for {
		r = l.s.Read()
		switch r {
		case '*':
			switch l.s.Read() {
			case '/':
				return
			case EOFRune:
				panic(&errs.SyntaxError{
					Location: l.s.Location(),
					Err:      errors.New("unexpected EOF"),
				})
			}
		case EOFRune:
			panic(&errs.SyntaxError{
				Location: l.s.Location(),
				Err:      errors.New("unexpected EOF"),
			})
		}
	}
}

// Consumes a single-line comment, eating until after the next line term.
func (l *Lexer) consumeSingleLineComment() {
	var r rune
	for {
		r = l.s.Read()
		if isLineTerm(r) || r == EOFRune {
			return
		}
	}
}

// Consumes an identifier.
func (l *Lexer) consumeIdentifier(typ TokenType) Token {
	r := l.s.Read()
	if !isIdentifierStart(r) {
		panic(&errs.SyntaxError{
			Location: l.s.Location(),
			Err:      fmt.Errorf("expected IdentifierStart, got %q", r),
		})
	}

	lit := &strings.Builder{}
	lit.WriteRune(r)
	for {
		r := l.s.Read()
		if !isIdentifierContinue(r) {
			l.s.Unread()
			s := lit.String()
			if typ == TokenIdentifier {
				if t, ok := strToKeywordType[s]; ok {
					return Token{Type: t, Literal: s}
				}
			}
			return Token{
				Type:    typ,
				Literal: s,
			}
		}
		lit.WriteRune(r)
	}
}

// Consumes binary digits.
func (l *Lexer) consumeBinaryPart(lit *strings.Builder) string {
	if lit == nil {
		lit = &strings.Builder{}
	}
	r := l.s.Read()

	if isBinaryDigit(r) {
		lit.WriteRune(r)
	} else {
		panic(&errs.SyntaxError{
			Location: l.s.Location(),
			Err:      fmt.Errorf("expected BinaryDigit, got %q", r),
		})
	}

	for {
		r = l.s.Read()
		if isBinaryDigit(r) {
			lit.WriteRune(r)
		} else if isNumericLiteralSeparator(r) {
			r = l.s.Read()
			if isBinaryDigit(r) {
				lit.WriteRune(r)
			} else {
				panic(&errs.SyntaxError{
					Location: l.s.Location(),
					Err:      fmt.Errorf("expected BinaryDigit, got %q", r),
				})
			}
		} else {
			l.s.Unread()
			break
		}
	}

	return lit.String()
}

func (l *Lexer) consumeOctalPart(lit *strings.Builder) string {
	if lit == nil {
		lit = &strings.Builder{}
	}
	r := l.s.Read()

	if isOctalDigit(r) {
		lit.WriteRune(r)
	} else {
		panic(&errs.SyntaxError{
			Location: l.s.Location(),
			Err:      fmt.Errorf("expected OctalDigit, got %q", r),
		})
	}

	for {
		r = l.s.Read()
		if isOctalDigit(r) {
			lit.WriteRune(r)
		} else if isNumericLiteralSeparator(r) {
			r = l.s.Read()
			if isOctalDigit(r) {
				lit.WriteRune(r)
			} else {
				panic(&errs.SyntaxError{
					Location: l.s.Location(),
					Err:      fmt.Errorf("expected OctalDigit, got %q", r),
				})
			}
		} else {
			l.s.Unread()
			break
		}
	}

	return lit.String()
}

func (l *Lexer) consumeHexPart(lit *strings.Builder) string {
	if lit == nil {
		lit = &strings.Builder{}
	}
	r := l.s.Read()

	if isHexDigit(r) {
		lit.WriteRune(r)
	} else {
		panic(&errs.SyntaxError{
			Location: l.s.Location(),
			Err:      fmt.Errorf("expected HexDigit, got %q", r),
		})
	}

	for {
		r = l.s.Read()
		if isHexDigit(r) {
			lit.WriteRune(r)
		} else if isNumericLiteralSeparator(r) {
			r = l.s.Read()
			if isHexDigit(r) {
				lit.WriteRune(r)
			} else {
				panic(&errs.SyntaxError{
					Location: l.s.Location(),
					Err:      fmt.Errorf("expected HexDigit, got %q", r),
				})
			}
		} else {
			l.s.Unread()
			break
		}
	}

	return lit.String()
}

func (l *Lexer) consumeDecimalPart(lit *strings.Builder) string {
	if lit == nil {
		lit = &strings.Builder{}
	}
	r := l.s.Read()

	if !isDecimalDigit(r) {
		panic(&errs.SyntaxError{
			Location: l.s.Location(),
			Err:      fmt.Errorf("expected DecimalDigit, got %q", r),
		})
	}
	lit.WriteRune(r)

	for {
		r = l.s.Read()
		if isDecimalDigit(r) {
			lit.WriteRune(r)
		} else if isNumericLiteralSeparator(r) {
			r = l.s.Read()
			if isDecimalDigit(r) {
				lit.WriteRune(r)
			} else {
				panic(&errs.SyntaxError{
					Location: l.s.Location(),
					Err:      fmt.Errorf("expected DecimalDigit, got %q", r),
				})
			}
		} else if r == '.' {
			lit.WriteRune(r)
			return l.consumeFractionalPart(lit)
		} else if isExponentIndicator(r) {
			for {
				r = l.s.Read()
				if isDecimalDigit(r) {
					lit.WriteRune(r)
				} else {
					l.s.Unread()
					break
				}
			}
			break
		} else {
			l.s.Unread()
			break
		}
	}

	return lit.String()
}

func (l *Lexer) consumeFractionalPart(lit *strings.Builder) string {
	if lit == nil {
		lit = &strings.Builder{}
	}
	r := l.s.Read()

	if isDecimalDigit(r) {
		lit.WriteRune(r)
	} else {
		panic(&errs.SyntaxError{
			Location: l.s.Location(),
			Err:      fmt.Errorf("expected DecimalDigit, got %q", r),
		})
	}

	for {
		r = l.s.Read()
		if isDecimalDigit(r) {
			lit.WriteRune(r)
		} else if isNumericLiteralSeparator(r) {
			r = l.s.Read()
			if isDecimalDigit(r) {
				lit.WriteRune(r)
			} else {
				panic(&errs.SyntaxError{
					Location: l.s.Location(),
					Err:      fmt.Errorf("expected DecimalDigit, got %q", r),
				})
			}
		} else {
			l.s.Unread()
			break
		}
	}

	r = l.s.Read()
	if !isExponentIndicator(r) {
		l.s.Unread()
		return lit.String()
	}
	lit.WriteRune(r)

	r = l.s.Read()
	if r != '+' && r != '-' && !isDecimalDigit(r) {
		panic(&errs.SyntaxError{
			Location: l.s.Location(),
			Err:      fmt.Errorf("expected DecimalDigit, +, or -, got %q", r),
		})
	}
	lit.WriteRune(r)

	for {
		r = l.s.Read()
		if isDecimalDigit(r) {
			lit.WriteRune(r)
		} else if isExponentIndicator(r) {
			for {
				r = l.s.Read()
				if isDecimalDigit(r) {
					lit.WriteRune(r)
				} else {
					l.s.Unread()
					break
				}
			}
			break
		} else {
			l.s.Unread()
			break
		}
	}

	return lit.String()
}

func (l *Lexer) consumeStringLiteral() Token {
	quo := l.s.Read()
	if quo != '\'' && quo != '"' {
		panic("unexpected string literal quote")
	}

	c := []rune{quo}
	for {
		r := l.s.Read()
		c = append(c, r)
		if r == quo {
			break
		}
		if r == '\\' {
			r = l.s.Read()
			c = append(c, r)
		}
		if r == EOFRune {
			panic(&errs.SyntaxError{
				Location: l.s.Location(),
				Err:      errors.New("unexpected EOF"),
			})
		}
	}

	return Token{
		Type:    TokenLiteralString,
		Literal: string(c),
	}
}

func (l *Lexer) consumeNextToken() Token {
	var r rune
	for {
		r = l.s.Read()
		if isLineTerm(r) {
			l.newLine = true
			continue
		}
		if isWhiteSpace(r) {
			continue
		}
		switch r {
		case '{':
			return Token{Type: TokenPunctuatorOpenBrace}
		case '(':
			return Token{Type: TokenPunctuatorOpenParen}
		case '[':
			return Token{Type: TokenPunctuatorOpenBracket}
		case ']':
			return Token{Type: TokenPunctuatorCloseBracket}
		case ')':
			return Token{Type: TokenPunctuatorCloseParen}
		case '}':
			return Token{Type: TokenPunctuatorCloseBrace}
		case '.':
			switch l.s.Read() {
			case '.':
				switch l.s.Read() {
				case '.':
					return Token{Type: TokenPunctuatorEllipsis}
				default:
					panic(&errs.SyntaxError{
						Location: l.s.Location(),
						Err:      fmt.Errorf("expected ., got %q", r),
					})
				}
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				l.s.Unread()
				lit := &strings.Builder{}
				lit.WriteRune(r)
				return Token{Type: TokenLiteralNumber, Literal: l.consumeFractionalPart(lit)}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorDot}
			}
		case '0':
			lit := &strings.Builder{}
			lit.WriteRune(r)
			r = l.s.Read()
			switch r {
			case 'n':
				return Token{Type: TokenLiteralNumber, Literal: "0n"}
			case 'b':
				lit.WriteRune(r)
				return Token{Type: TokenLiteralNumber, Literal: l.consumeBinaryPart(lit)}
			case 'B':
				lit.WriteRune(r)
				return Token{Type: TokenLiteralNumber, Literal: l.consumeBinaryPart(lit)}
			case 'o':
				lit.WriteRune(r)
				return Token{Type: TokenLiteralNumber, Literal: l.consumeOctalPart(lit)}
			case 'O':
				lit.WriteRune(r)
				return Token{Type: TokenLiteralNumber, Literal: l.consumeOctalPart(lit)}
			case 'x':
				lit.WriteRune(r)
				return Token{Type: TokenLiteralNumber, Literal: l.consumeHexPart(lit)}
			case 'X':
				lit.WriteRune(r)
				return Token{Type: TokenLiteralNumber, Literal: l.consumeHexPart(lit)}
			case '_':
				panic(&errs.SyntaxError{
					Location: l.s.Location(),
					Err:      fmt.Errorf("numeric separator can not be used after leading 0"),
				})
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				l.s.Unread()
				return Token{Type: TokenLiteralNumber, Literal: l.consumeDecimalPart(lit)}
			default:
				l.s.Unread()
				return Token{Type: TokenLiteralNumber, Literal: "0"}
			}
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.s.Unread()
			return Token{Type: TokenLiteralNumber, Literal: l.consumeDecimalPart(nil)}
		case ';':
			return Token{Type: TokenPunctuatorSemicolon}
		case ',':
			return Token{Type: TokenPunctuatorComma}
		case '<':
			switch l.s.Read() {
			case '<':
				switch l.s.Read() {
				case '<':
					switch l.s.Read() {
					case '=':
						return Token{Type: TokenPunctuatorUnsignedRShiftAssign}
					default:
						l.s.Unread()
						return Token{Type: TokenPunctuatorUnsignedRShift}
					}
				case '=':
					return Token{Type: TokenPunctuatorLShiftAssign}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorLShift}
				}
			case '=':
				return Token{Type: TokenPunctuatorLessThanEqual}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorLessThan}
			}
		case '>':
			switch l.s.Read() {
			case '>':
				switch l.s.Read() {
				case '>':
					switch l.s.Read() {
					case '=':
						return Token{Type: TokenPunctuatorUnsignedRShiftAssign}
					default:
						l.s.Unread()
						return Token{Type: TokenPunctuatorUnsignedRShift}
					}
				case '=':
					return Token{Type: TokenPunctuatorRShiftAssign}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorRShift}
				}
			case '=':
				return Token{Type: TokenPunctuatorGreaterThanEqual}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorGreaterThan}
			}
		case '=':
			switch l.s.Read() {
			case '=':
				switch l.s.Read() {
				case '=':
					return Token{Type: TokenPunctuatorStrictEqual}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorEqual}
				}
			case '>':
				return Token{Type: TokenPunctuatorFatArrow}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorAssign}
			}
		case '!':
			switch l.s.Read() {
			case '=':
				switch l.s.Read() {
				case '=':
					return Token{Type: TokenPunctuatorStrictNotEqual}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorNotEqual}
				}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorNot}
			}
		case '+':
			switch l.s.Read() {
			case '+':
				return Token{Type: TokenPunctuatorIncrement}
			case '=':
				return Token{Type: TokenPunctuatorPlusAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorPlus}
			}
		case '-':
			switch l.s.Read() {
			case '-':
				return Token{Type: TokenPunctuatorDecrement}
			case '=':
				return Token{Type: TokenPunctuatorMinusAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorMinus}
			}
		case '&':
			switch l.s.Read() {
			case '&':
				switch l.s.Read() {
				case '=':
					return Token{Type: TokenPunctuatorLogicalAndAssign}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorLogicalAnd}
				}
			case '=':
				return Token{Type: TokenPunctuatorBitAndAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorBitAnd}
			}
		case '|':
			switch l.s.Read() {
			case '|':
				switch l.s.Read() {
				case '=':
					return Token{Type: TokenPunctuatorLogicalOrAssign}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorLogicalOr}
				}
			case '=':
				return Token{Type: TokenPunctuatorBitOrAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorBitOr}
			}
		case '^':
			switch l.s.Read() {
			case '=':
				return Token{Type: TokenPunctuatorBitXorAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorBitXor}
			}
		case '~':
			return Token{Type: TokenPunctuatorBitNot}
		case '?':
			switch l.s.Read() {
			case '?':
				switch l.s.Read() {
				case '=':
					return Token{Type: TokenPunctuatorNullCoalesceAssign}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorNullCoalesce}
				}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorQuestionMark}
			}
		case ':':
			return Token{Type: TokenPunctuatorColon}
		case '*':
			switch l.s.Read() {
			case '*':
				switch l.s.Read() {
				case '=':
					return Token{Type: TokenPunctuatorExponentAssign}
				default:
					l.s.Unread()
					return Token{Type: TokenPunctuatorExponent}
				}
			case '=':
				return Token{Type: TokenPunctuatorMultAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorMult}
			}
		case '%':
			switch l.s.Read() {
			case '=':
				return Token{Type: TokenPunctuatorModAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorMod}
			}
		case '/':
			switch l.s.Read() {
			case '/':
				l.consumeSingleLineComment()
				continue
			case '*':
				l.consumeMultiLineComment()
				continue
			case '=':
				return Token{Type: TokenPunctuatorDivAssign}
			default:
				l.s.Unread()
				return Token{Type: TokenPunctuatorDiv}
			}
		case '"', '\'':
			l.s.Unread()
			return l.consumeStringLiteral()
		case '#':
			return l.consumeIdentifier(TokenPrivateIdentifier)
		case EOFRune:
			return Token{Type: TokenNone}
		default:
			if isIdentifierStart(r) {
				l.s.Unread()
				return l.consumeIdentifier(TokenIdentifier)
			}

			panic(&errs.SyntaxError{
				Location: l.s.Location(),
				Err:      fmt.Errorf("unexpected rune %q", r),
			})
		}
	}
}
