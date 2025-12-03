# JSSON Examples

This directory contains examples demonstrating JSSON features and use cases.

## Core vs Optional Features

**Core Features** (always available):
- Transpilation to JSON, YAML, TOML, TypeScript
- Variables and expressions
- Ranges and map expressions  
- Array templates
- String interpolation
- Conditionals
- Includes

**Optional Features** (opt-in):
- **Presets** - Use `@preset` syntax in your JSSON files when you need reusable configurations
- **Schema Validation** - Use `-schema` flag when you need to validate output against a JSON/YAML schema

## Directory Structure

```
examples/
├── schemas/                 # Generic validation schemas
│   ├── api-config.schema.json
│   └── database.schema.yaml
├── validation/              # Examples with schema validation
│   ├── api_config.jsson
│   └── database_config.jsson
├── use-cases/               # Domain-specific use cases
│   └── chaos-factory/       # Gaming macro configuration example
│       ├── chaos_macros.jsson
│       └── chaos_macro.schema.json
└── *.jsson                  # General feature examples
```

## Running Examples

### Basic Transpilation

```bash
# Transpile to JSON (default)
jsson -i demo.jsson

# Transpile to YAML
jsson -i demo.jsson -f yaml

# Transpile to TOML
jsson -i demo.jsson -f toml

# Transpile to TypeScript types
jsson -i demo.jsson -f typescript
```

### With Schema Validation

```bash
# Validate output against a JSON Schema
jsson -i validation/api_config.jsson -schema schemas/api-config.schema.json

# Validate YAML output against YAML Schema
jsson -i validation/database_config.jsson -f yaml -schema schemas/database.schema.yaml

# Validate only (don't output result)
jsson -i validation/api_config.jsson -schema schemas/api-config.schema.json -validate-only
```

## Features Demonstrated

### 1. Presets (`presets_example.jsson`)
Reusable configuration blocks that can be referenced and overridden.

```jsson
@preset "defaults" {
    timeout = 30
    retries = 3
}

config = @use "defaults" {
    timeout = 60  // Override specific values
}
```

### 2. Ranges and Generation (`ranges.jsson`)
Generate arrays from numeric or string ranges.

```jsson
// Numeric range
numbers = [1..10]

// With step
evens = [2..20..2]

// String range (for IPs, etc)
servers = ["server-01".."server-10"]
```

### 3. Map Expressions (`map.jsson`)
Transform arrays with map expressions.

```jsson
users = [1..5] -> (id) {
    id = id
    name = "User {id}"
    email = "user{id}@example.com"
}
```

### 4. Array Templates (`template.jsson`)
Generate structured data from tabular input.

```jsson
servers = [| (name, ip, port)
    "web-1",    "10.0.0.1", 80
    "web-2",    "10.0.0.2", 80
    "db-1",     "10.0.1.1", 5432
|]
```

### 5. Conditionals (`comparison-test.jsson`)
Ternary expressions for conditional values.

```jsson
env = "production"
debug = env == "development" ? true : false
```

### 6. String Interpolation
Template strings with variable substitution.

```jsson
name = "World"
greeting = "Hello, {name}!"       // Simple interpolation
template = `Hello, ${name}!`      // Template string style
```

## Schema Validation

JSSON supports validating transpiled output against JSON Schema (draft-07) or YAML Schema.

### Supported Validations
- Type checking (string, number, integer, boolean, array, object, null)
- Required properties
- Enum values
- Pattern matching (regex)
- Min/max for numbers
- MinLength/maxLength for strings
- MinItems/maxItems for arrays
- Unique items
- Additional properties control
- allOf, anyOf, oneOf, not
- if/then/else
- $ref references

### Custom JSSON Formats
JSSON extends JSON Schema with additional format validators:

- `identifier` - Valid identifier (starts with letter, alphanumeric + underscore)
- `kebab-case` - Lowercase with hyphens
- `snake_case` - Lowercase with underscores
- `camelCase` - Camel case format
- `PascalCase` - Pascal case format
- `semver` - Semantic versioning
- `duration` - Duration string (e.g., "500ms", "2h")
- `hex-color` - Hex color codes
- `port` - Valid port number
- `env-var` - Environment variable name format

Example in schema:
```json
{
    "type": "string",
    "$jsson_format": "semver"
}
```
