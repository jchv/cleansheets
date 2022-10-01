package ast

// ArrayExpression is a node containing an array literal.
type ArrayExpression struct {
	BaseNode
	Elements []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ArrayExpression) ESTree() interface{} {
	e := struct {
		Type     string        `json:"type"`
		Elements []interface{} `json:"elements"`
	}{
		Type:     "ArrayExpression",
		Elements: []interface{}{},
	}
	for _, elem := range n.Elements {
		e.Elements = append(e.Elements, estree(elem))
	}
	return e
}

// ContainsTemporalNodes returns true if the node contains any temporal
// children.
func (n ArrayExpression) ContainsTemporalNodes() bool {
	for _, elem := range n.Elements {
		if elem.ContainsTemporalNodes() {
			return true
		}
	}
	return false
}

// ConditionalExpression is the AST node for a conditional expression
// statement.
//
// Sometimes called the ternary operator, this expression is the equivalent to
// the 'if' statement, except as an expression. The test expression is
// evaluated, and then if it is truthy, the consequent is returned from the
// expression. If it is falsy, the alternate is returned instead.
//
// For example:
//
//	i == 1 ? 0 : 2
//
// Would be represented as:
//
//	ConditionalExpression{
//	    Test: BinaryExpression{
//	        Left: Identifier{Name: "i"},
//	        Right: NumberLiteral{Value: 1, ...},
//	    },
//	    Consequent: NumberLiteral{Value: 0, ...},
//	    Alternate: NumberLiteral{Value: 2, ...},
//	}
type ConditionalExpression struct {
	BaseNode

	Test       Node
	Consequent Node
	Alternate  Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ConditionalExpression) ESTree() interface{} {
	return struct {
		Type       string      `json:"type"`
		Test       interface{} `json:"test"`
		Alternate  interface{} `json:"alternate"`
		Consequent interface{} `json:"consequent"`
	}{
		Type:       "ConditionalExpression",
		Test:       estree(n.Test),
		Alternate:  estree(n.Alternate),
		Consequent: estree(n.Consequent),
	}
}

// FormalParameters stores function parameters.
type FormalParameters struct {
	Parameters    []BindingElement
	RestParameter string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n FormalParameters) ESTree() interface{} {
	e := []interface{}{}
	for _, elem := range n.Parameters {
		e = append(e, elem.ESTree())
	}
	if n.RestParameter != "" {
		e = append(e, struct {
			Type     string      `json:"type"`
			Argument interface{} `json:"argument"`
		}{
			Type:     "RestElement",
			Argument: estreeIdent(n.RestParameter),
		})
	}
	return e
}

// FunctionExpression stores a function expression.
type FunctionExpression struct {
	BaseNode
	ID         string
	Params     FormalParameters
	Body       Node
	Generator  bool
	Expression bool
	Async      bool
	Arrow      bool
}

// ESTree returns the corresponding ESTree representation for this node.
func (n FunctionExpression) ESTree() interface{} {
	typ := "FunctionExpression"
	if n.Arrow {
		typ = "ArrowFunctionExpression"
	}
	return struct {
		Type       string      `json:"type"`
		ID         interface{} `json:"id"`
		Params     interface{} `json:"params"`
		Body       interface{} `json:"body"`
		Generator  bool        `json:"generator"`
		Expression bool        `json:"expression"`
		Async      bool        `json:"async"`
	}{
		Type:       typ,
		ID:         estreeIdent(n.ID),
		Params:     n.Params.ESTree(),
		Body:       estree(n.Body),
		Generator:  n.Generator,
		Expression: n.Expression,
		Async:      n.Async,
	}
}

// Identifier is the node for an ECMAScript identifier expression.
//
// For example:
//
//	window
//
// Would be represented as:
//
//	Identifier{
//	    Name: "window",
//	}
type Identifier struct {
	BaseNode
	Name string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n Identifier) ESTree() interface{} {
	return struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{
		Type: "Identifier",
		Name: n.Name,
	}
}

// ThisExpression is a node for the ECMAScript `this` keyword.
type ThisExpression struct {
	BaseNode
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ThisExpression) ESTree() interface{} {
	return struct {
		Type string `json:"type"`
	}{
		Type: "ThisExpression",
	}
}

