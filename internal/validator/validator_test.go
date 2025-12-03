package validator

import (
	"testing"
)

func TestValidateSimpleObject(t *testing.T) {
	schema := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
			"age":  {Type: "number"},
		},
		Required: []string{"name"},
	}

	v := New()

	// Valid data
	validJSON := []byte(`{"name": "John", "age": 30}`)
	result := v.ValidateJSON(validJSON, schema)
	if !result.Valid {
		t.Errorf("Expected valid, got errors: %v", result.Errors)
	}

	// Missing required field
	invalidJSON := []byte(`{"age": 30}`)
	result = v.ValidateJSON(invalidJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for missing required field")
	}

	// Wrong type
	wrongTypeJSON := []byte(`{"name": 123, "age": 30}`)
	result = v.ValidateJSON(wrongTypeJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for wrong type")
	}
}

func TestValidateArray(t *testing.T) {
	minItems := 2
	maxItems := 5
	schema := &Schema{
		Type:     "array",
		MinItems: &minItems,
		MaxItems: &maxItems,
		Items: &Schema{
			Type: "string",
		},
	}

	v := New()

	// Valid array
	validJSON := []byte(`["a", "b", "c"]`)
	result := v.ValidateJSON(validJSON, schema)
	if !result.Valid {
		t.Errorf("Expected valid, got errors: %v", result.Errors)
	}

	// Too few items
	tooFewJSON := []byte(`["a"]`)
	result = v.ValidateJSON(tooFewJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for too few items")
	}

	// Too many items
	tooManyJSON := []byte(`["a", "b", "c", "d", "e", "f"]`)
	result = v.ValidateJSON(tooManyJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for too many items")
	}
}

func TestValidateEnum(t *testing.T) {
	schema := &Schema{
		Type: "string",
		Enum: []interface{}{"red", "green", "blue"},
	}

	v := New()

	// Valid enum value
	validJSON := []byte(`"red"`)
	result := v.ValidateJSON(validJSON, schema)
	if !result.Valid {
		t.Errorf("Expected valid, got errors: %v", result.Errors)
	}

	// Invalid enum value
	invalidJSON := []byte(`"yellow"`)
	result = v.ValidateJSON(invalidJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for non-enum value")
	}
}

func TestValidatePattern(t *testing.T) {
	schema := &Schema{
		Type:    "string",
		Pattern: "^[a-z]+$",
	}

	v := New()

	// Valid pattern
	validJSON := []byte(`"abc"`)
	result := v.ValidateJSON(validJSON, schema)
	if !result.Valid {
		t.Errorf("Expected valid, got errors: %v", result.Errors)
	}

	// Invalid pattern
	invalidJSON := []byte(`"ABC123"`)
	result = v.ValidateJSON(invalidJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for non-matching pattern")
	}
}

func TestValidateMinMax(t *testing.T) {
	min := float64(0)
	max := float64(100)
	schema := &Schema{
		Type:    "number",
		Minimum: &min,
		Maximum: &max,
	}

	v := New()

	// Valid range
	validJSON := []byte(`50`)
	result := v.ValidateJSON(validJSON, schema)
	if !result.Valid {
		t.Errorf("Expected valid, got errors: %v", result.Errors)
	}

	// Below minimum
	belowJSON := []byte(`-10`)
	result = v.ValidateJSON(belowJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for below minimum")
	}

	// Above maximum
	aboveJSON := []byte(`150`)
	result = v.ValidateJSON(aboveJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for above maximum")
	}
}

func TestValidateYAML(t *testing.T) {
	schema := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
			"enabled": {Type: "boolean"},
		},
		Required: []string{"name"},
	}

	v := New()

	yamlData := []byte(`
name: test
enabled: true
`)
	result := v.ValidateYAML(yamlData, schema)
	if !result.Valid {
		t.Errorf("Expected valid YAML, got errors: %v", result.Errors)
	}
}

