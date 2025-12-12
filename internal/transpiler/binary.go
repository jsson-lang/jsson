package transpiler

import (
	"fmt"
	ie "jsson/internal/errors"
)

// evalBinary evaluates a binary expression
func (t *Transpiler) evalBinary(left interface{}, op string, right interface{}) (interface{}, error) {
	// Prevent applying numeric/string operators directly to a RangeResult
	if _, ok := left.(RangeResult); ok {
		return nil, t.errMsg(fmt.Sprintf("cannot apply operator %q to a range — expand it or use in an array context", op))
	}
	if _, ok := right.(RangeResult); ok {
		return nil, t.errMsg(fmt.Sprintf("cannot apply operator %q to a range — expand it or use in an array context", op))
	}

	switch op {
	case "+":
		return t.evalAddition(left, right)
	case "-":
		return t.evalSubtraction(left, right)
	case "*":
		return t.evalMultiplication(left, right)
	case "/":
		return t.evalDivision(left, right)
	case "%":
		return t.evalModulo(left, right)
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
	case "&&":
		return t.isTruthy(left) && t.isTruthy(right), nil
	case "||":
		return t.isTruthy(left) || t.isTruthy(right), nil
	}
	return nil, t.errMsg(ie.UnsupportedBinaryOp(left, op, right))
}

// evalAddition handles + operator
func (t *Transpiler) evalAddition(left, right interface{}) (interface{}, error) {
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
	return nil, t.errMsg(ie.UnsupportedBinaryOp(left, "+", right))
}

// evalSubtraction handles - operator
func (t *Transpiler) evalSubtraction(left, right interface{}) (interface{}, error) {
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
	return nil, t.errMsg(ie.UnsupportedBinaryOp(left, "-", right))
}

// evalMultiplication handles * operator
func (t *Transpiler) evalMultiplication(left, right interface{}) (interface{}, error) {
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
	return nil, t.errMsg(ie.UnsupportedBinaryOp(left, "*", right))
}

// evalDivision handles / operator
func (t *Transpiler) evalDivision(left, right interface{}) (interface{}, error) {
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
	return nil, t.errMsg(ie.UnsupportedBinaryOp(left, "/", right))
}

// evalModulo handles % operator
func (t *Transpiler) evalModulo(left, right interface{}) (interface{}, error) {
	if lInt, okL := toInt64(left); okL {
		if rInt, okR := toInt64(right); okR {
			if rInt == 0 {
				return nil, t.errMsg(ie.ModuloByZero())
			}
			return lInt % rInt, nil
		}
	}
	return nil, t.errMsg(ie.UnsupportedBinaryOp(left, "%", right))
}
