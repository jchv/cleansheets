package ast

// estreeIdent returns an identifier node with the given string. Our AST does
// not use Identifier nodes in cases where it is unambiguous, so this function
// is useful for converting to estree.
func estreeIdent(ident string) interface{} {
	if ident == "" {
		return nil
	}
	return struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{
		Type: "Identifier",
		Name: ident,
	}
}

// estree returns the result of calling the ESTree method if the node is
// non-nil, or nil otherwise. This is useful since nil nodes may appear in many
// different structures.
func estree(node Node) interface{} {
	if node != nil {
		return node.ESTree()
	}
	return nil
}
