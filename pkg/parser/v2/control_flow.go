package v2

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// ControlFlowParser handles control flow statements (for, while, if, match)
type ControlFlowParser struct {
	// Control flow specific configuration
}

// NewControlFlowParser creates a new control flow parser
func NewControlFlowParser() *ControlFlowParser {
	return &ControlFlowParser{}
}

// ParseForStatement parses for loops with both <- and in syntax
func (cf *ControlFlowParser) ParseForStatement(p ParserInterface, token lexer.Token) ast.Statement {
	stmt := &ast.ForStatement{Token: token}
	
	if !p.ExpectPeek(lexer.IDENT) {
		return nil
	}
	
	stmt.Variable = &ast.Identifier{Token: p.GetCurrentToken(), Value: p.GetCurrentToken().Literal}
	
	// Check for <- or 'in'
	if p.PeekTokenIs(lexer.LARROW) {
		// for x <- iterable
		p.NextToken()
		stmt.IsInRange = false
	} else if p.PeekTokenIs(lexer.IN) {
		// for i in range
		p.NextToken()
		stmt.IsInRange = true
	} else {
		p.AddError("expected '<-' or 'in' after for variable")
		return nil
	}
	
	p.NextToken()
	stmt.Iterable = p.ParseExpression(LOWEST)
	
	if !p.ExpectPeek(lexer.LBRACE) {
		return nil
	}
	
	// Use block parser for the body
	blockParser := NewBlockParser()
	stmt.Body = blockParser.ParseBlockStatement(p, p.GetCurrentToken())
	
	// Consume the closing RBRACE
	if p.CurTokenIs(lexer.RBRACE) {
		p.NextToken()
	}
	
	return stmt
}

// ParseWhileStatement parses while loops
func (cf *ControlFlowParser) ParseWhileStatement(p ParserInterface, token lexer.Token) ast.Statement {
	stmt := &ast.WhileStatement{Token: token}
	
	if !p.ExpectPeek(lexer.LPAREN) {
		return nil
	}
	
	p.NextToken()
	stmt.Condition = p.ParseExpression(LOWEST)
	
	if !p.ExpectPeek(lexer.RPAREN) {
		return nil
	}
	
	if !p.ExpectPeek(lexer.LBRACE) {
		return nil
	}
	
	// Use block parser for the body
	blockParser := NewBlockParser()
	stmt.Body = blockParser.ParseBlockStatement(p, p.GetCurrentToken())
	
	// Consume the closing RBRACE
	if p.CurTokenIs(lexer.RBRACE) {
		p.NextToken()
	}
	
	return stmt
}

// ParseIfExpression parses if expressions
func (cf *ControlFlowParser) ParseIfExpression(p ParserInterface, token lexer.Token) ast.Expression {
	expression := &ast.IfExpression{Token: token}
	
	if !p.ExpectPeek(lexer.LPAREN) {
		return nil
	}
	
	p.NextToken()
	expression.Condition = p.ParseExpression(LOWEST)
	
	if !p.ExpectPeek(lexer.RPAREN) {
		return nil
	}
	
	if !p.ExpectPeek(lexer.LBRACE) {
		return nil
	}
	
	// Use block parser for consequence
	blockParser := NewBlockParser()
	expression.Consequence = blockParser.ParseBlockStatement(p, p.GetCurrentToken())
	
	// Consume the closing RBRACE
	if p.CurTokenIs(lexer.RBRACE) {
		p.NextToken()
	}
	
	// Handle optional else clause
	if p.PeekTokenIs(lexer.ELSE) {
		p.NextToken()
		
		if !p.ExpectPeek(lexer.LBRACE) {
			return nil
		}
		
		expression.Alternative = blockParser.ParseBlockStatement(p, p.GetCurrentToken())
		
		// Consume the closing RBRACE
		if p.CurTokenIs(lexer.RBRACE) {
			p.NextToken()
		}
	}
	
	return expression
}