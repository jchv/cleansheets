#!/bin/sh
set -e

mkdir -p ../dist
cp ./index.html ../dist

# Parser demo
mkdir -p ../dist/parser
GOOS=js GOARCH=wasm go build -o ../dist/parser/parser.wasm ./parser
cp ./parser/parser.html ../dist/parser/parser.html
cp ./parser/wasm_exec.js ../dist/parser/wasm_exec.js
