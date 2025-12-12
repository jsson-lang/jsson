// Package token re-exports jsson/internal/token for public use.
// This package provides token definitions for the JSSON lexer.
package token

import "jsson/internal/token"

// Re-export types
type Token = token.Token
type TokenType = token.TokenType

// Re-export constants
const (
	// Literals
	ILLEGAL = token.ILLEGAL
	EOF     = token.EOF
	IDENT   = token.IDENT
	INT     = token.INT
	FLOAT   = token.FLOAT
	STRING  = token.STRING

	// Operators
	ASSIGN   = token.ASSIGN
	DECLARE  = token.DECLARE
	PLUS     = token.PLUS
	MINUS    = token.MINUS
	ASTERISK = token.ASTERISK
	SLASH    = token.SLASH
	MODULO   = token.MODULO

	// Comparison
	EQ  = token.EQ
	NEQ = token.NEQ
	LT  = token.LT
	GT  = token.GT
	LTE = token.LTE
	GTE = token.GTE

	// Logical
	LAND = token.LAND
	LOR  = token.LOR

	// Delimiters
	COMMA    = token.COMMA
	COLON    = token.COLON
	LBRACE   = token.LBRACE
	RBRACE   = token.RBRACE
	LBRACKET = token.LBRACKET
	RBRACKET = token.RBRACKET
	LPAREN   = token.LPAREN
	RPAREN   = token.RPAREN
	DOT      = token.DOT
	RANGE    = token.RANGE
	QUESTION = token.QUESTION
	AT       = token.AT

	// Keywords
	INCLUDE  = token.INCLUDE
	TEMPLATE = token.TEMPLATE
	MAP      = token.MAP
	STEP     = token.STEP
	TRUE     = token.TRUE
	FALSE    = token.FALSE
	PRESET   = token.PRESET
	USE      = token.USE

	// Validators
	UUID     = token.UUID
	EMAIL    = token.EMAIL
	URL      = token.URL
	IPV4     = token.IPV4
	IPV6     = token.IPV6
	FILEPATH = token.FILEPATH
	DATE     = token.DATE
	DATETIME = token.DATETIME
	REGEX    = token.REGEX

	// String types
	RAWSTRING   = token.RAWSTRING
	TEMPLATESTR = token.TEMPLATESTR
)
