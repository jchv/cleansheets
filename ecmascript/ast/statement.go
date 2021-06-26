package ast

// BlockStatement is the AST node for a block.
type BlockStatement struct {
	BaseNode
	Body []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n BlockStatement) ESTree() interface{} {
	e := struct {
		Type string        `json:"type"`
		Body []interface{} `json:"body"`
	}{
		Type: "BlockStatement",
		Body: []interface{}{},
	}
	for _, stmt := range n.Body {
		e.Body = append(e.Body, estree(stmt))
	}
	return e
}

// EmptyStatement is the AST node for an empty expression statement.
type EmptyStatement struct {
	BaseNode
}

// ESTree returns the corresponding ESTree representation for this node.
func (n EmptyStatement) ESTree() interface{} {
	return struct {
		Type string `json:"type"`
	}{
		Type: "EmptyStatement",
	}
}

// ExpressionStatement is the AST node for an expression statement.
type ExpressionStatement struct {
	BaseNode

	Expression Node
	Directive  string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ExpressionStatement) ESTree() interface{} {
	return struct {
		Type       string      `json:"type"`
		Expression interface{} `json:"expression"`
		Directive  string      `json:"directive,omitempty"`
	}{
		Type:       "ExpressionStatement",
		Expression: estree(n.Expression),
		Directive:  n.Directive,
	}
}

// VarKind is an enumeration type for different kinds of variables.
type VarKind int

const (
	// VarDeclaration is the traditional var statement.
	VarDeclaration VarKind = iota

	// LetDeclaration is the newer declaration syntax that introduces temporal
	// dead zones.
	LetDeclaration

	// ConstDeclaration is like LetDeclaration, but the variable cannot be re-
	// assigned.
	ConstDeclaration
)

// estreeVarKindMap maps VarKind values to their corresponding ESTree strings.
var estreeVarKindMap = map[VarKind]string{
	VarDeclaration:   "var",
	LetDeclaration:   "let",
	ConstDeclaration: "const",
}

// VariableDeclaration is the AST node for a variable declaration statement.
type VariableDeclaration struct {
	BaseNode
	Declarations []VariableDeclarator
	Kind         VarKind
}

// ESTree returns the corresponding ESTree representation for this node.
func (n VariableDeclaration) ESTree() interface{} {
	e := struct {
		Type         string        `json:"type"`
		Declarations []interface{} `json:"declarations"`
		Kind         string        `json:"kind"`
	}{
		Type: "VariableDeclaration",
		Kind: estreeVarKindMap[n.Kind], // TODO
	}
	for _, decl := range n.Declarations {
		e.Declarations = append(e.Declarations, decl.ESTree())
	}
	return e
}

// VariableDeclarator contains one fragment of variable declaration.
type VariableDeclarator struct {
	// One and only one of ID, Pattern must be set.
	// - ID: var x;
	// - Pattern: var {x} = y;
	ID BindingPattern

	// Expression initializing the variable.
	// Optional if Identifier is set, non-optional otherwise.
	Init Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n VariableDeclarator) ESTree() interface{} {
	return struct {
		Type string      `json:"type"`
		ID   interface{} `json:"id"`
		Init interface{} `json:"init"`
	}{
		Type: "VariableDeclarator",
		ID:   n.ID.ESTree(),
		Init: estree(n.Init),
	}
}

type BindingPattern struct {
	// One and only one of Identifier, ObjectPattern, ArrayPattern must be set.
	// - Identifier: var x = y;
	// - ObjectPattern: var {x} = y;
	// - ArrayPattern: var [x] = y;
	Identifier    string
	ObjectPattern *ObjectBindingPattern
	ArrayPattern  *ArrayBindingPattern
}

// ESTree returns the corresponding ESTree representation for this node.
func (n BindingPattern) ESTree() interface{} {
	if n.Identifier != "" {
		return estreeIdent(n.Identifier)
	} else if n.ObjectPattern != nil {
		return n.ObjectPattern.ESTree()
	} else if n.ArrayPattern != nil {
		return n.ArrayPattern.ESTree()
	}
	return nil
}

// ObjectBindingPattern contains a full object binding pattern.
type ObjectBindingPattern struct {
	Properties []BindingProperty

	// Optional: rest pattern. e.g. {...a}
	RestElement string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ObjectBindingPattern) ESTree() interface{} {
	e := struct {
		Type       string        `json:"type"`
		Properties []interface{} `json:"properties"`
	}{
		Type:       "ObjectPattern",
		Properties: []interface{}{},
	}
	for _, p := range n.Properties {
		e.Properties = append(e.Properties, p.ESTree())
	}
	if n.RestElement != "" {
		e.Properties = append(e.Properties, struct {
			Type     string      `json:"type"`
			Argument interface{} `json:"argument"`
		}{
			Type:     "RestElement",
			Argument: estreeIdent(n.RestElement),
		})
	}
	return e
}

