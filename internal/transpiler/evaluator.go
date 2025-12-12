package transpiler

import (
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/token"
	"strings"
)

// evalExpression evaluates an AST expression and returns its value
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
	case *ast.ValidatorExpression:
		return t.generateValidatorValue(e)
	case *ast.PresetReference:
		return t.evalPresetReference(e, ctx)
	case *ast.MapExpression:
		return t.evalMapExpression(e, ctx)
	case *ast.InterpolatedString:
		return t.evalInterpolatedString(e, ctx)
	case *ast.Identifier:
		return t.evalIdentifier(e, ctx)
	case *ast.ObjectLiteral:
		return t.evalObjectLiteral(e, ctx)
	case *ast.ArrayLiteral:
		return t.evalArrayLiteral(e, ctx)
	case *ast.RangeExpression:
		return t.evalRangeExpression(e, ctx)
	case *ast.ArrayTemplate:
		return t.evalArrayTemplate(e, ctx)
	case *ast.BinaryExpression:
		return t.evalBinaryExpression(e, ctx)
	case *ast.ConditionalExpression:
		return t.evalConditionalExpression(e, ctx)
	case *ast.MemberExpression:
		return t.evalMemberExpression(e, ctx)
	default:
		return nil, t.errfNode(expr, "unknown expression type: %T", expr)
	}
}

// evalPresetReference evaluates a preset reference expression
func (t *Transpiler) evalPresetReference(e *ast.PresetReference, ctx map[string]interface{}) (interface{}, error) {
	presetName := e.Name.Value
	presetBody, exists := t.presetTable[presetName]
	if !exists {
		return nil, t.errfNode(e, "preset %q not found — define it with @preset %q { ... }", presetName, presetName)
	}

	// Evaluate the preset body to get base values
	baseVal, err := t.evalExpression(presetBody, ctx)
	if err != nil {
		return nil, err
	}

	baseObj, ok := baseVal.(map[string]interface{})
	if !ok {
		return nil, t.errfNode(e, "preset %q did not evaluate to an object", presetName)
	}

	// If there are overrides, merge them on top
	if e.Overrides != nil {
		overridesVal, err := t.evalExpression(e.Overrides, ctx)
		if err != nil {
			return nil, err
		}

		overridesObj, ok := overridesVal.(map[string]interface{})
		if !ok {
			return nil, t.errfNode(e, "preset overrides must be an object")
		}

		// Merge: overrides take precedence
		result := make(map[string]interface{})
		for k, v := range baseObj {
			result[k] = v
		}
		for k, v := range overridesObj {
			result[k] = v
		}
		return result, nil
	}

	// Return a copy to avoid mutation issues
	result := make(map[string]interface{})
	for k, v := range baseObj {
		result[k] = v
	}
	return result, nil
}

// evalMapExpression evaluates a map expression
func (t *Transpiler) evalMapExpression(e *ast.MapExpression, ctx map[string]interface{}) (interface{}, error) {
	// Evaluate the array to be mapped
	leftVal, err := t.evalExpression(e.Left, ctx)
	if err != nil {
		return nil, err
	}

	// Ensure it's an array (slice). Accept RangeResult (wraps generated slice) as well.
	var items []interface{}
	switch v := leftVal.(type) {
	case []interface{}:
		items = v
	case RangeResult:
		items = v.Values
	default:
		return nil, t.errfNode(e, "map target is not an array, it's a %T — gremlin is confused", leftVal)
	}

	var result []interface{}

	// Iterate and map
	for _, item := range items {
		// Create a new scope for the iteration
		// We copy the current context to allow access to outer variables
		newCtx := make(map[string]interface{})
		if ctx != nil {
			for k, v := range ctx {
				newCtx[k] = v
			}
		}
		// Bind the iterator variable
		newCtx[e.Iterator.Value] = item
		// Evaluate the body
		mappedVal, err := t.evalExpression(e.Body, newCtx)
		if err != nil {
			return nil, err
		}
		result = append(result, mappedVal)
	}
	return result, nil
}

