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
