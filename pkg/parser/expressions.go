package parser

import (
	"strconv"

	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// Expression parsing using Pratt parsing
func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.SEMICOLON) &&
		!p.peekTokenIs(lexer.RBRACE) && !p.peekTokenIs(lexer.RBRACKET) &&
		!p.peekTokenIs(lexer.COMMA) && !p.peekTokenIs(lexer.EOF) &&
		precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Prefix expression parsers
func (p *Parser) parseIdentifier() ast.Expression {
	// Check if this identifier is a known C function
	if p.cRegistry.IsFunction(p.curToken.Literal) {
		// For now, treat it as a regular identifier
		// In the future, we might want to create a special CFunctionIdentifier node
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := "could not parse " + p.curToken.Literal + " as integer"
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := "could not parse " + p.curToken.Literal + " as float"
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseWildcardExpression() ast.Expression {
	return &ast.WildcardExpression{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	// Check if this is a tuple literal
	if p.curTokenIs(lexer.RPAREN) {
		// Empty tuple ()
		return &ast.TupleLiteral{Token: p.curToken, Elements: []ast.Expression{}}
	}

	exp := p.parseExpression(LOWEST)

	// Check if there's a comma, indicating a tuple
	if p.peekTokenIs(lexer.COMMA) {
		return p.parseTupleLiteral(exp)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseTupleLiteral(firstElement ast.Expression) ast.Expression {
	tuple := &ast.TupleLiteral{Token: p.curToken}
	tuple.Elements = []ast.Expression{firstElement}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		if p.peekTokenIs(lexer.RPAREN) {
			// Trailing comma
			break
		}
		p.nextToken()
		tuple.Elements = append(tuple.Elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return tuple
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	// Check if this is a typed empty array like []int
	if p.peekTokenIs(lexer.RBRACKET) {
		p.nextToken() // consume ]

		// Check if followed by a type
		if p.isTypeToken(p.peekToken.Type) {
			p.nextToken() // move to the type
			// For now, we'll just parse this as an empty array
			// In the future, we might want to store the type information
			return array
		}

		// Just an empty array []
		return array
	}

	array.Elements = p.parseExpressionList(lexer.RBRACKET)
	return array
}

func (p *Parser) isTypeToken(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.INT_TYPE, lexer.LONG_TYPE, lexer.FLOAT_TYPE, lexer.DOUBLE_TYPE,
		lexer.BOOL_TYPE, lexer.STRING_TYPE, lexer.VOID_TYPE,
		lexer.I8_TYPE, lexer.I16_TYPE, lexer.I32_TYPE, lexer.I64_TYPE,
		lexer.U8_TYPE, lexer.U16_TYPE, lexer.U32_TYPE, lexer.U64_TYPE,
		lexer.F32_TYPE, lexer.F64_TYPE, lexer.BYTE_TYPE:
		return true
	default:
		return false
	}
}

// Infix expression parsers
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	// Handle right-associative operators (power **)
	if expression.Operator == "**" {
		expression.Right = p.parseExpression(precedence - 1)
	} else {
		expression.Right = p.parseExpression(precedence)
	}

	return expression
}

func (p *Parser) parseRangeExpression(left ast.Expression) ast.Expression {
	expression := &ast.RangeExpression{
		Token:     p.curToken,
		Start:     left,
		Inclusive: p.curToken.Type == lexer.DOTDOTEQ,
	}

	p.nextToken()

	// Check if there's an end expression
	if !p.curTokenIs(lexer.RBRACKET) && !p.curTokenIs(lexer.SEMICOLON) && !p.curTokenIs(lexer.RPAREN) {
		expression.End = p.parseExpression(LOWEST)
	}

	return expression
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: fn}
	args := []ast.Expression{}

	// Check if the next token is ')', meaning empty arguments
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken() // consume ')'
		exp.Arguments = args
		return exp
	}

	// Parse the first argument
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	// Parse additional arguments
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next argument
		args = append(args, p.parseExpression(LOWEST))
	}

	// Expect the closing ')'
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	exp.Arguments = args
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	// For now, treat dot as member access via infix expression
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(DOT)

	return expression
}

// Helper function to parse expression lists
func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	// If next token is the end token, we have an empty list
	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	// Move to the first argument
	p.nextToken()

	// Check if we immediately hit the end token after moving
	if p.curTokenIs(end) {
		return args
	}

	args = append(args, p.parseExpression(LOWEST))

	// Parse additional arguments separated by commas
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next argument
		args = append(args, p.parseExpression(LOWEST))
	}

	// Expect the end token
	if !p.expectPeek(end) {
		return nil
	}

	return args
}

