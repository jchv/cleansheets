package lexer

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func lexAll(s string) (t []Token) {
	l := NewLexer(NewScanner(strings.NewReader(s), nil))
	for {
		token := l.Lex()
		if token.Type == TokenNone {
			return t
		}
		t = append(t, token)
	}
}

func TestLex(t *testing.T) {
	tests := []struct {
		s string
		t []Token
	}{
		{"", nil},
		{
			"(1 + 1) / 2",
			[]Token{
				{Type: TokenPunctuatorOpenParen},
				{Type: TokenLiteralNumber, Literal: "1"},
				{Type: TokenPunctuatorPlus},
				{Type: TokenLiteralNumber, Literal: "1"},
				{Type: TokenPunctuatorCloseParen},
				{Type: TokenPunctuatorDiv},
				{Type: TokenLiteralNumber, Literal: "2"},
			},
		},
		{
			"async function* createAsyncIterable(syncIterable) { for (const elem of syncIterable) { yield elem; } }",
			[]Token{
				{Type: TokenKeywordAsync, Literal: "async"},
				{Type: TokenKeywordFunction, Literal: "function"},
				{Type: TokenPunctuatorMult},
				{Type: TokenIdentifier, Literal: "createAsyncIterable"},
				{Type: TokenPunctuatorOpenParen},
				{Type: TokenIdentifier, Literal: "syncIterable"},
				{Type: TokenPunctuatorCloseParen},
				{Type: TokenPunctuatorOpenBrace},
				{Type: TokenKeywordFor, Literal: "for"},
				{Type: TokenPunctuatorOpenParen},
				{Type: TokenKeywordConst, Literal: "const"},
				{Type: TokenIdentifier, Literal: "elem"},
				{Type: TokenKeywordOf, Literal: "of"},
				{Type: TokenIdentifier, Literal: "syncIterable"},
				{Type: TokenPunctuatorCloseParen},
				{Type: TokenPunctuatorOpenBrace},
				{Type: TokenKeywordYield, Literal: "yield"},
				{Type: TokenIdentifier, Literal: "elem"},
				{Type: TokenPunctuatorSemicolon},
				{Type: TokenPunctuatorCloseBrace},
				{Type: TokenPunctuatorCloseBrace},
			},
		},
		{
			`class Test {
				a() { console.log("Test A"); }
				#b() { console.log("Test B"); }
			}`,
			[]Token{
				{Type: TokenKeywordClass, Literal: "class"},
				{Type: TokenIdentifier, Literal: "Test"},
				{Type: TokenPunctuatorOpenBrace},
				{Type: TokenIdentifier, Literal: "a", NewLine: true},
				{Type: TokenPunctuatorOpenParen},
				{Type: TokenPunctuatorCloseParen},
				{Type: TokenPunctuatorOpenBrace},
				{Type: TokenIdentifier, Literal: "console"},
				{Type: TokenPunctuatorDot},
				{Type: TokenIdentifier, Literal: "log"},
				{Type: TokenPunctuatorOpenParen},
				{Type: TokenLiteralString, Literal: `"Test A"`},
				{Type: TokenPunctuatorCloseParen},
				{Type: TokenPunctuatorSemicolon},
				{Type: TokenPunctuatorCloseBrace},
				{Type: TokenPrivateIdentifier, Literal: "b", NewLine: true},
				{Type: TokenPunctuatorOpenParen},
				{Type: TokenPunctuatorCloseParen},
				{Type: TokenPunctuatorOpenBrace},
				{Type: TokenIdentifier, Literal: "console"},
				{Type: TokenPunctuatorDot},
				{Type: TokenIdentifier, Literal: "log"},
				{Type: TokenPunctuatorOpenParen},
				{Type: TokenLiteralString, Literal: `"Test B"`},
				{Type: TokenPunctuatorCloseParen},
				{Type: TokenPunctuatorSemicolon},
				{Type: TokenPunctuatorCloseBrace},
				{Type: TokenPunctuatorCloseBrace, NewLine: true},
			},
		},
	}

	for _, test := range tests {
		t.Run(strconv.Quote(test.s), func(t *testing.T) {
			result := lexAll(test.s)
			if !reflect.DeepEqual(result, test.t) {
				t.Errorf("lex(%q) = %v != %v", test.s, result, test.t)
			}
		})
	}
}
