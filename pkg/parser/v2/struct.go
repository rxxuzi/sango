package v2

import (
	"github.com/rxxuzi/sango/pkg/ast"
	"github.com/rxxuzi/sango/pkg/lexer"
)

// StructParser handles struct-related parsing operations
type StructParser struct {
	// Struct-specific configuration
}

// NewStructParser creates a new struct parser
func NewStructParser() *StructParser {
	return &StructParser{}
}

// ParseStructLiteral parses struct initialization: Struct { .field = value }
func (sp *StructParser) ParseStructLiteral(p ParserInterface, token lexer.Token) ast.Expression {
	lit := &ast.StructLiteral{Token: token}
	
	// We're already positioned at the first field token
	lit.Fields = sp.parseStructFields(p)
	
	if !p.ExpectPeek(lexer.RBRACE) {
		return nil
	}
	
	return lit
}

// parseStructFields parses field assignments in struct literals
func (sp *StructParser) parseStructFields(p ParserInterface) []*ast.StructField {
	fields := []*ast.StructField{}
	
	if p.CurTokenIs(lexer.RBRACE) {
		return fields
	}
	
	// Parse first field
	field := sp.parseStructField(p)
	if field != nil {
		fields = append(fields, field)
	}
	
	for p.PeekTokenIs(lexer.COMMA) {
		p.NextToken() // consume comma
		
		// Accept either IDENT or DOT for field names
		if p.PeekTokenIs(lexer.IDENT) {
			p.NextToken()
		} else if p.PeekTokenIs(lexer.DOT) {
			p.NextToken()
		} else {
			return nil
		}
		
		field := sp.parseStructField(p)
		if field != nil {
			fields = append(fields, field)
		}
	}
	
	return fields
}

// parseStructField parses individual field assignments
func (sp *StructParser) parseStructField(p ParserInterface) *ast.StructField {
	field := &ast.StructField{}
	
	if p.CurTokenIs(lexer.DOT) {
		// Handle .field = value syntax
		if !p.ExpectPeek(lexer.IDENT) {
			return nil
		}
		field.Name = &ast.Identifier{Token: p.GetCurrentToken(), Value: p.GetCurrentToken().Literal}
		
		if !p.ExpectPeek(lexer.ASSIGN) {
			return nil
		}
	} else if p.CurTokenIs(lexer.IDENT) {
		// Handle field: value syntax
		field.Name = &ast.Identifier{Token: p.GetCurrentToken(), Value: p.GetCurrentToken().Literal}
		
		if !p.ExpectPeek(lexer.COLON) {
			return nil
		}
	} else {
		return nil
	}
	
	p.NextToken()
	field.Value = p.ParseExpression(LOWEST)
	
	return field
}