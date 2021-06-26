#!/bin/sh
set -e

mkdir -p ./dist
cp ./web/index.html ../dist

# Parser demo
mkdir -p ./dist/parser
GOOS=js GOARCH=wasm go build -o ./dist/parser/parser.wasm ./web/parser
cp ./web/parser/parser.html ./dist/parser/parser.html
cp ./web/parser/wasm_exec.js ./dist/parser/wasm_exec.js
