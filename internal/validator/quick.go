package validator

import (
	"regexp"
	"strings"
)

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

// Inline Validators (for @uuid, @email, etc)

// ValidateUUID validates a UUID string
func ValidateUUID(value string) bool {
	pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// ValidateEmail validates an email string
func ValidateEmail(value string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// ValidateURL validates a URL string
func ValidateURL(value string) bool {
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// ValidateIPv4 validates an IPv4 address string
func ValidateIPv4(value string) bool {
	pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// ValidateIPv6 validates an IPv6 address string
func ValidateIPv6(value string) bool {
	pattern := `^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// ValidateFilePath validates a file path string
func ValidateFilePath(value string) bool {
	if value == "" {
		return false
	}
	invalidChars := []string{"\x00", "\n", "\r"}
	for _, char := range invalidChars {
		if strings.Contains(value, char) {
			return false
		}
	}
	return true
}

// ValidateDate validates a date string
func ValidateDate(value string) bool {
	patterns := []string{
		`^\d{4}-\d{2}-\d{2}$`,
		`^\d{2}/\d{2}/\d{4}$`,
		`^\d{2}-\d{2}-\d{4}$`,
	}
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, value); matched {
			return true
		}
	}
	return false
}

// ValidateDateTime validates a datetime string
func ValidateDateTime(value string) bool {
	patterns := []string{
		`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})?$`,
		`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`,
		`^\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}$`,
	}
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, value); matched {
			return true
		}
	}
	return false
}

// ValidateRegex validates a string against a regex pattern
func ValidateRegex(value string, pattern string) bool {
	if pattern == "" {
		return true
	}
	matched, err := regexp.MatchString(pattern, value)
	return err == nil && matched
}
