package transpiler

import (
	"encoding/json"
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"os"
	"path/filepath"
)

type Transpiler struct {
	program *ast.Program
	baseDir string
	// includeCache stores previously transpiled include outputs keyed by absolute path
	includeCache map[string]map[string]interface{}
	// inProgress marks includes currently being processed to detect cycles
	inProgress map[string]bool
	// mergeMode controls include merge behavior: "keep" (default), "overwrite", "error"
	mergeMode string
	// sourceFile is the path to the source file being transpiled (optional)
	sourceFile string
}

func New(program *ast.Program, baseDir string, mergeMode string, sourceFile string) *Transpiler {
	if mergeMode == "" {
		mergeMode = "keep"
	}
	return &Transpiler{program: program, baseDir: baseDir, includeCache: make(map[string]map[string]interface{}), inProgress: make(map[string]bool), mergeMode: mergeMode, sourceFile: sourceFile}
}

func (t *Transpiler) Transpile() ([]byte, error) {
	root := make(map[string]interface{})

	for _, stmt := range t.program.Statements {
		switch s := stmt.(type) {
		case *ast.AssignmentStatement:
			key := s.Name.Value
			val, err := t.evalExpression(s.Value, nil)
			if err != nil {
				return nil, err
			}
			root[key] = val
		case *ast.IncludeStatement:
			// Read included file and merge its output into root (do not overwrite existing keys)
			includePath := s.Path.Value

			// Resolve path relative to the current Transpiler baseDir when not absolute
			var includeAbs string
			if filepath.IsAbs(includePath) {
				includeAbs = filepath.Clean(includePath)
			} else {
				includeAbs = filepath.Clean(filepath.Join(t.baseDir, includePath))
			}

			// Detect cyclic include
			if t.inProgress[includeAbs] {
				return nil, t.errfNodeMsg(s, ie.CyclicInclude(includeAbs))
			}

			// If cached, use cached result
			if cached, ok := t.includeCache[includeAbs]; ok {
				for k, v := range cached {
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				}
				break
			}

			// Mark as in-progress
			t.inProgress[includeAbs] = true

			data, err := os.ReadFile(includeAbs)
			if err != nil {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "could not read include file %q — gremlin can't find it: %v", s.Path.Value, err)
			}

			l := lexer.New(string(data))
			l.SetSourceFile(includeAbs)
			p := parser.New(l)
			prog := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "parser errors in included file %q — wizard got confused: %v", s.Path.Value, p.Errors())
			}

			// Create a transpiler for the included program, setting its baseDir to the included file's dir
			incBase := filepath.Dir(includeAbs)
			incT := New(prog, incBase, t.mergeMode, includeAbs)
			// share cache and inProgress maps so nested includes use the same state
			incT.includeCache = t.includeCache
			incT.inProgress = t.inProgress

			incJSON, err := incT.Transpile()
			if err != nil {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "transpile error in included file %q: %v", s.Path.Value, err)
			}

			var incRoot map[string]interface{}
			if err := json.Unmarshal(incJSON, &incRoot); err != nil {
				t.inProgress[includeAbs] = false
				return nil, t.errfNode(s, "invalid json from include %q: %v", s.Path.Value, err)
			}

			// Cache result
			t.includeCache[includeAbs] = incRoot
			// Done processing
			t.inProgress[includeAbs] = false

			// Merge according to mergeMode
			for k, v := range incRoot {
				switch t.mergeMode {
				case "keep":
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				case "overwrite":
					root[k] = v
				case "error":
					if _, exists := root[k]; exists {
						return nil, t.errfNode(s, "include merge conflict for key %q from %s", k, includeAbs)
					}
					root[k] = v
				default:
					if _, exists := root[k]; !exists {
						root[k] = v
					}
				}
			}
		}
	}

	return json.MarshalIndent(root, "", "  ")
}

func (t *Transpiler) errf(format string, args ...interface{}) error {
	prefix := "Transpile gremlin:"
	if t != nil && t.sourceFile != "" {
		ctx := ie.FormatContext(t.sourceFile, 1, 1)
		return fmt.Errorf("%s %s — %s", prefix, ctx, fmt.Sprintf(format, args...))
	}
	return fmt.Errorf("%s — %s", prefix, fmt.Sprintf(format, args...))
}