// MemberExpression is a node for an ECMAScript member expression.
type MemberExpression struct {
	BaseNode
	Computed bool
	Object   Node
	Property Node
	Optional bool
}

// ESTree returns the corresponding ESTree representation for this node.
func (n MemberExpression) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Computed bool        `json:"computed"`
		Object   interface{} `json:"object"`
		Property interface{} `json:"property"`
		Optional bool        `json:"optional,omitempty"`
	}{
		Type:     "MemberExpression",
		Computed: n.Computed,
		Object:   estree(n.Object),
		Property: estree(n.Property),
		Optional: n.Optional,
	}
}

// ParenthesizedExpression is the node for an expression surrounded in parens.
type ParenthesizedExpression struct {
	BaseNode
	Expression Node
}

// ESTree returns the corresponding ESTree representation for this node.
// Because the ESTree AST does not store parenthetical expressions, this
// returns the underlying expression.
func (n ParenthesizedExpression) ESTree() interface{} {
	// ESTree does not retain parenthesis.
	// TODO: Maybe support Babel extension for extra data.
	return estree(n.Expression)
}

// SpreadElement is a node containing a spread operator.
type SpreadElement struct {
	BaseNode
	Argument Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n SpreadElement) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Argument interface{} `json:"argument"`
	}{
		Type:     "SpreadElement",
		Argument: estree(n.Argument),
	}
}

// CallExpression is a node containing a call expression.
type CallExpression struct {
	BaseNode
	Callee    Node
	Optional  bool
	Arguments []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n CallExpression) ESTree() interface{} {
	e := struct {
		Type      string        `json:"type"`
		Callee    interface{}   `json:"callee"`
		Optional  bool          `json:"optional,omitempty"`
		Arguments []interface{} `json:"arguments"`
	}{
		Type:      "CallExpression",
		Callee:    estree(n.Callee),
		Optional:  n.Optional,
		Arguments: []interface{}{},
	}
	for _, arg := range n.Arguments {
		e.Arguments = append(e.Arguments, estree(arg))
	}
	return e
}

// NewExpression is a node containing a new expression.
type NewExpression struct {
	BaseNode
	Callee    Node
	Arguments []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n NewExpression) ESTree() interface{} {
	e := struct {
		Type      string        `json:"type"`
		Callee    interface{}   `json:"callee"`
		Arguments []interface{} `json:"arguments"`
	}{
		Type:      "NewExpression",
		Callee:    estree(n.Callee),
		Arguments: []interface{}{},
	}
	for _, arg := range n.Arguments {
		e.Arguments = append(e.Arguments, estree(arg))
	}
	return e
}

// PropertyKind is an enumeration type for different kinds of properties.
type PropertyKind int

const (
	// InitProperty is the kind of the most basic property, which simply sets
	// a property value for a given key.
	InitProperty PropertyKind = iota

	// GetProperty is the kind of a getter property, which provides a getter
	// function to fetch the property key's value.
	GetProperty

	// SetProperty is the kind of a setter property, which provides a setter
	// function to handle values being written to the property key.
	SetProperty
)

// estreePropertyKindMap maps PropertyKind values to their corresponding ESTree
// strings.
var estreePropertyKindMap = map[PropertyKind]string{
	InitProperty: "init",
	GetProperty:  "get",
	SetProperty:  "set",
}

// Property stores a single property value in an object expression.
type Property struct {
	// Key specifies a property key. In the non-computed cases (e.g. {a: 1}),
	// it will be treated as a literal value for the property name. In the
	// computed cases, it will be treated as an expression that is evaluated
	// to calculate the property name.
	Key Node

	// Computed specifies whether the key node is a computed property key
	// (e.g. {["expression"]: 1}) -- in this case, the `Key` will be treated as
	// an expression to be evaluated. Otherwise, it will be read as an
	// Identifier or Literal node specifying a property name.
	Computed bool

	// Value is the node that represents the value for the property. In case of
	// init properties, the property is simply set to this value.
	Value Node

	// DestructureInit is set when parsing a destructuring assignment as a
	// property. This happpens when an ambiguity forces the parser to parse an
	// arrow function as an expression.
	DestructureInit Node

	// Method specifies whether or not this value is using the method shorthand
	// syntax. Note that this is only true in case of init properties; getter
	// and setter properties always have this field set to false.
	Method bool

	// Kind is the kind of property. Most properties are init properties.
	Kind PropertyKind
}

