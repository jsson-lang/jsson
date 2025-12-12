// Package ast re-exports jsson/internal/ast for public use.
// This package provides AST node definitions for JSSON.
package ast

import "jsson/internal/ast"

// Re-export types
type (
	Node       = ast.Node
	Statement  = ast.Statement
	Expression = ast.Expression

	Program               = ast.Program
	AssignmentStatement   = ast.AssignmentStatement
	VariableDeclaration   = ast.VariableDeclaration
	IncludeStatement      = ast.IncludeStatement
	PresetStatement       = ast.PresetStatement
	PresetReference       = ast.PresetReference
	Identifier            = ast.Identifier
	IntegerLiteral        = ast.IntegerLiteral
	FloatLiteral          = ast.FloatLiteral
	BooleanLiteral        = ast.BooleanLiteral
	StringLiteral         = ast.StringLiteral
	ValidatorExpression   = ast.ValidatorExpression
	InterpolatedString    = ast.InterpolatedString
	ObjectLiteral         = ast.ObjectLiteral
	ArrayLiteral          = ast.ArrayLiteral
	MapClause             = ast.MapClause
	BinaryExpression      = ast.BinaryExpression
	MemberExpression      = ast.MemberExpression
	ArrayTemplate         = ast.ArrayTemplate
	MapExpression         = ast.MapExpression
	RangeExpression       = ast.RangeExpression
	ConditionalExpression = ast.ConditionalExpression
)