// ArrayBindingPattern contains a full array binding pattern.
type ArrayBindingPattern struct {
	// Binding elements, i.e. [ a, b ]
	Elements []BindingElement

	// Optional.
	RestElement BindingPattern
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ArrayBindingPattern) ESTree() interface{} {
	e := struct {
		Type     string        `json:"type"`
		Elements []interface{} `json:"elements"`
	}{
		Type:     "ArrayPattern",
		Elements: []interface{}{},
	}
	for _, p := range n.Elements {
		e.Elements = append(e.Elements, p.ESTree())
	}
	rest := n.RestElement.ESTree()
	if rest != nil {
		e.Elements = append(e.Elements, struct {
			Type     string      `json:"type"`
			Argument interface{} `json:"argument"`
		}{
			Type:     "RestElement",
			Argument: rest,
		})
	}
	return e
}

// BindingProperty is a binding property in an object binding pattern.
type BindingProperty struct {
	// Property name. If BindingIdentifier is not specified, this is also the
	// BindingIdentifier.
	PropertyName string

	// Only one of BindingIdentifier and BindingPattern can be set.
	// - none: { PropertyName = Initializer }
	// - BindingIdentifier: { PropertyName: BindingIdentifier = Initializer }
	// - BindingPattern: { PropertyName: { BindingPattern } = Initializer }
	Value BindingPattern

	// Specifies the default value. Optional in all cases.
	Init Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n BindingProperty) ESTree() interface{} {
	k := estreeIdent(n.PropertyName)
	v, shorthand := n.Value.ESTree(), false
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
		Computed:  false, // TODO?
		Value:     v,
		Kind:      "init",
		Method:    false,
		Shorthand: shorthand,
	}
}

// BindingElement is a binding element in a binding pattern.
type BindingElement struct {
	// Only one of BindingIdentifier and Value can be set.
	// - none: [ , ]  (NOTE: Element is an elision.)
	// - BindingIdentifier: [ BindingIdentifier ]
	// - Value: [ { Value } ]
	Value BindingPattern

	// Specifies the default value. Optional in all cases.
	Init Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n BindingElement) ESTree() interface{} {
	e := n.Value.ESTree()
	if n.Init != nil {
		e = struct {
			Type  string      `json:"type"`
			Left  interface{} `json:"left"`
			Right interface{} `json:"right"`
		}{
			Type:  "AssignmentPattern",
			Left:  e,
			Right: estree(n.Init),
		}
	}
	return e
}

// ContinueStatement is a node containing a continue statement.
type ContinueStatement struct {
	BaseNode
	Label string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ContinueStatement) ESTree() interface{} {
	return struct {
		Type  string      `json:"type"`
		Label interface{} `json:"label"`
	}{
		Type:  "ContinueStatement",
		Label: estreeIdent(n.Label),
	}
}

// BreakStatement is a node containing a break statement.
type BreakStatement struct {
	BaseNode
	Label string
}

// ESTree returns the corresponding ESTree representation for this node.
func (n BreakStatement) ESTree() interface{} {
	return struct {
		Type  string      `json:"type"`
		Label interface{} `json:"label"`
	}{
		Type:  "BreakStatement",
		Label: estreeIdent(n.Label),
	}
}

// ReturnStatement is a node containing a return statement.
type ReturnStatement struct {
	BaseNode
	Argument Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ReturnStatement) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Argument interface{} `json:"argument"`
	}{
		Type:     "ReturnStatement",
		Argument: estree(n.Argument),
	}
}

// ThrowStatement is a node containing a throw statement.
type ThrowStatement struct {
	BaseNode
	Argument Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ThrowStatement) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Argument interface{} `json:"argument"`
	}{
		Type:     "ThrowStatement",
		Argument: estree(n.Argument),
	}
}

// IfStatement is a node containing an if statement.
type IfStatement struct {
	BaseNode
	Test       Node
	Consequent Node
	Alternate  Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n IfStatement) ESTree() interface{} {
	return struct {
		Type       string      `json:"type"`
		Test       interface{} `json:"test"`
		Consequent interface{} `json:"consequent"`
		Alternate  interface{} `json:"alternate"`
	}{
		Type:       "IfStatement",
		Test:       estree(n.Test),
		Consequent: estree(n.Consequent),
		Alternate:  estree(n.Alternate),
	}
}

// WhileStatement is a node containing an while statement.
type WhileStatement struct {
	BaseNode
	Test Node
	Body Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n WhileStatement) ESTree() interface{} {
	return struct {
		Type string      `json:"type"`
		Test interface{} `json:"test"`
		Body interface{} `json:"body"`
	}{
		Type: "WhileStatement",
		Test: estree(n.Test),
		Body: estree(n.Body),
	}
}

