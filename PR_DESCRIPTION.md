## 🔧 Bug Fixes & Improvements for JSSON v0.0.6

### 🐛 Critical Bug Fix: @preset/@use Support for All Output Formats

**Problem:** The `@preset` and `@use` syntax only worked when transpiling to JSON. When using YAML, TOML, or TypeScript output formats, the transpiler threw an error: `preset "X" not found`.

**Root Cause:** The `PresetStatement` case was missing from the YAML, TOML, and TypeScript transpilers. Only `transpiler.go` (JSON) had the code to register presets in the `presetTable`.

**Solution:** Added the `PresetStatement` case to all transpilers:
- `internal/transpiler/yaml.go`
- `internal/transpiler/toml.go`
- `internal/transpiler/typescript.go`

**Before:**
```bash
$ jsson -i config.jsson -f yaml
# Error: preset "s" not found — define it with @preset "s" { ... }
```

**After:**
```bash
$ jsson -i config.jsson -f yaml
# ✓ Works correctly with all presets expanded
```

---

### ✨ New Feature: TOML & TypeScript Schema Validation

**Added:** Support for TOML and TypeScript-formatted schemas in the validator.

**Implementation:**
- `LoadSchemaFromTOML()` - Parse TOML schema strings
- `LoadSchemaFromTypeScript()` - Parse TypeScript object declarations as schemas
- `LoadSchemaAuto()` - Enhanced to detect JSON, YAML, TOML, and TypeScript formats automatically
- HTTP API - Now accepts `schema_format: "toml"` and `schema_format: "typescript"` parameters

**Example TOML Schema:**

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

---

### 🧹 Cleanup: Removed Duplicate `pkg/` Directory

**Problem:** The repository had both `internal/` and `pkg/` directories with duplicated code. Only `internal/` was being used in imports.

**Solution:** Removed the unused `pkg/` directory to align with the upstream repository structure and avoid confusion.

---

### 📁 Files Changed

| File | Change |
|------|--------|
| `internal/transpiler/yaml.go` | Added `PresetStatement` case |
| `internal/transpiler/toml.go` | Added `PresetStatement` case |
| `internal/transpiler/typescript.go` | Added `PresetStatement` case |
| `internal/validator/validator.go` | Added TOML and TypeScript schema support |
| `cmd/jsson/main.go` | Added TOML and TypeScript schema format to HTTP API |
| `Dockerfile.chaos` | Removed reference to `pkg/` |
| `README.md` | Updated documentation with TOML schema examples |

---

### ✅ Testing Results

**Presets (@preset/@use)** now work with all output formats:

| Output Format | @preset/@use Support |
|---------------|----------------------|
| JSON | ✅ Working |
| YAML | ✅ **Fixed** |
| TOML | ✅ **Fixed** |
| TypeScript | ✅ **Fixed** |

**Schema Validation** now supports all combinations of input/output formats:

| Input Format | Output Support |
|--------------|-----------------|
| JSON Schema | ✅ JSON, YAML, TOML, TypeScript |
| YAML Schema | ✅ JSON, YAML, TOML, TypeScript |
| TOML Schema | ✅ JSON, YAML, TOML, TypeScript |
| TypeScript Schema | ✅ JSON, YAML, TOML, TypeScript |

**How it works:**

- **Presets:** Can be defined and used (`@preset` / `@use`) with all output formats (JSON, YAML, TOML, TypeScript)
- **Schema Validation:** Accepts schemas in any format (JSON/YAML/TOML/TypeScript) and validates output in any format because all formats share the same underlying data structure

---

### 🚀 Deployment

Changes have been deployed and tested on production server with the Discord bot successfully using presets with YAML output.
