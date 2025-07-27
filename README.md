# Sango

A minimal, extensible functional programming language that compiles to C.

## Features

- **Functional**: Immutable by default, first-class functions, pattern matching
- **Statically Typed**: With type inference  
- **C Interop**: Direct use of C libraries and macros
- **Minimal**: Small core language with extensibility
- **Flexible Syntax**: Optional semicolons, clean syntax inspired by C and Scala

## Quick Start

```bash
git clone https://github.com/rxxuzi/sango.git
cd sango
make build
./bin/sangoc examples/simple.sango
```

## Compilation Pipeline

```
.sango → [Lexer] → [Parser] → [Type Checker] → [Code Gen] → .c → [gcc] → executable
```

## Development

```bash
sangoc -l file.sango    # Tokenize
sangoc -p file.sango    # Parse AST  
sangoc file.sango       # Compile to binary
```

## Status

Currently implementing parser. Lexer complete, type checker and code generator planned.

## License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or distribute this software, either in source code form or as a compiled binary, for any purpose, commercial or non-commercial, and by any means.

For more information, please refer to http://unlicense.org/