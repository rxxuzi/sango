package ast

import (
	"bytes"
	"strings"

	"github.com/rxxuzi/sango/pkg/lexer"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is a node that doesn't produce a value
type Statement interface {
	Node
	statementNode()
}

// Expression is a node that produces a value
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents a variable or function name
type Identifier struct {
	Token lexer.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// ValStatement represents val x = expression
type ValStatement struct {
	Token lexer.Token // the token.VAL token
	Names []*Identifier
	Type  *TypeExpression // optional type annotation
	Value Expression
}

func (vs *ValStatement) statementNode()       {}
func (vs *ValStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *ValStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " ")
	names := []string{}
	for _, name := range vs.Names {
		names = append(names, name.String())
	}
	out.WriteString(strings.Join(names, ", "))
	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// VarStatement represents var x = expression
type VarStatement struct {
	Token lexer.Token // the token.VAR token
	Names []*Identifier
	Type  *TypeExpression // optional type annotation
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " ")
	names := []string{}
	for _, name := range vs.Names {
		names = append(names, name.String())
	}
	out.WriteString(strings.Join(names, ", "))
	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents return expression
type ReturnStatement struct {
	Token       lexer.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// AssignmentStatement represents identifier = expression
type AssignmentStatement struct {
	Token    lexer.Token // the assignment token
	Name     *Identifier
	Operator string // =, +=, -=, etc.
	Value    Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(" " + as.Operator + " ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement is a statement consisting of a single expression
type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// FunctionStatement represents top-level function definitions
type FunctionStatement struct {
	Token      lexer.Token // the 'def' token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *TypeExpression
	Body       Expression // can be BlockStatement or expression
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fs.TokenLiteral() + " ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fs.ReturnType != nil {
		out.WriteString(": ")
		out.WriteString(fs.ReturnType.String())
	}
	out.WriteString(" = ")
	if fs.Body != nil {
		out.WriteString(fs.Body.String())
	}
	return out.String()
}

// IncludeStatement represents include "header.h"
type IncludeStatement struct {
	Token lexer.Token // the 'include' token
	Path  string
}

func (is *IncludeStatement) statementNode()       {}
func (is *IncludeStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncludeStatement) String() string {
	return is.TokenLiteral() + " \"" + is.Path + "\""
}

// TypeStatement represents type aliases: type Name = Type
type TypeStatement struct {
	Token lexer.Token // the 'type' token
	Name  *Identifier
	Type  *TypeExpression
}

func (ts *TypeStatement) statementNode()       {}
func (ts *TypeStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TypeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ts.TokenLiteral() + " ")
	out.WriteString(ts.Name.String())
	out.WriteString(" = ")
	out.WriteString(ts.Type.String())
	return out.String()
}

// StructStatement represents struct definitions
type StructStatement struct {
	Token  lexer.Token // the 'struct' token
	Name   *Identifier
	Fields []*StructField
}

func (ss *StructStatement) statementNode()       {}
func (ss *StructStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StructStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ss.TokenLiteral() + " ")
	out.WriteString(ss.Name.String())
	out.WriteString(" { ")
	fields := []string{}
	for _, field := range ss.Fields {
		fields = append(fields, field.String())
	}
	out.WriteString(strings.Join(fields, "; "))
	out.WriteString(" }")
	return out.String()
}

// ImplStatement represents implementation blocks
type ImplStatement struct {
	Token        lexer.Token // the 'impl' token
	Type         *Identifier
	ReceiverInfo *ReceiverInfo // Parsed receiver type information
	Methods      []*FunctionStatement
}

func (is *ImplStatement) statementNode()       {}
func (is *ImplStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImplStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.TokenLiteral() + " ")
	out.WriteString(is.Type.String())
	out.WriteString(" { ")
	methods := []string{}
	for _, method := range is.Methods {
		methods = append(methods, method.String())
	}
	out.WriteString(strings.Join(methods, "; "))
	out.WriteString(" }")
	return out.String()
}

// DefineStatement represents C-style macro definitions
type DefineStatement struct {
	Token lexer.Token // the 'define' token
	Name  *Identifier
	Value string // Simple string for now
}

func (ds *DefineStatement) statementNode()       {}
func (ds *DefineStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DefineStatement) String() string {
	return ds.TokenLiteral() + " " + ds.Name.String() + " " + ds.Value
}

// ForStatement represents for loops: for x <- iterable { ... } or for i in range { ... }
type ForStatement struct {
	Token     lexer.Token // the 'for' token
	Variable  *Identifier
	Iterable  Expression
	Body      *BlockStatement
	IsInRange bool // true for 'for i in range', false for 'for x <- iterable'
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fs.TokenLiteral() + " ")
	out.WriteString(fs.Variable.String())
	if fs.IsInRange {
		out.WriteString(" in ")
	} else {
		out.WriteString(" <- ")
	}
	out.WriteString(fs.Iterable.String())
	out.WriteString(" ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// WhileStatement represents while loops: while condition { ... }
type WhileStatement struct {
	Token     lexer.Token // the 'while' token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ws.TokenLiteral() + " ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// DeferStatement represents defer statement: defer expr
type DeferStatement struct {
	Token      lexer.Token // the 'defer' token
	Expression Expression
}

func (ds *DeferStatement) statementNode()       {}
func (ds *DeferStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeferStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ds.TokenLiteral() + " ")
	if ds.Expression != nil {
		out.WriteString(ds.Expression.String())
	}
	return out.String()
}

// AssertStatement represents assert statement: assert expr
type AssertStatement struct {
	Token      lexer.Token // the 'assert' token
	Expression Expression
}

func (as *AssertStatement) statementNode()       {}
func (as *AssertStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssertStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.TokenLiteral() + " ")
	if as.Expression != nil {
		out.WriteString(as.Expression.String())
	}
	return out.String()
}

// IntegerLiteral represents an integer literal
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a floating point literal
type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string literal
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BooleanLiteral represents true or false
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }

