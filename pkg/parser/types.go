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
			p.nextToken() // move to element type
			type_expr.Array = true
			// Recursively parse element type (could be another array)
			elementType := p.parseTypeExpression()
			if elementType != nil {
				type_expr.Name = elementType.String()
				type_expr.ElementType = elementType
			}
		} else {
			// Fixed size array [N]T - for now treat as dynamic
			for !p.curTokenIs(lexer.RBRACKET) {
				p.nextToken()
			}
			p.nextToken() // consume ]
			p.nextToken() // move to element type
			type_expr.Array = true
			// Recursively parse element type
			elementType := p.parseTypeExpression()
			if elementType != nil {
				type_expr.Name = elementType.String()
				type_expr.ElementType = elementType
			}
		}
		return type_expr
	}

	// Handle parenthesized types - could be tuple or function parameters
	if p.curTokenIs(lexer.LPAREN) {
		return p.parseParenthesizedType()
	}
	
	// Handle record types { field: type, ... }
	if p.curTokenIs(lexer.LBRACE) {
		return p.parseRecordType()
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

// parseParenthesizedType handles both tuple types and function parameter lists
func (p *Parser) parseParenthesizedType() *ast.TypeExpression {
	type_expr := &ast.TypeExpression{Token: p.curToken}

	p.nextToken() // consume '('

	// Empty parentheses
	if p.curTokenIs(lexer.RPAREN) {
		p.nextToken()
		// Check if this is a function type: () -> ReturnType
		if p.curTokenIs(lexer.ARROW) {
			return p.parseFunctionType([]ast.TypeExpression{})
		}
		// Empty tuple
		return type_expr
	}

	// Parse first type
	var types []ast.TypeExpression
	firstType := p.parseTypeExpression()
	types = append(types, *firstType)

	// Parse remaining types
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume ','
		p.nextToken() // move to next type
		nextType := p.parseTypeExpression()
		types = append(types, *nextType)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	// Check what comes after parentheses
	if p.peekTokenIs(lexer.ARROW) {
		// This is a function type: (A, B) -> C
		p.nextToken() // consume '->'
		return p.parseFunctionType(types)
	}

	// This is a tuple type: (A, B, C)
	type_expr.Tuple = types
	return type_expr
}

// parseFunctionType parses function types: paramTypes -> returnType
func (p *Parser) parseFunctionType(paramTypes []ast.TypeExpression) *ast.TypeExpression {
	type_expr := &ast.TypeExpression{Token: p.curToken}

	p.nextToken() // move to return type
	returnType := p.parseTypeExpression()

	funcType := &ast.FunctionType{
		Parameters: paramTypes,
		ReturnType: returnType,
	}

	type_expr.Function = funcType
	return type_expr
}

// parseRecordType parses record types: { field: type, ... }
func (p *Parser) parseRecordType() *ast.TypeExpression {
	type_expr := &ast.TypeExpression{Token: p.curToken}
	type_expr.Record = &ast.RecordType{Fields: make(map[string]*ast.TypeExpression)}
	
	p.nextToken() // consume '{'
	
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		// Parse field name
		if !p.curTokenIs(lexer.IDENT) {
			p.addError("expected field name in record type")
			return nil
		}
		fieldName := p.curToken.Literal
		
		// Expect ':'
		if !p.expectPeek(lexer.COLON) {
			return nil
		}
		
		// Parse field type
		p.nextToken() // move to type
		fieldType := p.parseTypeExpression()
		if fieldType == nil {
			return nil
		}
		
		type_expr.Record.Fields[fieldName] = fieldType
		
		// Check for comma or end
		p.nextToken()
		if p.curTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma and continue
		} else if !p.curTokenIs(lexer.RBRACE) {
			p.addError("expected ',' or '}' in record type")
			return nil
		}
	}
	
	// Consume closing '}'
	if !p.curTokenIs(lexer.RBRACE) {
		p.addError("expected '}' to close record type")
		return nil
	}
	
	return type_expr
}
