package cinterop

import (
	"sync"
)

// FunctionRegistry manages C functions available in the current compilation context
type FunctionRegistry struct {
	mu        sync.RWMutex
	functions map[string]FunctionSignature // function name -> signature
	headers   map[string]bool             // track which headers have been included
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]FunctionSignature),
		headers:   make(map[string]bool),
	}
}

// IncludeHeader adds all functions from a header to the registry
func (r *FunctionRegistry) IncludeHeader(header string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Skip if already included
	if r.headers[header] {
		return
	}
	
	r.headers[header] = true
	
	// Add functions from the header
	if funcs := GetFunctionsForHeader(header); funcs != nil {
		for _, fn := range funcs {
			r.functions[fn.Name] = fn
		}
	}
}

// RegisterFunction manually registers a function
func (r *FunctionRegistry) RegisterFunction(fn FunctionSignature) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.functions[fn.Name] = fn
}

// LookupFunction checks if a function is available
func (r *FunctionRegistry) LookupFunction(name string) (FunctionSignature, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, ok := r.functions[name]
	return fn, ok
}

// IsFunction checks if a name is a registered C function
func (r *FunctionRegistry) IsFunction(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.functions[name]
	return ok
}

// GetAllFunctions returns all registered functions
func (r *FunctionRegistry) GetAllFunctions() map[string]FunctionSignature {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Create a copy to avoid concurrent modification
	result := make(map[string]FunctionSignature)
	for k, v := range r.functions {
		result[k] = v
	}
	return result
}