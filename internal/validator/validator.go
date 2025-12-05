/*
JSSON Validator
===============

Validates transpiled JSSON output against schemas in multiple formats.
Supports JSON Schema, YAML Schema, and custom JSSON validation rules.

Features:
  - JSON Schema Draft 7 validation
  - YAML Schema validation
  - TOML structure validation
  - Custom JSSON rules validation
  - Multi-format schema support
  - Detailed error reporting

Usage:
  validator := validator.New()
  result := validator.Validate(transpiledData, schema, "json")
*/
package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// ============================================================================
// Types
// ============================================================================

// ValidationResult contains the result of a validation
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
	Format   string            `json:"format"`
	Schema   string            `json:"schema_type"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Path       string `json:"path"`
	Message    string `json:"message"`
	SchemaPath string `json:"schema_path,omitempty"`
	Value      string `json:"value,omitempty"`
	Expected   string `json:"expected,omitempty"`
}

// Schema represents a validation schema
type Schema struct {
	Type       string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Properties map[string]*Schema     `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items      *Schema                `json:"items,omitempty" yaml:"items,omitempty"`
	Required   []string               `json:"required,omitempty" yaml:"required,omitempty"`
	Enum       []interface{}          `json:"enum,omitempty" yaml:"enum,omitempty"`
	Pattern    string                 `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Minimum    *float64               `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum    *float64               `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MinLength  *int                   `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength  *int                   `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinItems   *int                   `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems   *int                   `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	Const      interface{}            `json:"const,omitempty" yaml:"const,omitempty"`
	Default    interface{}            `json:"default,omitempty" yaml:"default,omitempty"`
	OneOf      []*Schema              `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf      []*Schema              `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	AllOf      []*Schema              `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	Not        *Schema                `json:"not,omitempty" yaml:"not,omitempty"`
	If         *Schema                `json:"if,omitempty" yaml:"if,omitempty"`
	Then       *Schema                `json:"then,omitempty" yaml:"then,omitempty"`
	Else       *Schema                `json:"else,omitempty" yaml:"else,omitempty"`
	Title      string                 `json:"title,omitempty" yaml:"title,omitempty"`
	Desc       string                 `json:"description,omitempty" yaml:"description,omitempty"`
	AdditionalProperties interface{}  `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	PatternProperties    map[string]*Schema `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	UniqueItems          *bool        `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	// JSSON-specific extensions
	JSSonPreset string `json:"$jsson_preset,omitempty" yaml:"$jsson_preset,omitempty"`
	JSSonFormat string `json:"$jsson_format,omitempty" yaml:"$jsson_format,omitempty"`
}

// Validator performs validation of transpiled output
type Validator struct {
	schemas     map[string]*Schema
	customRules map[string]ValidationRule
}

// ValidationRule is a custom validation function
type ValidationRule func(value interface{}, path string) []ValidationError

// ============================================================================
// Constructor
// ============================================================================

// New creates a new Validator instance
func New() *Validator {
	return &Validator{
		schemas:     make(map[string]*Schema),
		customRules: make(map[string]ValidationRule),
	}
}

// ============================================================================
// Schema Loading
// ============================================================================

// LoadSchemaFromJSON loads a schema from JSON string
func (v *Validator) LoadSchemaFromJSON(schemaJSON string) (*Schema, error) {
	var schema Schema
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}
	return &schema, nil
}

// LoadSchemaFromYAML loads a schema from YAML string
func (v *Validator) LoadSchemaFromYAML(schemaYAML string) (*Schema, error) {
	var schema Schema
	if err := yaml.Unmarshal([]byte(schemaYAML), &schema); err != nil {
		return nil, fmt.Errorf("failed to parse YAML schema: %w", err)
	}
	return &schema, nil
}

// LoadSchemaFromTOML loads a schema from TOML string
func (v *Validator) LoadSchemaFromTOML(schemaTOML string) (*Schema, error) {
	var schema Schema
	if _, err := toml.Decode(schemaTOML, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse TOML schema: %w", err)
	}
	return &schema, nil
}

// LoadSchemaAuto automatically detects format and loads schema
func (v *Validator) LoadSchemaAuto(schemaContent string) (*Schema, string, error) {
	// Try JSON first
	schemaContent = strings.TrimSpace(schemaContent)
	if strings.HasPrefix(schemaContent, "{") {
		schema, err := v.LoadSchemaFromJSON(schemaContent)
		if err == nil {
			return schema, "json", nil
		}
	}

	// Try TOML (check for typical TOML patterns: [section] or key = value without colon)
	if strings.Contains(schemaContent, "[") && strings.Contains(schemaContent, "]") &&
		!strings.Contains(schemaContent, ": ") {
		schema, err := v.LoadSchemaFromTOML(schemaContent)
		if err == nil {
			return schema, "toml", nil
		}
	}

	// Try YAML (default fallback since YAML is a superset of JSON)
	schema, err := v.LoadSchemaFromYAML(schemaContent)
	if err == nil {
		return schema, "yaml", nil
	}

	return nil, "", fmt.Errorf("could not parse schema as JSON, YAML, or TOML")
}

// RegisterSchema registers a named schema for reuse
func (v *Validator) RegisterSchema(name string, schema *Schema) {
	v.schemas[name] = schema
}

// RegisterCustomRule registers a custom validation rule
func (v *Validator) RegisterCustomRule(name string, rule ValidationRule) {
	v.customRules[name] = rule
}

// ============================================================================
// Validation Methods
// ============================================================================

// Validate validates data against a schema
// dataFormat can be: json, yaml, toml, typescript
func (v *Validator) Validate(data []byte, schema *Schema, dataFormat string) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Format: dataFormat,
		Schema: "json-schema",
	}

	// Parse data based on format
	var parsedData interface{}
	var err error

	switch strings.ToLower(dataFormat) {
	case "json":
		err = json.Unmarshal(data, &parsedData)
	case "yaml":
		err = yaml.Unmarshal(data, &parsedData)
	case "toml":
		// TOML needs special handling - parse as map
		parsedData, err = v.parseTOML(data)
	case "typescript", "ts":
		// TypeScript - extract the type definition and validate structure
		parsedData, err = v.parseTypeScript(data)
	default:
		// Try JSON as default
		err = json.Unmarshal(data, &parsedData)
	}

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    "$",
			Message: fmt.Sprintf("Failed to parse %s data: %v", dataFormat, err),
		})
		return result
	}

	// Normalize YAML maps (map[string]interface{} vs map[interface{}]interface{})
	parsedData = v.normalizeData(parsedData)

	// Validate against schema
	errors := v.validateValue(parsedData, schema, "$", "")
	if len(errors) > 0 {
		result.Valid = false
		result.Errors = errors
	}

	return result
}

// ValidateJSON is a convenience method for JSON validation
func (v *Validator) ValidateJSON(jsonData []byte, schema *Schema) *ValidationResult {
	return v.Validate(jsonData, schema, "json")
}

// ValidateYAML is a convenience method for YAML validation
func (v *Validator) ValidateYAML(yamlData []byte, schema *Schema) *ValidationResult {
	return v.Validate(yamlData, schema, "yaml")
}

// ValidateTOML is a convenience method for TOML validation
func (v *Validator) ValidateTOML(tomlData []byte, schema *Schema) *ValidationResult {
	return v.Validate(tomlData, schema, "toml")
}

// ValidateTypeScript is a convenience method for TypeScript validation
func (v *Validator) ValidateTypeScript(tsData []byte, schema *Schema) *ValidationResult {
	return v.Validate(tsData, schema, "typescript")
}

// ============================================================================
// Core Validation Logic
// ============================================================================

func (v *Validator) validateValue(value interface{}, schema *Schema, path string, schemaPath string) []ValidationError {
	var errors []ValidationError

	if schema == nil {
		return errors
	}

	// Handle allOf
	if len(schema.AllOf) > 0 {
		for i, subSchema := range schema.AllOf {
			subErrors := v.validateValue(value, subSchema, path, fmt.Sprintf("%s/allOf/%d", schemaPath, i))
			errors = append(errors, subErrors...)
		}
	}

	// Handle anyOf
	if len(schema.AnyOf) > 0 {
		anyValid := false
		for i, subSchema := range schema.AnyOf {
			subErrors := v.validateValue(value, subSchema, path, fmt.Sprintf("%s/anyOf/%d", schemaPath, i))
			if len(subErrors) == 0 {
				anyValid = true
				break
			}
		}
		if !anyValid {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Value does not match any of the schemas in anyOf",
				SchemaPath: schemaPath + "/anyOf",
			})
		}
	}

	// Handle oneOf
	if len(schema.OneOf) > 0 {
		matchCount := 0
		for _, subSchema := range schema.OneOf {
			subErrors := v.validateValue(value, subSchema, path, schemaPath)
			if len(subErrors) == 0 {
				matchCount++
			}
		}
		if matchCount != 1 {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    fmt.Sprintf("Value must match exactly one schema in oneOf, but matched %d", matchCount),
				SchemaPath: schemaPath + "/oneOf",
			})
		}
	}

	// Handle not
	if schema.Not != nil {
		notErrors := v.validateValue(value, schema.Not, path, schemaPath+"/not")
		if len(notErrors) == 0 {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Value should not match the schema in 'not'",
				SchemaPath: schemaPath + "/not",
			})
		}
	}

	// Handle if/then/else
	if schema.If != nil {
		ifErrors := v.validateValue(value, schema.If, path, schemaPath+"/if")
		if len(ifErrors) == 0 && schema.Then != nil {
			thenErrors := v.validateValue(value, schema.Then, path, schemaPath+"/then")
			errors = append(errors, thenErrors...)
		} else if len(ifErrors) > 0 && schema.Else != nil {
			elseErrors := v.validateValue(value, schema.Else, path, schemaPath+"/else")
			errors = append(errors, elseErrors...)
		}
	}

	// Handle const
	if schema.Const != nil {
		if !v.deepEqual(value, schema.Const) {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    fmt.Sprintf("Value must be exactly %v", schema.Const),
				SchemaPath: schemaPath + "/const",
				Value:      fmt.Sprintf("%v", value),
				Expected:   fmt.Sprintf("%v", schema.Const),
			})
		}
	}

	// Handle enum
	if len(schema.Enum) > 0 {
		found := false
		for _, enumVal := range schema.Enum {
			if v.deepEqual(value, enumVal) {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    fmt.Sprintf("Value must be one of: %v", schema.Enum),
				SchemaPath: schemaPath + "/enum",
				Value:      fmt.Sprintf("%v", value),
			})
		}
	}

	// Type validation
	if schema.Type != "" {
		typeErrors := v.validateType(value, schema, path, schemaPath)
		errors = append(errors, typeErrors...)
	}

	return errors
}

func (v *Validator) validateType(value interface{}, schema *Schema, path string, schemaPath string) []ValidationError {
	var errors []ValidationError

	actualType := v.getType(value)

	// Handle multiple types (e.g., ["string", "null"])
	types := strings.Split(schema.Type, ",")
	typeMatch := false
	for _, t := range types {
		t = strings.TrimSpace(t)
		if t == actualType || (t == "null" && value == nil) || (t == "integer" && actualType == "number" && v.isInteger(value)) {
			typeMatch = true
			break
		}
	}

	if !typeMatch {
		errors = append(errors, ValidationError{
			Path:       path,
			Message:    fmt.Sprintf("Expected type %s but got %s", schema.Type, actualType),
			SchemaPath: schemaPath + "/type",
			Value:      fmt.Sprintf("%v", value),
			Expected:   schema.Type,
		})
		return errors
	}

	// Type-specific validations
	switch actualType {
	case "object":
		objErrors := v.validateObject(value.(map[string]interface{}), schema, path, schemaPath)
		errors = append(errors, objErrors...)
	case "array":
		arrErrors := v.validateArray(value.([]interface{}), schema, path, schemaPath)
		errors = append(errors, arrErrors...)
	case "string":
		strErrors := v.validateString(value.(string), schema, path, schemaPath)
		errors = append(errors, strErrors...)
	case "number":
		numErrors := v.validateNumber(value, schema, path, schemaPath)
		errors = append(errors, numErrors...)
	}

	return errors
}

func (v *Validator) validateObject(obj map[string]interface{}, schema *Schema, path string, schemaPath string) []ValidationError {
	var errors []ValidationError

	// Check required properties
	for _, req := range schema.Required {
		if _, exists := obj[req]; !exists {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    fmt.Sprintf("Missing required property: %s", req),
				SchemaPath: schemaPath + "/required",
				Expected:   req,
			})
		}
	}

	// Validate each property
	for key, val := range obj {
		propPath := path + "." + key
		propSchemaPath := schemaPath + "/properties/" + key

		// Check if property is defined in schema
		if propSchema, exists := schema.Properties[key]; exists {
			propErrors := v.validateValue(val, propSchema, propPath, propSchemaPath)
			errors = append(errors, propErrors...)
		} else {
			// Check patternProperties
			matched := false
			if schema.PatternProperties != nil {
				for pattern, patternSchema := range schema.PatternProperties {
					if re, err := regexp.Compile(pattern); err == nil && re.MatchString(key) {
						propErrors := v.validateValue(val, patternSchema, propPath, schemaPath+"/patternProperties/"+pattern)
						errors = append(errors, propErrors...)
						matched = true
					}
				}
			}

			// Check additionalProperties
			if !matched && schema.AdditionalProperties != nil {
				switch ap := schema.AdditionalProperties.(type) {
				case bool:
					if !ap {
						errors = append(errors, ValidationError{
							Path:       propPath,
							Message:    fmt.Sprintf("Additional property '%s' is not allowed", key),
							SchemaPath: schemaPath + "/additionalProperties",
						})
					}
				case map[string]interface{}:
					// additionalProperties is a schema
					var addSchema Schema
					jsonBytes, _ := json.Marshal(ap)
					json.Unmarshal(jsonBytes, &addSchema)
					propErrors := v.validateValue(val, &addSchema, propPath, schemaPath+"/additionalProperties")
					errors = append(errors, propErrors...)
				}
			}
		}
	}

	return errors
}

func (v *Validator) validateArray(arr []interface{}, schema *Schema, path string, schemaPath string) []ValidationError {
	var errors []ValidationError

	// MinItems
	if schema.MinItems != nil && len(arr) < *schema.MinItems {
		errors = append(errors, ValidationError{
			Path:       path,
			Message:    fmt.Sprintf("Array must have at least %d items, got %d", *schema.MinItems, len(arr)),
			SchemaPath: schemaPath + "/minItems",
		})
	}

	// MaxItems
	if schema.MaxItems != nil && len(arr) > *schema.MaxItems {
		errors = append(errors, ValidationError{
			Path:       path,
			Message:    fmt.Sprintf("Array must have at most %d items, got %d", *schema.MaxItems, len(arr)),
			SchemaPath: schemaPath + "/maxItems",
		})
	}

	// UniqueItems
	if schema.UniqueItems != nil && *schema.UniqueItems {
		seen := make(map[string]bool)
		for i, item := range arr {
			key := fmt.Sprintf("%v", item)
			if seen[key] {
				errors = append(errors, ValidationError{
					Path:       fmt.Sprintf("%s[%d]", path, i),
					Message:    "Duplicate item in array where uniqueItems is required",
					SchemaPath: schemaPath + "/uniqueItems",
				})
			}
			seen[key] = true
		}
	}

	// Validate items
	if schema.Items != nil {
		for i, item := range arr {
			itemPath := fmt.Sprintf("%s[%d]", path, i)
			itemErrors := v.validateValue(item, schema.Items, itemPath, schemaPath+"/items")
			errors = append(errors, itemErrors...)
		}
	}

	return errors
}

func (v *Validator) validateString(str string, schema *Schema, path string, schemaPath string) []ValidationError {
	var errors []ValidationError

	// MinLength
	if schema.MinLength != nil && len(str) < *schema.MinLength {
		errors = append(errors, ValidationError{
			Path:       path,
			Message:    fmt.Sprintf("String must be at least %d characters, got %d", *schema.MinLength, len(str)),
			SchemaPath: schemaPath + "/minLength",
		})
	}

	// MaxLength
	if schema.MaxLength != nil && len(str) > *schema.MaxLength {
		errors = append(errors, ValidationError{
			Path:       path,
			Message:    fmt.Sprintf("String must be at most %d characters, got %d", *schema.MaxLength, len(str)),
			SchemaPath: schemaPath + "/maxLength",
		})
	}

	// Pattern
	if schema.Pattern != "" {
		re, err := regexp.Compile(schema.Pattern)
		if err != nil {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    fmt.Sprintf("Invalid regex pattern in schema: %s", schema.Pattern),
				SchemaPath: schemaPath + "/pattern",
			})
		} else if !re.MatchString(str) {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    fmt.Sprintf("String does not match pattern: %s", schema.Pattern),
				SchemaPath: schemaPath + "/pattern",
				Value:      str,
			})
		}
	}

	// JSSON Format validation
	if schema.JSSonFormat != "" {
		formatErrors := v.validateJSSonFormat(str, schema.JSSonFormat, path, schemaPath)
		errors = append(errors, formatErrors...)
	}

	return errors
}

func (v *Validator) validateNumber(value interface{}, schema *Schema, path string, schemaPath string) []ValidationError {
	var errors []ValidationError

	num := v.toFloat64(value)

	// Minimum
	if schema.Minimum != nil && num < *schema.Minimum {
		errors = append(errors, ValidationError{
			Path:       path,
			Message:    fmt.Sprintf("Value must be >= %v, got %v", *schema.Minimum, num),
			SchemaPath: schemaPath + "/minimum",
		})
	}

	// Maximum
	if schema.Maximum != nil && num > *schema.Maximum {
		errors = append(errors, ValidationError{
			Path:       path,
			Message:    fmt.Sprintf("Value must be <= %v, got %v", *schema.Maximum, num),
			SchemaPath: schemaPath + "/maximum",
		})
	}

	return errors
}

// ============================================================================
// Custom Format Validation (JSSON Extensions)
// ============================================================================

// validateJSSonFormat validates custom JSSON format extensions
// These are generic formats that can be used across any domain
func (v *Validator) validateJSSonFormat(str string, format string, path string, schemaPath string) []ValidationError {
	var errors []ValidationError

	switch format {
	// ---- Identifier Formats ----
	case "identifier":
		// Generic identifier: starts with letter, alphanumeric + underscores
		if matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid identifier format. Must start with letter and contain only alphanumeric characters and underscores",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "kebab-case":
		// kebab-case: lowercase letters and hyphens
		if matched, _ := regexp.MatchString(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid kebab-case format. Must be lowercase with hyphens (e.g., my-component-name)",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "snake_case":
		// snake_case: lowercase letters and underscores
		if matched, _ := regexp.MatchString(`^[a-z][a-z0-9]*(_[a-z0-9]+)*$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid snake_case format. Must be lowercase with underscores (e.g., my_variable_name)",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "camelCase":
		// camelCase: starts lowercase, no separators
		if matched, _ := regexp.MatchString(`^[a-z][a-zA-Z0-9]*$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid camelCase format (e.g., myVariableName)",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "PascalCase":
		// PascalCase: starts uppercase, no separators
		if matched, _ := regexp.MatchString(`^[A-Z][a-zA-Z0-9]*$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid PascalCase format (e.g., MyClassName)",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	// ---- Path Formats ----
	case "file-path":
		// File path (Unix or Windows style)
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_./-]+$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid file path format",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "url-path":
		// URL path segment
		if matched, _ := regexp.MatchString(`^/[a-zA-Z0-9_./-]*$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid URL path format. Must start with /",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	// ---- Version Formats ----
	case "semver":
		// Semantic versioning (e.g., 1.0.0, 2.1.3-beta.1)
		if matched, _ := regexp.MatchString(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?(\+[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid semantic version format (e.g., 1.0.0, 2.1.3-beta)",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	// ---- Duration Formats ----
	case "duration":
		// Duration string (e.g., 1h30m, 500ms, 2s)
		if matched, _ := regexp.MatchString(`^(\d+)(ms|s|m|h|d)$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid duration format. Use format like: 500ms, 30s, 5m, 2h, 1d",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "duration-ms":
		// Duration in milliseconds (numeric string)
		if _, err := strconv.Atoi(str); err != nil {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Duration must be a valid integer (milliseconds)",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	// ---- Color Formats ----
	case "hex-color":
		// Hex color (e.g., #fff, #ffffff, #ffffffff)
		if matched, _ := regexp.MatchString(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid hex color format (e.g., #fff, #ffffff)",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "rgb-color":
		// RGB color (e.g., rgb(255, 255, 255))
		if matched, _ := regexp.MatchString(`^rgb\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*\)$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid RGB color format (e.g., rgb(255, 128, 0))",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	// ---- Network Formats ----
	case "port":
		// Port number as string
		port, err := strconv.Atoi(str)
		if err != nil || port < 1 || port > 65535 {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid port number. Must be between 1 and 65535",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "host":
		// Hostname or IP
		if matched, _ := regexp.MatchString(`^([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$|^(\d{1,3}\.){3}\d{1,3}$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid host format. Must be a valid hostname or IP address",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	// ---- Code Formats ----
	case "env-var":
		// Environment variable name
		if matched, _ := regexp.MatchString(`^[A-Z][A-Z0-9_]*$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid environment variable name. Must be uppercase with underscores",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	case "template-var":
		// Template variable (e.g., ${VAR}, {{var}}, {var})
		if matched, _ := regexp.MatchString(`^(\$\{[a-zA-Z_][a-zA-Z0-9_]*\}|\{\{[a-zA-Z_][a-zA-Z0-9_]*\}\}|\{[a-zA-Z_][a-zA-Z0-9_]*\})$`, str); !matched {
			errors = append(errors, ValidationError{
				Path:       path,
				Message:    "Invalid template variable format. Use ${VAR}, {{var}}, or {var}",
				SchemaPath: schemaPath,
				Value:      str,
			})
		}

	default:
		// Unknown format - you can extend with custom rules
		// Check if there's a registered custom format validator
		// For now, we skip validation for unknown formats
	}

	return errors
}

// ============================================================================
// Helper Methods
// ============================================================================

func (v *Validator) getType(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch value.(type) {
	case bool:
		return "boolean"
	case string:
		return "string"
	case float64, float32, int, int32, int64:
		return "number"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

func (v *Validator) isInteger(value interface{}) bool {
	switch n := value.(type) {
	case float64:
		return n == float64(int64(n))
	case int, int32, int64:
		return true
	default:
		return false
	}
}

func (v *Validator) toFloat64(value interface{}) float64 {
	switch n := value.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}

func (v *Validator) deepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(v.normalizeData(a), v.normalizeData(b))
}

func (v *Validator) normalizeData(data interface{}) interface{} {
	switch d := data.(type) {
	case map[interface{}]interface{}:
		// Convert YAML's map[interface{}]interface{} to map[string]interface{}
		result := make(map[string]interface{})
		for k, val := range d {
			result[fmt.Sprintf("%v", k)] = v.normalizeData(val)
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range d {
			result[k] = v.normalizeData(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(d))
		for i, val := range d {
			result[i] = v.normalizeData(val)
		}
		return result
	default:
		return data
	}
}

// parseTOML parses TOML data into a map structure
func (v *Validator) parseTOML(data []byte) (interface{}, error) {
	// Simple TOML parser for validation purposes
	// For full TOML support, you'd use a library like pelletier/go-toml
	result := make(map[string]interface{})
	lines := strings.Split(string(data), "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			// Create nested structure for section
			parts := strings.Split(currentSection, ".")
			current := result
			for _, part := range parts {
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				if m, ok := current[part].(map[string]interface{}); ok {
					current = m
				}
			}
			continue
		}

		// Key = value
		if idx := strings.Index(line, "="); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])

			// Parse value
			parsedValue := v.parseTOMLValue(value)

			// Add to current section
			if currentSection == "" {
				result[key] = parsedValue
			} else {
				parts := strings.Split(currentSection, ".")
				current := result
				for _, part := range parts {
					if m, ok := current[part].(map[string]interface{}); ok {
						current = m
					}
				}
				current[key] = parsedValue
			}
		}
	}

	return result, nil
}

func (v *Validator) parseTOMLValue(value string) interface{} {
	// String
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return strings.Trim(value, "\"")
	}
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return strings.Trim(value, "'")
	}

	// Boolean
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// Array
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		inner := strings.Trim(value, "[]")
		if inner == "" {
			return []interface{}{}
		}
		parts := strings.Split(inner, ",")
		arr := make([]interface{}, len(parts))
		for i, part := range parts {
			arr[i] = v.parseTOMLValue(strings.TrimSpace(part))
		}
		return arr
	}

	// Number
	if n, err := strconv.ParseInt(value, 10, 64); err == nil {
		return float64(n) // Normalize to float64 for consistency
	}
	if n, err := strconv.ParseFloat(value, 64); err == nil {
		return n
	}

	return value
}

// parseTypeScript extracts structure from TypeScript type definitions
func (v *Validator) parseTypeScript(data []byte) (interface{}, error) {
	// For TypeScript validation, we extract the structure from type definitions
	// This is a simplified parser that handles common patterns
	content := string(data)
	result := make(map[string]interface{})

	// Look for exported const declarations
	// export const name: Type = { ... }
	constRegex := regexp.MustCompile(`export\s+const\s+(\w+)\s*[:\s=]`)
	matches := constRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			result[match[1]] = "exported_const"
		}
	}

	// Look for type/interface definitions
	typeRegex := regexp.MustCompile(`(?:export\s+)?(?:type|interface)\s+(\w+)`)
	typeMatches := typeRegex.FindAllStringSubmatch(content, -1)

	typesDefined := make([]string, 0)
	for _, match := range typeMatches {
		if len(match) >= 2 {
			typesDefined = append(typesDefined, match[1])
		}
	}

	if len(typesDefined) > 0 {
		result["_types"] = typesDefined
	}

	return result, nil
}

// ============================================================================
// Quick Validation Functions
// ============================================================================

// QuickValidate performs a quick validation with a JSON schema string
func QuickValidate(data []byte, schemaJSON string, dataFormat string) *ValidationResult {
	v := New()
	schema, err := v.LoadSchemaFromJSON(schemaJSON)
	if err != nil {
		return &ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Path: "$", Message: "Invalid schema: " + err.Error()}},
			Format: dataFormat,
		}
	}
	return v.Validate(data, schema, dataFormat)
}

// QuickValidateJSON validates JSON data against a JSON schema
func QuickValidateJSON(jsonData []byte, schemaJSON string) *ValidationResult {
	return QuickValidate(jsonData, schemaJSON, "json")
}

// QuickValidateYAML validates YAML data against a schema
func QuickValidateYAML(yamlData []byte, schemaYAML string) *ValidationResult {
	v := New()
	schema, err := v.LoadSchemaFromYAML(schemaYAML)
	if err != nil {
		return &ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Path: "$", Message: "Invalid schema: " + err.Error()}},
			Format: "yaml",
		}
	}
	return v.Validate(yamlData, schema, "yaml")
}
