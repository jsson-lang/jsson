package transpiler

import (
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
)

// errf creates an error with the transpiler prefix and source context
func (t *Transpiler) errf(format string, args ...interface{}) error {
	prefix := "Transpile gremlin:"
	if t != nil && t.sourceFile != "" {
		ctx := ie.FormatContext(t.sourceFile, 1, 1)
		return fmt.Errorf("%s %s — %s", prefix, ctx, fmt.Sprintf(format, args...))
	}
	return fmt.Errorf("%s — %s", prefix, fmt.Sprintf(format, args...))
}

// errfNode creates an error with the transpiler prefix and node position context
func (t *Transpiler) errfNode(node ast.Node, format string, args ...interface{}) error {
	prefix := "Transpile gremlin:"
	line, col := 0, 0
	if node != nil {
		line, col = node.Position()
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
	line, col := 0, 0
	if node != nil {
		line, col = node.Position()
	}

	if t != nil && t.sourceFile != "" {
		if line > 0 && col > 0 {
			ctx := ie.FormatContext(t.sourceFile, line, col)
			return fmt.Errorf("%s %s — %s", prefix, ctx, msg)
		}

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
