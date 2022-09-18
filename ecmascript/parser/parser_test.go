package parser

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jchv/cleansheets/ecmascript/ast"
	"github.com/jchv/cleansheets/ecmascript/lexer"
)

func ident(n string) ast.Identifier {
	return ast.Identifier{Name: n}
}

func assertTree(t *testing.T, input interface{}, expected ast.Node, opt ParseOptions, r ...bool) {
	var result ast.Node
	var err error

	todo := false
	if len(r) > 0 {
		todo = r[0]
	}

	switch b := input.(type) {
	case string:
		result, err = NewParser(lexer.NewLexer(lexer.NewScanner(strings.NewReader(b), nil))).Parse(opt)
	case []byte:
		result, err = NewParser(lexer.NewLexer(lexer.NewScanner(bytes.NewReader(b), nil))).Parse(opt)
	case io.Reader:
		result, err = NewParser(lexer.NewLexer(lexer.NewScanner(bufio.NewReader(b), nil))).Parse(opt)
	default:
		t.Fatalf("unsupported input type %t", b)
	}
	if err != nil && todo == false {
		t.Errorf("error parsing code: %v", err)
		return
	} else if err != nil && todo == true {
		t.Logf("todo: %v", err)
		return
	}

	ast.ClearSpans(result)
	if diff := cmp.Diff(expected, result, cmpopts.IgnoreUnexported(ast.BaseNode{})); diff != "" {
		if todo {
			t.Logf("todo: ast mismatch (-expected +result):\n%s", diff)
		} else {
			t.Errorf("ast mismatch (-expected +result):\n%s", diff)
		}
	} else {
		if todo {
			t.Error("assertion marked todo passed")
		}
	}
}

func TestParseImport(t *testing.T) {
	tests := []struct {
		s, e string
	}{
		{s: ``},

		// Import declarations.
		{s: `import 'react';`},
		{s: `import React from "react";`},
		{s: `import React, * as ReactNS from "react";`},
		{s: `import React, {Component,} from "react";`},
		{s: `import * as React from "react";`},
		{s: `import {Component as ReactComponent, useState} from "react";`},
		{s: `import React, { } from "react";`},

		// Import declarations with non-reserved keywords.
		{s: `import as, * as as from "reserved-never"; import as, {as as as} from "reserved-never";`},
		{s: `import async, * as async from "reserved-never"; import async, {async as async} from "reserved-never";`},
		{s: `import from, * as from from "reserved-never"; import from, {from as from} from "reserved-never";`},
		{s: `import get, * as get from "reserved-never"; import get, {get as get} from "reserved-never";`},
		{s: `import meta, * as meta from "reserved-never"; import meta, {meta as meta} from "reserved-never";`},
		{s: `import of, * as of from "reserved-never"; import of, {of as of} from "reserved-never";`},
		{s: `import set, * as set from "reserved-never"; import set, {set as set} from "reserved-never";`},
		{s: `import target, * as target from "reserved-never"; import target, {target as target} from "reserved-never";`},
		{s: `import await, * as await from "reserved-async"; import await, {await as await} from "reserved-async";`},
		{s: `import yield, * as yield from "reserved-generator"; import yield, {yield as yield} from "reserved-generator";`},

		// Import syntax errors.
		{s: `import`, e: "syntax error"},
		{s: `import React`, e: "syntax error"},
		{s: `import React from`, e: "syntax error"},
		{s: `import React from react;`, e: "syntax error"},
		{s: `import * as React, {Component}`, e: "syntax error"},
		{s: `import { Component, , } from "react";`, e: "syntax error"},
		{s: `import { Component as } from "react";`, e: "syntax error"},
		{s: `import { Component from "react";`, e: "syntax error"},
		{s: `import React, React from "react";`, e: "syntax error"},
		{s: `import {Component} "react";`, e: "syntax error"},
		{s: `import {,} "react";`, e: "syntax error"},

		// Variable declarations.
		{s: `var i, j, [k] = false, {l} = 0, [...m] = null, {...n} = undefined, {o: p} = this;`},

		// Expressions.
		{s: `window.alert`},
		{s: `window.localStorage.getItem`},
		{s: `window.[]`, e: "syntax error"},
		{s: `8 + 4 * 3`},
		{s: `4 * 3 + 8`},
		{s: `/[/]/`},
		{s: `/[\]/]/`},
	}

	for _, test := range tests {
		t.Run(strconv.Quote(test.s), func(t *testing.T) {
			_, err := NewParser(lexer.NewLexer(lexer.NewScanner(strings.NewReader(test.s), nil))).Parse(ParseOptions{Mode: ModuleMode})
			if test.e == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error to contain %v, got nil", test.e)
				} else if !strings.Contains(err.Error(), test.e) {
					t.Errorf("expected error to contain %v, got %v", test.e, err.Error())
				}
			}
		})
	}
}

func TestParseLibraries(t *testing.T) {
	tests := []string{"lodash-core-v4.17.15.min", "lodash-v4.17.15.min", "ramda-v0.25.0.min", "react-v17.0.2"}
	for _, test := range tests {
		jsFileName := "testdata/" + test + ".js"
		f, err := os.Open(jsFileName)
		if err != nil {
			t.Fatal(err)
		}
		r := bufio.NewReader(f)
		url, _ := url.Parse("file://" + jsFileName)
		_, err = NewParser(lexer.NewLexer(lexer.NewScanner(r, url))).Parse(ParseOptions{Mode: ScriptMode})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkParseReact(b *testing.B) {
	b.StopTimer()
	data, err := ioutil.ReadFile("testdata/react-v17.0.2.js")
	if err != nil {
		b.Fatal(err)
	}
	url, _ := url.Parse("file:///testdata/react-v17.0.2.js")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewParser(lexer.NewLexer(lexer.NewScanner(bytes.NewReader(data), url))).Parse(ParseOptions{Mode: ScriptMode})
		if err != nil {
			b.Fatal(err)
		}
	}
}
