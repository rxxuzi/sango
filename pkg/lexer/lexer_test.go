package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `val x = 5
var y = 10
val result = x + y

def add(x: int, y: int): int = x + y

// This is a comment
val name = "Sango"
val pi = 3.14159

if (x > 0) {
    println("positive")
} else {
    println("negative")
}

val a, b = (10, 20)
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{VAL, "val"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT, "5"},
		{VAR, "var"},
		{IDENT, "y"},
		{ASSIGN, "="},
		{INT, "10"},
		{VAL, "val"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{DEF, "def"},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COLON, ":"},
		{INT_TYPE, "int"},
		{COMMA, ","},
		{IDENT, "y"},
		{COLON, ":"},
		{INT_TYPE, "int"},
		{RPAREN, ")"},
		{COLON, ":"},
		{INT_TYPE, "int"},
		{ASSIGN, "="},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{VAL, "val"},
		{IDENT, "name"},
		{ASSIGN, "="},
		{STRING, "Sango"},
		{VAL, "val"},
		{IDENT, "pi"},
		{ASSIGN, "="},
		{FLOAT, "3.14159"},
		{IF, "if"},
		{LPAREN, "("},
		{IDENT, "x"},
		{GT, ">"},
		{INT, "0"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "println"},
		{LPAREN, "("},
		{STRING, "positive"},
		{RPAREN, ")"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{IDENT, "println"},
		{LPAREN, "("},
		{STRING, "negative"},
		{RPAREN, ")"},
		{RBRACE, "}"},
		{VAL, "val"},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{ASSIGN, "="},
		{LPAREN, "("},
		{INT, "10"},
		{COMMA, ","},
		{INT, "20"},
		{RPAREN, ")"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `+ - * / % **
== != < > <= >=
&& || !
& | ^ ~ << >>
+= -= *= /= %=
&= |= ^= <<= >>=
-> => <- .. ..=
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{PLUS, "+"},
		{MINUS, "-"},
		{ASTERISK, "*"},
		{SLASH, "/"},
		{PERCENT, "%"},
		{POWER, "**"},
		{EQ, "=="},
		{NEQ, "!="},
		{LT, "<"},
		{GT, ">"},
		{LEQ, "<="},
		{GEQ, ">="},
		{AND, "&&"},
		{OR, "||"},
		{NOT, "!"},
		{AMPERSAND, "&"},
		{PIPE, "|"},
		{CARET, "^"},
		{TILDE, "~"},
		{LSHIFT, "<<"},
		{RSHIFT, ">>"},
		{PLUSASSIGN, "+="},
		{MINUSASSIGN, "-="},
		{ASTERISKASSIGN, "*="},
		{SLASHASSIGN, "/="},
		{PERCENTASSIGN, "%="},
		{AMPERSANDASSIGN, "&="},
		{PIPEASSIGN, "|="},
		{CARETASSIGN, "^="},
		{LSHIFTASSIGN, "<<="},
		{RSHIFTASSIGN, ">>="},
		{ARROW, "->"},
		{DARROW, "=>"},
		{LARROW, "<-"},
		{DOTDOT, ".."},
		{DOTDOTEQ, "..="},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestTypes(t *testing.T) {
	input := `int long float double bool string void
i8 i16 i32 i64 u8 u16 u32 u64 f32 f64 byte
[]int [][]float
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{INT_TYPE, "int"},
		{LONG_TYPE, "long"},
		{FLOAT_TYPE, "float"},
		{DOUBLE_TYPE, "double"},
		{BOOL_TYPE, "bool"},
		{STRING_TYPE, "string"},
		{VOID_TYPE, "void"},
		{I8_TYPE, "i8"},
		{I16_TYPE, "i16"},
		{I32_TYPE, "i32"},
		{I64_TYPE, "i64"},
		{U8_TYPE, "u8"},
		{U16_TYPE, "u16"},
		{U32_TYPE, "u32"},
		{U64_TYPE, "u64"},
		{F32_TYPE, "f32"},
		{F64_TYPE, "f64"},
		{BYTE_TYPE, "byte"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{INT_TYPE, "int"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{FLOAT_TYPE, "float"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestKeywords(t *testing.T) {
	input := `def val var if else match type struct impl
return true false for in while break continue
defer sizeof include import define null
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{DEF, "def"},
		{VAL, "val"},
		{VAR, "var"},
		{IF, "if"},
		{ELSE, "else"},
		{MATCH, "match"},
		{TYPE, "type"},
		{STRUCT, "struct"},
		{IMPL, "impl"},
		{RETURN, "return"},
		{TRUE, "true"},
		{FALSE, "false"},
		{FOR, "for"},
		{IN, "in"},
		{WHILE, "while"},
		{BREAK, "break"},
		{CONTINUE, "continue"},
		{DEFER, "defer"},
		{SIZEOF, "sizeof"},
		{INCLUDE, "include"},
		{IMPORT, "import"},
		{DEFINE, "define"},
		{NULL, "null"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
