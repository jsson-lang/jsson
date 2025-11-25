# JSSON - V0.0.3

[![JSSON Banner](https://i.postimg.cc/xTPzcqNQ/og.png)](https://postimg.cc/zVVVBQzk2)

**JSON Simplified Object Notation** - A human-friendly syntax that transpiles to JSON.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![VS Code Extension](https://img.shields.io/badge/VS%20Code-Extension-blue)](https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson)

---

## 📑 Table of Contents

- [What is JSSON?](#what-is-jsson)
- [Why JSSON?](#why-jsson)
- [Quick Start](#quick-start)
- [Features](#features)
- [Examples](#examples)
- [Installation](#installation)
- [Documentation](#documentation)
- [Real-World Use Cases](#real-world-use-cases)
- [VS Code Extension](#vs-code-extension)
- [Contributing](#contributing)
- [License](#license)

---

## 🎯 What is JSSON?

JSSON is a **transpiler** that converts human-friendly syntax into standard JSON. It eliminates the pain points of writing JSON manually while maintaining full JSON compatibility.

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

| Pain Point               | JSSON Solution                      |
| ------------------------ | ----------------------------------- |
| 😤 Quotes everywhere     | ✅ No quotes needed for keys        |
| 🐛 Trailing comma errors | ✅ No commas required               |
| 📋 Repetitive data       | ✅ Templates for arrays             |
| 🔢 Manual ranges         | ✅ Built-in range syntax (`1..100`) |
| 📁 Scattered configs     | ✅ File includes                    |
| 🔄 Copy-paste errors     | ✅ Map transformations              |

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
jsson -i config.jsson
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

### Arrays & Objects

Full support for nested structures:

```jsson
config {
  methods = ["GET", "POST", "PUT"]
  nested {
    items = [1, 2, 3]
  }
}
```

---

## 📚 Examples

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

**More examples** in [`examples/real-world/`](./examples/real-world/)

---

## Installation

### CLI Tool

**From Source:**

```bash
git clone https://github.com/carlosedujs/jsson
cd jsson
go build -o jsson ./cmd/jsson
```

**Usage:**

```bash
jsson -i input.jsson              # Output to stdout
jsson -i input.jsson -o output.json  # Output to file
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

**Full documentation available at:** [JSSON Docs](https://jsson-docs.vercel.app/)

- [Getting Started](https://jsson-docs.vercel.app/guides/getting-started/)
- [Syntax Reference](https://jsson-docs.vercel.app/reference/syntax/)
- [Templates Guide](https://jsson-docs.vercel.app/guides/templates/)
- [Advanced Patterns](https://jsson-docs.vercel.app/guides/advanced-patterns/)
- [CLI Usage](https://jsson-docs.vercel.app/cli/usage/)

---

## Real-World Use Cases

JSSON excels at:

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

- **Documentation**: https://jsson-docs.vercel.app/
- **VS Code Extension**: https://marketplace.visualstudio.com/items?itemName=carlosedujs.jsson
- **GitHub**: https://github.com/carlosedujs/jsson
- **Issues**: https://github.com/carlosedujs/jsson/issues

---

<div align="center">

**Made with ❤️ by [Carlos Eduardo](https://github.com/carlosedujs)**

**Enjoy coding with JSSON!** 🚀

</div>
