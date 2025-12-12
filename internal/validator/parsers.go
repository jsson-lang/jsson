package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// parseTOML parses a simple TOML string into a map
func parseTOML(tomlStr string) (map[string]any, error) {
	result := make(map[string]any)
	currentSection := result

	lines := strings.Split(tomlStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.Trim(line, "[]")
			parts := strings.Split(sectionName, ".")

			currentSection = result
			for _, part := range parts {
				if _, exists := currentSection[part]; !exists {
					currentSection[part] = make(map[string]any)
				}
				if nextSection, ok := currentSection[part].(map[string]any); ok {
					currentSection = nextSection
				}
			}
			continue
		}

		// Key-value pair
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				currentSection[key] = parseTOMLValue(value)
			}
		}
	}

	return result, nil
}

// parseTOMLValue parses a TOML value string
func parseTOMLValue(value string) any {
	// String (quoted)
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

	// Number
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Array
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		inner := strings.Trim(value, "[]")
		if inner == "" {
			return []any{}
		}
		// Simple split by comma (doesn't handle nested arrays well)
		parts := strings.Split(inner, ",")
		arr := make([]any, 0, len(parts))
		for _, part := range parts {
			arr = append(arr, parseTOMLValue(strings.TrimSpace(part)))
		}
		return arr
	}

	// Default to string
	return value
}

// parseTypeScript parses TypeScript const/interface to extract data
func parseTypeScript(tsStr string) (map[string]any, error) {
	result := make(map[string]any)

	// Match const declarations: const name = { ... }
	constPattern := regexp.MustCompile(`(?s)const\s+(\w+)\s*=\s*(\{[^}]+\})`)
	matches := constPattern.FindAllStringSubmatch(tsStr, -1)

	for _, match := range matches {
		if len(match) == 3 {
			varName := match[1]
			objStr := match[2]

			// Parse simple object literal
			obj := parseTypeScriptObject(objStr)
			result[varName] = obj
		}
	}

	// Match export default
	defaultPattern := regexp.MustCompile(`(?s)export\s+default\s+(\{[^}]+\})`)
	if defaultMatch := defaultPattern.FindStringSubmatch(tsStr); len(defaultMatch) == 2 {
		result["default"] = parseTypeScriptObject(defaultMatch[1])
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid TypeScript data found")
	}

	return result, nil
}

// parseTypeScriptObject parses a simple TypeScript object literal
func parseTypeScriptObject(objStr string) map[string]any {
	result := make(map[string]any)

	// Remove braces
	inner := strings.Trim(objStr, "{}")
	inner = strings.TrimSpace(inner)

	// Split by comma or newline
	lines := strings.FieldsFunc(inner, func(r rune) bool {
		return r == ',' || r == '\n'
	})

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split by colon
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes from key
			key = strings.Trim(key, "\"'")

			// Parse value
			result[key] = parseTypeScriptValue(value)
		}
	}

	return result
}

// parseTypeScriptValue parses a TypeScript value
func parseTypeScriptValue(value string) any {
	// Remove trailing comma if present
	value = strings.TrimSuffix(value, ",")
	value = strings.TrimSpace(value)

	// String
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return strings.Trim(value, "\"")
	}
	if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
		return strings.Trim(value, "'")
	}
	if strings.HasPrefix(value, "`") && strings.HasSuffix(value, "`") {
		return strings.Trim(value, "`")
	}

	// Boolean
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// Null/undefined
	if value == "null" || value == "undefined" {
		return nil
	}

	// Number
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Default to string
	return value
}
