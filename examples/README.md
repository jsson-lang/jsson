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
├── current/                 # ✅ Working examples (current version)
│   ├── simple_config.jsson
│   ├── invalid_config.jsson
│   └── README.md
├── planned/                 # 🚧 Roadmap examples (future syntax)
│   ├── database_config.jsson
│   ├── api_config.jsson
│   └── README.md
├── schemas/                 # Generic validation schemas
│   ├── api-config.schema.json
│   └── database.schema.yaml
├── real-world/              # Real-world use cases
├── use-cases/               # Domain-specific examples
└── *.jsson                  # General feature examples
```

> **⚠️ Important**: Examples in `planned/` use **unimplemented syntax** from the roadmap.
> They will **not compile** with the current version. See `ROADMAP.md` for details.

## Running Examples

### Current Features (Working Now)

```bash
# Simple configuration example
jsson -i current/simple_config.jsson

# Different output formats
jsson -i current/simple_config.jsson -f yaml
jsson -i current/simple_config.jsson -f toml
jsson -i current/simple_config.jsson -f ts
```

### Planned Features (Roadmap)

⚠️ **These will NOT work** - they use unimplemented syntax:

```bash
# ❌ Will fail - uses list comprehension and @use
jsson -i planned/database_config.jsson

# ❌ Will fail - uses list comprehension and flatten
jsson -i planned/api_config.jsson
```

See `planned/README.md` and `ROADMAP.md` for details on planned features.

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