func TestLoadSchemaFromJSON(t *testing.T) {
	v := New()

	schemaJSON := `{
		"type": "object",
		"properties": {
			"id": {"type": "number"},
			"name": {"type": "string"}
		},
		"required": ["id"]
	}`

	schema, err := v.LoadSchemaFromJSON(schemaJSON)
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	if schema.Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", schema.Type)
	}

	if len(schema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(schema.Properties))
	}

	if len(schema.Required) != 1 || schema.Required[0] != "id" {
		t.Error("Expected required field 'id'")
	}
}

func TestValidateAnyOf(t *testing.T) {
	schema := &Schema{
		AnyOf: []*Schema{
			{Type: "string"},
			{Type: "number"},
		},
	}

	v := New()

	// String is valid
	stringJSON := []byte(`"hello"`)
	result := v.ValidateJSON(stringJSON, schema)
	if !result.Valid {
		t.Errorf("Expected string to be valid, got errors: %v", result.Errors)
	}

	// Number is valid
	numberJSON := []byte(`42`)
	result = v.ValidateJSON(numberJSON, schema)
	if !result.Valid {
		t.Errorf("Expected number to be valid, got errors: %v", result.Errors)
	}

	// Boolean is invalid (neither string nor number)
	boolJSON := []byte(`true`)
	result = v.ValidateJSON(boolJSON, schema)
	if result.Valid {
		t.Error("Expected boolean to be invalid for anyOf[string, number]")
	}
}

func TestValidateNestedObject(t *testing.T) {
	minActions := 1
	schema := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"macro": {
				Type: "object",
				Properties: map[string]*Schema{
					"id":   {Type: "string"},
					"name": {Type: "string"},
					"actions": {
						Type:     "array",
						MinItems: &minActions,
						Items: &Schema{
							Type: "object",
							Properties: map[string]*Schema{
								"type":  {Type: "string"},
								"key":   {Type: "string"},
								"delay": {Type: "number"},
							},
							Required: []string{"type"},
						},
					},
				},
				Required: []string{"id", "name", "actions"},
			},
		},
	}

	v := New()

	validJSON := []byte(`{
		"macro": {
			"id": "macro_1",
			"name": "Test Macro",
			"actions": [
				{"type": "click", "key": "KEY_A", "delay": 100},
				{"type": "hold", "key": "KEY_CTRL"}
			]
		}
	}`)

	result := v.ValidateJSON(validJSON, schema)
	if !result.Valid {
		t.Errorf("Expected valid nested object, got errors: %v", result.Errors)
	}

	// Missing required nested field
	invalidJSON := []byte(`{
		"macro": {
			"id": "macro_1",
			"actions": []
		}
	}`)

	result = v.ValidateJSON(invalidJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for missing required 'name' field")
	}
}

func TestQuickValidate(t *testing.T) {
	schemaJSON := `{"type": "string", "minLength": 3}`
	
	// Valid
	result := QuickValidateJSON([]byte(`"hello"`), schemaJSON)
	if !result.Valid {
		t.Errorf("Expected valid, got errors: %v", result.Errors)
	}

	// Invalid
	result = QuickValidateJSON([]byte(`"hi"`), schemaJSON)
	if result.Valid {
		t.Error("Expected invalid for string too short")
	}
}

func TestJSSonFormatValidation(t *testing.T) {
	schema := &Schema{
		Type:       "string",
		JSSonFormat: "macro-id",
	}

	v := New()

	// Valid macro ID
	validJSON := []byte(`"my_macro_1"`)
	result := v.ValidateJSON(validJSON, schema)
	if !result.Valid {
		t.Errorf("Expected valid macro-id, got errors: %v", result.Errors)
	}

	// Invalid macro ID (starts with number)
	invalidJSON := []byte(`"123_macro"`)
	result = v.ValidateJSON(invalidJSON, schema)
	if result.Valid {
		t.Error("Expected invalid for macro-id starting with number")
	}
}
