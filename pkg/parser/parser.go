package parser

import (
	"fmt"

	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/cinterop"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// Precedence levels for operator precedence parsing
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

// Operator precedence table
var precedences = map[lexer.TokenType]Precedence{
	// Assignment operators (right-associative, lowest precedence)
	lexer.ASSIGN:          ASSIGN,
	lexer.PLUSASSIGN:      ASSIGN,
	lexer.MINUSASSIGN:     ASSIGN,
	lexer.ASTERISKASSIGN:  ASSIGN,
	lexer.SLASHASSIGN:     ASSIGN,
	lexer.PERCENTASSIGN:   ASSIGN,
	lexer.AMPERSANDASSIGN: ASSIGN,
	lexer.PIPEASSIGN:      ASSIGN,
	lexer.CARETASSIGN:     ASSIGN,
	lexer.LSHIFTASSIGN:    ASSIGN,
	lexer.RSHIFTASSIGN:    ASSIGN,

	// Logical operators
	lexer.OR:  OR,
	lexer.AND: AND,

	// Bitwise operators
	lexer.PIPE:      BITOR,
	lexer.CARET:     BITXOR,
	lexer.AMPERSAND: BITAND,

	// Comparison operators
	lexer.EQ:  EQUALS,
	lexer.NEQ: EQUALS,
	lexer.LT:  LESSGREATER,
	lexer.GT:  LESSGREATER,
	lexer.LEQ: LESSGREATER,
	lexer.GEQ: LESSGREATER,

	// Shift operators
	lexer.LSHIFT: SHIFT,
	lexer.RSHIFT: SHIFT,

	// Arithmetic operators
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTERISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
	lexer.PERCENT:  PRODUCT,
	lexer.POWER:    POWER,

	// Access operators
	lexer.LPAREN:   CALL,
	lexer.LBRACKET: INDEX,
	lexer.LBRACE:   CALL,
	lexer.DOT:      DOT,

	// Range operators
	lexer.DOTDOT:   LESSGREATER,
	lexer.DOTDOTEQ: LESSGREATER,
}

// Parser represents the Sango parser
type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token

	errors []string

	// Parsing functions
	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn

	// C function registry
	cRegistry *cinterop.FunctionRegistry

	// Bracket tracking stack for proper nesting
	bracketStack []lexer.TokenType
	
	// V2 parsers
	v2 *V2Parsers
}

// Function types for Pratt parsing
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:            l,
		errors:       []string{},
		cRegistry:    cinterop.NewFunctionRegistry(),
		bracketStack: []lexer.TokenType{},
	}

	// Initialize prefix parse functions
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NULL, p.parseNullLiteral)
	p.registerPrefix(lexer.UNDERSCORE, p.parseWildcardExpression)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.TILDE, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseBraceExpression)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.MATCH, p.parseMatchExpression)
	p.registerPrefix(lexer.DEF, p.parseFunctionLiteral)
	p.registerPrefix(lexer.DOT, p.parseDotFieldExpression)
	p.registerPrefix(lexer.SIZEOF, p.parseSizeofExpression)
	
	// Register primitive type tokens as prefix parsers
	p.registerPrefix(lexer.INT_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.LONG_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.FLOAT_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.DOUBLE_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.BOOL_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.STRING_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.VOID_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.I8_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.I16_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.I32_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.I64_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.U8_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.U16_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.U32_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.U64_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.F32_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.F64_TYPE, p.parseTypeIdentifier)
	p.registerPrefix(lexer.BYTE_TYPE, p.parseTypeIdentifier)

	// Initialize infix parse functions
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.PERCENT, p.parseInfixExpression)
	p.registerInfix(lexer.POWER, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NEQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LEQ, p.parseInfixExpression)
	p.registerInfix(lexer.GEQ, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.AMPERSAND, p.parseInfixExpression)
	p.registerInfix(lexer.PIPE, p.parseInfixExpression)
	p.registerInfix(lexer.CARET, p.parseInfixExpression)
	p.registerInfix(lexer.LSHIFT, p.parseInfixExpression)
	p.registerInfix(lexer.RSHIFT, p.parseInfixExpression)
	p.registerInfix(lexer.DOTDOT, p.parseRangeExpression)
	p.registerInfix(lexer.DOTDOTEQ, p.parseRangeExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.LBRACE, p.parseStructConstructorExpression)
	p.registerInfix(lexer.DOT, p.parseDotExpression)

	// Initialize v2 parsers
	p.initV2Parsers()

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Error handling
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d:%d",
		t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found at line %d:%d",
		t, p.curToken.Line, p.curToken.Column)
	p.errors = append(p.errors, msg)
}

