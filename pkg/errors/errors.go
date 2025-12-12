// Package errors re-exports jsson/internal/errors for public use.
// This package provides error formatting and messages for JSSON.
package errors

import "jsson/internal/errors"

// Re-export functions
var (
	FormatContext      = errors.FormatContext
	LexerError         = errors.LexerError
	ParserError        = errors.ParserError
	TranspilerError    = errors.TranspilerError
	UnterminatedString = errors.UnterminatedString
	IllegalCharacter   = errors.IllegalCharacter
	ExpectedToken      = errors.ExpectedToken
)
