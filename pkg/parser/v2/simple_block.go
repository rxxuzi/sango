package v2

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// SimpleBlockParser handles block parsing with a straightforward approach
type SimpleBlockParser struct{}

// NewSimpleBlockParser creates a new simple block parser
func NewSimpleBlockParser() *SimpleBlockParser {
	return &SimpleBlockParser{}
}

// ParseBlockStatement parses { statements... } with robust token management
func (sbp *SimpleBlockParser) ParseBlockStatement(p ParserInterface, token lexer.Token) *ast.BlockStatement {
	block := &ast.BlockStatement{Token: token}
	block.Statements = []ast.Statement{}

	p.NextToken() // consume '{'

	for !p.CurTokenIs(lexer.RBRACE) && !p.CurTokenIs(lexer.EOF) {
		// Remember position before parsing
		prevToken := p.GetCurrentToken()
		
		// Parse each statement
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		
		// Token advancement strategy:
		// If statement parsing didn't advance, we need to advance manually
		if prevToken.Type == p.GetCurrentToken().Type && 
		   prevToken.Literal == p.GetCurrentToken().Literal {
			// Statement parsing didn't advance - advance manually
			if !p.CurTokenIs(lexer.RBRACE) && !p.CurTokenIs(lexer.EOF) {
				p.NextToken()
			}
		}
		
		// Skip semicolons if present
		if p.CurTokenIs(lexer.SEMICOLON) {
			p.NextToken()
		}
	}

	return block
}

// isNextStatement checks if we're positioned at the start of a statement
func (sbp *SimpleBlockParser) isNextStatement(p ParserInterface) bool {
	return p.CurTokenIs(lexer.DEF) ||
		   p.CurTokenIs(lexer.VAL) ||
		   p.CurTokenIs(lexer.VAR) ||
		   p.CurTokenIs(lexer.IF) ||
		   p.CurTokenIs(lexer.FOR) ||
		   p.CurTokenIs(lexer.WHILE) ||
		   p.CurTokenIs(lexer.RETURN) ||
		   p.CurTokenIs(lexer.MATCH) ||
		   p.CurTokenIs(lexer.LBRACE) ||
		   p.CurTokenIs(lexer.IDENT) // Could be start of expression statement
}