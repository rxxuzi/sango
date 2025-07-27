package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rxxuzi/sango/pkg/lexer"
	"github.com/rxxuzi/sango/pkg/parser"
)

const VERSION = "v0.1.5"

type CompileMode int

const (
	ModeLexOnly CompileMode = iota
	ModeParseOnly
)

type Config struct {
	mode        CompileMode
	inputFile   string
	showHelp    bool
	showVersion bool
}

func main() {
	config := parseArgs()

	if config.showVersion {
		fmt.Printf("sangoc %s - Sango Compiler (Lexer/Parser Only)\n", VERSION)
		return
	}

	if config.showHelp {
		showHelp()
		return
	}

	if config.inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: No input file specified\n")
		showUsage()
		os.Exit(1)
	}

	// Check if input file exists
	if _, err := os.Stat(config.inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist\n", config.inputFile)
		os.Exit(1)
	}

	// Check file extension
	if filepath.Ext(config.inputFile) != ".sango" {
		fmt.Fprintf(os.Stderr, "Error: Input file must have .sango extension\n")
		os.Exit(1)
	}

	// Read source file
	source, err := ioutil.ReadFile(config.inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", config.inputFile, err)
		os.Exit(1)
	}

	// Execute compilation based on mode
	switch config.mode {
	case ModeLexOnly:
		lexOnly(string(source), config.inputFile)
	case ModeParseOnly:
		parseOnly(string(source), config.inputFile)
	}
}

func parseArgs() Config {
	var config Config

	// Define flags
	lexFlag := flag.Bool("l", false, "Lexical analysis only - show tokens")
	parseFlag := flag.Bool("p", false, "Parse only - show AST")
	versionFlag := flag.Bool("v", false, "Show version")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	config.showHelp = *helpFlag
	config.showVersion = *versionFlag

	// Determine mode
	modeCount := 0
	if *lexFlag {
		config.mode = ModeLexOnly
		modeCount++
	}
	if *parseFlag {
		config.mode = ModeParseOnly
		modeCount++
	}

	if modeCount > 1 {
		fmt.Fprintf(os.Stderr, "Error: Multiple modes specified. Use only one of -l or -p\n")
		os.Exit(1)
	}

	if modeCount == 0 && !config.showHelp && !config.showVersion {
		fmt.Fprintf(os.Stderr, "Error: No mode specified. Use -l for lexing or -p for parsing\n")
		showUsage()
		os.Exit(1)
	}

	// Get input file
	args := flag.Args()
	if len(args) > 0 {
		config.inputFile = args[0]
	}

	return config
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: sangoc [options] <file.sango>\n")
	fmt.Fprintf(os.Stderr, "Try 'sangoc -h' for more information.\n")
}

func showHelp() {
	fmt.Printf(`sangoc %s - Sango Compiler (Lexer/Parser Development Version)

Usage:
  sangoc -l <file.sango>                 Lexical analysis only - show tokens
  sangoc -p <file.sango>                 Parse only - show AST
  sangoc -v                              Show version
  sangoc -h                              Show this help

Options:
  -l    Perform lexical analysis only and display tokens
  -p    Perform parsing only and display AST
  -v    Display version information
  -h    Display this help message

Examples:
  sangoc -l hello.sango                  # Show tokens
  sangoc -p hello.sango                  # Show AST

Note: This is a development version focused on lexer and parser implementation.
Type checking, code generation, and compilation are not yet implemented.

`, VERSION)
}

// Lexical analysis only
func lexOnly(source, filename string) {
	fmt.Printf("=== Lexical Analysis of %s ===\n", filename)

	l := lexer.New(source)

	for {
		tok := l.NextToken()
		if tok.Type == lexer.EOF {
			break
		}
		fmt.Printf("%s\n", tok)
	}
}

// Parse only
func parseOnly(source, filename string) {
	fmt.Printf("=== Parsing %s ===\n", filename)

	l := lexer.New(source)
	p := parser.New(l)

	program := p.ParseProgram()

	// Check for parser errors
	errors := p.Errors()
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Parser errors:\n")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "  %s\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("AST:\n%s\n", program.String())
}