// evalInterpolatedString evaluates an interpolated string
func (t *Transpiler) evalInterpolatedString(e *ast.InterpolatedString, ctx map[string]interface{}) (interface{}, error) {
	var result strings.Builder
	for _, part := range e.Parts {
		switch p := part.(type) {
		case string:
			result.WriteString(p)
		case ast.Expression:
			if ident, ok := p.(*ast.Identifier); ok {
				found := false
				if ctx != nil {
					_, found = ctx[ident.Value]
				}

				if !found {
					if e.Token.Type == token.TEMPLATESTR {
						result.WriteString("${")
						result.WriteString(ident.Value)
						result.WriteString("}")
					} else {
						// Raw string uses {var}
						result.WriteString("{")
						result.WriteString(ident.Value)
						result.WriteString("}")
					}
					continue
				}
			}

			val, err := t.evalExpression(p, ctx)
			if err != nil {
				return nil, err
			}
			result.WriteString(fmt.Sprintf("%v", val))
		}
	}
	return result.String(), nil
}

// evalIdentifier evaluates an identifier
func (t *Transpiler) evalIdentifier(e *ast.Identifier, ctx map[string]interface{}) (interface{}, error) {
	// Variable lookup: check context (local) first, then symbol table (global)
	if ctx != nil {
		if val, ok := ctx[e.Value]; ok {
			return val, nil
		}
	}
	if val, ok := t.symbolTable[e.Value]; ok {
		return val, nil
	}
	return e.Value, nil
}

// evalObjectLiteral evaluates an object literal
func (t *Transpiler) evalObjectLiteral(e *ast.ObjectLiteral, ctx map[string]interface{}) (interface{}, error) {
	obj := make(map[string]interface{})

	localCtx := make(map[string]interface{})
	if ctx != nil {
		for k, v := range ctx {
			localCtx[k] = v
		}
	}

	for _, decl := range e.Declarations {
		val, err := t.evalExpression(decl.Value, localCtx)
		if err != nil {
			return nil, err
		}
		localCtx[decl.Name.Value] = val
	}

	// Evaluate properties using local context
	for _, key := range e.Keys {
		valExpr := e.Properties[key]
		if valExpr == nil {
			continue
		}
		val, err := t.evalExpression(valExpr, localCtx)
		if err != nil {
			return nil, err
		}
		obj[key] = val
	}
	return obj, nil
}

// evalArrayLiteral evaluates an array literal
func (t *Transpiler) evalArrayLiteral(e *ast.ArrayLiteral, ctx map[string]interface{}) (interface{}, error) {
	arr := make([]interface{}, 0, len(e.Elements))
	for _, el := range e.Elements {
		val, err := t.evalExpression(el, ctx)
		if err != nil {
			return nil, err
		}
		// If the element is a RangeResult, flatten it
		if rr, ok := val.(RangeResult); ok {
			arr = append(arr, rr.Values...)
		} else {
			arr = append(arr, val)
		}
	}
	return arr, nil
}

// evalRangeExpression evaluates a range expression
func (t *Transpiler) evalRangeExpression(e *ast.RangeExpression, ctx map[string]interface{}) (interface{}, error) {
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

	return t.evalIntegerRange(sInt, eInt, stepV, e)
}

