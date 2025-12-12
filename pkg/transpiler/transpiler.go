// Package transpiler re-exports jsson/internal/transpiler for public use.
// This package provides the transpiler that converts JSSON to JSON/YAML/TOML/TypeScript.
package transpiler

import "jsson/internal/transpiler"

// Re-export types
type Transpiler = transpiler.Transpiler
type RangeResult = transpiler.RangeResult

// Re-export functions
var New = transpiler.New

// Note: Transpiler methods like TranspileToYAML(), TranspileToTOML(),
// and TranspileToTypeScript() are available directly on the Transpiler type.
