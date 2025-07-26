package cinterop

// StandardCFunctions defines commonly used C standard library functions
// These are automatically available when including the corresponding headers
var StandardCFunctions = map[string][]FunctionSignature{
	"stdio.h": {
		{Name: "printf", ReturnType: "int", Variadic: true},
		{Name: "fprintf", ReturnType: "int", Variadic: true},
		{Name: "sprintf", ReturnType: "int", Variadic: true},
		{Name: "scanf", ReturnType: "int", Variadic: true},
		{Name: "fopen", ReturnType: "*FILE", Args: []Argument{{Type: "*char"}, {Type: "*char"}}},
		{Name: "fclose", ReturnType: "int", Args: []Argument{{Type: "*FILE"}}},
		{Name: "fread", ReturnType: "size_t", Args: []Argument{{Type: "*void"}, {Type: "size_t"}, {Type: "size_t"}, {Type: "*FILE"}}},
		{Name: "fwrite", ReturnType: "size_t", Args: []Argument{{Type: "*void"}, {Type: "size_t"}, {Type: "size_t"}, {Type: "*FILE"}}},
		{Name: "puts", ReturnType: "int", Args: []Argument{{Type: "*char"}}},
		{Name: "getchar", ReturnType: "int"},
		{Name: "putchar", ReturnType: "int", Args: []Argument{{Type: "int"}}},
	},
	"stdlib.h": {
		{Name: "malloc", ReturnType: "*void", Args: []Argument{{Type: "size_t"}}},
		{Name: "free", ReturnType: "void", Args: []Argument{{Type: "*void"}}},
		{Name: "realloc", ReturnType: "*void", Args: []Argument{{Type: "*void"}, {Type: "size_t"}}},
		{Name: "exit", ReturnType: "void", Args: []Argument{{Type: "int"}}},
		{Name: "atoi", ReturnType: "int", Args: []Argument{{Type: "*char"}}},
		{Name: "atof", ReturnType: "double", Args: []Argument{{Type: "*char"}}},
		{Name: "rand", ReturnType: "int"},
		{Name: "srand", ReturnType: "void", Args: []Argument{{Type: "unsigned int"}}},
	},
	"math.h": {
		{Name: "sqrt", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "pow", ReturnType: "double", Args: []Argument{{Type: "double"}, {Type: "double"}}},
		{Name: "sin", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "cos", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "tan", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "ceil", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "floor", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "abs", ReturnType: "int", Args: []Argument{{Type: "int"}}},
		{Name: "fabs", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "log", ReturnType: "double", Args: []Argument{{Type: "double"}}},
		{Name: "exp", ReturnType: "double", Args: []Argument{{Type: "double"}}},
	},
	"string.h": {
		{Name: "strlen", ReturnType: "size_t", Args: []Argument{{Type: "*char"}}},
		{Name: "strcpy", ReturnType: "*char", Args: []Argument{{Type: "*char"}, {Type: "*char"}}},
		{Name: "strncpy", ReturnType: "*char", Args: []Argument{{Type: "*char"}, {Type: "*char"}, {Type: "size_t"}}},
		{Name: "strcat", ReturnType: "*char", Args: []Argument{{Type: "*char"}, {Type: "*char"}}},
		{Name: "strcmp", ReturnType: "int", Args: []Argument{{Type: "*char"}, {Type: "*char"}}},
		{Name: "strchr", ReturnType: "*char", Args: []Argument{{Type: "*char"}, {Type: "int"}}},
		{Name: "strstr", ReturnType: "*char", Args: []Argument{{Type: "*char"}, {Type: "*char"}}},
		{Name: "memcpy", ReturnType: "*void", Args: []Argument{{Type: "*void"}, {Type: "*void"}, {Type: "size_t"}}},
		{Name: "memset", ReturnType: "*void", Args: []Argument{{Type: "*void"}, {Type: "int"}, {Type: "size_t"}}},
	},
}

// FunctionSignature represents a C function signature
type FunctionSignature struct {
	Name       string
	ReturnType string
	Args       []Argument
	Variadic   bool // for functions like printf that accept variable arguments
}

// Argument represents a function argument
type Argument struct {
	Name string // optional
	Type string
}

// GetFunctionsForHeader returns all function signatures for a given header
func GetFunctionsForHeader(header string) []FunctionSignature {
	if funcs, ok := StandardCFunctions[header]; ok {
		return funcs
	}
	return nil
}