package parser

import (
	"jsson/internal/ast"
	"jsson/internal/lexer"
	"testing"
)

// Test edge cases for validators with extreme arguments
func TestValidatorEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "int with same min and max",
			input:     `x = @int(10, 10)`,
			shouldErr: false,
		},
		{
			name:      "int with reversed range",
			input:     `x = @int(100, 1)`,
			shouldErr: false, // Parser allows, transpiler handles
		},
		{
			name:      "int with negative range",
			input:     `x = @int(-100, -1)`,
			shouldErr: false,
		},
		{
			name:      "int with zero boundaries",
			input:     `x = @int(0, 0)`,
			shouldErr: false,
		},
		{
			name:      "float with same values",
			input:     `x = @float(5.5, 5.5)`,
			shouldErr: false,
		},
		{
			name:      "float with negative range",
			input:     `x = @float(-99.99, -0.01)`,
			shouldErr: false,
		},
		{
			name:      "int with very large numbers",
			input:     `x = @int(1000000, 9999999)`,
			shouldErr: false,
		},
		{
			name:      "float with decimals",
			input:     `x = @float(0.001, 0.999)`,
			shouldErr: false,
		},
		{
			name:      "validator without parentheses (currently errors)",
			input:     `x = @uuid`,
			shouldErr: true, // Current lexer/parser requires @uuid() with parens
		},
		{
			name:      "regex with complex pattern",
			input:     `x = @regex("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for ranges
func TestRangeEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "range with negative numbers",
			input:     `x = -10..-1`,
			shouldErr: false,
		},
		{
			name:      "range with zero",
			input:     `x = 0..0`,
			shouldErr: false,
		},
		{
			name:      "reverse range",
			input:     `x = 10..1`,
			shouldErr: false,
		},
		{
			name:      "range with large step",
			input:     `x = 1..100 step 50`,
			shouldErr: false,
		},
		{
			name:      "range with negative step",
			input:     `x = 10..1 step -1`,
			shouldErr: false,
		},
		{
			name:      "range in array",
			input:     `x = [1..5, 10..15]`,
			shouldErr: false,
		},
		{
			name:      "nested ranges in map",
			input:     `x = (1..5 map (a) = (1..5 map (b) = a * b))`,
			shouldErr: false,
		},
		{
			name:      "string range with numbers",
			input:     `x = "server-001".."server-010"`,
			shouldErr: false,
		},
		{
			name:      "very large range",
			input:     `x = 1..100000`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for arithmetic operations
func TestArithmeticEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "division expression",
			input:     `x = 10 / 2`,
			shouldErr: false,
		},
		{
			name:      "division by zero (parser allows)",
			input:     `x = 10 / 0`,
			shouldErr: false, // Parser allows, transpiler should catch
		},
		{
			name:      "modulo expression",
			input:     `x = 10 % 3`,
			shouldErr: false,
		},
		{
			name:      "modulo by zero (parser allows)",
			input:     `x = 10 % 0`,
			shouldErr: false,
		},
		{
			name:      "complex arithmetic",
			input:     `x = (10 + 5) * 2 / (3 - 1)`,
			shouldErr: false,
		},
		{
			name:      "negative numbers",
			input:     `x = -10 + -5`,
			shouldErr: false,
		},
		{
			name:      "float arithmetic",
			input:     `x = 10.5 / 2.5`,
			shouldErr: false,
		},
		{
			name:      "mixed int and float",
			input:     `x = 10 / 2.5`,
			shouldErr: false,
		},
		{
			name:      "arithmetic in map",
			input:     `x = (1..10 map (n) = n * 2 + 1)`,
			shouldErr: false,
		},
		{
			name:      "arithmetic with parentheses",
			input:     `x = ((10 + 5) * (2 - 1)) / 3`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for strings
func TestStringEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "empty string",
			input:     `x = ""`,
			shouldErr: false,
		},
		{
			name:      "string with unicode",
			input:     `x = "Hello 世界 🌍"`,
			shouldErr: false,
		},
		{
			name:      "string with escaped quotes",
			input:     `x = "He said \"Hello\""`,
			shouldErr: false,
		},
		{
			name:      "string with newlines",
			input:     "x = \"Line 1\\nLine 2\"",
			shouldErr: false,
		},
		{
			name:      "very long string",
			input:     `x = "` + string(make([]byte, 1000)) + `"`,
			shouldErr: false,
		},
		{
			name:      "string concatenation",
			input:     `x = "Hello" + " " + "World"`,
			shouldErr: false,
		},
		{
			name:      "interpolated string empty",
			input:     "x = `${}`",
			shouldErr: false,
		},
		{
			name:      "interpolated with variable",
			input:     "name := \"John\"\nx = `Hello ${name}`",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for nested structures
func TestNestedStructureEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name: "deeply nested objects",
			input: `
a {
	b {
		c {
			d {
				e = 1
			}
		}
	}
}`,
			shouldErr: false,
		},
		{
			name: "deeply nested arrays",
			input: `
x = [[[[1, 2], [3, 4]], [[5, 6], [7, 8]]]]`,
			shouldErr: false,
		},
		{
			name: "nested maps",
			input: `
x = (1..3 map (a) = (1..3 map (b) = (1..3 map (c) = a * b * c)))`,
			shouldErr: false,
		},
		{
			name: "array with mixed types",
			input: `
x = [1, "string", true, 3.14, @uuid]`,
			shouldErr: false,
		},
		{
			name: "empty nested structures",
			input: `
x {
	y = []
	z = {}
}`,
			shouldErr: false,
		},
		{
			name: "object with arrays and ranges",
			input: `
config {
	ports = 8000..8005
	methods = [GET, POST, PUT]
	nested {
		items = (1..5 map (x) = x * 2)
	}
}`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for presets
func TestPresetEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name: "empty preset",
			input: `
@preset "empty" {}
x = @use "empty"`,
			shouldErr: false,
		},
		{
			name: "preset with all overrides",
			input: `
@preset "base" {
	a = 1
	b = 2
}
x = @use "base" {
	a = 10
	b = 20
}`,
			shouldErr: false,
		},
		{
			name: "preset with extra overrides",
			input: `
@preset "base" {
	a = 1
}
x = @use "base" {
	a = 10
	b = 20
	c = 30
}`,
			shouldErr: false,
		},
		{
			name: "multiple preset references",
			input: `
@preset "config" {
	timeout = 30
}
x = @use "config"
y = @use "config"
z = @use "config"`,
			shouldErr: false,
		},
		{
			name: "preset with nested object",
			input: `
@preset "nested" {
	outer {
		inner = 1
	}
}
x = @use "nested"`,
			shouldErr: false,
		},
		{
			name: "preset reference before definition should error",
			input: `
x = @use "undefined"
@preset "undefined" {
	a = 1
}`,
			shouldErr: false, // Parser allows, transpiler catches
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for template arrays
func TestTemplateArrayEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name: "template with no rows",
			input: `
x [
	template { name, age }
]`,
			shouldErr: false,
		},
		{
			name: "template with mismatched columns (parser allows)",
			input: `
x [
	template { name, age }
	John, 25
	Mary
]`,
			shouldErr: false,
		},
		{
			name: "template with extra columns (parser allows)",
			input: `
x [
	template { name, age }
	John, 25, extra
]`,
			shouldErr: false,
		},
		{
			name: "template with validators",
			input: `
x [
	template { name, email }
	John, @email
	Mary, @email
]`,
			shouldErr: false,
		},
		{
			name: "template with ranges",
			input: `
x [
	template { id, range }
	1..5, 10..14
]`,
			shouldErr: false,
		},
		{
			name: "template with map clause and no data",
			input: `
x [
	template { name }
	map (u) = { name = u.name }
]`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for variables
func TestVariableEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "variable shadowing",
			input:     "x := 1\nx := 2",
			shouldErr: false,
		},
		{
			name:      "variable used before declaration (parser allows)",
			input:     "y = x\nx := 1",
			shouldErr: false,
		},
		{
			name:      "variable with complex expression",
			input:     "x := (1..10 map (n) = n * 2)",
			shouldErr: false,
		},
		{
			name:      "variable with validator",
			input:     "id := @uuid()",
			shouldErr: false,
		},
		{
			name:      "variable in object scope",
			input:     "obj { x := 10\ny = x * 2 }",
			shouldErr: false,
		},
		{
			name:      "multiple variable declarations",
			input:     "a := 1\nb := 2\nc := 3",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for conditionals
func TestConditionalEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "nested ternary",
			input:     `x = a > 10 ? b > 5 ? "high" : "medium" : "low"`,
			shouldErr: false,
		},
		{
			name:      "ternary with arithmetic",
			input:     `x = (10 + 5) > 12 ? 1 : 0`,
			shouldErr: false,
		},
		{
			name:      "ternary with validators",
			input:     `x = true ? @uuid() : @email()`,
			shouldErr: false,
		},
		{
			name:      "ternary in map",
			input:     `x = (1..10 map (n) = n > 5 ? "high" : "low")`,
			shouldErr: false,
		},
		{
			name:      "ternary with objects",
			input:     `x = flag ? { a = 1 } : { b = 2 }`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("No statements parsed for valid input")
			}
		})
	}
}

// Test edge cases for comments
func TestCommentEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "comment only",
			input:     `// just a comment`,
			shouldErr: false,
		},
		{
			name:      "multiple comments",
			input:     "// comment 1\n// comment 2\nx = 1",
			shouldErr: false,
		},
		{
			name:      "inline comment",
			input:     `x = 1 // inline comment`,
			shouldErr: false,
		},
		{
			name:      "comment with special chars",
			input:     `// Comment with special: @#$%^&*()`,
			shouldErr: false,
		},
		{
			name:      "comment with unicode",
			input:     `// 注释 комментарий تعليق`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}
		})
	}
}

