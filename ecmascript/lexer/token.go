package lexer

import (
	"fmt"
	"strconv"
)

// TokenType is an enumeration of possible token types.
type TokenType int

//go:generate go run golang.org/x/tools/cmd/stringer -type=TokenType

// These are all of the possible token types.
const (
	TokenNone TokenType = iota

	TokenIdentifier
	TokenPrivateIdentifier

	// Keywords
	TokenKeywordAs
	TokenKeywordAsync
	TokenKeywordAwait
	TokenKeywordBreak
	TokenKeywordCase
	TokenKeywordCatch
	TokenKeywordClass
	TokenKeywordConst
	TokenKeywordContinue
	TokenKeywordDebugger
	TokenKeywordDefault
	TokenKeywordDelete
	TokenKeywordDo
	TokenKeywordElse
	TokenKeywordEnum
	TokenKeywordExport
	TokenKeywordExtends
	TokenKeywordFalse
	TokenKeywordFinally
	TokenKeywordFor
	TokenKeywordFrom
	TokenKeywordFunction
	TokenKeywordGet
	TokenKeywordIf
	TokenKeywordImplements
	TokenKeywordImport
	TokenKeywordIn
	TokenKeywordInstanceOf
	TokenKeywordInterface
	TokenKeywordLet
	TokenKeywordNew
	TokenKeywordNull
	TokenKeywordMeta
	TokenKeywordOf
	TokenKeywordPackage
	TokenKeywordPrivate
	TokenKeywordProtected
	TokenKeywordPublic
	TokenKeywordReturn
	TokenKeywordSet
	TokenKeywordStatic
	TokenKeywordSuper
	TokenKeywordSwitch
	TokenKeywordTarget
	TokenKeywordThis
	TokenKeywordThrow
	TokenKeywordTrue
	TokenKeywordTry
	TokenKeywordTypeOf
	TokenKeywordVar
	TokenKeywordVoid
	TokenKeywordWhile
	TokenKeywordWith
	TokenKeywordYield

	// Punctuators
	TokenPunctuatorOptionalChain
	TokenPunctuatorOpenBrace
	TokenPunctuatorOpenParen
	TokenPunctuatorOpenBracket
	TokenPunctuatorCloseBracket
	TokenPunctuatorCloseParen
	TokenPunctuatorCloseBrace
	TokenPunctuatorDot
	TokenPunctuatorEllipsis
	TokenPunctuatorSemicolon
	TokenPunctuatorComma
	TokenPunctuatorLessThan
	TokenPunctuatorGreaterThan
	TokenPunctuatorLessThanEqual
	TokenPunctuatorGreaterThanEqual
	TokenPunctuatorEqual
	TokenPunctuatorNotEqual
	TokenPunctuatorStrictEqual
	TokenPunctuatorStrictNotEqual
	TokenPunctuatorPlus
	TokenPunctuatorMinus
	TokenPunctuatorMult
	TokenPunctuatorDiv
	TokenPunctuatorMod
	TokenPunctuatorExponent
	TokenPunctuatorIncrement
	TokenPunctuatorDecrement
	TokenPunctuatorLShift
	TokenPunctuatorRShift
	TokenPunctuatorUnsignedRShift
	TokenPunctuatorBitAnd
	TokenPunctuatorBitOr
	TokenPunctuatorBitXor
	TokenPunctuatorNot
	TokenPunctuatorBitNot
	TokenPunctuatorLogicalAnd
	TokenPunctuatorLogicalOr
	TokenPunctuatorNullCoalesce
	TokenPunctuatorQuestionMark
	TokenPunctuatorColon
	TokenPunctuatorAssign
	TokenPunctuatorPlusAssign
	TokenPunctuatorMinusAssign
	TokenPunctuatorMultAssign
	TokenPunctuatorDivAssign
	TokenPunctuatorModAssign
	TokenPunctuatorExponentAssign
	TokenPunctuatorLShiftAssign
	TokenPunctuatorRShiftAssign
	TokenPunctuatorUnsignedRShiftAssign
	TokenPunctuatorBitAndAssign
	TokenPunctuatorBitOrAssign
	TokenPunctuatorBitXorAssign
	TokenPunctuatorLogicalAndAssign
	TokenPunctuatorLogicalOrAssign
	TokenPunctuatorNullCoalesceAssign
	TokenPunctuatorFatArrow

	// Literals
	TokenLiteralNumber
	TokenLiteralString
	TokenLiteralRegExp
	TokenLiteralTemplate
)

