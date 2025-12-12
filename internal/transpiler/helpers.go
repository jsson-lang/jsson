package transpiler

import ie "jsson/internal/errors"

// toFloat converts a value to float64, returning true if the original was a float
func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), false
	case int:
		return float64(v), false
	}
	return 0, false
}

// toInt64 converts a value to int64
func toInt64(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	}
	return 0, false
}

// isTruthy determines if a value is truthy in JSSON
func (t *Transpiler) isTruthy(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case int64:
		return v != 0
	case float64:
		return v != 0.0
	case string:
		return v != ""
	case nil:
		return false
	default:
		return true
	}
}

// compareEqual checks equality between two values
func (t *Transpiler) compareEqual(left, right interface{}) bool {
	// Handle mixed numeric types
	lFloat, lIsFloat := toFloat(left)
	rFloat, rIsFloat := toFloat(right)
	if lIsFloat && rIsFloat {
		return lFloat == rFloat
	}

	// Handle same types
	switch l := left.(type) {
	case int64:
		if r, ok := right.(int64); ok {
			return l == r
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l == r
		}
	case string:
		if r, ok := right.(string); ok {
			return l == r
		}
	case bool:
		if r, ok := right.(bool); ok {
			return l == r
		}
	}
	return false
}

// compareLess checks if left < right
func (t *Transpiler) compareLess(left, right interface{}) (bool, error) {
	// Handle mixed numeric types
	lFloat, lIsFloat := toFloat(left)
	rFloat, rIsFloat := toFloat(right)
	if lIsFloat && rIsFloat {
		return lFloat < rFloat, nil
	}

	switch l := left.(type) {
	case int64:
		if r, ok := right.(int64); ok {
			return l < r, nil
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l < r, nil
		}
	case string:
		if r, ok := right.(string); ok {
			return l < r, nil
		}
	}
	return false, t.errMsg(ie.UnsupportedComparison(left, right))
}
