package v2

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/rxxuzi/sango/pkg/ast"
	sangoLexer "github.com/rxxuzi/sango/pkg/lexer"
)

// BlockParser は { f() } のようなブロック式を解析するための専用パーサー
type BlockParser struct {
	parser *participle.Parser[BlockStmt]
}

// BlockStmt は Participle v2 用のブロック文定義
type BlockStmt struct {
	Statements []SimpleStmt `"{"  @@* "}"`
}

// SimpleStmt は基本的な文を表現
type SimpleStmt struct {
	ExpressionStmt *ExprStmt `@@`
}

// ExprStmt は式文を表現
type ExprStmt struct {
	Expression *SimpleExpr `@@`
}

// SimpleExpr は基本的な式を表現
type SimpleExpr struct {
	FuncCall   *FuncCall `@@`
	Identifier *string   `| @Ident`
	Integer    *string   `| @Int`
	String     *string   `| @String`
}

// FuncCall は関数呼び出しを表現
type FuncCall struct {
	Name      string       `@Ident`
	Arguments []SimpleExpr `"(" (@@ ("," @@)*)? ")"`
}

// Sango lexer definition (simplified for block parsing)
var blockLexer = lexer.MustSimple([]lexer.SimpleRule{
	// Comments
	{Name: "comment", Pattern: `//.*|/\*([^*]|\*[^/])*\*/`},

	// Delimiters
	{Name: "(", Pattern: `\(`},
	{Name: ")", Pattern: `\)`},
	{Name: "{", Pattern: `\{`},
	{Name: "}", Pattern: `\}`},
	{Name: ",", Pattern: `,`},
	{Name: ";", Pattern: `;`},

	// Literals
	{Name: "Int", Pattern: `\d+`},
	{Name: "String", Pattern: `"(?:[^"\\]|\\.)*"`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},

	// Whitespace
	{Name: "whitespace", Pattern: `\s+`},
})

// NewBlockParser creates a new block parser
func NewBlockParser() (*BlockParser, error) {
	parser, err := participle.Build[BlockStmt](
		participle.Lexer(blockLexer),
		participle.Elide("whitespace", "comment"),
		participle.Unquote("String"),
	)
	if err != nil {
		return nil, err
	}

	return &BlockParser{parser: parser}, nil
}

// ParseBlockExpression はブロック式を解析して元のAST形式に変換
func (bp *BlockParser) ParseBlockExpression(source string) (*ast.BlockStatement, error) {
	blockStmt, err := bp.parser.ParseString("", source)
	if err != nil {
		return nil, err
	}

	// Participle v2の結果を元のAST形式に変換
	result := &ast.BlockStatement{
		Token:      sangoLexer.Token{Type: sangoLexer.LBRACE, Literal: "{"},
		Statements: []ast.Statement{},
	}

	for _, stmt := range blockStmt.Statements {
		if stmt.ExpressionStmt != nil {
			expr := stmt.ExpressionStmt.Expression

			var astExpr ast.Expression
			if expr.FuncCall != nil {
				// 関数呼び出しをAST形式に変換
				args := []ast.Expression{}
				for _, arg := range expr.FuncCall.Arguments {
					args = append(args, convertToASTExpression(arg))
				}

				astExpr = &ast.CallExpression{
					Token:     sangoLexer.Token{Type: sangoLexer.IDENT, Literal: expr.FuncCall.Name},
					Function:  &ast.Identifier{Value: expr.FuncCall.Name},
					Arguments: args,
				}
			} else {
				astExpr = convertToASTExpression(*expr)
			}

			result.Statements = append(result.Statements, &ast.ExpressionStatement{
				Token:      sangoLexer.Token{Type: sangoLexer.IDENT},
				Expression: astExpr,
			})
		}
	}

	return result, nil
}

// convertToASTExpression は SimpleExpr を ast.Expression に変換
func convertToASTExpression(expr SimpleExpr) ast.Expression {
	if expr.Identifier != nil {
		return &ast.Identifier{
			Token: sangoLexer.Token{Type: sangoLexer.IDENT, Literal: *expr.Identifier},
			Value: *expr.Identifier,
		}
	}
	if expr.Integer != nil {
		return &ast.IntegerLiteral{
			Token: sangoLexer.Token{Type: sangoLexer.INT, Literal: *expr.Integer},
			Value: 0, // TODO: parse the actual value
		}
	}
	if expr.String != nil {
		return &ast.StringLiteral{
			Token: sangoLexer.Token{Type: sangoLexer.STRING, Literal: *expr.String},
			Value: *expr.String,
		}
	}
	if expr.FuncCall != nil {
		args := []ast.Expression{}
		for _, arg := range expr.FuncCall.Arguments {
			args = append(args, convertToASTExpression(arg))
		}

		return &ast.CallExpression{
			Token:     sangoLexer.Token{Type: sangoLexer.IDENT, Literal: expr.FuncCall.Name},
			Function:  &ast.Identifier{Value: expr.FuncCall.Name},
			Arguments: args,
		}
	}

	// Fallback
	return &ast.Identifier{Value: "unknown"}
}
