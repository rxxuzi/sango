package v2

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// BlockParser handles block statement and expression parsing
type BlockParser struct {
	// Block-specific configuration
}

// NewBlockParser creates a new block parser
func NewBlockParser() *BlockParser {
	return &BlockParser{}
}

// ParseBlockStatement parses { statements... } blocks
func (bp *BlockParser) ParseBlockStatement(p ParserInterface, token lexer.Token) *ast.BlockStatement {
	block := &ast.BlockStatement{Token: token}
	block.Statements = []ast.Statement{}

	p.NextToken() // consume '{'

	for !p.CurTokenIs(lexer.RBRACE) && !p.CurTokenIs(lexer.EOF) {
		// Look ahead to detect if this might be the last expression in the block
		if bp.isLastExpressionInBlock(p) {
			// Parse as expression statement for block return value
			expr := p.ParseExpression(LOWEST)
			if expr != nil {
				exprStmt := &ast.ExpressionStatement{
					Token:      p.GetCurrentToken(),
					Expression: expr,
				}
				block.Statements = append(block.Statements, exprStmt)
			}
			// After parsing the last expression, we should be at RBRACE
			break
		}
		
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		
		// Smart token advancement for block statements
		if !p.CurTokenIs(lexer.RBRACE) && !p.CurTokenIs(lexer.EOF) {
			if p.CurTokenIs(lexer.SEMICOLON) {
				p.NextToken() // consume semicolon
			} else if !bp.isStatementStart(p) {
				p.NextToken()
			}
		}
	}

	// DO NOT consume the closing RBRACE here - let the caller handle it
	return block
}

// ParseBlockExpression parses blocks that return values
func (bp *BlockParser) ParseBlockExpression(p ParserInterface, token lexer.Token) ast.Expression {
	// For now, delegate to block statement parsing
	// In the future, we might want different logic for expression blocks
	return bp.ParseBlockStatement(p, token)
}

// isStatementStart checks if current token starts a new statement
func (bp *BlockParser) isStatementStart(p ParserInterface) bool {
	return p.CurTokenIs(lexer.DEF) ||
		   p.CurTokenIs(lexer.VAL) ||
		   p.CurTokenIs(lexer.VAR) ||
		   p.CurTokenIs(lexer.IF) ||
		   p.CurTokenIs(lexer.FOR) ||
		   p.CurTokenIs(lexer.WHILE) ||
		   p.CurTokenIs(lexer.RETURN) ||
		   p.CurTokenIs(lexer.MATCH)
}

// isLastExpressionInBlock detects if current token is likely the last expression in a block
func (bp *BlockParser) isLastExpressionInBlock(p ParserInterface) bool {
	// If current token is an identifier and the next token is RBRACE,
	// this is likely a return expression
	if p.CurTokenIs(lexer.IDENT) && p.PeekTokenIs(lexer.RBRACE) {
		return true
	}
	
	// If current token starts an expression but is not a statement keyword,
	// and there's a RBRACE somewhere ahead, it might be a return expression
	if !bp.isStatementStart(p) && bp.hasRBraceAhead(p) {
		return true
	}
	
	return false
}

// hasRBraceAhead checks if there's an RBRACE token coming up (simple heuristic)
func (bp *BlockParser) hasRBraceAhead(p ParserInterface) bool {
	// Simple check: if peek token is RBRACE, we're likely at the end
	return p.PeekTokenIs(lexer.RBRACE)
}