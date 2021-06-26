package ast

// ModuleNode is the AST node for a module.
type ModuleNode struct {
	BaseNode
	Body []Node
}

// ESTree returns the corresponding ESTree representation for this node.
func (n ModuleNode) ESTree() interface{} {
	e := struct {
		Type       string        `json:"type"`
		Body       []interface{} `json:"body"`
		SourceType string        `json:"sourceType"`
	}{
		Type:       "Program",
		SourceType: "module",
	}
	for _, stmt := range n.Body {
		e.Body = append(e.Body, estree(stmt))
	}
	return e
}

// ImportDeclNode is the AST node for an import declaration.
type ImportDeclNode struct {
	BaseNode

	// Possible combinations:
	// - none:
	//       import "react";
	// - DefaultBinding:
	//       import React from "react";
	// - DefaultBinding + NameSpace
	//       import React, * as ReactNS from "react";
	// - DefaultBinding + NamedImports
	//       import React, {Component} from "react";
	// - NameSpace
	//       import * as React from "react";
	// - NamedImports
	//       import {Component as ReactComponent, useState} from "react";

	// Default binding, e.g. import React from "react";
	DefaultBinding *ImportDefaultBinding

	// Namespace binding, e.g. import * as React from "react";
	NameSpace *NameSpaceImport

	// Named imports, e.g. import {Component as ReactComponent} from "react";
	NamedImports []NamedImport

	// Module to import; string literal.
	Module string
}

func (n ImportDeclNode) ESTree() interface{} {
	panic("unimplemented")
}

type ImportDefaultBinding struct {
	Identifier string
}

type NameSpaceImport struct {
	Identifier string
}

type NamedImport struct {
	Identifier string
	AsBinding  string
}
