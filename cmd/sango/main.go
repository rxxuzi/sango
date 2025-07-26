package main

import (
	"fmt"
	"os"
)

const VERSION = "v0.1.1"

func main() {
	fmt.Printf("Sango REPL/Interpreter %s\n", VERSION)
	fmt.Printf("Interactive Sango development environment (under development)\n")
	fmt.Printf("For compilation, use 'sangoc' instead.\n")
	fmt.Printf("\nUsage:\n")
	fmt.Printf("  sangoc <file.sango>     # Compile Sango programs\n")
	fmt.Printf("  sangoc -h               # Show compiler help\n")
	
	if len(os.Args) > 1 {
		fmt.Printf("\nNote: File compilation has moved to 'sangoc'.\n")
		fmt.Printf("Try: sangoc %s\n", os.Args[1])
	}
}