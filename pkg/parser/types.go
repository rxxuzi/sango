package parser

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// Type expression parsing
func (p *Parser) parseTypeExpression() *ast.TypeExpression {
	type_expr := &ast.TypeExpression{Token: p.curToken}

	// Handle array types []
	if p.curTokenIs(lexer.LBRACKET) {
		if p.peekTokenIs(lexer.RBRACKET) {
			// Dynamic array []T
			p.nextToken() // consume ]
			p.nextToken() // move to type
			type_expr.Array = true
			type_expr.Name = p.curToken.Literal
		} else {
			// Fixed size array [N]T - for now treat as dynamic
			for !p.curTokenIs(lexer.RBRACKET) {
				p.nextToken()
			}
			p.nextToken() // consume ]
			p.nextToken() // move to type
			type_expr.Array = true
			type_expr.Name = p.curToken.Literal
		}
		return type_expr
	}

	// Handle tuple types (A, B, C)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()
		for !p.curTokenIs(lexer.RPAREN) && !p.curTokenIs(lexer.EOF) {
			tupleType := p.parseTypeExpression()
			type_expr.Tuple = append(type_expr.Tuple, *tupleType)
			if p.peekTokenIs(lexer.COMMA) {
				p.nextToken()
			}
			p.nextToken()
		}
		return type_expr
	}

	// Handle function types (A, B) -> C
	if p.peekTokenIs(lexer.ARROW) {
		funcType := &ast.FunctionType{}
		// Simple case: single parameter type
		funcType.Parameters = []ast.TypeExpression{{Token: p.curToken, Name: p.curToken.Literal}}
		p.nextToken() // consume ->
		p.nextToken() // move to return type
		funcType.ReturnType = p.parseTypeExpression()
		type_expr.Function = funcType
		return type_expr
	}

	// Simple type name
	type_expr.Name = p.curToken.Literal
	return type_expr
}

// Function parameter parsing
func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	identifiers := []*ast.Parameter{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// Parse first parameter
	param := &ast.Parameter{}
	param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken()
		p.nextToken()
		param.Type = p.parseTypeExpression()
	}

	identifiers = append(identifiers, param)

	// Parse remaining parameters
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()

		param := &ast.Parameter{}
		param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		// Optional type annotation
		if p.peekTokenIs(lexer.COLON) {
			p.nextToken()
			p.nextToken()
			param.Type = p.parseTypeExpression()
		}

		identifiers = append(identifiers, param)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return identifiers
}