// evalArrayTemplate evaluates an array template expression
func (t *Transpiler) evalArrayTemplate(e *ast.ArrayTemplate, ctx map[string]interface{}) (interface{}, error) {
	result := make([]interface{}, 0, len(e.Rows))
	keys := e.Template.Keys

	// Detect if this is an implicit template (single field matching map parameter)
	isImplicitTemplate := false
	if e.Map != nil && len(keys) == 1 && keys[0] == e.Map.Param.Value {
		isImplicitTemplate = true
	}

	for _, row := range e.Rows {
		// First, evaluate all expressions in the row
		evaluatedRow := make([]interface{}, len(row))
		for i, expr := range row {
			val, err := t.evalExpression(expr, ctx)
			if err != nil {
				return nil, err
			}
			if rr, ok := val.(RangeResult); ok {
				evaluatedRow[i] = rr.Values
			} else {
				evaluatedRow[i] = val
			}
		}

		// Check if we have ranges that need zipping
		hasArrays := false
		minArrayLength := -1

		for _, val := range evaluatedRow {
			if arr, ok := val.([]interface{}); ok {
				isObjectArray := false
				if len(arr) > 0 {
					if _, isMap := arr[0].(map[string]interface{}); isMap {
						isObjectArray = true
					}
				}

				if !isObjectArray {
					hasArrays = true
					if minArrayLength == -1 || len(arr) < minArrayLength {
						minArrayLength = len(arr)
					}
				}
			}
		}

		// Range Zipping: if we have arrays, zip them
		if hasArrays && minArrayLength > 0 {
			for idx := 0; idx < minArrayLength; idx++ {
				var itemValue interface{}

				if isImplicitTemplate {
					if arr, ok := evaluatedRow[0].([]interface{}); ok {
						itemValue = arr[idx]
					} else {
						itemValue = evaluatedRow[0]
					}
				} else {
					rowObj := make(map[string]interface{})
					for i, val := range evaluatedRow {
						if i >= len(keys) {
							break
						}
						key := keys[i]

						if arr, ok := val.([]interface{}); ok {
							rowObj[key] = arr[idx]
						} else {
							rowObj[key] = val
						}
					}
					itemValue = rowObj
				}

				// Apply Map Clause if present
				if e.Map != nil {
					mapCtx := make(map[string]interface{})
					for k, v := range ctx {
						mapCtx[k] = v
					}
					mapCtx[e.Map.Param.Value] = itemValue

					mappedVal, err := t.evalExpression(e.Map.Body, mapCtx)
					if err != nil {
						return nil, err
					}
					result = append(result, mappedVal)
				} else {
					result = append(result, itemValue)
				}
			}
		} else {
			// No zipping needed
			var itemValue interface{}

			if isImplicitTemplate {
				itemValue = evaluatedRow[0]
			} else {
				rowObj := make(map[string]interface{})
				for i, val := range evaluatedRow {
					if i >= len(keys) {
						break
					}
					key := keys[i]
					rowObj[key] = val
				}
				itemValue = rowObj
			}

			// Apply Map Clause if present
			if e.Map != nil {
				mapCtx := make(map[string]interface{})
				for k, v := range ctx {
					mapCtx[k] = v
				}
				mapCtx[e.Map.Param.Value] = itemValue

				mappedVal, err := t.evalExpression(e.Map.Body, mapCtx)
				if err != nil {
					return nil, err
				}
				result = append(result, mappedVal)
			} else {
				result = append(result, itemValue)
			}
		}
	}
	return result, nil
}

// evalBinaryExpression evaluates a binary expression
func (t *Transpiler) evalBinaryExpression(e *ast.BinaryExpression, ctx map[string]interface{}) (interface{}, error) {
	left, err := t.evalExpression(e.Left, ctx)
	if err != nil {
		return nil, err
	}
	right, err := t.evalExpression(e.Right, ctx)
	if err != nil {
		return nil, err
	}

	return t.evalBinary(left, e.Operator, right)
}

// evalConditionalExpression evaluates a ternary conditional expression
func (t *Transpiler) evalConditionalExpression(e *ast.ConditionalExpression, ctx map[string]interface{}) (interface{}, error) {
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
}

// evalMemberExpression evaluates a member access expression
func (t *Transpiler) evalMemberExpression(e *ast.MemberExpression, ctx map[string]interface{}) (interface{}, error) {
	leftVal, err := t.evalExpression(e.Left, ctx)
	if err != nil {
		return nil, err
	}

	// Handle map access
	if obj, ok := leftVal.(map[string]interface{}); ok {
		if val, ok := obj[e.Property.Value]; ok {
			return val, nil
		}
		return nil, t.errfNode(e, "property %q not found — gremlin searched everywhere", e.Property.Value)
	}
	return nil, t.errfNodeMsg(e, ie.NotAnObject())
}
