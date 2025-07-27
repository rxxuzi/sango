package parser

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// Temporary compatibility functions until we fully migrate to v2 parsers

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// Use fixed block parser that properly handles RBRACE consumption
	return p.parseBlockStatementFixed()
}

func (p *Parser) parseBlockExpressionFromBrace(token lexer.Token) ast.Expression {
	// Use fixed block parser for expressions too
	return p.parseBlockStatementFixed()
}