func (t *Transpiler) errfNode(node ast.Node, format string, args ...interface{}) error {
	prefix := "Transpile gremlin:"
	var line, col int
	if node != nil {
		switch n := node.(type) {
		case *ast.AssignmentStatement:
			line, col = n.Token.Line, n.Token.Column
		case *ast.IncludeStatement:
			line, col = n.Token.Line, n.Token.Column
		case *ast.IntegerLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.FloatLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.StringLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.Identifier:
			line, col = n.Token.Line, n.Token.Column
		case *ast.ObjectLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.ArrayLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.RangeExpression:
			line, col = n.Token.Line, n.Token.Column
		case *ast.ArrayTemplate:
			line, col = n.Token.Line, n.Token.Column
		case *ast.MapClause:
			line, col = n.Token.Line, n.Token.Column
		case *ast.BinaryExpression:
			line, col = n.Token.Line, n.Token.Column
		case *ast.MemberExpression:
			line, col = n.Token.Line, n.Token.Column
		default:
			line, col = 0, 0
		}
	}

	if t != nil && t.sourceFile != "" {
		if line > 0 && col > 0 {
			ctx := ie.FormatContext(t.sourceFile, line, col)
			return fmt.Errorf("%s %s — %s", prefix, ctx, fmt.Sprintf(format, args...))
		}
		// fallback to file-only context
		ctx := ie.FormatContext(t.sourceFile, 1, 1)
		return fmt.Errorf("%s %s — %s", prefix, ctx, fmt.Sprintf(format, args...))
	}
	return fmt.Errorf("%s — %s", prefix, fmt.Sprintf(format, args...))
}

// errfNodeMsg formats an already-formatted error message with node context
func (t *Transpiler) errfNodeMsg(node ast.Node, msg string) error {
	prefix := "Transpile gremlin:"
	var line, col int
	if node != nil {
		switch n := node.(type) {
		case *ast.AssignmentStatement:
			line, col = n.Token.Line, n.Token.Column
		case *ast.IncludeStatement:
			line, col = n.Token.Line, n.Token.Column
		case *ast.IntegerLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.FloatLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.StringLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.Identifier:
			line, col = n.Token.Line, n.Token.Column
		case *ast.ObjectLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.ArrayLiteral:
			line, col = n.Token.Line, n.Token.Column
		case *ast.RangeExpression:
			line, col = n.Token.Line, n.Token.Column
		case *ast.ArrayTemplate:
			line, col = n.Token.Line, n.Token.Column
		case *ast.MapClause:
			line, col = n.Token.Line, n.Token.Column
		case *ast.BinaryExpression:
			line, col = n.Token.Line, n.Token.Column
		case *ast.MemberExpression:
			line, col = n.Token.Line, n.Token.Column
		default:
			line, col = 0, 0
		}
	}

	if t != nil && t.sourceFile != "" {
		if line > 0 && col > 0 {
			ctx := ie.FormatContext(t.sourceFile, line, col)
			return fmt.Errorf("%s %s — %s", prefix, ctx, msg)
		}
		// fallback to file-only context
		ctx := ie.FormatContext(t.sourceFile, 1, 1)
		return fmt.Errorf("%s %s — %s", prefix, ctx, msg)
	}
	return fmt.Errorf("%s — %s", prefix, msg)
}

// errMsg formats an already-formatted error message with context
func (t *Transpiler) errMsg(msg string) error {
	prefix := "Transpile gremlin:"
	if t != nil && t.sourceFile != "" {
		ctx := ie.FormatContext(t.sourceFile, 1, 1)
		return fmt.Errorf("%s %s — %s", prefix, ctx, msg)
	}
	return fmt.Errorf("%s — %s", prefix, msg)
}

