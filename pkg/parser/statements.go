package parser

import (
	"bytes"

	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// Statement parsing
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.VAL:
		return p.parseValStatement()
	case lexer.VAR:
		return p.parseVarStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.DEF:
		return p.parseFunctionStatement()
	case lexer.TYPE:
		return p.parseTypeStatement()
	case lexer.STRUCT:
		return p.parseStructStatement()
	case lexer.IMPL:
		return p.parseImplStatement()
	case lexer.INCLUDE:
		return p.parseIncludeStatement()
	case lexer.IMPORT:
		return p.parseImportStatement()
	case lexer.DEFINE:
		return p.parseDefineStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.DEFER:
		return p.parseDeferStatement()
	case lexer.ASSERT:
		return p.parseAssertStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseValStatement() *ast.ValStatement {
	stmt := &ast.ValStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	// Parse multiple identifiers for tuple destructuring
	stmt.Names = []*ast.Identifier{}
	stmt.Names = append(stmt.Names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	// Check for tuple destructuring (val a, b, c = ...)
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.Names = append(stmt.Names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	// Check for optional type annotation (val x: Type = ...)
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume ':'
		p.nextToken() // move to type
		stmt.Type = p.parseTypeExpression()
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	// Parse multiple identifiers for tuple destructuring
	stmt.Names = []*ast.Identifier{}
	stmt.Names = append(stmt.Names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	// Check for tuple destructuring (var a, b, c = ...)
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.Names = append(stmt.Names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	// Check for optional type annotation (var x: Type = ...)
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume ':'
		p.nextToken() // move to type
		stmt.Type = p.parseTypeExpression()
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	// Check if this is an assignment statement
	if p.curTokenIs(lexer.IDENT) && p.isAssignmentOperator(p.peekToken.Type) {
		return p.parseAssignmentStatement()
	}

	// Otherwise, it's a regular expression statement
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	// After expression parsing, we should be positioned correctly
	// Skip semicolon if present
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) isAssignmentOperator(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.ASSIGN, lexer.PLUSASSIGN, lexer.MINUSASSIGN, lexer.ASTERISKASSIGN,
		lexer.SLASHASSIGN, lexer.PERCENTASSIGN, lexer.AMPERSANDASSIGN,
		lexer.PIPEASSIGN, lexer.CARETASSIGN, lexer.LSHIFTASSIGN, lexer.RSHIFTASSIGN:
		return true
	default:
		return false
	}
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}

	// Parse the identifier being assigned to
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Move to the assignment operator
	p.nextToken()
	stmt.Operator = p.curToken.Literal

	// Parse the value expression
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Block statement parsing
// parseBlockStatement moved to block_parser.go

// Enhanced error recovery and synchronization
func (p *Parser) synchronize() {
	p.nextToken()

	for !p.curTokenIs(lexer.EOF) {
		if p.curTokenIs(lexer.SEMICOLON) {
			return
		}

		switch p.peekToken.Type {
		case lexer.DEF, lexer.VAL, lexer.VAR, lexer.TYPE, lexer.STRUCT, lexer.IMPL, lexer.RETURN:
			return
		}

		p.nextToken()
	}
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) recoverFromError() {
	p.synchronize()
}

// Statement implementations (basic versions for completeness)
func (p *Parser) parseTypeStatement() ast.Statement {
	stmt := &ast.TypeStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Parse type directly without expecting ASSIGN token
	// Supports "type Name Type" syntax (no equals sign)
	p.nextToken()
	stmt.Type = p.parseTypeExpression()

	return stmt
}

func (p *Parser) parseStructStatement() ast.Statement {
	stmt := &ast.StructStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	// Parse struct fields
	stmt.Fields = []*ast.StructField{}

	p.nextToken()
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		field := p.parseStructFieldDefinition()
		if field != nil {
			stmt.Fields = append(stmt.Fields, field)
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseStructFieldDefinition() *ast.StructField {
	if !p.curTokenIs(lexer.IDENT) {
		return nil
	}

	field := &ast.StructField{}
	field.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()
	fieldType := p.parseTypeExpression()

	// For struct definitions, we store the type as a simple identifier expression
	// This is a simplification for now
	field.Value = &ast.Identifier{Token: p.curToken, Value: fieldType.String()}

	return field
}

func (p *Parser) parseImplStatement() ast.Statement {
	stmt := &ast.ImplStatement{Token: p.curToken}

	p.nextToken() // Move past 'impl'

	// Handle pointer types: *, &
	isPointer := false
	isReference := false
	
	if p.curTokenIs(lexer.ASTERISK) {
		isPointer = true
		p.nextToken() // Move past '*'
	} else if p.curTokenIs(lexer.AMPERSAND) {
		isReference = true
		p.nextToken() // Move past '&'
	}

	if !p.curTokenIs(lexer.IDENT) {
		p.addError("expected type name after impl")
		return nil
	}

	// Create type identifier with receiver type info
	typeName := p.curToken.Literal
	fullTypeName := typeName
	if isPointer {
		fullTypeName = "*" + typeName
	} else if isReference {
		fullTypeName = "&" + typeName
	}
	
	stmt.Type = &ast.Identifier{Token: p.curToken, Value: fullTypeName}
	
	// Set receiver info
	stmt.ReceiverInfo = &ast.ReceiverInfo{
		TypeName: typeName,
	}
	if isPointer {
		stmt.ReceiverInfo.Type = ast.PointerReceiver
	} else if isReference {
		stmt.ReceiverInfo.Type = ast.ReferenceReceiver
	} else {
		stmt.ReceiverInfo.Type = ast.ValueReceiver
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	// Parse impl methods directly using main parser
	stmt.Methods = []*ast.FunctionStatement{}

	p.nextToken() // consume '{'
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		// Skip semicolons and newlines
		if p.curTokenIs(lexer.SEMICOLON) {
			p.nextToken()
			continue
		}

		if p.curTokenIs(lexer.DEF) {
			method := p.parseFunctionStatement()
			if method != nil {
				stmt.Methods = append(stmt.Methods, method)
			}
			// After parsing function, we should be positioned correctly
			// Only advance if we're not at a logical boundary
			if !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.DEF) && !p.curTokenIs(lexer.EOF) {
				p.nextToken()
			}
		} else {
			// If we're not at a DEF token and not at the end, advance
			p.nextToken()
		}
	}

	// Consume the closing RBRACE
	if p.curTokenIs(lexer.RBRACE) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIncludeStatement() ast.Statement {
	stmt := &ast.IncludeStatement{Token: p.curToken}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}

	stmt.Path = p.curToken.Literal

	// Register C functions from the included header
	p.cRegistry.IncludeHeader(stmt.Path)

	return stmt
}

func (p *Parser) parseImportStatement() ast.Statement {
	// Sango module imports: import "module.sango"
	p.addError("import statements not fully implemented yet")
	p.recoverFromError()
	return nil
}

func (p *Parser) parseDefineStatement() ast.Statement {
	stmt := &ast.DefineStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Simple implementation: collect everything until end of line as value
	var value bytes.Buffer
	for !p.peekTokenIs(lexer.EOF) && p.peekToken.Line == p.curToken.Line {
		p.nextToken()
		value.WriteString(p.curToken.Literal)
		if !p.peekTokenIs(lexer.EOF) && p.peekToken.Line == p.curToken.Line {
			value.WriteString(" ")
		}
	}

	stmt.Value = value.String()
	return stmt
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	// Check for return type
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken()
		p.nextToken()
		stmt.ReturnType = p.parseTypeExpression()
	}

	// Parse function body - must have '=' before body
	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	
	// Handle function body - could be expression or block
	if p.curTokenIs(lexer.LBRACE) {
		stmt.Body = p.parseBlockStatement()
		// parseBlockStatement handles RBRACE consumption
	} else {
		stmt.Body = p.parseExpression(0)
	}

	return stmt
}

func (p *Parser) parseForStatement() ast.Statement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Variable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check for <- or 'in'
	if p.peekTokenIs(lexer.LARROW) {
		// for x <- iterable
		p.nextToken()
		stmt.IsInRange = false
	} else if p.peekTokenIs(lexer.IN) {
		// for i in range
		p.nextToken()
		stmt.IsInRange = true
	} else {
		p.addError("expected '<-' or 'in' after for variable")
		return nil
	}

	p.nextToken()
	stmt.Iterable = p.parseExpression(0)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	// parseBlockStatementFixed already handles RBRACE consumption
	return stmt
}

func (p *Parser) parseWhileStatement() ast.Statement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(0)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	// parseBlockStatementFixed already handles RBRACE consumption
	return stmt
}

func (p *Parser) parseDeferStatement() ast.Statement {
	stmt := &ast.DeferStatement{Token: p.curToken}

	p.nextToken()
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssertStatement() ast.Statement {
	stmt := &ast.AssertStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Expression = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
