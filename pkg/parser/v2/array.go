package v2

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// ArrayParser handles array-related parsing operations
type ArrayParser struct {
	// We can add array-specific configuration here if needed
}

// NewArrayParser creates a new array parser
func NewArrayParser() *ArrayParser {
	return &ArrayParser{}
}

// ParseArrayLiteral parses array literals [1, 2, 3]
func (ap *ArrayParser) ParseArrayLiteral(p ParserInterface, token lexer.Token) ast.Expression {
	lit := &ast.ArrayLiteral{Token: token}

	if p.PeekTokenIs(lexer.RBRACKET) {
		p.NextToken()
		return lit
	}

	p.NextToken()
	lit.Elements = append(lit.Elements, p.ParseExpression(LOWEST))

	for p.PeekTokenIs(lexer.COMMA) {
		p.NextToken()
		p.NextToken()
		lit.Elements = append(lit.Elements, p.ParseExpression(LOWEST))
	}

	if !p.ExpectPeek(lexer.RBRACKET) {
		return nil
	}

	return lit
}

// ParseIndexExpression parses array[index] expressions
func (ap *ArrayParser) ParseIndexExpression(p ParserInterface, left ast.Expression, token lexer.Token) ast.Expression {
	exp := &ast.IndexExpression{Token: token, Left: left}

	p.NextToken()
	exp.Index = p.ParseExpression(LOWEST)

	if !p.ExpectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

// ParseSliceExpression parses array[start..end] expressions
func (ap *ArrayParser) ParseSliceExpression(p ParserInterface, left ast.Expression, token lexer.Token) ast.Expression {
	// TODO: Implement slice expression parsing
	// This would handle patterns like arr[1..3], arr[..5], arr[2..]
	return nil
}