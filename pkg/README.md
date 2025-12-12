# JSSON Public API (`pkg/`)

This directory contains the **public API** for using JSSON as a Go library in your projects.

## Architecture

The `pkg/` directory **re-exports** the internal packages from `internal/`. This ensures:
- Single source of truth (code lives in `internal/`)
- Stable public API for external consumers
- No code duplication or divergence

## Purpose

The `pkg/` directory provides a stable, importable API for external projects that want to:
- Parse JSSON syntax programmatically
- Transpile JSSON to JSON/YAML/TOML/TypeScript
- Validate JSSON documents
- Build tools on top of JSSON

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "jsson/pkg/lexer"
    "jsson/pkg/parser"
    "jsson/pkg/transpiler"
)

func main() {
    // JSSON source code
    source := `
        user {
            name = "Carlos"
            age = 25
            active = yes
            email = @email
        }
    `
    
    // Parse
    l := lexer.New(source)
    p := parser.New(l)
    program := p.ParseProgram()
    
    if len(p.Errors()) > 0 {
        for _, err := range p.Errors() {
            fmt.Println("Error:", err)
        }
        return
    }
    
    // Transpile to JSON
    tr := transpiler.New(program, "", "keep", "")
    output, err := tr.Transpile()
    if err != nil {
        fmt.Println("Transpile error:", err)
        return
    }
    
    fmt.Println(string(output))
}
```

### Transpile to Different Formats

```go
// JSON (default)
jsonOutput, _ := tr.Transpile()

// YAML
yamlOutput, _ := tr.TranspileToYAML()

// TOML
tomlOutput, _ := tr.TranspileToTOML()

// TypeScript
tsOutput, _ := tr.TranspileToTypeScript()
```

### With Schema Validation

```go
import "jsson/pkg/validator"

v := validator.New()
schema, _, err := v.LoadSchemaAuto(schemaContent)
if err != nil {
    // handle error
}

result := v.Validate(output, schema, "json")
if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Validation error: %s\n", err.Message)
    }
}
```

## Available Packages

- **`pkg/lexer`** - Tokenize JSSON source code
- **`pkg/parser`** - Parse tokens into AST
- **`pkg/transpiler`** - Convert AST to JSON/YAML/TOML/TypeScript
- **`pkg/validator`** - Validate output against schemas
- **`pkg/ast`** - Abstract Syntax Tree definitions
- **`pkg/token`** - Token type definitions
- **`pkg/errors`** - Error handling utilities

## Features Supported (v0.0.6)

- Boolean literals: `yes`/`no`, `on`/`off`
- Validators: `@uuid`, `@email`, `@url`, `@ipv4`, `@ipv6`, `@filepath`, `@date`, `@datetime`, `@regex`
- Presets: `@preset "name" { ... }` and `@use "name"`
- Ranges: `1..100`
- Maps: `(1..10 map (x) = x * 2)`
- Templates
- Variables: `x := 10`
- Include system
- Streaming mode for large datasets
- Multi-format output (JSON, YAML, TOML, TypeScript)
- Schema validation

## Versioning

The `pkg/` API follows the main JSSON version. This is currently **v0.0.6**.

Breaking changes will be documented in the main CHANGELOG.

## ⚠️ Note

The `pkg/` directory is synchronized with `internal/` but may lag behind by one minor version for stability purposes. For the absolute latest features, use the command-line tool directly.

## Documentation

For full JSSON language documentation, see the main [README.md](https://github.com/jsson-lang/jsson#readme).

For API documentation, run:
```bash
go doc jsson/pkg/parser
go doc jsson/pkg/transpiler
```

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass: `go test ./pkg/...`
2. Code is formatted: `go fmt ./pkg/...`
3. API remains backwards compatible when possible

## License

MIT License - See [LICENSE](../LICENSE) for details.
