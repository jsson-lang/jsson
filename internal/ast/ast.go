package ast

import (
	"bytes"
	"jsson/internal/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

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

// Assignment: name = "value"
type AssignmentStatement struct {
	Token token.Token // the token.IDENT
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	return out.String()
}

// Identifier
type Identifier struct {
	Token token.Token // the token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Literals
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// Object: { key = value }
type ObjectLiteral struct {
	Token      token.Token           // '{'
	Properties map[string]Expression // Simplificado por enquanto mesmo, as chaves são strings mesmo
	Keys       []string              // Para manter a ordem das chaves
}

func (o *ObjectLiteral) expressionNode()      {}
func (o *ObjectLiteral) TokenLiteral() string { return o.Token.Literal }
func (o *ObjectLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, key := range o.Keys {
		out.WriteString(key)
		if val := o.Properties[key]; val != nil {
			out.WriteString(" = ")
			out.WriteString(val.String())
		}
		out.WriteString(", ")
	}
	out.WriteString(" }")
	return out.String()
}

// Array: [ 1, 2, 3 ]
type ArrayLiteral struct {
	Token    token.Token // '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, el := range al.Elements {
		out.WriteString(el.String())
		if i < len(al.Elements)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("]")
	return out.String()
}

// MapClause: map (x) = { ... }
type MapClause struct {
	Token token.Token    // "map"
	Param *Identifier    // "x"
	Body  *ObjectLiteral // "{ ... }"
}

func (mc *MapClause) expressionNode()      {}
func (mc *MapClause) TokenLiteral() string { return mc.Token.Literal }
func (mc *MapClause) String() string {
	return "map (" + mc.Param.String() + ") = " + mc.Body.String()
}

// BinaryExpression: x + y
type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) TokenLiteral() string { return be.Operator }
func (be *BinaryExpression) String() string {
	return "(" + be.Left.String() + " " + be.Operator + " " + be.Right.String() + ")"
}

// MemberExpression: item.path
type MemberExpression struct {
	Token    token.Token // The '.' token
	Left     Expression
	Property *Identifier
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
	return me.Left.String() + "." + me.Property.String()
}

// ArrayTemplate: users [ template { name, age } ... ]
type ArrayTemplate struct {
	Token    token.Token // The identifier token before '['
	Name     *Identifier
	Template *ObjectLiteral // The template definition
	Map      *MapClause     // Optional map clause
	Rows     [][]Expression // The data rows
}

func (at *ArrayTemplate) expressionNode()      {}
func (at *ArrayTemplate) TokenLiteral() string { return at.Token.Literal }
func (at *ArrayTemplate) String() string {
	var out bytes.Buffer
	if at.Name != nil {
		out.WriteString(at.Name.String())
		out.WriteString(" ")
	}
	out.WriteString("[ template ")
	if at.Template != nil {
		out.WriteString(at.Template.String())
	}
	if at.Map != nil {
		out.WriteString(" ")
		out.WriteString(at.Map.String())
	}
	out.WriteString(" ... ]")
	return out.String()
}

// RangeExpression: start .. end [ step N ]
type RangeExpression struct {
	Token token.Token
	Start Expression
	End   Expression
	Step  Expression // optional
}

func (re *RangeExpression) expressionNode()      {}
func (re *RangeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RangeExpression) String() string {
	if re.Step != nil {
		return re.Start.String() + ".." + re.End.String() + " step " + re.Step.String()
	}
	return re.Start.String() + ".." + re.End.String()
}

// IncludeStatement: include "file.jsson"
type IncludeStatement struct {
	Token token.Token // the 'include' token
	Path  *StringLiteral
}

func (is *IncludeStatement) statementNode()       {}
func (is *IncludeStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncludeStatement) String() string {
	return "include " + is.Path.String()
}

// ConditionalExpression: condition ? consequence : alternative
type ConditionalExpression struct {
	Token       token.Token // The '?' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ce *ConditionalExpression) expressionNode()      {}
func (ce *ConditionalExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ConditionalExpression) String() string {
	return "(" + ce.Condition.String() + " ? " + ce.Consequence.String() + " : " + ce.Alternative.String() + ")"
}
