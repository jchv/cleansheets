package ast

// ScriptNode is the AST node for a script.
type ScriptNode struct {
	BaseNode
	Body []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ScriptNode) ESTree() interface{} {
	e := struct {
		Type       string        `json:"type"`
		Body       []interface{} `json:"body"`
		SourceType string        `json:"sourceType"`
	}{
		Type:       "Program",
		SourceType: "script",
	}
	for _, stmt := range n.Body {
		e.Body = append(e.Body, estree(stmt))
	}
	return e
}