// ESTree returns the corresponding ESTree representation for this node.
func (n Property) ESTree() interface{} {
	k := estree(n.Key)
	v, shorthand := estree(n.Value), false
	if v == nil {
		v, shorthand = k, true
	}
	return struct {
		Type      string      `json:"type"`
		Key       interface{} `json:"key"`
		Computed  bool        `json:"computed"`
		Value     interface{} `json:"value"`
		Kind      string      `json:"kind"`
		Method    bool        `json:"method"`
		Shorthand bool        `json:"shorthand"`
	}{
		Type:      "Property",
		Key:       k,
		Computed:  n.Computed,
		Value:     v,
		Kind:      estreePropertyKindMap[n.Kind],
		Method:    n.Method,
		Shorthand: shorthand,
	}
}

// ObjectExpression is a node containing an object literal.
type ObjectExpression struct {
	BaseNode
	Properties []Property
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ObjectExpression) ESTree() interface{} {
	e := struct {
		Type       string        `json:"type"`
		Properties []interface{} `json:"properties"`
	}{
		Type:       "ObjectExpression",
		Properties: []interface{}{},
	}
	for _, elem := range n.Properties {
		e.Properties = append(e.Properties, elem.ESTree())
	}
	return e
}

// ContainsTemporalNodes returns true if the node contains any temporal
// children.
func (n ObjectExpression) ContainsTemporalNodes() bool {
	for _, prop := range n.Properties {
		if prop.Key.ContainsTemporalNodes() || prop.Value.ContainsTemporalNodes() {
			return true
		}
	}
	return false
}

// SequenceExpression is a node containing expressions separated with the comma
// operator.
//
// For example:
//
//	1, 2, 3
//
// Would be represented as:
//
//	    SequenceExpression{
//		       Expressions: []Node{
//	            NumberLiteral{Value: 1, ...},
//	            NumberLiteral{Value: 2, ...},
//	            NumberLiteral{Value: 3, ...},
//	        },
//	    }
type SequenceExpression struct {
	BaseNode
	Expressions []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n SequenceExpression) ESTree() interface{} {
	e := struct {
		Type        string        `json:"type"`
		Expressions []interface{} `json:"expressions"`
	}{
		Type:        "SequenceExpression",
		Expressions: []interface{}{},
	}
	for _, expr := range n.Expressions {
		e.Expressions = append(e.Expressions, estree(expr))
	}
	return e
}

// ContainsTemporalNodes returns true if the node contains any temporal
// children.
func (n SequenceExpression) ContainsTemporalNodes() bool {
	for _, expr := range n.Expressions {
		if expr.ContainsTemporalNodes() {
			return true
		}
	}
	return false
}

// ClassExpression is the AST node that corresponds to an ECMAscript
// class expression.
//
// For example:
//
//     class { }
//
// Would be represented as:
//
//     ClassExpression{
// 	       ID: "",
//         SuperClass: "",
//         Body: ClassBody{},
//     }
type ClassExpression struct {
	BaseNode
	ID         string
	SuperClass Node
	Body       []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ClassExpression) ESTree() interface{} {
	e := struct {
		Type       string      `json:"type"`
		ID         interface{} `json:"id"`
		SuperClass interface{} `json:"params"`
		Body       struct {
			Type string        `json:"type"`
			Body []interface{} `json:"body"`
		} `json:"body"`
	}{
		Type:       "ClassExpression",
		ID:         estreeIdent(n.ID),
		SuperClass: estree(n.SuperClass),
	}

	e.Body.Type = "ClassBody"
	for _, elem := range n.Body {
		e.Body.Body = append(e.Body.Body, estree(elem))
	}

	return e
}
