package lexer

import "unicode"

var whitespace = map[rune]struct{}{
	'\u0009': {}, '\u000b': {}, '\u000c': {},
	'\u0020': {}, '\u00a0': {}, '\u1680': {},
	'\u2000': {}, '\u2001': {}, '\u2002': {},
	'\u2003': {}, '\u2004': {}, '\u2005': {},
	'\u2006': {}, '\u2007': {}, '\u2008': {},
	'\u2009': {}, '\u200a': {}, '\u202f': {},
	'\u205f': {}, '\u3000': {}, '\ufeff': {},
}

var lineterms = map[rune]struct{}{
	'\u000a': {}, '\u000d': {},
	'\u2028': {}, '\u2029': {},
}

func isWhiteSpace(r rune) bool {
	_, ok := whitespace[r]
	return ok
}

func isLineTerm(r rune) bool {
	_, ok := lineterms[r]
	return ok
}

func isIdentifierStart(r rune) bool {
	return (r == '$' || r == '_' ||
		(unicode.In(r, unicode.L, unicode.Nl, unicode.Other_ID_Start) &&
			!unicode.In(r, unicode.Pattern_Syntax, unicode.Pattern_White_Space)))
}

func isIdentifierContinue(r rune) bool {
	return (r == '$' || r == '_' || r == 0x200C || r == 0x200D ||
		(unicode.In(r, unicode.L, unicode.Nl, unicode.Other_ID_Start, unicode.Mn,
			unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue) &&
			!unicode.In(r, unicode.Pattern_Syntax, unicode.Pattern_White_Space)))
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}

func isExponentIndicator(r rune) bool {
	return r == 'e' || r == 'E'
}

func isNumericLiteralSeparator(r rune) bool {
	return r == '_'
}

// EncodeUTF16 encodes a UTF-8 string as a UTF-16 string.
func EncodeUTF16(s string) []uint16 {
	n := len(s)
	for _, v := range s {
		if v >= 0x10000 {
			n++
		}
	}

	a := make([]uint16, n)
	n = 0

	for _, v := range s {
		switch {
		case 0 <= v && v < 0xd800, 0xe000 <= v && v < 0x10000:
			a[n] = uint16(v)
			n++

		case 0x10000 <= v && v <= 0x10ffff:
			v -= 0x10000
			a[n] = uint16(0xd800 + (v>>10)&0x3ff)
			a[n+1] = uint16(0xdc00 + v&0x3ff)
			n += 2

		default:
			a[n] = uint16(0xfffd)
			n++
		}
	}

	return a[:n]
}
