package v2

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// Precedence levels for v2 parsers
type Precedence int

const (
	_ Precedence = iota
	LOWEST
	ASSIGN      // = += -= *= /= %= &= |= ^= <<= >>=
	OR          // ||
	AND         // &&
	BITOR       // |
	BITXOR      // ^
	BITAND      // &
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SHIFT       // << >>
	SUM         // + -
	PRODUCT     // * / %
	POWER       // **
	PREFIX      // -x !x ~x
	POSTFIX     // x++ x--
	CALL        // myFunction(X)
	INDEX       // array[index]
	DOT         // obj.field
)

// ParserInterface defines the interface that v2 parsers expect from the main parser
type ParserInterface interface {
	// Token management
	CurTokenIs(lexer.TokenType) bool
	PeekTokenIs(lexer.TokenType) bool
	NextToken()
	ExpectPeek(lexer.TokenType) bool
	GetCurrentToken() lexer.Token
	GetPeekToken() lexer.Token
	
	// Expression parsing
	ParseExpression(Precedence) ast.Expression
	ParseStatement() ast.Statement
	
	// Error handling
	AddError(string)
}