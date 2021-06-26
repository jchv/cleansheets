package parser

import "github.com/jchv/cleansheets/ecmascript/lexer"

// reservedType specifies what contexts a keyword could also be a valid
// identifier in.
type reservedType int

const (
	// reservedNever specifies that a keyword is not reserved; can always be an
	// identifier.
	reservedNever reservedType = iota

	// reservedAsync specifies that a keyword is reserved in async contexts;
	// can be an identifier outside of async contexts.
	reservedAsync

	// reservedGenerator specifies that a keyword is reserved in generators;
	// can be an identifier outside of generators.
	reservedGenerator

	// reservedStrict specifies that a keyword is reserved in strict contexts;
	// can be an identifier outside of strict contexts.
	reservedStrict

	// reservedAlways specifies that a keyword is reserved in all contexts; can
	// never be an identifier.
	reservedAlways
)

// reservedWords specifies the reservation state for keyword tokens.
var reservedWords = map[lexer.TokenType]reservedType{
	lexer.TokenKeywordAs:     reservedNever,
	lexer.TokenKeywordAsync:  reservedNever,
	lexer.TokenKeywordFrom:   reservedNever,
	lexer.TokenKeywordGet:    reservedNever,
	lexer.TokenKeywordMeta:   reservedNever,
	lexer.TokenKeywordOf:     reservedNever,
	lexer.TokenKeywordSet:    reservedNever,
	lexer.TokenKeywordTarget: reservedNever,

	lexer.TokenKeywordAwait: reservedAsync,
	lexer.TokenKeywordYield: reservedGenerator,

	lexer.TokenKeywordImplements: reservedStrict,
	lexer.TokenKeywordInterface:  reservedStrict,
	lexer.TokenKeywordLet:        reservedStrict,
	lexer.TokenKeywordPackage:    reservedStrict,
	lexer.TokenKeywordPrivate:    reservedStrict,
	lexer.TokenKeywordProtected:  reservedStrict,
	lexer.TokenKeywordPublic:     reservedStrict,
	lexer.TokenKeywordStatic:     reservedStrict,

	lexer.TokenKeywordBreak:      reservedAlways,
	lexer.TokenKeywordCase:       reservedAlways,
	lexer.TokenKeywordCatch:      reservedAlways,
	lexer.TokenKeywordClass:      reservedAlways,
	lexer.TokenKeywordConst:      reservedAlways,
	lexer.TokenKeywordContinue:   reservedAlways,
	lexer.TokenKeywordDebugger:   reservedAlways,
	lexer.TokenKeywordDefault:    reservedAlways,
	lexer.TokenKeywordDelete:     reservedAlways,
	lexer.TokenKeywordDo:         reservedAlways,
	lexer.TokenKeywordElse:       reservedAlways,
	lexer.TokenKeywordEnum:       reservedAlways,
	lexer.TokenKeywordExport:     reservedAlways,
	lexer.TokenKeywordExtends:    reservedAlways,
	lexer.TokenKeywordFalse:      reservedAlways,
	lexer.TokenKeywordFinally:    reservedAlways,
	lexer.TokenKeywordFor:        reservedAlways,
	lexer.TokenKeywordFunction:   reservedAlways,
	lexer.TokenKeywordIf:         reservedAlways,
	lexer.TokenKeywordImport:     reservedAlways,
	lexer.TokenKeywordIn:         reservedAlways,
	lexer.TokenKeywordInstanceOf: reservedAlways,
	lexer.TokenKeywordNew:        reservedAlways,
	lexer.TokenKeywordNull:       reservedAlways,
	lexer.TokenKeywordReturn:     reservedAlways,
	lexer.TokenKeywordSuper:      reservedAlways,
	lexer.TokenKeywordSwitch:     reservedAlways,
	lexer.TokenKeywordThis:       reservedAlways,
	lexer.TokenKeywordThrow:      reservedAlways,
	lexer.TokenKeywordTrue:       reservedAlways,
	lexer.TokenKeywordTry:        reservedAlways,
	lexer.TokenKeywordTypeOf:     reservedAlways,
	lexer.TokenKeywordVar:        reservedAlways,
	lexer.TokenKeywordVoid:       reservedAlways,
	lexer.TokenKeywordWhile:      reservedAlways,
	lexer.TokenKeywordWith:       reservedAlways,
}
