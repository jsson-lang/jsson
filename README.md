# JSSON - V0.0.6

[![JSSON Banner](https://i.postimg.cc/yx4C3YqC/og.png)](https://postimg.cc/WFnHQVb5)

**JavaScript Simplified Object Notation** - A human-friendly syntax that transpiles to JSON, YAML, TOML, and TypeScript.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![VS Code Extension](https://img.shields.io/badge/VS%20Code-Extension-blue)](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson)

---

## 📑 Table of Contents

- [What is JSSON?](#what-is-jsson)
- [Why JSSON?](#why-jsson)
- [Quick Start](#quick-start)
- [Features](#features)
- [What's New in v0.0.6](#whats-new-in-v006)
- [Multi-Format Output](#multi-format-output)
- [Examples](#examples)
- [HTTP Server](#http-server)
- [Installation](#installation)
- [Documentation](#documentation)
- [LLM-Optimized Documentation](#llm-optimized-documentation)
- [Real-World Use Cases](#real-world-use-cases)
- [VS Code Extension](#vs-code-extension)
- [Contributing](#contributing)
- [License](#license)

---

## 🎯 What is JSSON?

JSSON is a **transpiler** that converts human-friendly syntax into standard configuration formats. It eliminates the pain points of writing JSON manually while maintaining full compatibility.

**JSSON Input:**

```jsson
users [
  template { name, age }

  João, 19
  Maria, 25
  Pedro, 30
]

ports = 8080..8085
```

**JSON Output:**

```json
{
  "users": [
    { "name": "João", "age": 19 },
    { "name": "Maria", "age": 25 },
    { "name": "Pedro", "age": 30 }
  ],
  "ports": [8080, 8081, 8082, 8083, 8084, 8085]
}
```

---

## Why JSSON?

| Pain Point               | JSSON Solution                       |
| ------------------------ | ------------------------------------ |
| 😤 Quotes everywhere     | ✅ No quotes needed for keys         |
| 🐛 Trailing comma errors | ✅ No commas required                |
| 📋 Repetitive data       | ✅ Templates for arrays              |
| 🔢 Manual ranges         | ✅ Built-in range syntax (`1..100`)  |
| 📁 Scattered configs     | ✅ File includes                     |
| 🔄 Copy-paste errors     | ✅ Map transformations and variables |

---

## Quick Start

### 1. Install the CLI

```bash
# Download from releases
# Or build from source
go build -o jsson ./cmd/jsson
```

### 2. Create a `.jsson` file

```jsson
// config.jsson
app {
  name = "My App"
  version = "1.0.0"
  ports = 3000..3005
}
```

### 3. Transpile to JSON

```bash
jsson -i config.jsson > config.json

# Or start HTTP server for API access
jsson serve
```

**Output:**

```json
{
  "app": {
    "name": "My App",
    "version": "1.0.0",
    "ports": [3000, 3001, 3002, 3003, 3004, 3005]
  }
}
```

---

## Features

#### Variable Declarations 🔧

Declare reusable variables with `:=` to avoid repetition and keep your configs DRY:

```jsson
// Declare once, use everywhere
api_url := "https://api.example.com"
timeout := 5000
max_retries := 3

production {
  url = api_url + "/v1"
  timeout = timeout
  retries = max_retries
}

staging {
  url = api_url + "/dev"
  timeout = timeout * 2
  retries = max_retries
}
```

Variables support:

- **Global scope**: Declared at root level, accessible everywhere
- **Local scope**: Declared inside objects, scoped to that object
- **Shadowing**: Inner scopes can override outer variables
- **Not in output**: Variables are internal-only, never appear in final JSON

#### Nested Map Transformations 🔄

Map transformations can now be nested inside other maps for multi-level data pipelines:

```jsson
// Generate a multiplication table
table = (1..5 map (row) = (1..5 map (col) = row * col))

// Product variants (all combinations)
products = (["S", "M", "L"] map (size) = (
  ["Red", "Blue"] map (color) = {
    sku = size + "-" + color
    price = 29.99
  }
))
```

#### Nested Arrays 📦

Full support for multi-dimensional arrays:

```jsson
// 2D Matrix
matrix = [
  [ 1, 2, 3 ],
  [ 4, 5, 6 ],
  [ 7, 8, 9 ]
]

// Generated with nested maps
grid = (0..2 map (y) = (0..2 map (x) = [x, y]))
```

#### Universal Ranges 🔢

Ranges now work **everywhere** expressions are allowed:

```jsson
// Inside arrays
numbers = [ 1..5, 10..15, 20..25 ]

// Inside map arguments
data = (0..999 map (x) = { id = x, value = x * 2 })

// Large-scale generation
bigData = 0..9999  // 10,000 items!
```

#### Stable Arithmetic ➗

Division (`/`) and modulo (`%`) now work everywhere:

```jsson
hash = "uid-" + (user.id * 91 % 17)
halves = (0..10 map (x) = x / 2)
```

---

### Clean Syntax

- No quotes for keys
- No trailing commas
- Comments with `//`
- Bare identifiers

### Templates

Generate arrays from structured data:

```jsson
products [
  template { id, name, price }

  1, "Laptop", 999.99
  2, "Mouse", 29.99
  3, "Keyboard", 79.99
]
```

### Ranges

Generate sequences effortlessly:

```jsson
numbers = 1..100
ports = 8080..8090
```

### Map Transformations

Transform data with map clauses:

```jsson
users [
  template { name, age }

  map (u) = {
    name = u.name
    age = u.age
    isAdult = u.age >= 18
    category = u.age >= 18 ? "adult" : "minor"
  }

  João, 25
  Maria, 16
]
```

### 📁 File Includes

Modularize your configurations:

```jsson
include "database.jsson"
include "api-config.jsson"
```

### Arithmetic & Logic

- Operators: `+`, `-`, `*`, `/`, `%`
- Comparisons: `==`, `!=`, `>`, `<`, `>=`, `<=`
- Ternary: `condition ? true : false`

### 🚀 Streaming Support

Handle large datasets efficiently with streaming mode:

```bash
# Enable streaming for large data
jsson -i large-data.jsson --stream > output.json

# Auto-enable for ranges > 10,000 items (default)
jsson -i data.jsson > output.json

# Custom threshold
jsson -i data.jsson --stream-threshold 5000 > output.json
```

**Benefits:**

- **Memory efficient**: Reduces memory usage from ~500MB to <50MB for 100k items
- **Scalable**: Process millions of items without OOM errors
- **Automatic**: Smart threshold detection enables streaming when needed

**Perfect for:**

- Generating 100k+ records
- Large matrix/grid data (100x100+)
- Database seeding with millions of rows
- Memory-constrained environments

### Arrays & Objects

Full support for nested structures:

```jsson
config {
  methods = [ GET, POST, PUT ]
  nested {
    items = [ 1, 2, 3 ]
  }
  // Nested arrays (v0.0.5+)
  matrix = [
    [ 1, 2 ],
    [ 3, 4 ]
  ]
}
```

---

## Multi-Format Output

JSSON isn't just for JSON anymore. Transpile to your favorite format:

### YAML (Infrastructure)

```bash
jsson -i config.jsson -f yaml > config.yaml
```

### TOML (Configuration)

```bash
jsson -i config.jsson -f toml > config.toml
```

### TypeScript (Frontend)

Generates `as const` objects and type definitions!

```bash
jsson -i config.jsson -f ts > config.ts
```

---

## 🆕 What's New in v0.0.6

### Presets (Optional)

Define reusable configuration blocks with `@preset` and use them with `@use`. **Presets work with ALL output formats** (JSON, YAML, TOML, TypeScript):

```jsson
// Define reusable presets
@preset "api-defaults" {
    timeout = 30
    retries = 3
    cache = true
}

@preset "logging" {
    level = "info"
    format = "json"
}

// Use preset as-is
dev_api = @use "api-defaults"

// Use preset with overrides
prod_api = @use "api-defaults" {
    timeout = 60
    retries = 5
}

// Combine in nested structures
service {
    name = "my-service"
    api = @use "api-defaults"
    logging = @use "logging" {
        level = "debug"
    }
}
```

**JSON Output:**
```json
{
  "dev_api": { "timeout": 30, "retries": 3, "cache": true },
  "prod_api": { "timeout": 60, "retries": 5, "cache": true },
  "service": {
    "name": "my-service",
    "api": { "timeout": 30, "retries": 3, "cache": true },
    "logging": { "level": "debug", "format": "json" }
  }
}
```

**YAML Output (`-f yaml`):**
```yaml
dev_api:
  timeout: 30
  retries: 3
  cache: true
prod_api:
  timeout: 60
  retries: 5
  cache: true
service:
  name: my-service
  api:
    timeout: 30
    retries: 3
    cache: true
  logging:
    level: debug
    format: json
```

**TOML Output (`-f toml`):**
```toml
[dev_api]
timeout = 30
retries = 3
cache = true

[prod_api]
timeout = 60
retries = 5
cache = true

[service]
name = "my-service"

[service.api]
timeout = 30
retries = 3
cache = true

[service.logging]
level = "debug"
format = "json"
```
```

### Schema Validation (Optional)

Validate transpiled output against JSON Schema, YAML Schema, or TOML Schema:

```bash
# Validate output against a JSON schema
jsson -i config.jsson -schema schema.json

# Validate against YAML schema
jsson -i config.jsson -schema schema.yaml

# Validate against TOML schema
jsson -i config.jsson -schema schema.toml

# Validate only (don't output result)
jsson -i config.jsson -schema schema.json -validate-only

# Works with any output format + any schema format
jsson -i config.jsson -f yaml -schema schema.yaml
jsson -i config.jsson -f toml -schema schema.toml
```

**Supported Schema Formats:**

| Schema Format | File Extension | Description |
|---------------|----------------|-------------|
| JSON Schema | `.json` | Standard JSON Schema Draft 7 |
| YAML Schema | `.yaml`, `.yml` | YAML-formatted schema (same structure as JSON Schema) |
| TOML Schema | `.toml` | TOML-formatted schema (same structure as JSON Schema) |

**Example JSON Schema (`schema.json`):**
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "name": { "type": "string", "minLength": 1 },
    "port": { "type": "integer", "minimum": 1, "maximum": 65535 }
  },
  "required": ["name", "port"]
}
```

**Example YAML Schema (`schema.yaml`):**
```yaml
type: object
properties:
  name:
    type: string
    minLength: 1
  port:
    type: integer
    minimum: 1
    maximum: 65535
required:
  - name
  - port
```

**Example TOML Schema (`schema.toml`):**
```toml
type = "object"
required = ["name", "port"]

[properties.name]
type = "string"
minLength = 1

[properties.port]
type = "integer"
minimum = 1
maximum = 65535
```

**Validation Output:**
```
✓ Validation passed against schema
{ "name": "my-app", "port": 8080 }
✓ Compiled in 0.5ms
```

**On Error:**
```
❌ Validation failed against schema (json format):
  • $.port: Value must be >= 1, got -1
  • $.name: Missing required property
```

### HTTP Server (Unified Binary)

Run JSSON as a microservice using the `serve` subcommand:

```bash
# Start server (default port 8090)
jsson serve

# Custom port
jsson serve -port 3000

# Endpoints:
# POST /transpile        - Transpile JSSON
# POST /validate         - Validate syntax
# POST /validate-schema  - Validate against schema
# GET  /health          - Health check
# GET  /version         - Version info
```

**Example Request:**
```bash
curl -X POST http://localhost:8090/transpile \
  -H "Content-Type: application/json" \
  -d '{"source": "name = \"test\"\nport = 8080", "format": "json"}'
```

**Response:**
```json
{
  "success": true,
  "output": { "name": "test", "port": 8080 },
  "format": "json",
  "transpile_time_ms": 0.05
}
```

### Custom Format Validators

JSSON extends JSON Schema with additional format validators via `$jsson_format`:

| Format | Description | Example |
|--------|-------------|---------|
| `identifier` | Valid identifier | `myVar_123` |
| `kebab-case` | Lowercase with hyphens | `my-component` |
| `snake_case` | Lowercase with underscores | `my_variable` |
| `camelCase` | Camel case | `myVariable` |
| `PascalCase` | Pascal case | `MyClass` |
| `semver` | Semantic version | `1.0.0`, `2.1.3-beta` |
| `duration` | Duration string | `500ms`, `2h`, `1d` |
| `hex-color` | Hex color | `#fff`, `#ffffff` |
| `port` | Valid port (1-65535) | `8080` |
| `env-var` | Environment variable | `MY_VAR` |

---

## 📚 Examples

### Matrix Generation (v0.0.5)

Generate 2D matrices using nested maps:

```jsson
// Multiplication table
table = (1..10 map (row) = (1..10 map (col) = row * col))

// Coordinate grid with objects
grid = (0..99 map (y) = (0..99 map (x) = {
  x = x
  y = y
  id = y * 100 + x
}))
// Generates 10,000 coordinates!
```

### Product Variants (v0.0.5)

Generate all combinations for e-commerce:

```jsson
products = (["S", "M", "L", "XL"] map (size) = (
  ["Black", "White", "Navy", "Gray"] map (color) = {
    sku = size + "-" + color
    size = size
    color = color
    price = 29.99
    inStock = true
  }
))
// Generates 16 product variants automatically!
```

### Large-Scale Test Data (v0.0.5)

Generate thousands of records effortlessly:

```jsson
testUsers = (0..9999 map (id) = {
  id = id
  username = "user_" + id
  email = "user" + id + "@test.com"
  active = (id % 2) == 0
  tier = id < 1000 ? "bronze" : id < 5000 ? "silver" : "gold"
})
// 10,000 users generated!
```

### Geographic Coordinates

Generate millions of coordinate records:

```jsson
cityGrid [
  template { id, zone }

  map (point) = {
    id = "grid-" + point.id
    lat = -23.5505 + (point.id / 10) * 0.01
    lon = -46.6333 + (point.id % 10) * 0.01
    zone = point.zone
  }

  0..9999, "urban"  // 10,000 points!
]
```

### Kubernetes Deployments

Multi-environment infrastructure:

```jsson
deployments [
  template { app, env, replicas }

  map (d) = {
    name = d.app + "-" + d.env
    replicas = d.replicas
    image = "registry/" + d.app + ":" + (d.env == "prod" ? "stable" : "latest")
    resources = d.env == "prod" ? "high" : "low"
  }

  "api", "prod", 5
  "api", "staging", 2
  "web", "prod", 3
]
```

**More examples** in [`examples/`](./examples/)

---

## 🌐 HTTP Server

JSSON includes a built-in HTTP server via the `serve` subcommand:

### Starting the Server

```bash
# Build JSSON (single binary includes CLI + Server)
go build -o jsson ./cmd/jsson

# Start on default port (8090)
jsson serve

# Start on custom port
jsson serve -port 3000

# With CORS disabled
jsson serve -port 8090 -cors=false
```

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/transpile` | POST | Transpile JSSON to JSON/YAML/TOML/TypeScript |
| `/validate` | POST | Validate JSSON syntax |
| `/validate-schema` | POST | Validate output against schema (optional) |
| `/health` | GET | Health check |
| `/version` | GET | Version info |

### Example: Transpile

```bash
curl -X POST http://localhost:8090/transpile \
  -H "Content-Type: application/json" \
  -d '{
    "source": "name = \"my-app\"\nport = 8080\nenabled = true",
    "format": "json"
  }'
```

**Response:**

```json
{
  "success": true,
  "output": {
    "name": "my-app",
    "port": 8080,
    "enabled": true
  },
  "format": "json",
  "transpile_time_ms": 0.05
}
```

### Example: Validate with Schema

```bash
curl -X POST http://localhost:8090/validate-schema \
  -H "Content-Type: application/json" \
  -d '{
    "source": "name = \"test\"\nport = 8080",
    "schema": "{\"type\":\"object\",\"required\":[\"name\",\"port\"]}",
    "output_format": "json"
  }'
```

---

## Installation

### CLI Tool

**From Source:**

```bash
git clone https://github.com/carlosedujs/jsson
cd jsson
go build -o jsson ./cmd/jsson
```

**CLI Usage:**

```bash
# Basic transpilation
jsson -i input.jsson                    # JSON (default)
jsson -i input.jsson -f yaml            # YAML
jsson -i input.jsson -f toml            # TOML
jsson -i input.jsson -f ts              # TypeScript

# With schema validation (optional)
jsson -i input.jsson -schema schema.json

# Validate only
jsson -i input.jsson -schema schema.json -validate-only

# Streaming for large datasets
jsson -i large.jsson --stream

# Start HTTP server
jsson serve                             # Port 8090 (default)
jsson serve -port 3000                  # Custom port

# Help
jsson help
jsson version
```

### VS Code Extension

Install from the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson):

1. Open VS Code
2. Press `Ctrl+P` (or `Cmd+P` on Mac)
3. Type: `ext install carlosedujs.jsson`
4. Press Enter

Or search for "JSSON" in the Extensions view (`Ctrl+Shift+X`).

---

## 📖 Documentation

**Full documentation available at:** [JSSON Docs](https://docs.jssonlang.tech/)

- [Getting Started](https://docs.jssonlang.tech/guides/getting-started/)
- [Syntax Reference](https://docs.jssonlang.tech/reference/syntax/)
- [Templates Guide](https://docs.jssonlang.tech/guides/templates/)
- [Advanced Patterns](https://docs.jssonlang.tech/guides/advanced-patterns/)
- [CLI Usage](https://docs.jssonlang.tech/cli/usage/)

### 🤖 LLM-Optimized Documentation

JSSON provides LLM-friendly documentation at `/llms.txt/` for AI assistants and language models:

- **Clean Text Format**: All documentation converted to `.txt` without HTML/JSX noise
- **Structured Index**: Easy navigation with organized sections
- **Auto-Generated**: Updated on every build from source `.mdx` files

**Access:**

- Index: [https://docs.jssonlang.tech/llms.txt/index.txt](https://docs.jssonlang.tech/llms.txt/index.txt)
- Example: [https://docs.jssonlang.tech/llms.txt/guides/getting-started.txt](https://docs.jssonlang.tech/llms.txt/guides/getting-started.txt)

Perfect for:

- AI coding assistants (Copilot, Cursor, Claude)
- LLM context injection
- Automated documentation queries
- RAG (Retrieval-Augmented Generation) systems

---

## Real-World Use Cases

JSSON excels at:

- **🔢 Matrix Generation**: 2D/3D grids, multiplication tables, coordinate systems (v0.0.5)
- **🛍️ E-commerce**: Product variants, SKU generation, pricing matrices (v0.0.5)
- **🧪 Test Data**: Generate thousands of realistic records with patterns (v0.0.5)
- **🗺️ Geographic Data**: Generate millions of coordinate records efficiently
- **☸️ Infrastructure as Code**: Kubernetes configs, Terraform, CloudFormation
- **🌐 API Configurations**: Gateway routes, rate limiting, CORS policies
- **🌍 Internationalization**: Multi-language translation files
- **🚀 Feature Flags**: Environment-specific configuration management
- **💾 Database Seeding**: Generate realistic test data with relationships

**See real examples** in [`examples/real-world/`](./examples/real-world/):

- Geographic coordinates (10,000+ points)
- Kubernetes deployments (multi-environment)
- API gateway configuration
- i18n translations (4 languages)
- Feature flags (prod/staging/dev)
- Database seed data (200+ records)

---

## VS Code Extension

The JSSON VS Code extension provides:

- ✨ **Syntax Highlighting** for all JSSON keywords
- 🎯 **Auto-closing** brackets and braces
- 💬 **Comment support** (`//`)
- 🎨 **Color coding** for strings, numbers, operators
- 📝 **Language configuration** for better editing

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## License

MIT License - see the [LICENSE](LICENSE) file for details.

---

## Links

- **Documentation**: https://docs.jssonlang.tech/
- **VS Code Extension**: https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson
- **GitHub**: https://github.com/carlosedujs/jsson
- **Issues**: https://github.com/carlosedujs/jsson/issues

---

<div align="center">

**Made with ❤️ by [Carlos Eduardo](https://github.com/carlosedujs)**

**Enjoy coding with JSSON!** 🚀

</div>
