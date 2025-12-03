package errors

import "fmt"

// LexerError formats a lexer error with the "Lexer goblin" prefix and fun messaging
func LexerError(sourceFile string, line, col int, format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	if sourceFile != "" {
		ctx := FormatContext(sourceFile, line, col)
		return fmt.Sprintf("Lexer goblin: %s — %s", ctx, msg)
	}
	ctx := fmt.Sprintf("%d:%d", line, col)
	return fmt.Sprintf("Lexer goblin: %s — %s", ctx, msg)
}

// ParserError formats a parser error with the "Syntax wizard" prefix and fun messaging
func ParserError(sourceFile string, line, col int, format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	if sourceFile != "" {
		ctx := FormatContext(sourceFile, line, col)
		return fmt.Sprintf("Syntax wizard: %s — %s", ctx, msg)
	}
	ctx := fmt.Sprintf("%d:%d", line, col)
	return fmt.Sprintf("Syntax wizard: %s — %s", ctx, msg)
}

// TranspilerError formats a transpiler error with the "Transpile gremlin" prefix and fun messaging
func TranspilerError(sourceFile string, line, col int, format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	if sourceFile != "" {
		ctx := FormatContext(sourceFile, line, col)
		return fmt.Sprintf("Transpile gremlin: %s — %s", ctx, msg)
	}
	ctx := fmt.Sprintf("%d:%d", line, col)
	return fmt.Sprintf("Transpile gremlin: %s — %s", ctx, msg)
}

// Common fun error messages for reuse

// UnterminatedString returns a fun message for unterminated string literals
func UnterminatedString() string {
	return "found an endless string (missing closing quote)"
}

// IllegalCharacter returns a fun message for illegal characters
func IllegalCharacter(ch rune) string {
	return fmt.Sprintf("stumbled upon a strange character: %q", ch)
}

// ExpectedToken returns a fun message for expected tokens
func ExpectedToken(expected, got string) string {
	return fmt.Sprintf("expected %s but found %s instead", expected, got)
}

// MissingClosingBrace returns a fun message for missing closing braces
func MissingClosingBrace() string {
	return "expected '}' — wizard can't find the closing brace"
}

// MissingClosingParen returns a fun message for missing closing parentheses
func MissingClosingParen() string {
	return "expected ')' — wizard needs balanced parentheses"
}

// PropertyNotFound returns a fun message for missing properties
func PropertyNotFound(prop string) string {
	return fmt.Sprintf("property %q not found — gremlin searched everywhere", prop)
}

// NotAnObject returns a fun message when expecting an object
func NotAnObject() string {
	return "left side of '.' is not an object — gremlin expected a map"
}

// CyclicInclude returns a fun message for cyclic includes
func CyclicInclude(path string) string {
	return fmt.Sprintf("cyclic include detected: %s — gremlin is going in circles!", path)
}

// RangeBoundsNotIntegers returns a fun message for non-integer range bounds
func RangeBoundsNotIntegers(start, end interface{}) string {
	return fmt.Sprintf("range bounds must be integers: %v .. %v — gremlin can't count with those", start, end)
}

// StepNotInteger returns a fun message for non-integer step values
func StepNotInteger(step interface{}) string {
	return fmt.Sprintf("step must be integer: %v — gremlin needs whole numbers to step", step)
}

// StepCannotBeZero returns a fun message for zero step values
func StepCannotBeZero() string {
	return "step cannot be 0 — gremlin would be stuck forever!"
}

// UnsupportedBinaryOp returns a fun message for unsupported binary operations
func UnsupportedBinaryOp(left interface{}, op string, right interface{}) string {
	return fmt.Sprintf("unsupported binary operation: %v %s %v — gremlin doesn't know how to do that math", left, op, right)
}

// IncludePathExpected returns a fun message when include needs a path
func IncludePathExpected() string {
	return "expected a path string after include — wizard needs directions"
}

// IntegerTooSpicy returns a fun message for unparseable integers
func IntegerTooSpicy(literal string) string {
	return fmt.Sprintf("could not parse %q as integer — maybe it's too spicy for me", literal)
}

// ExpectedIdentifierAfterDot returns a fun message for member access errors
func ExpectedIdentifierAfterDot() string {
	return "expected identifier after '.' — maybe use letters, not emojis"
}

// DivisionByZero returns a fun message for division by zero
func DivisionByZero() string {
	return "division by zero — even gremlins can't divide by nothing!"
}

// ModuloByZero returns a fun message for modulo by zero
func ModuloByZero() string {
	return "modulo by zero — gremlins are confused!"
}

// MissingColonInTernary returns a fun message for missing colon in ternary
func MissingColonInTernary() string {
	return "expected ':' in ternary expression — wizard needs both ? and :"
}

// UnsupportedComparison returns a fun message for unsupported comparisons
func UnsupportedComparison(left, right interface{}) string {
	return fmt.Sprintf("can't compare %v and %v — gremlin doesn't know how", left, right)
}
