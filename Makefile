# Sango Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=sango
BINARY_DIR=bin

# Source directories
CMD_DIR=./cmd/sango
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

.PHONY: all build test clean fmt vet deps install uninstall runtime example

all: build

# Build the Sango compiler
build: deps runtime
	@echo "Building Sango compiler..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Build complete: $(BINARY_DIR)/$(BINARY_NAME)"

# Build runtime library
runtime:
	@echo "Building runtime library..."
	@mkdir -p $(BINARY_DIR)
	$(CC) $(CFLAGS) -c $(RUNTIME_DIR)/sango.c -o $(BINARY_DIR)/sango.o
	@ar rcs $(BINARY_DIR)/libsango.a $(BINARY_DIR)/sango.o
	@echo "Runtime library built: $(BINARY_DIR)/libsango.a"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v $(PKG_DIR)

# Run lexer test
test-lexer: build
	@echo "Testing lexer with examples/test.sango..."
	./$(BINARY_DIR)/$(BINARY_NAME) examples/test.sango

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

# Install the compiler
install: build
	@echo "Installing Sango to $(BINDIR)..."
	@mkdir -p $(BINDIR)
	@cp $(BINARY_DIR)/$(BINARY_NAME) $(BINDIR)/
	@mkdir -p $(LIBDIR)
	@cp $(BINARY_DIR)/libsango.a $(LIBDIR)/
	@mkdir -p $(INCDIR)
	@cp $(RUNTIME_DIR)/sango.h $(INCDIR)/
	@echo "Installation complete"

# Uninstall the compiler
uninstall:
	@echo "Uninstalling Sango..."
	@rm -f $(BINDIR)/$(BINARY_NAME)
	@rm -rf $(LIBDIR)
	@rm -rf $(INCDIR)
	@echo "Uninstallation complete"

# Development helpers
run: build
	./$(BINARY_DIR)/$(BINARY_NAME)

# Run a specific example
example: build
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make example FILE=examples/hello.sango"; \
	else \
		./$(BINARY_DIR)/$(BINARY_NAME) $(FILE) -o $(BINARY_DIR)/example && \
		echo "Compiled to $(BINARY_DIR)/example"; \
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
	@echo "  make build       - Build the Sango compiler"
	@echo "  make test        - Run all tests"
	@echo "  make test-lexer  - Test lexer with examples/test.sango"
	@echo "  make fmt         - Format Go code"
	@echo "  make vet         - Run go vet"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make install     - Install Sango compiler"
	@echo "  make uninstall   - Uninstall Sango compiler"
	@echo "  make example FILE=<file> - Compile a Sango file"
	@echo "  make quick       - Quick test (fmt, vet, build, test-lexer)"
	@echo "  make help        - Show this help message"