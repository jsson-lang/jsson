package ast

import (
	"bytes"
	"jsson/internal/token"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
	// Position returns the line and column of the node in the source file
	Position() (line, col int)
}

// Statement is a node that represents a statement
type Statement interface {
	Node
	statementNode()
}

// Expression is a node that represents an expression
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

func (p *Program) Position() (int, int) {
	if len(p.Statements) > 0 {
		return p.Statements[0].Position()
	}
	return 0, 0
}

// Assignment: name = "value"
type AssignmentStatement struct {
	Token token.Token // the token.IDENT
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) Position() (int, int) { return as.Token.Line, as.Token.Column }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	return out.String()
}

// VariableDeclaration: name := value
type VariableDeclaration struct {
	Token token.Token // the ':=' token
	Name  *Identifier
	Value Expression
}

func (vd *VariableDeclaration) statementNode()       {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Literal }
func (vd *VariableDeclaration) Position() (int, int) { return vd.Token.Line, vd.Token.Column }
func (vd *VariableDeclaration) String() string {
	var out bytes.Buffer
	out.WriteString(vd.Name.String())
	out.WriteString(" := ")
	if vd.Value != nil {
		out.WriteString(vd.Value.String())
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
func (i *Identifier) Position() (int, int) { return i.Token.Line, i.Token.Column }
func (i *Identifier) String() string       { return i.Value }

// Literals
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) Position() (int, int) { return il.Token.Line, il.Token.Column }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) Position() (int, int) { return fl.Token.Line, fl.Token.Column }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) Position() (int, int) { return b.Token.Line, b.Token.Column }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }

type StringLiteral struct {
	Token     token.Token
	Value     string
	IsRaw     bool
	Validator *ValidatorExpression
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) Position() (int, int) { return sl.Token.Line, sl.Token.Column }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ValidatorExpression struct {
	Token   token.Token
	Type    string
	Pattern string
	Args    []interface{} // For validators like @int(min, max), @float(min, max)
}

func (ve *ValidatorExpression) expressionNode()      {}
func (ve *ValidatorExpression) TokenLiteral() string { return ve.Token.Literal }
func (ve *ValidatorExpression) Position() (int, int) { return ve.Token.Line, ve.Token.Column }
func (ve *ValidatorExpression) String() string {
	if ve.Pattern != "" {
		return "@" + ve.Type + "(\"" + ve.Pattern + "\")"
	}
	if len(ve.Args) > 0 {
		return "@" + ve.Type + "(...)"
	}
	return "@" + ve.Type
}

// InterpolatedString represents a raw string with {var} interpolations
// Example: """Hello {user.name}, balance: {user.balance * 10}"""
type InterpolatedString struct {
	Token token.Token
	Parts []interface{} // alternating string and Expression
}

func (is *InterpolatedString) expressionNode()      {}
func (is *InterpolatedString) TokenLiteral() string { return is.Token.Literal }
func (is *InterpolatedString) Position() (int, int) { return is.Token.Line, is.Token.Column }
func (is *InterpolatedString) String() string {
	var out bytes.Buffer
	for _, part := range is.Parts {
		switch p := part.(type) {
		case string:
			out.WriteString(p)
		case Expression:
			out.WriteString("{")
			out.WriteString(p.String())
			out.WriteString("}")
		}
	}
	return out.String()
}

// Object: { key = value }
type ObjectLiteral struct {
	Token        token.Token            // '{'
	Declarations []*VariableDeclaration // Local variables (key := value)
	Properties   map[string]Expression  // Properties (key = value)
	Keys         []string               // Para manter a ordem das chaves
}

func (o *ObjectLiteral) expressionNode()      {}
func (o *ObjectLiteral) TokenLiteral() string { return o.Token.Literal }
func (o *ObjectLiteral) Position() (int, int) { return o.Token.Line, o.Token.Column }
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
func (al *ArrayLiteral) Position() (int, int) { return al.Token.Line, al.Token.Column }
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
func (mc *MapClause) Position() (int, int) { return mc.Token.Line, mc.Token.Column }
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
func (be *BinaryExpression) Position() (int, int) { return be.Token.Line, be.Token.Column }
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
func (me *MemberExpression) Position() (int, int) { return me.Token.Line, me.Token.Column }
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
func (at *ArrayTemplate) Position() (int, int) { return at.Token.Line, at.Token.Column }
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

type MapExpression struct {
	Token    token.Token // The 'map' token
	Left     Expression  // The array being mapped
	Iterator *Identifier // The variable name (e.g. 't')
	Body     Expression  // The transformation body
}

func (me *MapExpression) expressionNode()      {}
func (me *MapExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MapExpression) Position() (int, int) { return me.Token.Line, me.Token.Column }
func (me *MapExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(me.Left.String())
	out.WriteString(" map (")
	out.WriteString(me.Iterator.String())
	out.WriteString(") = ")
	out.WriteString(me.Body.String())
	out.WriteString(")")
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
func (re *RangeExpression) Position() (int, int) { return re.Token.Line, re.Token.Column }
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
func (is *IncludeStatement) Position() (int, int) { return is.Token.Line, is.Token.Column }
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
func (ce *ConditionalExpression) Position() (int, int) { return ce.Token.Line, ce.Token.Column }
func (ce *ConditionalExpression) String() string {
	return "(" + ce.Condition.String() + " ? " + ce.Consequence.String() + " : " + ce.Alternative.String() + ")"
}

// PresetStatement: @preset "name" { ... }
// Defines a reusable configuration preset
type PresetStatement struct {
	Token token.Token    // The '@' token
	Name  *StringLiteral // Preset name
	Body  *ObjectLiteral // Preset contents
}

func (ps *PresetStatement) statementNode()       {}
func (ps *PresetStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PresetStatement) Position() (int, int) { return ps.Token.Line, ps.Token.Column }
func (ps *PresetStatement) String() string {
	var out bytes.Buffer
	out.WriteString("@preset ")
	out.WriteString(ps.Name.String())
	out.WriteString(" ")
	out.WriteString(ps.Body.String())
	return out.String()
}

// PresetReference: @"name" or @"name" { overrides }
// References and optionally extends a preset
type PresetReference struct {
	Token     token.Token    // The '@' token
	Name      *StringLiteral // Preset name to reference
	Overrides *ObjectLiteral // Optional overrides
}

func (pr *PresetReference) expressionNode()      {}
func (pr *PresetReference) TokenLiteral() string { return pr.Token.Literal }
func (pr *PresetReference) Position() (int, int) { return pr.Token.Line, pr.Token.Column }
func (pr *PresetReference) String() string {
	var out bytes.Buffer
	out.WriteString("@")
	out.WriteString(pr.Name.String())
	if pr.Overrides != nil {
		out.WriteString(" ")
		out.WriteString(pr.Overrides.String())
	}
	return out.String()
}