// NullLiteral represents null
type NullLiteral struct {
	Token lexer.Token
}

func (n *NullLiteral) expressionNode()      {}
func (n *NullLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NullLiteral) String() string       { return "null" }

// WildcardExpression represents _ (wildcard pattern)
type WildcardExpression struct {
	Token lexer.Token
}

func (w *WildcardExpression) expressionNode()      {}
func (w *WildcardExpression) TokenLiteral() string { return w.Token.Literal }
func (w *WildcardExpression) String() string       { return "_" }

// PrefixExpression represents !expression or -expression
type PrefixExpression struct {
	Token    lexer.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents left operator right
type InfixExpression struct {
	Token    lexer.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// BlockStatement represents { statements }
// In Sango, blocks can be both statements and expressions
type BlockStatement struct {
	Token      lexer.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) expressionNode()      {} // Blocks can also be expressions in Sango
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

// IfExpression represents if (condition) { consequence } else { alternative }
type IfExpression struct {
	Token       lexer.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// FunctionLiteral represents def name(params): type = body
type FunctionLiteral struct {
	Token      lexer.Token // The 'def' token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *TypeExpression
	Body       Expression // Can be a BlockStatement or a single expression
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	if fl.Name != nil {
		out.WriteString(" " + fl.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fl.ReturnType != nil {
		out.WriteString(": ")
		out.WriteString(fl.ReturnType.String())
	}
	out.WriteString(" = ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// Parameter represents a function parameter with optional type
type Parameter struct {
	Name *Identifier
	Type *TypeExpression
}

func (p *Parameter) String() string {
	if p.Type != nil {
		return p.Name.String() + ": " + p.Type.String()
	}
	return p.Name.String()
}

// TypeExpression represents a type annotation
type TypeExpression struct {
	Token       lexer.Token
	Name        string
	Array       bool               // true if []Type
	ElementType *TypeExpression    // for array element type
	Tuple       []TypeExpression   // for tuple types (A, B, C)
	Function    *FunctionType      // for function types (A, B) -> C
	Record      *RecordType        // for record types { field: type }
}

func (te *TypeExpression) expressionNode()      {}
func (te *TypeExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TypeExpression) String() string {
	if te.Array {
		if te.ElementType != nil {
			return "[]" + te.ElementType.String()
		}
		return "[]" + te.Name
	}
	if len(te.Tuple) > 0 {
		types := []string{}
		for _, t := range te.Tuple {
			types = append(types, t.String())
		}
		return "(" + strings.Join(types, ", ") + ")"
	}
	if te.Function != nil {
		return te.Function.String()
	}
	if te.Record != nil {
		return te.Record.String()
	}
	return te.Name
}

// FunctionType represents (A, B) -> C
type FunctionType struct {
	Parameters []TypeExpression
	ReturnType *TypeExpression
}

// RecordType represents { field: type, ... }
type RecordType struct {
	Fields map[string]*TypeExpression
}

func (rt *RecordType) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	fields := []string{}
	for name, fieldType := range rt.Fields {
		fields = append(fields, name + ": " + fieldType.String())
	}
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString(" }")
	return out.String()
}

func (ft *FunctionType) String() string {
	params := []string{}
	for _, p := range ft.Parameters {
		params = append(params, p.String())
	}
	return "(" + strings.Join(params, ", ") + ") -> " + ft.ReturnType.String()
}

// CallExpression represents function(arguments)
type CallExpression struct {
	Token     lexer.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// BuiltinFunctionCall represents builtin functions like printf
type BuiltinFunctionCall struct {
	Token     lexer.Token // The function name token
	Name      string      // Function name (e.g., "printf")
	Arguments []Expression
}

func (bfc *BuiltinFunctionCall) expressionNode()      {}
func (bfc *BuiltinFunctionCall) TokenLiteral() string { return bfc.Token.Literal }
func (bfc *BuiltinFunctionCall) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range bfc.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(bfc.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// ArrayLiteral represents [1, 2, 3]
type ArrayLiteral struct {
	Token    lexer.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression represents array[index] or array[start..end]
type IndexExpression struct {
	Token lexer.Token // The '[' token
	Left  Expression
	Index Expression // Can be a single index or a RangeExpression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// RangeExpression represents start..end or start..=end
type RangeExpression struct {
	Token     lexer.Token // The '..' or '..=' token
	Start     Expression  // Can be nil for ..end
	End       Expression  // Can be nil for start..
	Inclusive bool        // true for ..=
}

func (re *RangeExpression) expressionNode()      {}
func (re *RangeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RangeExpression) String() string {
	var out bytes.Buffer
	if re.Start != nil {
		out.WriteString(re.Start.String())
	}
	if re.Inclusive {
		out.WriteString("..=")
	} else {
		out.WriteString("..")
	}
	if re.End != nil {
		out.WriteString(re.End.String())
	}
	return out.String()
}

// TupleLiteral represents (a, b, c)
type TupleLiteral struct {
	Token    lexer.Token // the '(' token
	Elements []Expression
}

func (tl *TupleLiteral) expressionNode()      {}
func (tl *TupleLiteral) TokenLiteral() string { return tl.Token.Literal }
func (tl *TupleLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range tl.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(")")
	return out.String()
}

// StructField represents a field in a struct literal: name: value
type StructField struct {
	Name  *Identifier
	Value Expression
}

func (sf *StructField) String() string {
	return sf.Name.String() + ": " + sf.Value.String()
}

// StructLiteral represents Point { x: 1, y: 2 }
type StructLiteral struct {
	Token  lexer.Token // the '{' token
	Name   *Identifier // optional struct name
	Fields []*StructField
}

func (sl *StructLiteral) expressionNode()      {}
func (sl *StructLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StructLiteral) String() string {
	var out bytes.Buffer
	if sl.Name != nil {
		out.WriteString(sl.Name.String())
		out.WriteString(" ")
	}
	out.WriteString("{ ")
	fields := []string{}
	for _, field := range sl.Fields {
		fields = append(fields, field.String())
	}
	out.WriteString(strings.Join(fields, ", "))
	out.WriteString(" }")
	return out.String()
}

// MatchExpression represents match expr { patterns }
type MatchExpression struct {
	Token lexer.Token // the 'match' token
	Value Expression
	Cases []*MatchCase
}

func (me *MatchExpression) expressionNode()      {}
func (me *MatchExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MatchExpression) String() string {
	var out bytes.Buffer
	out.WriteString(me.TokenLiteral() + " ")
	out.WriteString(me.Value.String())
	out.WriteString(" { ")
	cases := []string{}
	for _, c := range me.Cases {
		cases = append(cases, c.String())
	}
	out.WriteString(strings.Join(cases, "; "))
	out.WriteString(" }")
	return out.String()
}

// MatchCase represents pattern => expression, with optional guard
type MatchCase struct {
	Pattern Expression
	Guard   Expression // Optional guard condition (if clause)
	Value   Expression
}

func (mc *MatchCase) String() string {
	if mc.Guard != nil {
		return mc.Pattern.String() + " if " + mc.Guard.String() + " => " + mc.Value.String()
	}
	return mc.Pattern.String() + " => " + mc.Value.String()
}
