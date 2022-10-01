package ast

// FunctionDeclaration is the AST node that corresponds to an ECMAscript
// function declaration.
//
// For example:
//
//     function a() { }
//
// Would be represented as:
//
//     FunctionDeclaration{
// 	       ID: "a",
//         Params: FormalParameters{},
//         Body: BlockStatement{},
//     }
type FunctionDeclaration struct {
	BaseNode
	ID         string
	Params     FormalParameters
	Body       BlockStatement
	Generator  bool
	Expression bool
	Async      bool
}

// ESTree returns the corresponding ESTree representation for this node.
func (n FunctionDeclaration) ESTree() interface{} {
	return struct {
		Type       string      `json:"type"`
		ID         interface{} `json:"id"`
		Params     interface{} `json:"params"`
		Body       interface{} `json:"body"`
		Generator  bool        `json:"generator"`
		Expression bool        `json:"expression"`
		Async      bool        `json:"async"`
	}{
		Type:       "FunctionDeclaration",
		ID:         estreeIdent(n.ID),
		Params:     n.Params.ESTree(),
		Body:       estree(n.Body),
		Generator:  n.Generator,
		Expression: n.Expression,
		Async:      n.Async,
	}
}

// ClassDeclaration is the AST node that corresponds to an ECMAscript
// class declaration.
//
// For example:
//
//     class a { }
//
// Would be represented as:
//
//     ClassDeclaration{
// 	       ID: "a",
//         SuperClass: "",
//         Body: ClassBody{},
//     }
type ClassDeclaration struct {
	BaseNode
	ID         string
	SuperClass Node
	Body       []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ClassDeclaration) ESTree() interface{} {
	e := struct {
		Type       string      `json:"type"`
		ID         interface{} `json:"id"`
		SuperClass interface{} `json:"params"`
		Body       struct {
			Type string        `json:"type"`
			Body []interface{} `json:"body"`
		} `json:"body"`
	}{
		Type:       "ClassDeclaration",
		ID:         estreeIdent(n.ID),
		SuperClass: estree(n.SuperClass),
	}

	e.Body.Type = "ClassBody"
	for _, elem := range n.Body {
		e.Body.Body = append(e.Body.Body, estree(elem))
	}

	return e
}

type MethodKind int

const (
	Method MethodKind = iota
	GetMethod
	SetMethod
)

// estreeMethodKindMap maps MethodKind values to their corresponding ESTree strings.
var estreeMethodKindMap = map[MethodKind]string{
	Method:    "method",
	GetMethod: "get",
	SetMethod: "set",
}

// MethodDefinition represents a method in a class body.
type MethodDefinition struct {
	BaseNode
	Key      Node
	Computed bool
	Value    FunctionExpression
	Kind     MethodKind
	Static   bool
}

// ESTree returns the corresponding ESTree representation for this node.
func (n MethodDefinition) ESTree() interface{} {
	return struct {
		Type     string      `json:"type"`
		Key      interface{} `json:"key"`
		Computed bool        `json:"computed"`
		Value    interface{} `json:"value"`
		Kind     string      `json:"kind"`
		Static   bool        `json:"static"`
	}{
		Type:     "MethodDefinition",
		Key:      estree(n.Key),
		Computed: n.Computed,
		Value:    estree(n.Value),
		Kind:     estreeMethodKindMap[n.Kind],
		Static:   n.Static,
	}
}
