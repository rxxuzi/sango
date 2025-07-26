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

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
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
	array.Elements = p.parseExpressionList(lexer.RBRACKET)
	return array
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
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
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

	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

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
	// TODO: Implement match expression parsing with patterns
	// This requires extending the AST to include match expressions and patterns
	return nil
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

// parseStructLiteral parses struct literals: { name: value, ... } or Type { name: value, ... }
func (p *Parser) parseStructLiteral() ast.Expression {
	lit := &ast.StructLiteral{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) && !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	// Check if it's an empty struct literal
	if p.curTokenIs(lexer.RBRACE) {
		return lit
	}

	// Parse first field
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