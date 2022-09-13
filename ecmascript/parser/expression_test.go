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

func TestArrowFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.FunctionExpression
		todo     bool
	}{
		{
			name:     "arrow function with no parameters",
			input:    "() => {}",
			expected: ast.FunctionExpression{Body: ast.BlockStatement{}, Arrow: true},
		},
		{
			name:     "arrow function with no parameters, async",
			input:    "async () => {}",
			expected: ast.FunctionExpression{Body: ast.BlockStatement{}, Arrow: true, Async: true},
		},
		{
			name:  "arrow function with parameter bare",
			input: "x => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with parameter bare, async",
			input: "async x => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
				},
				Body:  ast.BlockStatement{},
				Async: true,
				Arrow: true,
			},
		},
		{
			name:  "arrow function with parameter returning parameter",
			input: "x => x",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
				},
				Body:  ast.Identifier{Name: "x"},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with parameter returning parameter, async",
			input: "async x => x",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
				},
				Body:  ast.Identifier{Name: "x"},
				Async: true,
				Arrow: true,
			},
		},
		{
			name:  "arrow function with parameter parenthesized",
			input: "(x) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with parameter parenthesized, async",
			input: "async (x) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
				},
				Body:  ast.BlockStatement{},
				Async: true,
				Arrow: true,
			},
		},
		{
			name:  "arrow function with multiple parameters",
			input: "(x, y) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
						{Value: ast.BindingPattern{Identifier: "y"}},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with multiple parameters, async",
			input: "async (x, y) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
						{Value: ast.BindingPattern{Identifier: "y"}},
					},
				},
				Body:  ast.BlockStatement{},
				Async: true,
				Arrow: true,
			},
		},
		{
			name:  "arrow function with rest parameter",
			input: "(x, ...y) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
					RestParameter: "y",
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with rest parameter, async",
			input: "async (x, ...y) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{Value: ast.BindingPattern{Identifier: "x"}},
					},
					RestParameter: "y",
				},
				Body:  ast.BlockStatement{},
				Async: true,
				Arrow: true,
			},
		},
		{
			name:  "arrow function with default parameter",
			input: "(x = 1) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{Identifier: "x"},
							Init: ast.NumberLiteral{
								Value: 1,
								Raw:   "1",
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with default parameter, async",
			input: "async (x = 1) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{Identifier: "x"},
							Init: ast.NumberLiteral{
								Value: 1,
								Raw:   "1",
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Async: true,
				Arrow: true,
			},
		},
		{
			name:  "arrow function with object destructuring parameter",
			input: "({x}) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ObjectPattern: &ast.ObjectBindingPattern{
									Properties: []ast.BindingProperty{
										{PropertyName: "x"},
									},
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with object destructuring parameter and default",
			input: "({x = 1}) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ObjectPattern: &ast.ObjectBindingPattern{
									Properties: []ast.BindingProperty{
										{
											PropertyName: "x",
											Init: ast.NumberLiteral{
												Value: 1,
												Raw:   "1",
											},
										},
									},
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with object destructuring parameter, renamed with default",
			input: "({x: y = 1}) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ObjectPattern: &ast.ObjectBindingPattern{
									Properties: []ast.BindingProperty{
										{
											PropertyName: "x",
											Value: ast.BindingPattern{
												Identifier: "y",
											},
											Init: ast.NumberLiteral{
												Value: 1,
												Raw:   "1",
											},
										},
									},
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with object destructuring parameter and rest",
			input: "({x, ...y}) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ObjectPattern: &ast.ObjectBindingPattern{
									Properties: []ast.BindingProperty{
										{PropertyName: "x"},
									},
									RestElement: "y",
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with object destructuring parameter and default and rest",
			input: "({x = 1, ...y}) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ObjectPattern: &ast.ObjectBindingPattern{
									Properties: []ast.BindingProperty{
										{
											PropertyName: "x",
											Init: ast.NumberLiteral{
												Value: 1,
												Raw:   "1",
											},
										},
									},
									RestElement: "y",
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with object destructuring parameter and default and rest and other parameter",
			input: "({x = 1, ...y}, z) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ObjectPattern: &ast.ObjectBindingPattern{
									Properties: []ast.BindingProperty{
										{
											PropertyName: "x",
											Init: ast.NumberLiteral{
												Value: 1,
												Raw:   "1",
											},
										},
									},
									RestElement: "y",
								},
							},
						},
						{
							Value: ast.BindingPattern{
								Identifier: "z",
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with object destructuring parameter and default and rest and other parameter and rest",
			input: "({x = 1, ...y}, z, ...w) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ObjectPattern: &ast.ObjectBindingPattern{
									Properties: []ast.BindingProperty{
										{
											PropertyName: "x",
											Init: ast.NumberLiteral{
												Value: 1,
												Raw:   "1",
											},
										},
									},
									RestElement: "y",
								},
							},
						},
						{
							Value: ast.BindingPattern{
								Identifier: "z",
							},
						},
					},
					RestParameter: "w",
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with array destructuring parameter",
			input: "([x]) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ArrayPattern: &ast.ArrayBindingPattern{
									Elements: []ast.BindingElement{
										{
											Value: ast.BindingPattern{
												Identifier: "x",
											},
										},
									},
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with array destructuring parameter and default",
			input: "([x = 1]) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ArrayPattern: &ast.ArrayBindingPattern{
									Elements: []ast.BindingElement{
										{
											Value: ast.BindingPattern{
												Identifier: "x",
											},
											Init: ast.NumberLiteral{
												Value: 1,
												Raw:   "1",
											},
										},
									},
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with array destructuring parameter and rest",
			input: "([x, ...y]) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{{
						Value: ast.BindingPattern{
							ArrayPattern: &ast.ArrayBindingPattern{
								Elements: []ast.BindingElement{{
									Value: ast.BindingPattern{Identifier: "x"},
								}},
								RestElement: ast.BindingPattern{Identifier: "y"},
							},
						},
					}},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
		},
		{
			name:  "arrow function with array destructuring parameter and elided element",
			input: "([x, , y]) => {}",
			expected: ast.FunctionExpression{
				Params: ast.FormalParameters{
					Parameters: []ast.BindingElement{
						{
							Value: ast.BindingPattern{
								ArrayPattern: &ast.ArrayBindingPattern{
									Elements: []ast.BindingElement{
										{
											Value: ast.BindingPattern{
												Identifier: "x",
											},
										},
										{},
										{
											Value: ast.BindingPattern{
												Identifier: "y",
											},
										},
									},
								},
							},
						},
					},
				},
				Body:  ast.BlockStatement{},
				Arrow: true,
			},
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
			}, ParseOptions{Mode: ModuleMode}, test.todo)
			if test.todo {
				t.Skip("TODO")
			}
		})
	}
}
