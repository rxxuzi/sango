# Sango Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary names
SANGO_BINARY=sango
SANGOC_BINARY=sangoc
BINARY_DIR=bin

# Source directories
SANGO_CMD_DIR=./cmd/sango
SANGOC_CMD_DIR=./cmd/sangoc
PKG_DIR=./pkg/...

# C compiler parameters
CC=gcc
CFLAGS=-O2 -Wall -Wextra -std=c11
RUNTIME_DIR=runtime

# Installation directory
PREFIX=/usr/local
BINDIR=$(PREFIX)/bin
LIBDIR=$(PREFIX)/lib/sango
INCDIR=$(PREFIX)/include/sango

.PHONY: all build build-sango build-sangoc test clean fmt vet deps install uninstall runtime example

all: build

# Build both compilers
build: deps runtime build-sango build-sangoc

# Build the Sango REPL/interpreter
build-sango: deps runtime
	@echo "Building Sango REPL/interpreter..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(SANGO_BINARY) $(SANGO_CMD_DIR)
	@echo "Build complete: $(BINARY_DIR)/$(SANGO_BINARY)"

# Build the Sango compiler
build-sangoc: deps runtime
	@echo "Building Sango compiler (sangoc)..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(SANGOC_BINARY) $(SANGOC_CMD_DIR)
	@echo "Build complete: $(BINARY_DIR)/$(SANGOC_BINARY)"

# Build runtime library (disabled for lexer/parser development)
runtime:
	@echo "Runtime library disabled for lexer/parser development"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v $(PKG_DIR)

# Run lexer test
test-lexer: build-sangoc
	@echo "Testing lexer with examples/test.sango..."
	./$(BINARY_DIR)/$(SANGOC_BINARY) -l examples/test.sango

# Test parser
test-parser: build-sangoc
	@echo "Testing parser with examples/test.sango..."
	./$(BINARY_DIR)/$(SANGOC_BINARY) -p examples/test.sango

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f *.o *.a
	@echo "Clean complete"

# Install the compilers
install: build
	@echo "Installing Sango to $(BINDIR)..."
	@mkdir -p $(BINDIR)
	@cp $(BINARY_DIR)/$(SANGO_BINARY) $(BINDIR)/
	@cp $(BINARY_DIR)/$(SANGOC_BINARY) $(BINDIR)/
	@mkdir -p $(LIBDIR)
	@cp $(BINARY_DIR)/libsango.a $(LIBDIR)/
	@mkdir -p $(INCDIR)
	@cp $(RUNTIME_DIR)/sango.h $(INCDIR)/
	@echo "Installation complete"

# Uninstall the compilers
uninstall:
	@echo "Uninstalling Sango..."
	@rm -f $(BINDIR)/$(SANGO_BINARY)
	@rm -f $(BINDIR)/$(SANGOC_BINARY)
	@rm -rf $(LIBDIR)
	@rm -rf $(INCDIR)
	@echo "Uninstallation complete"

# Development helpers
run-sango: build-sango
	./$(BINARY_DIR)/$(SANGO_BINARY)

run-sangoc: build-sangoc
	./$(BINARY_DIR)/$(SANGOC_BINARY) -h

# Run a specific example
example: build-sangoc
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make example FILE=examples/hello.sango"; \
	else \
		./$(BINARY_DIR)/$(SANGOC_BINARY) $(FILE) && \
		echo "Compiled $(FILE)"; \
	fi

# Generate documentation
docs:
	@echo "Generating documentation..."
	@godoc -http=:6060

# Quick test for development
quick: fmt vet build test-lexer

# Initialize project
init:
	@echo "Initializing Sango project..."
	@mkdir -p cmd/sango pkg/lexer pkg/parser pkg/ast pkg/semantic pkg/codegen
	@mkdir -p runtime stdlib examples tests docs
	@echo "Project structure created"

# Show help
help:
	@echo "Sango Makefile commands:"
	@echo "  make build         - Build both sango and sangoc"
	@echo "  make build-sango   - Build sango REPL/interpreter only"
	@echo "  make build-sangoc  - Build sangoc compiler only"
	@echo "  make test          - Run all tests"
	@echo "  make test-lexer    - Test lexer with examples/test.sango"
	@echo "  make test-parser   - Test parser with examples/test.sango"
	@echo "  make fmt           - Format Go code"
	@echo "  make vet           - Run go vet"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install       - Install both binaries"
	@echo "  make uninstall     - Uninstall both binaries"
	@echo "  make example FILE=<file> - Compile a Sango file with sangoc"
	@echo "  make run-sangoc    - Show sangoc help"
	@echo "  make quick         - Quick test (fmt, vet, build, test-lexer)"
	@echo "  make help          - Show this help message"