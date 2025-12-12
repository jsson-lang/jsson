package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Validate validates data against a schema
func (v *Validator) Validate(data []byte, schema *Schema, format string) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
		Format: format,
	}

	var parsedData any
	var err error

	switch format {
	case "json":
		err = json.Unmarshal(data, &parsedData)
	case "yaml":
		err = yaml.Unmarshal(data, &parsedData)
		if err == nil {
			parsedData = normalizeData(parsedData)
		}
	case "toml":
		parsedData, err = parseTOML(string(data))
	case "typescript", "ts":
		parsedData, err = parseTypeScript(string(data))
	default:
		err = json.Unmarshal(data, &parsedData)
		if err != nil {
			err = yaml.Unmarshal(data, &parsedData)
			if err == nil {
				parsedData = normalizeData(parsedData)
			}
		}
	}

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    "$",
			Message: fmt.Sprintf("Failed to parse %s: %s", format, err.Error()),
		})
		return result
	}

	v.validateValue(parsedData, schema, "$", result)
	return result
}

// ValidateJSON validates JSON data against a schema
func (v *Validator) ValidateJSON(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "json")
}

// ValidateYAML validates YAML data against a schema
func (v *Validator) ValidateYAML(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "yaml")
}

// ValidateTOML validates TOML data against a schema
func (v *Validator) ValidateTOML(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "toml")
}

// ValidateTypeScript validates TypeScript data against a schema
func (v *Validator) ValidateTypeScript(data []byte, schema *Schema) *ValidationResult {
	return v.Validate(data, schema, "typescript")
}

// validateValue recursively validates a value against a schema
func (v *Validator) validateValue(value any, schema *Schema, path string, result *ValidationResult) {
	if schema == nil {
		return
	}

	// Check const
	if schema.Const != nil {
		if !deepEqual(value, schema.Const) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Value must be equal to constant %v", schema.Const),
				Value:   value,
			})
			return
		}
	}

	// Check enum
	if len(schema.Enum) > 0 {
		found := false
		for _, enumVal := range schema.Enum {
			if deepEqual(value, enumVal) {
				found = true
				break
			}
		}
		if !found {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Value must be one of: %v", schema.Enum),
				Value:   value,
			})
			return
		}
	}

	// Check oneOf
	if len(schema.OneOf) > 0 {
		validCount := 0
		for _, subSchema := range schema.OneOf {
			subResult := &ValidationResult{Valid: true, Errors: []ValidationError{}}
			v.validateValue(value, subSchema, path, subResult)
			if subResult.Valid {
				validCount++
			}
		}
		if validCount != 1 {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: "Value must match exactly one of the oneOf schemas",
				Value:   value,
			})
			return
		}
	}

	// Check anyOf
	if len(schema.AnyOf) > 0 {
		validAny := false
		for _, subSchema := range schema.AnyOf {
			subResult := &ValidationResult{Valid: true, Errors: []ValidationError{}}
			v.validateValue(value, subSchema, path, subResult)
			if subResult.Valid {
				validAny = true
				break
			}
		}
		if !validAny {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: "Value must match at least one of the anyOf schemas",
				Value:   value,
			})
			return
		}
	}

	// Check allOf
	if len(schema.AllOf) > 0 {
		for _, subSchema := range schema.AllOf {
			v.validateValue(value, subSchema, path, result)
		}
	}

	// Check not
	if schema.Not != nil {
		subResult := &ValidationResult{Valid: true, Errors: []ValidationError{}}
		v.validateValue(value, schema.Not, path, subResult)
		if subResult.Valid {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: "Value must not match the 'not' schema",
				Value:   value,
			})
			return
		}
	}

	// Check type
	if schema.Type != "" {
		v.validateType(value, schema, path, result)
	}
}

// validateType validates the type of a value
func (v *Validator) validateType(value any, schema *Schema, path string, result *ValidationResult) {
	actualType := getType(value)

	// Handle null
	if value == nil {
		if schema.Type != "null" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Expected type '%s', got 'null'", schema.Type),
				Value:   value,
			})
		}
		return
	}

	// Type validation
	validType := false
	switch schema.Type {
	case "string":
		validType = actualType == "string"
	case "number":
		validType = actualType == "number" || actualType == "integer"
	case "integer":
		validType = actualType == "integer" || (actualType == "number" && isInteger(value))
	case "boolean":
		validType = actualType == "boolean"
	case "array":
		validType = actualType == "array"
	case "object":
		validType = actualType == "object"
	case "null":
		validType = value == nil
	}

	if !validType {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Expected type '%s', got '%s'", schema.Type, actualType),
			Value:   value,
		})
		return
	}

	// Type-specific validations
	switch schema.Type {
	case "string":
		v.validateString(value.(string), schema, path, result)
	case "number", "integer":
		v.validateNumber(value, schema, path, result)
	case "array":
		v.validateArray(value, schema, path, result)
	case "object":
		v.validateObject(value, schema, path, result)
	}
}

