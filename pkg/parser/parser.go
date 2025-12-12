// Package parser re-exports jsson/internal/parser for public use.
// This package provides the parser for JSSON syntax.
package parser

import "jsson/internal/parser"

// Re-export types
type Parser = parser.Parser

// Re-export functions
var New = parser.New