func (t *Transpiler) evalExpression(expr ast.Expression, ctx map[string]interface{}) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return e.Value, nil
	case *ast.FloatLiteral:
		return e.Value, nil
	case *ast.BooleanLiteral:
		return e.Value, nil
	case *ast.StringLiteral:
		return e.Value, nil
	case *ast.Identifier:
		// Variable lookup
		if ctx != nil {
			if val, ok := ctx[e.Value]; ok {
				return val, nil
			}
		}
		// Se não for encontrado, talvez retorne como string? Ou um erro? Slr KKKKK, se funciona não mexe.
		// For now, return as string to support bare strings that parsed as Identifiers
		return e.Value, nil
	case *ast.ObjectLiteral:
		obj := make(map[string]interface{})
		// Use Keys to preserve order? JSON map is unordered.
		// But we evaluate in order.
		for _, key := range e.Keys {
			valExpr := e.Properties[key]
			if valExpr == nil {
				// Key only, ignore? Or null?
				continue
			}
			val, err := t.evalExpression(valExpr, ctx)
			if err != nil {
				return nil, err
			}
			obj[key] = val
		}
		return obj, nil
	case *ast.ArrayLiteral:
		arr := make([]interface{}, 0, len(e.Elements))
		for _, el := range e.Elements {
			val, err := t.evalExpression(el, ctx)
			if err != nil {
				return nil, err
			}
			// If the element evaluated to a slice (e.g., a range), flatten it into the array
			if slice, ok := val.([]interface{}); ok {
				arr = append(arr, slice...)
			} else {
				arr = append(arr, val)
			}
		}
		return arr, nil

	case *ast.RangeExpression:
		// Evaluate start, end and optional step
		startV, err := t.evalExpression(e.Start, ctx)
		if err != nil {
			return nil, err
		}
		endV, err := t.evalExpression(e.End, ctx)
		if err != nil {
			return nil, err
		}

		var stepV interface{}
		if e.Step != nil {
			stepV, err = t.evalExpression(e.Step, ctx)
			if err != nil {
				return nil, err
			}
		}

		// Check if both start and end are strings (String Range)
		if startStr, ok1 := startV.(string); ok1 {
			if endStr, ok2 := endV.(string); ok2 {
				// String Range: find numeric suffix and increment it
				return t.evalStringRange(startStr, endStr, stepV, e)
			}
		}

		// Integer Range (original behavior)
		sInt, ok1 := startV.(int64)
		eInt, ok2 := endV.(int64)
		if !ok1 || !ok2 {
			return nil, t.errfNodeMsg(e, ie.RangeBoundsNotIntegers(startV, endV))
		}

		step := int64(1)
		if stepV != nil {
			if st, ok := stepV.(int64); ok {
				step = st
			} else {
				return nil, t.errfNodeMsg(e, ie.StepNotInteger(stepV))
			}
		} else {
			if sInt > eInt {
				step = -1
			}
		}

		if step == 0 {
			return nil, t.errfNodeMsg(e, ie.StepCannotBeZero())
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
		return res, nil
	case *ast.ArrayTemplate:
		result := make([]interface{}, 0, len(e.Rows))
		keys := e.Template.Keys

		for _, row := range e.Rows {
			// First, evaluate all expressions in the row
			evaluatedRow := make([]interface{}, len(row))
			for i, expr := range row {
				val, err := t.evalExpression(expr, ctx)
				if err != nil {
					return nil, err
				}
				evaluatedRow[i] = val
			}

			// Check if we have ranges that need zipping
			// If we have arrays, we zip them up to the shortest length
			hasArrays := false
			minArrayLength := -1

			for _, val := range evaluatedRow {
				if arr, ok := val.([]interface{}); ok {
					hasArrays = true
					if minArrayLength == -1 || len(arr) < minArrayLength {
						minArrayLength = len(arr)
					}
				}
			}

			// Range Zipping: if we have arrays, zip them
			if hasArrays && minArrayLength > 0 {
				// Create multiple objects, one for each index up to minArrayLength
				for idx := 0; idx < minArrayLength; idx++ {
					rowObj := make(map[string]interface{})
					for i, val := range evaluatedRow {
						if i >= len(keys) {
							break
						}
						key := keys[i]

						// If it's an array, take the element at idx
						if arr, ok := val.([]interface{}); ok {
							rowObj[key] = arr[idx]
						} else {
							// If it's not an array, use the same value for all rows
							rowObj[key] = val
						}
					}

					// Apply Map Clause if present
					if e.Map != nil {
						mapCtx := make(map[string]interface{})
						for k, v := range ctx {
							mapCtx[k] = v
						}
						mapCtx[e.Map.Param.Value] = rowObj

						mappedVal, err := t.evalExpression(e.Map.Body, mapCtx)
						if err != nil {
							return nil, err
						}
						result = append(result, mappedVal)
					} else {
						result = append(result, rowObj)
					}
				}
			} else {
				// No zipping needed
				rowObj := make(map[string]interface{})
				for i, val := range evaluatedRow {
					if i >= len(keys) {
						break
					}
					key := keys[i]
					rowObj[key] = val
				}

				// Apply Map Clause if present
				if e.Map != nil {
					mapCtx := make(map[string]interface{})
					for k, v := range ctx {
						mapCtx[k] = v
					}
					mapCtx[e.Map.Param.Value] = rowObj

					mappedVal, err := t.evalExpression(e.Map.Body, mapCtx)
					if err != nil {
						return nil, err
					}
					result = append(result, mappedVal)
				} else {
					result = append(result, rowObj)
				}
			}
		}
		return result, nil
	case *ast.BinaryExpression:
		left, err := t.evalExpression(e.Left, ctx)
		if err != nil {
			return nil, err
		}
		right, err := t.evalExpression(e.Right, ctx)
		if err != nil {
			return nil, err
		}

		return t.evalBinary(left, e.Operator, right)
	case *ast.ConditionalExpression:
		condition, err := t.evalExpression(e.Condition, ctx)
		if err != nil {
			return nil, err
		}

		// Convert to boolean
		isTruthy := t.isTruthy(condition)

		if isTruthy {
			return t.evalExpression(e.Consequence, ctx)
		} else {
			return t.evalExpression(e.Alternative, ctx)
		}
	case *ast.MemberExpression:
		left, err := t.evalExpression(e.Left, ctx)
		if err != nil {
			return nil, err
		}

		if obj, ok := left.(map[string]interface{}); ok {
			prop := e.Property.Value
			if val, ok := obj[prop]; ok {
				return val, nil
			}
			return nil, t.errfNodeMsg(e, ie.PropertyNotFound(prop))
		}
		return nil, t.errfNodeMsg(e, ie.NotAnObject())
	default:
		return nil, t.errfNode(expr, "unknown expression type: %T", expr)
	}
}