// DoWhileStatement is a node containing an do/while statement.
type DoWhileStatement struct {
	BaseNode
	Body Node
	Test Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n DoWhileStatement) ESTree() interface{} {
	return struct {
		Type string      `json:"type"`
		Test interface{} `json:"test"`
		Body interface{} `json:"body"`
	}{
		Type: "DoWhileStatement",
		Test: estree(n.Test),
		Body: estree(n.Body),
	}
}

// ForStatement is a node containing a for statement.
type ForStatement struct {
	BaseNode
	Init   Node
	Test   Node
	Update Node
	Body   Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ForStatement) ESTree() interface{} {
	return struct {
		Type   string      `json:"type"`
		Init   interface{} `json:"init"`
		Test   interface{} `json:"test"`
		Update interface{} `json:"update"`
		Body   interface{} `json:"body"`
	}{
		Type:   "ForStatement",
		Init:   estree(n.Init),
		Test:   estree(n.Test),
		Update: estree(n.Update),
		Body:   estree(n.Body),
	}
}

// ForInStatement is a node containing a for in statement.
type ForInStatement struct {
	BaseNode
	Left  Node
	Right Node
	Body  Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ForInStatement) ESTree() interface{} {
	return struct {
		Type  string      `json:"type"`
		Each  bool        `json:"each"`
		Left  interface{} `json:"left"`
		Right interface{} `json:"right"`
		Body  interface{} `json:"body"`
	}{
		Type:  "ForInStatement",
		Each:  false,
		Left:  estree(n.Left),
		Right: estree(n.Right),
		Body:  estree(n.Body),
	}
}

// ForOfStatement is a node containing a for of statement.
type ForOfStatement struct {
	BaseNode
	Left  Node
	Right Node
	Body  Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ForOfStatement) ESTree() interface{} {
	return struct {
		Type  string      `json:"type"`
		Left  interface{} `json:"left"`
		Right interface{} `json:"right"`
		Body  interface{} `json:"body"`
	}{
		Type:  "ForOfStatement",
		Left:  estree(n.Left),
		Right: estree(n.Right),
		Body:  estree(n.Body),
	}
}

// SwitchStatement is a node containing a switch statement.
type SwitchStatement struct {
	BaseNode
	Discriminant Node
	Cases        []SwitchCase
}

// ESTree returns the corresponding ESTree representation for this node.
func (n SwitchStatement) ESTree() interface{} {
	e := struct {
		Type         string        `json:"type"`
		Discriminant interface{}   `json:"discriminant"`
		Cases        []interface{} `json:"cases"`
	}{
		Type:         "SwitchStatement",
		Discriminant: estree(n.Discriminant),
		Cases:        []interface{}{},
	}
	for _, stmt := range n.Cases {
		e.Cases = append(e.Cases, stmt.ESTree())
	}
	return e
}

// SwitchCase contains an individual switch case.
type SwitchCase struct {
	Test       Node
	Consequent []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n SwitchCase) ESTree() interface{} {
	e := struct {
		Type       string        `json:"type"`
		Test       interface{}   `json:"test"`
		Consequent []interface{} `json:"consequent"`
	}{
		Type:       "SwitchCase",
		Test:       estree(n.Test),
		Consequent: []interface{}{},
	}
	for _, stmt := range n.Consequent {
		e.Consequent = append(e.Consequent, estree(stmt))
	}
	return e
}

// LabeledStatement is a node containing an ECMAScript labelled statement.
type LabeledStatement struct {
	BaseNode
	Label string
	Body  Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n LabeledStatement) ESTree() interface{} {
	return struct {
		Type  string      `json:"type"`
		Label interface{} `json:"label"`
		Body  interface{} `json:"body"`
	}{
		Type:  "LabeledStatement",
		Label: estreeIdent(n.Label),
		Body:  estree(n.Body),
	}
}

// TryStatement is a node containing an ECMAScript try/catch/finally statement.
type TryStatement struct {
	BaseNode
	Block     Node
	Handler   Node
	Finalizer Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n TryStatement) ESTree() interface{} {
	return struct {
		Type      string      `json:"type"`
		Block     interface{} `json:"block"`
		Handler   interface{} `json:"handler"`
		Finalizer interface{} `json:"finalizer"`
	}{
		Type:      "TryStatement",
		Block:     estree(n.Block),
		Handler:   estree(n.Handler),
		Finalizer: estree(n.Finalizer),
	}
}

// CatchClause is a node representing the ECMAScript catch clause.
type CatchClause struct {
	BaseNode
	Param BindingPattern
	Body  Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n CatchClause) ESTree() interface{} {
	return struct {
		Type  string      `json:"type"`
		Param interface{} `json:"param"`
		Body  interface{} `json:"body"`
	}{
		Type:  "CatchClause",
		Param: n.Param.ESTree(),
		Body:  estree(n.Body),
	}
}
