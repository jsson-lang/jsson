package validator

import (
	"reflect"
)

// getType returns the JSON type of a value
func getType(value any) string {
	if value == nil {
		return "null"
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if f == float64(int64(f)) {
			return "integer"
		}
		return "number"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map:
		return "object"
	default:
		return "unknown"
	}
}

// isInteger checks if a value is an integer
func isInteger(value any) bool {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32:
		return v == float32(int64(v))
	case float64:
		return v == float64(int64(v))
	default:
		return false
	}
}

// toFloat64 converts a numeric value to float64
func toFloat64(value any) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}

// deepEqual compares two values for equality
func deepEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// normalizeData converts YAML-specific types to JSON-compatible types
func normalizeData(data any) any {
	switch v := data.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, value := range v {
			result[key] = normalizeData(value)
		}
		return result
	case map[any]any:
		result := make(map[string]any)
		for key, value := range v {
			if strKey, ok := key.(string); ok {
				result[strKey] = normalizeData(value)
			}
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = normalizeData(item)
		}
		return result
	default:
		return data
	}
}