func (t *Transpiler) evalBinary(left interface{}, op string, right interface{}) (interface{}, error) {
	switch op {
	case "+":
		// String concatenation
		if lStr, ok := left.(string); ok {
			return lStr + fmt.Sprintf("%v", right), nil
		}
		if rStr, ok := right.(string); ok {
			return fmt.Sprintf("%v", left) + rStr, nil
		}
		// Numeric addition
		lFloat, lIsFloat := toFloat(left)
		rFloat, rIsFloat := toFloat(right)
		if lIsFloat || rIsFloat {
			return lFloat + rFloat, nil
		}
		if lInt, ok := left.(int64); ok {
			if rInt, ok := right.(int64); ok {
				return lInt + rInt, nil
			}
		}
	case "-":
		// Numeric subtraction
		lFloat, lIsFloat := toFloat(left)
		rFloat, rIsFloat := toFloat(right)
		if lIsFloat || rIsFloat {
			return lFloat - rFloat, nil
		}
		if lInt, ok := left.(int64); ok {
			if rInt, ok := right.(int64); ok {
				return lInt - rInt, nil
			}
		}
	case "*":
		// Numeric multiplication
		lFloat, lIsFloat := toFloat(left)
		rFloat, rIsFloat := toFloat(right)
		if lIsFloat || rIsFloat {
			return lFloat * rFloat, nil
		}
		if lInt, ok := left.(int64); ok {
			if rInt, ok := right.(int64); ok {
				return lInt * rInt, nil
			}
		}
	case "/":
		// Numeric division
		lFloat, lIsFloat := toFloat(left)
		rFloat, rIsFloat := toFloat(right)
		if lIsFloat || rIsFloat {
			if rFloat == 0 {
				return nil, t.errMsg(ie.DivisionByZero())
			}
			return lFloat / rFloat, nil
		}
		if lInt, ok := left.(int64); ok {
			if rInt, ok := right.(int64); ok {
				if rInt == 0 {
					return nil, t.errMsg(ie.DivisionByZero())
				}
				return lInt / rInt, nil
			}
		}
	case "%":
		// Int modulo
		if lInt, ok := left.(int64); ok {
			if rInt, ok := right.(int64); ok {
				if rInt == 0 {
					return nil, t.errMsg(ie.ModuloByZero())
				}
				return lInt % rInt, nil
			}
		}
	case "==":
		return t.compareEqual(left, right), nil
	case "!=":
		return !t.compareEqual(left, right), nil
	case "<":
		return t.compareLess(left, right)
	case ">":
		return t.compareLess(right, left)
	case "<=":
		eq := t.compareEqual(left, right)
		if eq {
			return true, nil
		}
		return t.compareLess(left, right)
	case ">=":
		eq := t.compareEqual(left, right)
		if eq {
			return true, nil
		}
		return t.compareLess(right, left)
	}
	return nil, t.errMsg(ie.UnsupportedBinaryOp(left, op, right))
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	}
	return 0, false
}

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

// wtf???
