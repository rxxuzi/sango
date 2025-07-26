package lexer

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	NEWLINE

	// Identifiers and literals
	IDENT  // x, add, foo
	INT    // 123
	FLOAT  // 123.45
	STRING // "hello"

	// Operators
	PLUS       // +
	MINUS      // -
	ASTERISK   // *
	SLASH      // /
	PERCENT    // %
	POWER      // **
	EQ         // ==
	NEQ        // !=
	LT         // <
	GT         // >
	LEQ        // <=
	GEQ        // >=
	ASSIGN     // =
	PLUSASSIGN // +=
	MINUSASSIGN // -=
	ASTERISKASSIGN // *=
	SLASHASSIGN // /=
	PERCENTASSIGN // %=
	AMPERSAND  // &
	PIPE       // |
	CARET      // ^
	TILDE      // ~
	LSHIFT     // <<
	RSHIFT     // >>
	AMPERSANDASSIGN // &=
	PIPEASSIGN // |=
	CARETASSIGN // ^=
	LSHIFTASSIGN // <<=
	RSHIFTASSIGN // >>=
	AND        // &&
	OR         // ||
	NOT        // !
	ARROW      // ->
	DARROW     // =>
	LARROW     // <-
	DOTDOT     // ..
	DOTDOTEQ   // ..=
	AT         // @

	// Delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .
	UNDERSCORE // _

	// Keywords
	DEF      // def
	VAL      // val
	VAR      // var
	IF       // if
	ELSE     // else
	MATCH    // match
	TYPE     // type
	STRUCT   // struct
	IMPL     // impl
	RETURN   // return
	TRUE     // true
	FALSE    // false
	FOR      // for
	IN       // in
	WHILE    // while
	BREAK    // break
	CONTINUE // continue
	DEFER    // defer
	SIZEOF   // sizeof
	INCLUDE  // include
	IMPORT   // import
	DEFINE   // define
	NULL     // null

	// Basic types
	INT_TYPE    // int
	LONG_TYPE   // long
	FLOAT_TYPE  // float
	DOUBLE_TYPE // double
	BOOL_TYPE   // bool
	STRING_TYPE // string
	VOID_TYPE   // void
	
	// Detailed types
	I8_TYPE   // i8
	I16_TYPE  // i16
	I32_TYPE  // i32
	I64_TYPE  // i64
	U8_TYPE   // u8
	U16_TYPE  // u16
	U32_TYPE  // u32
	U64_TYPE  // u64
	F32_TYPE  // f32
	F64_TYPE  // f64
	BYTE_TYPE // byte
)

var tokenStrings = map[TokenType]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	NEWLINE: "NEWLINE",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	PLUS:           "+",
	MINUS:          "-",
	ASTERISK:       "*",
	SLASH:          "/",
	PERCENT:        "%",
	POWER:          "**",
	EQ:             "==",
	NEQ:            "!=",
	LT:             "<",
	GT:             ">",
	LEQ:            "<=",
	GEQ:            ">=",
	ASSIGN:         "=",
	PLUSASSIGN:     "+=",
	MINUSASSIGN:    "-=",
	ASTERISKASSIGN: "*=",
	SLASHASSIGN:    "/=",
	PERCENTASSIGN:  "%=",
	AMPERSAND:      "&",
	PIPE:           "|",
	CARET:          "^",
	TILDE:          "~",
	LSHIFT:         "<<",
	RSHIFT:         ">>",
	AMPERSANDASSIGN: "&=",
	PIPEASSIGN:     "|=",
	CARETASSIGN:    "^=",
	LSHIFTASSIGN:   "<<=",
	RSHIFTASSIGN:   ">>=",
	AND:            "&&",
	OR:             "||",
	NOT:            "!",
	ARROW:          "->",
	DARROW:         "=>",
	LARROW:         "<-",
	DOTDOT:         "..",
	DOTDOTEQ:       "..=",
	AT:             "@",

	LPAREN:     "(",
	RPAREN:     ")",
	LBRACE:     "{",
	RBRACE:     "}",
	LBRACKET:   "[",
	RBRACKET:   "]",
	COMMA:      ",",
	SEMICOLON:  ";",
	COLON:      ":",
	DOT:        ".",
	UNDERSCORE: "_",

	DEF:      "def",
	VAL:      "val",
	VAR:      "var",
	IF:       "if",
	ELSE:     "else",
	MATCH:    "match",
	TYPE:     "type",
	STRUCT:   "struct",
	IMPL:     "impl",
	RETURN:   "return",
	TRUE:     "true",
	FALSE:    "false",
	FOR:      "for",
	IN:       "in",
	WHILE:    "while",
	BREAK:    "break",
	CONTINUE: "continue",
	DEFER:    "defer",
	SIZEOF:   "sizeof",
	INCLUDE:  "include",
	IMPORT:   "import",
	DEFINE:   "define",
	NULL:     "null",

	INT_TYPE:    "int",
	LONG_TYPE:   "long",
	FLOAT_TYPE:  "float",
	DOUBLE_TYPE: "double",
	BOOL_TYPE:   "bool",
	STRING_TYPE: "string",
	VOID_TYPE:   "void",

	I8_TYPE:   "i8",
	I16_TYPE:  "i16",
	I32_TYPE:  "i32",
	I64_TYPE:  "i64",
	U8_TYPE:   "u8",
	U16_TYPE:  "u16",
	U32_TYPE:  "u32",
	U64_TYPE:  "u64",
	F32_TYPE:  "f32",
	F64_TYPE:  "f64",
	BYTE_TYPE: "byte",
}

func (t TokenType) String() string {
	if s, ok := tokenStrings[t]; ok {
		return s
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

var keywords = map[string]TokenType{
	"def":      DEF,
	"val":      VAL,
	"var":      VAR,
	"if":       IF,
	"else":     ELSE,
	"match":    MATCH,
	"type":     TYPE,
	"struct":   STRUCT,
	"impl":     IMPL,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"for":      FOR,
	"in":       IN,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
	"defer":    DEFER,
	"sizeof":   SIZEOF,
	"include":  INCLUDE,
	"import":   IMPORT,
	"define":   DEFINE,
	"null":     NULL,
	
	// Basic types
	"int":    INT_TYPE,
	"long":   LONG_TYPE,
	"float":  FLOAT_TYPE,
	"double": DOUBLE_TYPE,
	"bool":   BOOL_TYPE,
	"string": STRING_TYPE,
	"void":   VOID_TYPE,
	
	// Detailed types
	"i8":     I8_TYPE,
	"i16":    I16_TYPE,
	"i32":    I32_TYPE,
	"i64":    I64_TYPE,
	"u8":     U8_TYPE,
	"u16":    U16_TYPE,
	"u32":    U32_TYPE,
	"u64":    U64_TYPE,
	"f32":    F32_TYPE,
	"f64":    F64_TYPE,
	"byte":   BYTE_TYPE,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// NewToken creates a new token
func NewToken(tokenType TokenType, literal string, line, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("{%s %q %d:%d}", t.Type, t.Literal, t.Line, t.Column)
}