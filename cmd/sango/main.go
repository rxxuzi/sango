package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rxxuzi/sango/pkg/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename.sango> [options]\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	
	// Read the source file
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Create lexer
	l := lexer.New(string(source))
	
	// For now, just tokenize and print tokens
	fmt.Printf("=== Lexing %s ===\n", filename)
	for {
		tok := l.NextToken()
		if tok.Type == lexer.EOF {
			break
		}
		fmt.Printf("%s\n", tok)
	}
}