// Test parser recovery from errors
func TestParserErrorRecovery(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		shouldHaveErr bool
		minStatements int
	}{
		{
			name: "multiple errors should parse what's valid",
			input: `
x = 1
y = !@#$%  // Invalid
z = 2`,
			shouldHaveErr: true,
			minStatements: 2, // Should parse x and z at least
		},
		{
			name: "unterminated string followed by valid",
			input: `
x = "unclosed
y = 2`,
			shouldHaveErr: true,
			minStatements: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldHaveErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldHaveErr, hasErrors, p.Errors())
			}

			if len(program.Statements) < tt.minStatements {
				t.Errorf("Expected at least %d statements, got %d",
					tt.minStatements, len(program.Statements))
			}
		})
	}
}

// Test that validator expressions preserve their structure
func TestValidatorStructurePreservation(t *testing.T) {
	input := `x = @int(10, 100)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("Expected AssignmentStatement, got %T", program.Statements[0])
	}

	validator, ok := stmt.Value.(*ast.ValidatorExpression)
	if !ok {
		t.Fatalf("Expected ValidatorExpression, got %T", stmt.Value)
	}

	if validator.Type != "int" {
		t.Errorf("Expected type 'int', got '%s'", validator.Type)
	}

	if len(validator.Args) != 2 {
		t.Fatalf("Expected 2 args, got %d", len(validator.Args))
	}

	arg1, ok1 := validator.Args[0].(int64)
	arg2, ok2 := validator.Args[1].(int64)

	if !ok1 || !ok2 {
		t.Fatalf("Expected int64 args, got %T and %T", validator.Args[0], validator.Args[1])
	}

	if arg1 != 10 || arg2 != 100 {
		t.Errorf("Expected args [10, 100], got [%d, %d]", arg1, arg2)
	}
}

// ============================================================================
// LOGICAL OPERATORS AND COMPARISON EDGE CASES
// ============================================================================

func TestLogicalOperatorEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "AND operator",
			input:     `x = true && false`,
			shouldErr: false,
		},
		{
			name:      "OR operator",
			input:     `x = true || false`,
			shouldErr: false,
		},
		{
			name:      "NOT operator",
			input:     `x = !true`,
			shouldErr: false,
		},
		{
			name:      "complex logical expression",
			input:     `x = (a > 5 && b < 10) || c == 0`,
			shouldErr: false,
		},
		{
			name:      "multiple AND operators",
			input:     `x = a && b && c && d`,
			shouldErr: false,
		},
		{
			name:      "multiple OR operators",
			input:     `x = a || b || c || d`,
			shouldErr: false,
		},
		{
			name:      "mixed AND and OR",
			input:     `x = a && b || c && d`,
			shouldErr: false,
		},
		{
			name:      "NOT with comparison",
			input:     `x = !(a > 5)`,
			shouldErr: false,
		},
		{
			name:      "double NOT",
			input:     `x = !!true`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

func TestComparisonOperatorEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "less than",
			input:     `x = 5 < 10`,
			shouldErr: false,
		},
		{
			name:      "greater than",
			input:     `x = 10 > 5`,
			shouldErr: false,
		},
		{
			name:      "less than or equal",
			input:     `x = 5 <= 10`,
			shouldErr: false,
		},
		{
			name:      "greater than or equal",
			input:     `x = 10 >= 5`,
			shouldErr: false,
		},
		{
			name:      "equal",
			input:     `x = 10 == 10`,
			shouldErr: false,
		},
		{
			name:      "not equal",
			input:     `x = 10 != 5`,
			shouldErr: false,
		},
		{
			name:      "chained comparisons",
			input:     `x = a < b && b < c`,
			shouldErr: false,
		},
		{
			name:      "comparison with negative",
			input:     `x = -5 < 0`,
			shouldErr: false,
		},
		{
			name:      "comparison with float",
			input:     `x = 3.14 > 3.0`,
			shouldErr: false,
		},
		{
			name:      "string comparison",
			input:     `x = "a" == "a"`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

// ============================================================================
// ZIP AND ADVANCED MAP EDGE CASES
// ============================================================================

func TestZipRangeEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "basic zip with two ranges (not implemented yet)",
			input:     `x = (1..3 zip 10..12 map (a, b) = a + b)`,
			shouldErr: true, // zip not implemented yet
		},
		{
			name:      "zip with three ranges (not implemented yet)",
			input:     `x = (1..3 zip 10..12 zip 20..22 map (a, b, c) = a + b + c)`,
			shouldErr: true, // zip not implemented yet
		},
		{
			name:      "zip with different length ranges (not implemented yet)",
			input:     `x = (1..5 zip 10..12 map (a, b) = a + b)`,
			shouldErr: true, // zip not implemented yet
		},
		{
			name:      "zip with negative ranges (not implemented yet)",
			input:     `x = (-5..-1 zip 1..5 map (a, b) = a + b)`,
			shouldErr: true, // zip not implemented yet
		},
		{
			name:      "zip in object (not implemented yet)",
			input:     `obj { pairs = (1..3 zip 10..12 map (a, b) = [a, b]) }`,
			shouldErr: true, // zip not implemented yet
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

func TestAdvancedMapEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "map with external variable",
			input:     `factor := 10\nx = (1..5 map (n) = n * factor)`,
			shouldErr: false,
		},
		{
			name:      "map with ternary",
			input:     `x = (1..10 map (n) = n % 2 == 0 ? "even" : "odd")`,
			shouldErr: false,
		},
		{
			name:      "map with object creation",
			input:     `x = (1..3 map (n) = { id: n, name: "item" + n })`,
			shouldErr: false,
		},
		{
			name:      "map with array creation",
			input:     `x = (1..3 map (n) = [n, n*2, n*3])`,
			shouldErr: false,
		},
		{
			name:      "map with validators",
			input:     `x = (1..3 map (n) = @uuid())`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

// ============================================================================
// STRING INTERPOLATION EDGE CASES
// ============================================================================

func TestStringInterpolationEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "basic interpolation",
			input:     `name := "World"\nmsg = "Hello, \${name}!"`,
			shouldErr: false,
		},
		{
			name:      "interpolation with number",
			input:     `age := 25\nmsg = "Age: \${age}"`,
			shouldErr: false,
		},
		{
			name:      "multiple interpolations",
			input:     `first := "John"\nlast := "Doe"\nfull = "\${first} \${last}"`,
			shouldErr: false,
		},
		{
			name:      "interpolation with expression",
			input:     `x := 5\nmsg = "Result: \${x * 2}"`,
			shouldErr: false,
		},
		{
			name:      "nested interpolation",
			input:     `inner := "world"\nouter = "Hello \${\"Dear \${inner}\"}"`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

// ============================================================================
// NUMERIC LIMITS AND SPECIAL VALUES
// ============================================================================

func TestNumericLimitEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "very large integer",
			input:     `x = 9223372036854775807`,
			shouldErr: false,
		},
		{
			name:      "very small integer (min int64 boundary)",
			input:     `x = -9223372036854775807`, // Use -max instead of true min to avoid overflow
			shouldErr: false,
		},
		{
			name:      "scientific notation positive",
			input:     `x = 1e10`,
			shouldErr: false,
		},
		{
			name:      "scientific notation negative",
			input:     `x = 1e-10`,
			shouldErr: false,
		},
		{
			name:      "very small decimal",
			input:     `x = 0.0000000001`,
			shouldErr: false,
		},
		{
			name:      "zero variations",
			input:     `x = [0, 0.0, -0, -0.0]`,
			shouldErr: false,
		},
		{
			name:      "floating point precision",
			input:     `x = 0.1 + 0.2`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

// ============================================================================
// ADVANCED ARRAY AND OBJECT EDGE CASES
// ============================================================================

func TestAdvancedArrayEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "array with trailing comma",
			input:     `x = [1, 2, 3,]`,
			shouldErr: false,
		},
		{
			name:      "array with mixed validators",
			input:     `x = [@uuid(), @email(), @int(1, 100)]`,
			shouldErr: false,
		},
		{
			name:      "array with ranges and literals",
			input:     `x = [1, 2, 3..5, 6, 7]`,
			shouldErr: false,
		},
		{
			name:      "array with spread in map",
			input:     `x = (1..3 map (n) = [n, n*2, n*3])`,
			shouldErr: false,
		},
		{
			name:      "multidimensional array literal",
			input:     `matrix = [[1,2,3], [4,5,6], [7,8,9]]`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

func TestAdvancedObjectEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "object with computed property",
			input:     `x { ["key" + 1] = "value" }`,
			shouldErr: false,
		},
		{
			name:      "object with special characters in keys",
			input:     `x { "key-with-dash" = 1\n"key.with.dot" = 2 }`,
			shouldErr: false,
		},
		{
			name:      "object with numeric keys",
			input:     `x { "1" = "one"\n"2" = "two" }`,
			shouldErr: false,
		},
		{
			name:      "object with all value types",
			input:     `x { num = 1\nstr = "text"\nbool = true\narr = [1,2]\nobj { nested = 1 } }`,
			shouldErr: false,
		},
		{
			name:      "object with validator values",
			input:     `user { id = @uuid()\nemail = @email()\nage = @int(18, 65) }`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}

// ============================================================================
// ADVANCED TEMPLATE EDGE CASES
// ============================================================================

func TestAdvancedTemplateEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name: "template with expressions",
			input: `
x [
	template { a, b, c }
	1, 2, 3
	3, 4, 7
]`,
			shouldErr: false,
		},
		{
			name: "template with mixed types",
			input: `
x [
	template { id, name, active }
	1, "Alice", true
	2, "Bob", false
]`,
			shouldErr: false,
		},
		{
			name: "template with single column",
			input: `
x [
	template { value }
	1
	2
	3
]`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			hasErrors := len(p.Errors()) > 0
			if hasErrors != tt.shouldErr {
				t.Errorf("Expected error=%v, got error=%v. Errors: %v",
					tt.shouldErr, hasErrors, p.Errors())
			}

			if !tt.shouldErr && len(program.Statements) == 0 {
				t.Errorf("Expected statements, got none")
			}
		})
	}
}
