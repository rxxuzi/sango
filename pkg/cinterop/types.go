package cinterop

// TypeMapping maps C types to Sango types
var TypeMapping = map[string]string{
	// Basic integer types
	"int":           "int",
	"long":          "long",
	"short":         "i16",
	"char":          "i8",
	"unsigned int":  "u32",
	"unsigned long": "u64",
	"unsigned char": "u8",
	"size_t":        "u64",
	
	// Floating point types
	"float":  "float",
	"double": "double",
	
	// Pointer types
	"*void":   "*void",
	"*char":   "*u8",    // C strings
	"*int":    "*int",
	"*float":  "*float",
	"*double": "*double",
	"*FILE":   "*void",  // FILE is typically an opaque pointer
	
	// Special types
	"void": "void",
	"bool": "bool",
}

// ConvertCTypeToSango converts a C type string to its Sango equivalent
func ConvertCTypeToSango(cType string) string {
	if sangoType, ok := TypeMapping[cType]; ok {
		return sangoType
	}
	// If not found, return the original type
	// This allows for user-defined types
	return cType
}