package parser

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
	"github.com/rxxuzi/sango/pkg/parser/v2"
)

// Adapter methods to make Parser implement v2.ParserInterface

func (p *Parser) CurTokenIs(t lexer.TokenType) bool {
	return p.curTokenIs(t)
}

func (p *Parser) PeekTokenIs(t lexer.TokenType) bool {
	return p.peekTokenIs(t)
}

func (p *Parser) NextToken() {
	p.nextToken()
}

func (p *Parser) ExpectPeek(t lexer.TokenType) bool {
	return p.expectPeek(t)
}

func (p *Parser) GetCurrentToken() lexer.Token {
	return p.curToken
}

func (p *Parser) GetPeekToken() lexer.Token {
	return p.peekToken
}

func (p *Parser) ParseExpression(precedence v2.Precedence) ast.Expression {
	// Convert v2.Precedence to our internal Precedence
	return p.parseExpression(Precedence(precedence))
}

func (p *Parser) ParseStatement() ast.Statement {
	return p.parseStatement()
}

func (p *Parser) AddError(msg string) {
	p.addError(msg)
}

// V2 parser instances
type V2Parsers struct {
	Array       *v2.ArrayParser
	Block       *v2.BlockParser
	SimpleBlock *v2.SimpleBlockParser
	Struct      *v2.StructParser
	ControlFlow *v2.ControlFlowParser
}

// Initialize v2 parsers
func (p *Parser) initV2Parsers() {
	p.v2 = &V2Parsers{
		Array:       v2.NewArrayParser(),
		Block:       v2.NewBlockParser(),
		SimpleBlock: v2.NewSimpleBlockParser(),
		Struct:      v2.NewStructParser(),
		ControlFlow: v2.NewControlFlowParser(),
	}
}