// validateString validates a string value
func (v *Validator) validateString(value string, schema *Schema, path string, result *ValidationResult) {
	// Check minLength
	if schema.MinLength != nil && len(value) < *schema.MinLength {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("String length must be at least %d", *schema.MinLength),
			Value:   value,
		})
	}

	// Check maxLength
	if schema.MaxLength != nil && len(value) > *schema.MaxLength {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("String length must be at most %d", *schema.MaxLength),
			Value:   value,
		})
	}

	// Check pattern
	if schema.Pattern != "" {
		matched, err := regexp.MatchString(schema.Pattern, value)
		if err != nil || !matched {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("String does not match pattern '%s'", schema.Pattern),
				Value:   value,
			})
		}
	}

	// Check format
	if schema.Format != "" {
		v.validateJSSonFormat(value, schema.Format, path, result)
	}

	// Check JSSON-specific format
	if schema.JSSonFormat != "" {
		v.validateJSSonFormat(value, schema.JSSonFormat, path, result)
	}
}

// validateNumber validates a numeric value
func (v *Validator) validateNumber(value any, schema *Schema, path string, result *ValidationResult) {
	num := toFloat64(value)

	// Check minimum
	if schema.Minimum != nil && num < *schema.Minimum {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Number must be at least %v", *schema.Minimum),
			Value:   value,
		})
	}

	// Check maximum
	if schema.Maximum != nil && num > *schema.Maximum {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Number must be at most %v", *schema.Maximum),
			Value:   value,
		})
	}
}

// validateArray validates an array value
func (v *Validator) validateArray(value any, schema *Schema, path string, result *ValidationResult) {
	// Handle both []any and []interface{}
	var arr []any
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Slice {
		arr = make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			arr[i] = rv.Index(i).Interface()
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: "Expected array type",
			Value:   value,
		})
		return
	}

	// Check minItems
	if schema.MinItems != nil && len(arr) < *schema.MinItems {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Array must have at least %d items", *schema.MinItems),
			Value:   value,
		})
	}

	// Check maxItems
	if schema.MaxItems != nil && len(arr) > *schema.MaxItems {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("Array must have at most %d items", *schema.MaxItems),
			Value:   value,
		})
	}

	// Check uniqueItems
	if schema.UniqueItems {
		seen := make(map[string]bool)
		for i, item := range arr {
			key := fmt.Sprintf("%v", item)
			if seen[key] {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Path:    fmt.Sprintf("%s[%d]", path, i),
					Message: "Array items must be unique",
					Value:   item,
				})
			}
			seen[key] = true
		}
	}

	// Validate items
	if schema.Items != nil {
		for i, item := range arr {
			itemPath := fmt.Sprintf("%s[%d]", path, i)
			v.validateValue(item, schema.Items, itemPath, result)
		}
	}
}

// validateObject validates an object value
func (v *Validator) validateObject(value any, schema *Schema, path string, result *ValidationResult) {
	obj, ok := value.(map[string]any)
	if !ok {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: "Expected object type",
			Value:   value,
		})
		return
	}

	// Check required properties
	for _, req := range schema.Required {
		if _, exists := obj[req]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("Missing required property '%s'", req),
			})
		}
	}

	// Check additionalProperties
	if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
		for key := range obj {
			if _, defined := schema.Properties[key]; !defined {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Path:    fmt.Sprintf("%s.%s", path, key),
					Message: fmt.Sprintf("Additional property '%s' is not allowed", key),
					Value:   obj[key],
				})
			}
		}
	}

	// Validate properties
	for key, propSchema := range schema.Properties {
		if propValue, exists := obj[key]; exists {
			propPath := fmt.Sprintf("%s.%s", path, key)
			v.validateValue(propValue, propSchema, propPath, result)
		}
	}
}
