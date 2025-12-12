package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// New creates a new Validator instance
func New() *Validator {
	return &Validator{
		schemas:     make(map[string]*Schema),
		customRules: make(map[string]ValidationRule),
	}
}

// LoadSchemaFromJSON loads a schema from a JSON string
func (v *Validator) LoadSchemaFromJSON(schemaJSON string) (*Schema, error) {
	var schema Schema
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}
	return &schema, nil
}

// LoadSchemaFromYAML loads a schema from a YAML string
func (v *Validator) LoadSchemaFromYAML(schemaYAML string) (*Schema, error) {
	var schema Schema
	if err := yaml.Unmarshal([]byte(schemaYAML), &schema); err != nil {
		return nil, fmt.Errorf("failed to parse YAML schema: %w", err)
	}
	return &schema, nil
}

// LoadSchemaAuto auto-detects format and loads schema
func (v *Validator) LoadSchemaAuto(schemaStr string) (*Schema, string, error) {
	trimmed := strings.TrimSpace(schemaStr)
	if strings.HasPrefix(trimmed, "{") {
		schema, err := v.LoadSchemaFromJSON(schemaStr)
		return schema, "json", err
	}
	schema, err := v.LoadSchemaFromYAML(schemaStr)
	return schema, "yaml", err
}

// RegisterSchema registers a schema with a name for reuse
func (v *Validator) RegisterSchema(name string, schema *Schema) {
	v.schemas[name] = schema
}

// GetSchema retrieves a registered schema by name
func (v *Validator) GetSchema(name string) (*Schema, bool) {
	schema, ok := v.schemas[name]
	return schema, ok
}

// RegisterCustomRule registers a custom validation rule
func (v *Validator) RegisterCustomRule(name string, rule ValidationRule) {
	v.customRules[name] = rule
}
