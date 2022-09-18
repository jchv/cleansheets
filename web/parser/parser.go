//go:build js
// +build js

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"github.com/jchv/cleansheets/ecmascript/lexer"
	"github.com/jchv/cleansheets/ecmascript/parser"
)

func main() {
	c := make(chan struct{}, 0)
	js.Global().Call("parserLoaded", js.FuncOf(ParseES))
	<-c
}

func ParseES(this js.Value, p []js.Value) interface{} {
	n, err := parser.NewParser(lexer.NewLexer(lexer.NewScanner(strings.NewReader(p[0].String()), nil))).Parse(parser.ParseOptions{Mode: parser.ScriptMode})
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	w := &strings.Builder{}
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	err = e.Encode(n.ESTree())
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{"result": w.String()}
}