// Complex expression parsing
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// Handle else clause
	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseMatchExpression() ast.Expression {
	expr := &ast.MatchExpression{Token: p.curToken}

	p.nextToken()
	expr.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expr.Cases = []*ast.MatchCase{}

	p.nextToken()
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		matchCase := p.parseMatchCase()
		if matchCase != nil {
			expr.Cases = append(expr.Cases, matchCase)
		}
		p.nextToken()
	}

	return expr
}

func (p *Parser) parseMatchCase() *ast.MatchCase {
	matchCase := &ast.MatchCase{}

	// Parse pattern (left side of =>)
	matchCase.Pattern = p.parseExpression(LOWEST)

	// Check for guard clause (if condition)
	if p.peekTokenIs(lexer.IF) {
		p.nextToken() // consume 'if'
		p.nextToken() // move to the guard expression
		matchCase.Guard = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(lexer.DARROW) { // => token
		return nil
	}

	p.nextToken()
	matchCase.Value = p.parseExpression(LOWEST)

	return matchCase
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Check if function has a name (def name(...) vs def(...))
	if p.peekTokenIs(lexer.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken()
		p.nextToken()
		lit.ReturnType = p.parseTypeExpression()
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()

	// Function body can be a single expression or block
	if p.curTokenIs(lexer.LBRACE) {
		lit.Body = p.parseBlockStatement()
	} else {
		lit.Body = p.parseExpression(LOWEST)
	}

	return lit
}

// parseBraceExpression parses either struct literals or block statements based on content
func (p *Parser) parseBraceExpression() ast.Expression {
	token := p.curToken

	// Push the opening brace onto the bracket stack
	p.pushBracket(lexer.LBRACE)

	// Look ahead to determine if this is a struct literal or block statement
	if p.peekTokenIs(lexer.RBRACE) {
		// Empty braces - treat as empty struct literal
		p.nextToken()
		p.popBracket() // pop the matching '{'
		return &ast.StructLiteral{Token: token, Fields: []*ast.StructField{}}
	}

	// Look for identifier followed by colon (struct literal pattern)
	if p.peekTokenIs(lexer.IDENT) {
		// Save current position for backtracking
		currentPos := p.curToken
		peekPos := p.peekToken

		p.nextToken() // move to IDENT
		if p.peekTokenIs(lexer.COLON) {
			// This is a struct literal: { name: value }
			result := p.parseStructLiteralFromBrace(token)
			p.popBracket() // pop the matching '{'
			return result
		}

		// Not a struct literal, restore position and parse as block
		p.curToken = currentPos
		p.peekToken = peekPos
	}

	// Parse as block statement
	result := p.parseBlockStatement()
	p.popBracket() // pop the matching '{'
	return result
}

// parseStructLiteralFromBrace parses struct literals starting from after '{'
func (p *Parser) parseStructLiteralFromBrace(token lexer.Token) ast.Expression {
	lit := &ast.StructLiteral{Token: token}

	// We're already positioned at the first identifier
	lit.Fields = p.parseStructFields()

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return lit
}

func (p *Parser) parseStructFields() []*ast.StructField {
	fields := []*ast.StructField{}

	if p.curTokenIs(lexer.RBRACE) {
		return fields
	}

	// Parse first field
	field := p.parseStructField()
	if field != nil {
		fields = append(fields, field)
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		field := p.parseStructField()
		if field != nil {
			fields = append(fields, field)
		}
	}

	return fields
}

func (p *Parser) parseStructField() *ast.StructField {
	field := &ast.StructField{}

	if !p.curTokenIs(lexer.IDENT) {
		return nil
	}

	field.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()
	field.Value = p.parseExpression(LOWEST)

	return field
}
