// Package lexer re-exports jsson/internal/lexer for public use.
// This package provides the lexical analyzer for JSSON.
package lexer

import "jsson/internal/lexer"

// Re-export types
type Lexer = lexer.Lexer

// Re-export functions
var New = lexer.New
