package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/jchv/cleansheets/ecmascript/lexer"
	"github.com/jchv/cleansheets/ecmascript/parser"
)

func main() {
	flag.Parse()

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	for i, filename := range flag.Args() {
		// Write separator if multiple files.
		if i != 0 {
			os.Stdout.Write([]byte("\n---\n"))
		}

		// Open file for reading and create a buffered reader.
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Could not open file for reading: %q", filename)
		}
		defer func(file *os.File) {
			if err := file.Close(); err != nil {
				log.Printf("Warning: Error closing file: %v", err)
			}
		}(file)
		reader := bufio.NewReader(file)

		// Try to calculate a file URL.
		absname, err := filepath.Abs(filename)
		if err != nil {
			absname = filename
		}
		url := &url.URL{}
		url.Scheme = "file"
		url.Path = absname
		log.Printf("Parsing %q...", url)

		// Parse script.
		script, err := parser.NewParser(lexer.NewLexer(lexer.NewScanner(reader, url))).Parse(parser.ParseOptions{Mode: parser.ScriptMode})
		if err != nil {
			log.Fatalf("Could not parse ECMAscript file %q: %v", filename, err)
		}

		// Output ESTree AST.
		err = encoder.Encode(script.ESTree())
		if err != nil {
			log.Fatalf("Error while encoding ESTree AST: %v", err)
		}
	}
}
