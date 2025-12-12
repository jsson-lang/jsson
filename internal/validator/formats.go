package validator

import (
	"regexp"
	"strings"
)

// validateJSSonFormat validates JSSON-specific format validators
func (v *Validator) validateJSSonFormat(value string, format string, path string, result *ValidationResult) {
	valid := true
	message := ""

	switch format {
	case "email":
		pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid email format"

	case "uri", "url":
		pattern := `^https?://[^\s/$.?#].[^\s]*$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid URL format"

	case "uuid":
		pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid UUID format"

	case "date":
		patterns := []string{
			`^\d{4}-\d{2}-\d{2}$`,
			`^\d{2}/\d{2}/\d{4}$`,
			`^\d{2}-\d{2}-\d{4}$`,
		}
		valid = false
		for _, pattern := range patterns {
			if matched, _ := regexp.MatchString(pattern, value); matched {
				valid = true
				break
			}
		}
		message = "Invalid date format"

	case "date-time", "datetime":
		patterns := []string{
			`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})?$`,
			`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`,
			`^\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}$`,
		}
		valid = false
		for _, pattern := range patterns {
			if matched, _ := regexp.MatchString(pattern, value); matched {
				valid = true
				break
			}
		}
		message = "Invalid datetime format"

	case "time":
		pattern := `^([01]?[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid time format"

	case "ipv4":
		pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid IPv4 format"

	case "ipv6":
		pattern := `^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid IPv6 format"

	case "hostname":
		pattern := `^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid hostname format"

	case "file-path", "filepath":
		valid = value != "" && !strings.Contains(value, "\x00")
		message = "Invalid file path format"

	case "semver":
		pattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-((0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid semantic version format"

	case "hex-color", "hexcolor":
		pattern := `^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid hex color format"

	case "rgb-color", "rgbcolor":
		pattern := `^rgb\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid RGB color format"

	case "port":
		pattern := `^([1-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid port number (must be 1-65535)"

	case "host":
		// Can be hostname or IP
		hostnamePattern := `^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
		ipv4Pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
		matchedHostname, _ := regexp.MatchString(hostnamePattern, value)
		matchedIP, _ := regexp.MatchString(ipv4Pattern, value)
		valid = matchedHostname || matchedIP
		message = "Invalid host format"

	case "env-var", "envvar":
		pattern := `^[A-Z_][A-Z0-9_]*$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid environment variable name format"

	case "template-var", "templatevar":
		pattern := `^\{\{[a-zA-Z_][a-zA-Z0-9_]*\}\}$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid template variable format"

	case "json-pointer":
		pattern := `^(/([^/~]|~[01])*)*$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid JSON pointer format"

	case "regex":
		_, err := regexp.Compile(value)
		valid = err == nil
		message = "Invalid regex pattern"

	case "base64":
		pattern := `^[A-Za-z0-9+/]*={0,2}$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched && len(value)%4 == 0
		message = "Invalid base64 format"

	case "phone":
		pattern := `^[\+]?[(]?[0-9]{1,3}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,4}[-\s\.]?[0-9]{1,9}$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid phone number format"

	case "credit-card", "creditcard":
		// Remove spaces and dashes
		cleaned := strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", "")
		pattern := `^\d{13,19}$`
		matched, _ := regexp.MatchString(pattern, cleaned)
		valid = matched
		message = "Invalid credit card format"

	case "slug":
		pattern := `^[a-z0-9]+(-[a-z0-9]+)*$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid slug format"

	case "alpha":
		pattern := `^[a-zA-Z]+$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Must contain only alphabetic characters"

	case "alphanumeric":
		pattern := `^[a-zA-Z0-9]+$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Must contain only alphanumeric characters"

	case "macro-id":
		// JSSON macro identifier: starts with letter/underscore, contains letters, numbers, underscores
		pattern := `^[a-zA-Z_][a-zA-Z0-9_]*$`
		matched, _ := regexp.MatchString(pattern, value)
		valid = matched
		message = "Invalid macro-id format (must start with letter or underscore)"

	default:
		// Unknown format - allow by default
		return
	}

	if !valid {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: message,
			Value:   value,
		})
	}
}

// ValidateFormat checks if a value matches a specific format (exported helper)
func ValidateFormat(value string, format string) bool {
	v := New()
	result := &ValidationResult{Valid: true, Errors: []ValidationError{}}
	v.validateJSSonFormat(value, format, "$", result)
	return result.Valid
}

// Common format validation helpers for direct use

// IsValidEmail checks if a string is a valid email
func IsValidEmail(value string) bool {
	return ValidateFormat(value, "email")
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(value string) bool {
	return ValidateFormat(value, "url")
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(value string) bool {
	return ValidateFormat(value, "uuid")
}

// IsValidIPv4 checks if a string is a valid IPv4 address
func IsValidIPv4(value string) bool {
	return ValidateFormat(value, "ipv4")
}

// IsValidIPv6 checks if a string is a valid IPv6 address
func IsValidIPv6(value string) bool {
	return ValidateFormat(value, "ipv6")
}

// IsValidDate checks if a string is a valid date
func IsValidDate(value string) bool {
	return ValidateFormat(value, "date")
}

// IsValidDateTime checks if a string is a valid datetime
func IsValidDateTime(value string) bool {
	return ValidateFormat(value, "datetime")
}

// IsValidSemVer checks if a string is a valid semantic version
func IsValidSemVer(value string) bool {
	return ValidateFormat(value, "semver")
}

// IsValidHexColor checks if a string is a valid hex color
func IsValidHexColor(value string) bool {
	return ValidateFormat(value, "hex-color")
}

// IsValidPort checks if a string is a valid port number
func IsValidPort(value string) bool {
	return ValidateFormat(value, "port")
}
