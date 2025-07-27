package parser

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// Fixed block parsing that properly handles token advancement

func (p *Parser) parseBlockStatementFixed() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken() // consume '{'

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		// Skip any semicolons
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
			continue
		}
		
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		
		// Consistent token advancement - but be smarter about it
		if !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
			// If we're at a statement start token, don't advance
			if p.isStatementStartToken() {
				continue
			}
			p.nextToken()
		}
	}

	// Consume the closing RBRACE
	if p.curTokenIs(lexer.RBRACE) {
		p.nextToken()
	}

	return block
}

// isStatementStartToken checks if current token can start a statement
func (p *Parser) isStatementStartToken() bool {
	return p.curTokenIs(lexer.DEF) ||
		   p.curTokenIs(lexer.VAL) ||
		   p.curTokenIs(lexer.VAR) ||
		   p.curTokenIs(lexer.IF) ||
		   p.curTokenIs(lexer.FOR) ||
		   p.curTokenIs(lexer.WHILE) ||
		   p.curTokenIs(lexer.RETURN) ||
		   p.curTokenIs(lexer.MATCH) ||
		   p.curTokenIs(lexer.DEFER) ||
		   p.curTokenIs(lexer.ASSERT) ||
		   (p.curTokenIs(lexer.IDENT) && p.isAssignmentOperator(p.peekToken.Type))
}