func (p *Parser) unexpectedTokenError(expected string) {
	msg := fmt.Sprintf("unexpected token %s, expected %s at line %d:%d",
		p.curToken.Type, expected, p.curToken.Line, p.curToken.Column)
	p.errors = append(p.errors, msg)
}

// Token management
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Precedence helpers
func (p *Parser) peekPrecedence() Precedence {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() Precedence {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// Bracket stack management
func (p *Parser) pushBracket(bracket lexer.TokenType) {
	p.bracketStack = append(p.bracketStack, bracket)
}

func (p *Parser) popBracket() lexer.TokenType {
	if len(p.bracketStack) == 0 {
		return lexer.ILLEGAL
	}
	bracket := p.bracketStack[len(p.bracketStack)-1]
	p.bracketStack = p.bracketStack[:len(p.bracketStack)-1]
	return bracket
}

func (p *Parser) peekBracket() lexer.TokenType {
	if len(p.bracketStack) == 0 {
		return lexer.ILLEGAL
	}
	return p.bracketStack[len(p.bracketStack)-1]
}

func (p *Parser) isInParentheses() bool {
	return p.peekBracket() == lexer.LPAREN
}

// shouldStopExpression determines if expression parsing should stop based on context
func (p *Parser) shouldStopExpression() bool {
	// If we see a closing brace, only stop if we're not inside parentheses
	if p.peekTokenIs(lexer.RBRACE) {
		// If we're inside parentheses, don't stop for braces
		return !p.isInParentheses()
	}
	return false
}

// Parser function registration
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// ParseProgram parses the entire program and validates structure
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		
		// Only advance token if we're not at a natural statement boundary
		// Some statements (like functions with blocks) handle their own token advancement
		if !p.curTokenIs(lexer.EOF) && 
		   !p.curTokenIs(lexer.DEF) && 
		   !p.curTokenIs(lexer.STRUCT) && 
		   !p.curTokenIs(lexer.TYPE) && 
		   !p.curTokenIs(lexer.IMPL) &&
		   !p.curTokenIs(lexer.FOR) &&
		   !p.curTokenIs(lexer.WHILE) {
			p.nextToken()
		}
	}

	// Validate program structure for executables
	p.validateProgram(program)

	return program
}

// validateProgram checks if the program has required structure for executable Sango programs
func (p *Parser) validateProgram(program *ast.Program) {
	hasMainFunction := false

	// Check for main function
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.ExpressionStatement:
			if fn, ok := s.Expression.(*ast.FunctionLiteral); ok {
				if fn.Name != nil && fn.Name.Value == "main" {
					hasMainFunction = true

					// Validate main function signature
					if len(fn.Parameters) > 1 {
						p.addError("main function can have at most one parameter (args: []string)")
					}

					// Check if main has proper signature
					if len(fn.Parameters) == 1 {
						param := fn.Parameters[0]
						if param.Type == nil {
							p.addError("main function parameter should have type []string")
						}
					}

					break
				}
			}
		}
	}

	// Note: main function requirement is now optional
	// The compiler (sangoc) will handle main function requirements based on compilation mode
	// For now, we just validate main function signature if it exists
	_ = hasMainFunction // suppress unused variable warning
}

// isExecutableProgram determines if this looks like an executable program vs library/test code
func (p *Parser) isExecutableProgram(program *ast.Program) bool {
	// Simple heuristic: if it has function definitions or more than 3 statements,
	// it's likely meant to be executable
	functionCount := 0
	for _, stmt := range program.Statements {
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			if _, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
				functionCount++
			}
		}
	}

	// If it has multiple functions or is a substantial program, require main
	return functionCount > 1 || len(program.Statements) > 5
}