var strToKeywordType = map[string]TokenType{
	"as":         TokenKeywordAs,
	"async":      TokenKeywordAsync,
	"await":      TokenKeywordAwait,
	"break":      TokenKeywordBreak,
	"case":       TokenKeywordCase,
	"catch":      TokenKeywordCatch,
	"class":      TokenKeywordClass,
	"const":      TokenKeywordConst,
	"continue":   TokenKeywordContinue,
	"debugger":   TokenKeywordDebugger,
	"default":    TokenKeywordDefault,
	"delete":     TokenKeywordDelete,
	"do":         TokenKeywordDo,
	"else":       TokenKeywordElse,
	"enum":       TokenKeywordEnum,
	"export":     TokenKeywordExport,
	"extends":    TokenKeywordExtends,
	"false":      TokenKeywordFalse,
	"finally":    TokenKeywordFinally,
	"for":        TokenKeywordFor,
	"from":       TokenKeywordFrom,
	"function":   TokenKeywordFunction,
	"get":        TokenKeywordGet,
	"if":         TokenKeywordIf,
	"implements": TokenKeywordImplements,
	"import":     TokenKeywordImport,
	"in":         TokenKeywordIn,
	"instanceof": TokenKeywordInstanceOf,
	"interface":  TokenKeywordInterface,
	"let":        TokenKeywordLet,
	"meta":       TokenKeywordMeta,
	"new":        TokenKeywordNew,
	"null":       TokenKeywordNull,
	"of":         TokenKeywordOf,
	"package":    TokenKeywordPackage,
	"private":    TokenKeywordPrivate,
	"protected":  TokenKeywordProtected,
	"public":     TokenKeywordPublic,
	"return":     TokenKeywordReturn,
	"set":        TokenKeywordSet,
	"static":     TokenKeywordStatic,
	"super":      TokenKeywordSuper,
	"switch":     TokenKeywordSwitch,
	"target":     TokenKeywordTarget,
	"this":       TokenKeywordThis,
	"throw":      TokenKeywordThrow,
	"true":       TokenKeywordTrue,
	"try":        TokenKeywordTry,
	"typeof":     TokenKeywordTypeOf,
	"var":        TokenKeywordVar,
	"void":       TokenKeywordVoid,
	"while":      TokenKeywordWhile,
	"with":       TokenKeywordWith,
	"yield":      TokenKeywordYield,
}

// Token represents an ECMAScript lexical token.
type Token struct {
	Type    TokenType
	Literal string
	NewLine bool
}

// ReToken represents an ECMAScript regular expression token.
type ReToken struct {
	Token
	Pattern string
	Flags   string
}

// String implements the Stringer interface.
func (t Token) String() string {
	if t.Literal == "" {
		if t.NewLine {
			return fmt.Sprintf("{ Type: %s, NewLine: true }", t.Type)
		} else {
			return fmt.Sprintf("{ Type: %s }", t.Type)
		}
	}
	if t.NewLine {
		return fmt.Sprintf("{ Type: %s, Literal: %q, NewLine: true }", t.Type, t.Literal)
	} else {
		return fmt.Sprintf("{ Type: %s, Literal: %q }", t.Type, t.Literal)
	}
}

