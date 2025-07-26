# Sango

A minimal, extensible functional programming language that compiles to C.

## Features

- **Functional**: Immutable by default, first-class functions, pattern matching
- **Statically Typed**: With type inference
- **C Interop**: Direct use of C libraries and macros
- **Minimal**: Small core language with extensibility
- **Flexible Syntax**: Optional semicolons, clean syntax inspired by C and Scala

## Quick Start

```sango
// Hello World
def main() : int = {
    println("Hello, Sango!")
    0
}

// Factorial with pattern matching
def factorial(n: int) : int = match n {
    0 => 1
    n => n * factorial(n - 1)
}

// Using C libraries
include "math.h"
def distance(x: float, y: float) : float = 
    sqrt(x * x + y * y)
```

## Installation

```bash
# Clone the repository
git clone https://github.com/rxxuzi/sango.git
cd sango

# Build the compilers
make

# Use the Sango compiler (sangoc)
./bin/sangoc examples/simple.sango      # Compile to executable 
./bin/sangoc -h                          # Show help
./bin/sangoc -v                          # Show version

# Development and testing
./bin/sangoc -l examples/simple.sango   # Show tokens
./bin/sangoc -p examples/simple.sango   # Show AST
```

## Compilation Pipeline

```
.sango → [Lexer] → [Parser] → [Type Checker] → [Code Gen] → .c → [gcc] → executable
```

## Language Guide

### Basic Types

- `int`: 64-bit integer
- `float`: 64-bit floating point
- `bool`: Boolean
- `string`: UTF-8 string
- `void`: Unit type (like void)

### Function Definition

```sango
// Basic function
def add(x: int, y: int) : int = x + y

// With type inference (void return)
def identity(x: int) = x

// Multiple statements
def compute(x: int) : int = {
    val y = x * 2
    val z = y + 1
    z
}
```

### Pattern Matching

```sango
type Option[T] = 
    | Some(T)
    | None

def unwrap(opt: Option[int], default: int) : int = match opt {
    Some(x) => x
    None => default
}
```

### C Interoperability

```sango
// Include C headers
include "stdio.h"
include "stdlib.h"

// Use C macros
define BUFFER_SIZE 1024

// Call C functions (future feature)
extern def malloc(size: int) : ptr[void]
extern def free(ptr: ptr[void])

// Use C types (future feature)
type FILE = extern struct
```

## Building from Source

### Requirements

- Go 1.21 or later
- GCC or Clang
- Make

### Build Steps

```bash
make build         # Build both sango and sangoc
make build-sangoc  # Build compiler only
make test          # Run tests
make install       # Install to /usr/local/bin
```

## sangoc Compiler Usage

```bash
# Full compilation (requires main function)
sangoc program.sango                    # Compile to 'program'
sangoc -o myapp program.sango          # Compile to 'myapp'

# Development and debugging
sangoc -l program.sango                # Lexical analysis - show tokens
sangoc -p program.sango                # Parse only - show AST
sangoc -t program.sango                # Type check only (TODO)
sangoc -e program.sango                # Emit C code (TODO)

# Version and help
sangoc -v                              # Show version (v0.1.2)
sangoc -h                              # Show help
```

## License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

For more information, please refer to <http://unlicense.org/>