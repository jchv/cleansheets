package parser

import (
	"testing"

	"github.com/jchv/cleansheets/ecmascript/ast"
)

func TestObjectLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.Property
	}{
		{
			"init, non-computed name",
			"{ property: null }",
			ast.Property{
				Kind:  ast.InitProperty,
				Key:   ident("property"),
				Value: ast.NullLiteral{},
			},
		},
		{
			"init, string key",
			"{ \"property\": null }",
			ast.Property{
				Kind:  ast.InitProperty,
				Key:   ast.StringLiteral{Value: "property", Raw: "\"property\""},
				Value: ast.NullLiteral{},
			},
		},
		{
			"init, number key",
			"{ 0: null }",
			ast.Property{
				Kind:  ast.InitProperty,
				Key:   ast.NumberLiteral{Value: 0, Raw: "0"},
				Value: ast.NullLiteral{},
			},
		},
		{
			"getter, non-computed name",
			"{ get property() {} }",
			ast.Property{
				Kind:  ast.GetProperty,
				Key:   ident("property"),
				Value: ast.FunctionExpression{Body: ast.BlockStatement{}},
			},
		},
		{
			"setter, non-computed name",
			"{ set property() {} }",
			ast.Property{
				Kind:  ast.SetProperty,
				Key:   ident("property"),
				Value: ast.FunctionExpression{Body: ast.BlockStatement{}},
			},
		},
		{
			"method, non-computed name",
			"{ property() {} }",
			ast.Property{
				Kind:   ast.InitProperty,
				Key:    ident("property"),
				Value:  ast.FunctionExpression{Body: ast.BlockStatement{}},
				Method: true,
			},
		},
		{
			"generator method, non-computed name",
			"{ *property() {} }",
			ast.Property{
				Kind: ast.InitProperty,
				Key:  ident("property"),
				Value: ast.FunctionExpression{
					Body:      ast.BlockStatement{},
					Generator: true,
				},
				Method: true,
			},
		},
		{
			"async method, non-computed name",
			"{ async property() {} }",
			ast.Property{
				Kind: ast.InitProperty,
				Key:  ident("property"),
				Value: ast.FunctionExpression{
					Body:  ast.BlockStatement{},
					Async: true,
				},
				Method: true,
			},
		},
		{
			"async generator method, non-computed name",
			"{ async* property() {} }",
			ast.Property{
				Kind: ast.InitProperty,
				Key:  ident("property"),
				Value: ast.FunctionExpression{
					Body:      ast.BlockStatement{},
					Async:     true,
					Generator: true,
				},
				Method: true,
			},
		},
		{
			"init, computed name",
			"{ ['property']: null }",
			ast.Property{
				Kind:     ast.InitProperty,
				Key:      ast.StringLiteral{Value: "property", Raw: "'property'"},
				Computed: true,
				Value:    ast.NullLiteral{},
			},
		},
		{
			"getter, computed name",
			"{ get ['property']() {} }",
			ast.Property{
				Kind:     ast.GetProperty,
				Key:      ast.StringLiteral{Value: "property", Raw: "'property'"},
				Computed: true,
				Value:    ast.FunctionExpression{Body: ast.BlockStatement{}},
			},
		},
		{
			"setter, computed name",
			"{ set ['property']() {} }",
			ast.Property{
				Kind:     ast.SetProperty,
				Key:      ast.StringLiteral{Value: "property", Raw: "'property'"},
				Computed: true,
				Value:    ast.FunctionExpression{Body: ast.BlockStatement{}},
			},
		},
		{
			"method, computed name",
			"{ ['property']() {} }",
			ast.Property{
				Kind:     ast.InitProperty,
				Key:      ast.StringLiteral{Value: "property", Raw: "'property'"},
				Computed: true,
				Value:    ast.FunctionExpression{Body: ast.BlockStatement{}},
				Method:   true,
			},
		},
		{
			"generator method, computed name",
			"{ *['property']() {} }",
			ast.Property{
				Kind:     ast.InitProperty,
				Key:      ast.StringLiteral{Value: "property", Raw: "'property'"},
				Computed: true,
				Value: ast.FunctionExpression{
					Body:      ast.BlockStatement{},
					Generator: true,
				},
				Method: true,
			},
		},
		{
			"async method, computed name",
			"{ async ['property']() {} }",
			ast.Property{
				Kind:     ast.InitProperty,
				Key:      ast.StringLiteral{Value: "property", Raw: "'property'"},
				Computed: true,
				Value: ast.FunctionExpression{
					Body:  ast.BlockStatement{},
					Async: true,
				},
				Method: true,
			},
		},
		{
			"async generator method, computed name",
			"{ async* ['property']() {} }",
			ast.Property{
				Kind:     ast.InitProperty,
				Key:      ast.StringLiteral{Value: "property", Raw: "'property'"},
				Computed: true,
				Value: ast.FunctionExpression{
					Body:      ast.BlockStatement{},
					Async:     true,
					Generator: true,
				},
				Method: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertTree(t, test.input, ast.ObjectExpression{
				Properties: []ast.Property{test.expected},
			}, ParseOptions{Mode: ExpressionMode})
		})
	}
}

func TestRegexpLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.RegExpLiteral
	}{
		{
			"character class containing slash",
			"/[/]/",
			ast.RegExpLiteral{Pattern: "[/]", Raw: "/[/]/"},
		},
		{
			"character class containing left bracket and slash",
			`/[\]/]/`,
			ast.RegExpLiteral{Pattern: `[\]/]`, Raw: `/[\]/]/`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertTree(t, test.input, ast.ModuleNode{
				Body: []ast.Node{
					ast.ExpressionStatement{
						Expression: test.expected,
					},
				},
			}, ParseOptions{Mode: ModuleMode})
		})
	}
}