// Source returns the corresponding source code for a token.
func (t Token) Source() string {
	switch t.Type {
	case TokenIdentifier:
		return t.Literal
	case TokenPrivateIdentifier:
		return "#" + t.Literal
	case
		// Keywords
		TokenKeywordAs, TokenKeywordAsync, TokenKeywordAwait,
		TokenKeywordBreak, TokenKeywordCase, TokenKeywordCatch,
		TokenKeywordClass, TokenKeywordConst, TokenKeywordContinue,
		TokenKeywordDebugger, TokenKeywordDefault, TokenKeywordDelete,
		TokenKeywordDo, TokenKeywordElse, TokenKeywordEnum,
		TokenKeywordExport, TokenKeywordExtends, TokenKeywordFalse,
		TokenKeywordFinally, TokenKeywordFor, TokenKeywordFrom,
		TokenKeywordFunction, TokenKeywordGet, TokenKeywordIf,
		TokenKeywordImplements, TokenKeywordImport, TokenKeywordIn,
		TokenKeywordInstanceOf, TokenKeywordInterface, TokenKeywordLet,
		TokenKeywordNew, TokenKeywordNull, TokenKeywordMeta,
		TokenKeywordOf, TokenKeywordPackage, TokenKeywordPrivate,
		TokenKeywordProtected, TokenKeywordPublic, TokenKeywordReturn,
		TokenKeywordSet, TokenKeywordStatic, TokenKeywordSuper,
		TokenKeywordSwitch, TokenKeywordTarget, TokenKeywordThis,
		TokenKeywordThrow, TokenKeywordTrue, TokenKeywordTry,
		TokenKeywordTypeOf, TokenKeywordVar, TokenKeywordVoid,
		TokenKeywordWhile, TokenKeywordWith, TokenKeywordYield,
		// Literals
		TokenLiteralNumber, TokenLiteralString, TokenLiteralRegExp,
		TokenLiteralTemplate:
		return t.Literal
	case TokenPunctuatorOptionalChain:
		return ".?"
	case TokenPunctuatorOpenBrace:
		return "{"
	case TokenPunctuatorOpenParen:
		return "("
	case TokenPunctuatorOpenBracket:
		return "["
	case TokenPunctuatorCloseBracket:
		return "]"
	case TokenPunctuatorCloseParen:
		return ")"
	case TokenPunctuatorCloseBrace:
		return "}"
	case TokenPunctuatorDot:
		return "."
	case TokenPunctuatorEllipsis:
		return "..."
	case TokenPunctuatorSemicolon:
		return ";"
	case TokenPunctuatorComma:
		return ","
	case TokenPunctuatorLessThan:
		return "<"
	case TokenPunctuatorGreaterThan:
		return ">"
	case TokenPunctuatorLessThanEqual:
		return "<="
	case TokenPunctuatorGreaterThanEqual:
		return ">="
	case TokenPunctuatorEqual:
		return "=="
	case TokenPunctuatorNotEqual:
		return "==="
	case TokenPunctuatorStrictEqual:
		return "==="
	case TokenPunctuatorStrictNotEqual:
		return "!=="
	case TokenPunctuatorPlus:
		return "+"
	case TokenPunctuatorMinus:
		return "-"
	case TokenPunctuatorMult:
		return "*"
	case TokenPunctuatorDiv:
		return "/"
	case TokenPunctuatorMod:
		return "%"
	case TokenPunctuatorExponent:
		return "**"
	case TokenPunctuatorIncrement:
		return "++"
	case TokenPunctuatorDecrement:
		return "--"
	case TokenPunctuatorLShift:
		return "<<"
	case TokenPunctuatorRShift:
		return ">>"
	case TokenPunctuatorUnsignedRShift:
		return ">>>"
	case TokenPunctuatorBitAnd:
		return "&"
	case TokenPunctuatorBitOr:
		return "|"
	case TokenPunctuatorBitXor:
		return "^"
	case TokenPunctuatorNot:
		return "!"
	case TokenPunctuatorBitNot:
		return "~"
	case TokenPunctuatorLogicalAnd:
		return "&&"
	case TokenPunctuatorLogicalOr:
		return "||"
	case TokenPunctuatorNullCoalesce:
		return "??"
	case TokenPunctuatorQuestionMark:
		return "?"
	case TokenPunctuatorColon:
		return ":"
	case TokenPunctuatorAssign:
		return "="
	case TokenPunctuatorPlusAssign:
		return "+="
	case TokenPunctuatorMinusAssign:
		return "-="
	case TokenPunctuatorMultAssign:
		return "*="
	case TokenPunctuatorDivAssign:
		return "/="
	case TokenPunctuatorModAssign:
		return "*="
	case TokenPunctuatorExponentAssign:
		return "**="
	case TokenPunctuatorLShiftAssign:
		return "<<="
	case TokenPunctuatorRShiftAssign:
		return ">>="
	case TokenPunctuatorUnsignedRShiftAssign:
		return ">>>="
	case TokenPunctuatorBitAndAssign:
		return "&="
	case TokenPunctuatorBitOrAssign:
		return "|="
	case TokenPunctuatorBitXorAssign:
		return "^="
	case TokenPunctuatorLogicalAndAssign:
		return "&&="
	case TokenPunctuatorLogicalOrAssign:
		return "||="
	case TokenPunctuatorNullCoalesceAssign:
		return "??="
	case TokenPunctuatorFatArrow:
		return "=>"
	}
	return t.Type.String()
}

// StringConstant returns the parsed value for a string constant.
func (t Token) StringConstant() string {
	if t.Type != TokenLiteralString {
		panic("expected string literal token")
	}

	// TODO: actual string parsing :)
	return t.Literal[1 : len(t.Literal)-1]
}

// NumberConstant returns the parsed value for a numeric constant.
func (t Token) NumberConstant() float64 {
	// TODO: lexer should be parsing numbers accurately
	v, err := strconv.ParseFloat(t.Literal, 64)
	if err == nil {
		return v
	} else {
		v, err := strconv.ParseInt(t.Literal, 0, 64)
		if err == nil {
			return float64(v)
		} else {
			panic(err)
		}
	}
}
