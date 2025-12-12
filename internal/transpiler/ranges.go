package transpiler

import (
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
)

// evalStringRange handles ranges of strings with numeric suffixes (e.g., IP addresses)
// Example: "192.168.1.100".."192.168.1.109" generates ["192.168.1.100", "192.168.1.101", ...]
func (t *Transpiler) evalStringRange(start, end string, stepV interface{}, node ast.Node) (interface{}, error) {
	// Find the numeric suffix in both strings
	// We'll look for the last sequence of digits
	var startPrefix, endPrefix string
	var startNum, endNum int64
	var foundStart, foundEnd bool

	// Extract numeric suffix from start
	for i := len(start) - 1; i >= 0; i-- {
		if start[i] < '0' || start[i] > '9' {
			// Found non-digit, extract number after this position
			if i < len(start)-1 {
				startPrefix = start[:i+1]
				numStr := start[i+1:]
				if n, err := fmt.Sscanf(numStr, "%d", &startNum); n == 1 && err == nil {
					foundStart = true
				}
			}
			break
		}
		if i == 0 {
			// Entire string is a number
			startPrefix = ""
			if n, err := fmt.Sscanf(start, "%d", &startNum); n == 1 && err == nil {
				foundStart = true
			}
			break
		}
	}

	// Extract numeric suffix from end
	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < '0' || end[i] > '9' {
			if i < len(end)-1 {
				endPrefix = end[:i+1]
				numStr := end[i+1:]
				if n, err := fmt.Sscanf(numStr, "%d", &endNum); n == 1 && err == nil {
					foundEnd = true
				}
			}
			break
		}
		if i == 0 {
			endPrefix = ""
			if n, err := fmt.Sscanf(end, "%d", &endNum); n == 1 && err == nil {
				foundEnd = true
			}
			break
		}
	}

	if !foundStart || !foundEnd {
		return nil, t.errfNode(node, "string range requires numeric suffix in both start and end (e.g., \"192.168.1.100\"..\"192.168.1.109\")")
	}

	if startPrefix != endPrefix {
		return nil, t.errfNode(node, "string range prefixes must match (start: %q, end: %q)", startPrefix, endPrefix)
	}

	// Determine step
	step := int64(1)
	if stepV != nil {
		if st, ok := stepV.(int64); ok {
			step = st
		} else {
			return nil, t.errfNode(node, "step must be an integer for string ranges")
		}
	} else {
		if startNum > endNum {
			step = -1
		}
	}

	if step == 0 {
		return nil, t.errfNode(node, "step cannot be zero")
	}

	// Calculate number of digits in original (for zero-padding)
	startStr := start[len(startPrefix):]
	padding := len(startStr)

	// Generate range
	res := make([]interface{}, 0)
	if step > 0 {
		for i := startNum; i <= endNum; i += step {
			// Format with zero-padding if original had it
			if padding > 1 && startStr[0] == '0' {
				res = append(res, fmt.Sprintf("%s%0*d", startPrefix, padding, i))
			} else {
				res = append(res, fmt.Sprintf("%s%d", startPrefix, i))
			}
		}
	} else {
		for i := startNum; i >= endNum; i += step {
			if padding > 1 && startStr[0] == '0' {
				res = append(res, fmt.Sprintf("%s%0*d", startPrefix, padding, i))
			} else {
				res = append(res, fmt.Sprintf("%s%d", startPrefix, i))
			}
		}
	}

	return res, nil
}

// evalIntegerRange evaluates an integer range expression
func (t *Transpiler) evalIntegerRange(sInt, eInt int64, stepV interface{}, node ast.Node) (interface{}, error) {
	step := int64(1)
	if stepV != nil {
		if st, ok := stepV.(int64); ok {
			step = st
		} else {
			return nil, t.errfNodeMsg(node, ie.StepNotInteger(stepV))
		}
	} else {
		if sInt > eInt {
			step = -1
		}
	}

	if step == 0 {
		return nil, t.errfNodeMsg(node, ie.StepCannotBeZero())
	}

	res := make([]interface{}, 0)
	if step > 0 {
		for i := sInt; i <= eInt; i += step {
			res = append(res, i)
		}
	} else {
		for i := sInt; i >= eInt; i += step {
			res = append(res, i)
		}
	}
	return RangeResult{Values: res}, nil
}
