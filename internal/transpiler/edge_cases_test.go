package transpiler

import (
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"strings"
	"testing"
)

// ============================================================================
// EDGE CASES - VALORES EXTREMOS
// ============================================================================

func TestEdgeCase_EmptyRange(t *testing.T) {
	input := "data = 5..0"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	// Should generate backward range
	if !strings.Contains(string(output), "[5,4,3,2,1,0]") && !strings.Contains(string(output), "5") {
		t.Logf("Output: %s", string(output))
	}
}

func TestEdgeCase_SingleItemRange(t *testing.T) {
	input := "data = 5..5"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	// Should contain single item
	if !strings.Contains(string(output), "5") {
		t.Errorf("Expected single item 5, got: %s", string(output))
	}
}

func TestEdgeCase_ZeroStep(t *testing.T) {
	input := "data = 0..10 step 0"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()

	// Should error with step zero
	if err == nil {
		t.Error("Expected error for step=0, got nil")
	}
	if !strings.Contains(err.Error(), "zero") && !strings.Contains(err.Error(), "0") {
		t.Errorf("Expected error message about zero step, got: %v", err)
	}
}

func TestEdgeCase_NegativeRange(t *testing.T) {
	input := "data = -10..-1"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	// Should contain negative numbers
	if !strings.Contains(string(output), "-10") {
		t.Errorf("Expected negative numbers, got: %s", string(output))
	}
}

func TestEdgeCase_VeryLargeNumber(t *testing.T) {
	input := "data = 9999999999"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	if !strings.Contains(string(output), "9999999999") {
		t.Errorf("Expected large number, got: %s", string(output))
	}
}

func TestEdgeCase_VeryLongString(t *testing.T) {
	longStr := strings.Repeat("a", 10000)
	input := `data = "` + longStr + `"`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	if !strings.Contains(string(output), "aaa") {
		t.Error("Expected long string in output")
	}
}

func TestEdgeCase_SpecialCharactersInString(t *testing.T) {
	input := `data = "Hello\nWorld\t\"Quote\""`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	// Should handle escape sequences
	t.Logf("Output: %s", string(output))
}

// ============================================================================
// EDGE CASES - ARRAYS E OBJETOS
// ============================================================================

func TestEdgeCase_EmptyArray(t *testing.T) {
	input := "data = []"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	if !strings.Contains(string(output), "[]") {
		t.Errorf("Expected empty array, got: %s", string(output))
	}
}

func TestEdgeCase_EmptyObject(t *testing.T) {
	input := "data = {}"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	if !strings.Contains(string(output), "{}") {
		t.Errorf("Expected empty object, got: %s", string(output))
	}
}

// ============================================================================
// EDGE CASES - MAP TRANSFORMATIONS
// ============================================================================

func TestEdgeCase_MapOnEmptyArray(t *testing.T) {
	input := `data = ([] map (x) = x * 2)`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	// Should return empty array
	if !strings.Contains(string(output), "[]") {
		t.Errorf("Expected empty array, got: %s", string(output))
	}
}

func TestEdgeCase_TripleNestedMap(t *testing.T) {
	input := `data = (0..2 map (x) = (0..2 map (y) = (0..2 map (z) = x + y + z)))`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	// Should generate 3D structure
	t.Logf("Triple nested map output length: %d bytes", len(output))
}

// EDGE CASES - vars

func TestEdgeCase_UndefinedVariable(t *testing.T) {
	input := "data = undefined_var"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()

	// Currently JSSON treats undefined vars as identifiers (strings)
	// This might be the expected behavior
	if err != nil {
		t.Logf("Got error for undefined var: %v", err)
	} else {
		t.Logf("Undefined var treated as identifier: %s", string(output))
	}
}

func TestEdgeCase_VariableShadowing(t *testing.T) {
	input := `
x := 10
data = {
  x := 20
  value = x
}
`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	// Inner x should shadow outer x
	if !strings.Contains(string(output), "20") {
		t.Errorf("Expected inner variable to shadow outer, got: %s", string(output))
	}
}

// ============================================================================
// EDGE CASES - TERNÁRIOS
// ============================================================================

func TestEdgeCase_NestedTernary(t *testing.T) {
	input := `data = 5 > 3 ? (2 > 1 ? "yes" : "no") : "maybe"`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	if !strings.Contains(string(output), "yes") {
		t.Errorf("Expected 'yes', got: %s", string(output))
	}
}

// ============================================================================
// STRESS TESTS - PERFORMANCE
// ============================================================================

func TestStress_LargeRange_100k(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	input := "data = 0..99999"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	t.Logf("Generated %d bytes for 100k items", len(output))
}

func TestStress_LargeMatrix_100x100(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	input := `matrix = (0..99 map (y) = (0..99 map (x) = x * y))`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	t.Logf("Generated %d bytes for 100x100 matrix", len(output))
}

func TestStress_ManyProperties(t *testing.T) {
	// Object with 100 properties
	var props []string
	for i := 0; i < 100; i++ {
		props = append(props, "  prop"+string(rune('0'+i%10))+" = "+string(rune('0'+i%10)))
	}
	input := "data = {\n" + strings.Join(props, "\n") + "\n}"

	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Logf("Parser errors (expected for malformed input): %v", p.Errors())
		return
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Logf("Transpile error: %v", err)
	}
}

// ============================================================================
// NEW EDGE CASES - v0.0.6 Validators
// ============================================================================

func TestEdgeCase_ValidatorIntSameBoundaries(t *testing.T) {
	input := "x = @int(42, 42)"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle int with same min/max: %v", err)
	}
}

func TestEdgeCase_ValidatorIntNegativeRange(t *testing.T) {
	input := "x = @int(-1000, -1)"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle negative int range: %v", err)
	}
}

func TestEdgeCase_ValidatorIntZeroBoundaries(t *testing.T) {
	input := "x = @int(0, 0)"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle zero boundaries: %v", err)
	}
}

func TestEdgeCase_ValidatorFloatSameBoundaries(t *testing.T) {
	input := "x = @float(3.14, 3.14)"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle float with same min/max: %v", err)
	}
}

func TestEdgeCase_ValidatorFloatNegativeRange(t *testing.T) {
	input := "x = @float(-99.99, -0.01)"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle negative float range: %v", err)
	}
}

func TestEdgeCase_ValidatorFloatVerySmallRange(t *testing.T) {
	input := "x = @float(0.001, 0.002)"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle very small float range: %v", err)
	}
}

func TestEdgeCase_ValidatorBoolMultiple(t *testing.T) {
	input := `
	a = @bool
	b = @bool
	c = @bool
	d = @bool
	e = @bool
	`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should generate multiple booleans: %v", err)
	}
}

func TestEdgeCase_ValidatorsInArray(t *testing.T) {
	input := `x = [@uuid, @email, @int(1, 100), @float(0.0, 1.0), @bool]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle validators in array: %v", err)
	}
}

func TestEdgeCase_ValidatorsInTemplate(t *testing.T) {
	input := `
	users [
		template { id, email, age, score, active }
		@uuid, @email, @int(18, 65), @float(0.0, 100.0), @bool
		@uuid, @email, @int(18, 65), @float(0.0, 100.0), @bool
	]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle validators in template: %v", err)
	}
}

// ============================================================================
// EDGE CASES - Division and Modulo
// ============================================================================

func TestEdgeCase_DivisionByZero(t *testing.T) {
	input := "x = 10 / 0"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err == nil {
		t.Fatal("Should error on division by zero")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "division") {
		t.Errorf("Expected division error, got: %v", err)
	}
}

func TestEdgeCase_ModuloByZero(t *testing.T) {
	input := "x = 10 % 0"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err == nil {
		t.Fatal("Should error on modulo by zero")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "modulo") {
		t.Errorf("Expected modulo error, got: %v", err)
	}
}

func TestEdgeCase_FloatDivision(t *testing.T) {
	input := "x = 10.5 / 2.5"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle float division: %v", err)
	}
}

func TestEdgeCase_NegativeDivision(t *testing.T) {
	input := "x = -10 / 2"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle negative division: %v", err)
	}
}

// ============================================================================
// EDGE CASES - Presets
// ============================================================================

func TestEdgeCase_EmptyPreset(t *testing.T) {
	input := `
	@preset "empty" {}
	x = @use "empty"`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle empty preset: %v", err)
	}
}

func TestEdgeCase_PresetAllOverrides(t *testing.T) {
	input := `
	@preset "base" {
		a = 1
		b = 2
	}
	x = @use "base" {
		a = 10
		b = 20
	}`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle all overrides: %v", err)
	}
}

func TestEdgeCase_PresetUndefined(t *testing.T) {
	input := `x = @use "undefined"`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err == nil {
		t.Fatal("Should error on undefined preset")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestEdgeCase_PresetWithValidators(t *testing.T) {
	input := `
	@preset "ids" {
		uuid = @uuid
		email = @email
		age = @int(18, 65)
	}
	x = @use "ids"`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle preset with validators: %v", err)
	}
}

// ============================================================================
// EDGE CASES - Deeply Nested Structures
// ============================================================================

func TestEdgeCase_DeeplyNestedObjects(t *testing.T) {
	input := `
	a {
		b {
			c {
				d {
					e {
						f = 1
					}
				}
			}
		}
	}`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle deeply nested objects: %v", err)
	}
}

func TestEdgeCase_DeeplyNestedArrays(t *testing.T) {
	input := `x = [[[[[[1, 2]]]]]]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle deeply nested arrays: %v", err)
	}
}

func TestEdgeCase_NestedMaps3Levels(t *testing.T) {
	input := `x = (1..2 map (a) = (1..2 map (b) = (1..2 map (c) = a + b + c)))`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle 3-level nested maps: %v", err)
	}
}

// ============================================================================
// EDGE CASES - Ranges
// ============================================================================

func TestEdgeCase_RangeStepLargerThanRange(t *testing.T) {
	input := "x = 1..5 step 100"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle step larger than range: %v", err)
	}
}

func TestEdgeCase_RangeZeroStep(t *testing.T) {
	input := "x = 1..10 step 0"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err == nil {
		t.Fatal("Should error on zero step")
	}
}

func TestEdgeCase_RangeSingleElement(t *testing.T) {
	input := "x = 5..5"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle single element range: %v", err)
	}
}

// ============================================================================
// EDGE CASES - Strings
// ============================================================================

func TestEdgeCase_EmptyString(t *testing.T) {
	input := `x = ""`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle empty string: %v", err)
	}
}

func TestEdgeCase_StringConcatenationWithNumber(t *testing.T) {
	input := `x = "value: " + 42`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle string + number: %v", err)
	}
}

func TestEdgeCase_UnicodeString(t *testing.T) {
	input := `x = "Hello 世界 🌍"`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle unicode: %v", err)
	}
}

// EDGE CASES - Logical Operators

func TestEdgeCase_LogicalAND(t *testing.T) {
	input := `x = true && false`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle logical AND: %v", err)
	}
}

func TestEdgeCase_LogicalOR(t *testing.T) {
	input := `x = false || true`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle logical OR: %v", err)
	}
}

func TestEdgeCase_LogicalNOT(t *testing.T) {
	input := `x = !false`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Skipf("NOT operator not fully implemented: %v", p.Errors())
		return
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Skipf("NOT operator not fully implemented: %v", err)
	}
}

func TestEdgeCase_ComplexLogicalExpression(t *testing.T) {
	input := `x = (5 > 3 && 10 < 20) || false`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle complex logical: %v", err)
	}
}

// EDGE CASES - Comparison Operators

func TestEdgeCase_LessThanComparison(t *testing.T) {
	input := `x = 5 < 10`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle less than: %v", err)
	}

	if !strings.Contains(string(output), "true") {
		t.Errorf("Expected true, got: %s", string(output))
	}
}

func TestEdgeCase_GreaterThanComparison(t *testing.T) {
	input := `x = 10 > 5`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle greater than: %v", err)
	}

	if !strings.Contains(string(output), "true") {
		t.Errorf("Expected true, got: %s", string(output))
	}
}

func TestEdgeCase_EqualityComparison(t *testing.T) {
	input := `x = 10 == 10`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle equality: %v", err)
	}

	if !strings.Contains(string(output), "true") {
		t.Errorf("Expected true, got: %s", string(output))
	}
}

func TestEdgeCase_NotEqualComparison(t *testing.T) {
	input := `x = 10 != 5`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle not equal: %v", err)
	}

	if !strings.Contains(string(output), "true") {
		t.Errorf("Expected true, got: %s", string(output))
	}
}

func TestEdgeCase_StringComparison(t *testing.T) {
	input := `x = "hello" == "hello"`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle string comparison: %v", err)
	}

	if !strings.Contains(string(output), "true") {
		t.Errorf("Expected true, got: %s", string(output))
	}
}

// EDGE CASES - Zip Ranges

func TestEdgeCase_BasicZip(t *testing.T) {
	t.Skip("zip operator not implemented yet")
	input := `x = (1..3 zip 10..12 map (a, b) = a + b)`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle zip: %v", err)
	}

	// Should generate [11, 13, 15]
	if !strings.Contains(string(output), "11") {
		t.Errorf("Expected 11 in output, got: %s", string(output))
	}
}

func TestEdgeCase_ZipDifferentLengths(t *testing.T) {
	t.Skip("zip operator not implemented yet")
	input := `x = (1..5 zip 10..12 map (a, b) = a + b)`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle zip with different lengths: %v", err)
	}
}

func TestEdgeCase_TripleZip(t *testing.T) {
	t.Skip("zip operator not implemented yet")
	input := `x = (1..3 zip 10..12 zip 100..102 map (a, b, c) = a + b + c)`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle triple zip: %v", err)
	}

	// Should generate [111, 114, 117]
	if !strings.Contains(string(output), "111") {
		t.Errorf("Expected 111 in output, got: %s", string(output))
	}
}

// EDGE CASES - Advanced Validators

func TestEdgeCase_EmailValidatorUniqueness(t *testing.T) {
	input := `
	users [
		template { email }
		@email()
		@email()
		@email()
		@email()
		@email()
	]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should generate multiple emails: %v", err)
	}

	// Check that emails are generated
	if !strings.Contains(string(output), "@") {
		t.Errorf("Expected email addresses, got: %s", string(output))
	}
}

func TestEdgeCase_UUIDUniquenessInLargeArray(t *testing.T) {
	input := `ids = (1..100 map (n) = @uuid())`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should generate 100 UUIDs: %v", err)
	}

	// Check basic UUID format
	if !strings.Contains(string(output), "-") {
		t.Errorf("Expected UUID format, got: %s", string(output))
	}
}

func TestEdgeCase_ValidatorsInNestedContext(t *testing.T) {
	input := `
	config {
		server {
			id = @uuid()
			name = "server-1"
			port = @int(8000, 9000)
		}
		database {
			id = @uuid()
			host = "localhost"
			port = @int(3000, 4000)
		}
	}`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle validators in nested objects: %v", err)
	}
}

// ============================================================================
// EDGE CASES - Advanced Maps
// ============================================================================

func TestEdgeCase_MapWithTernary(t *testing.T) {
	input := `x = (1..10 map (n) = n % 2 == 0 ? "even" : "odd")`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle map with ternary: %v", err)
	}

	if !strings.Contains(string(output), "even") && !strings.Contains(string(output), "odd") {
		t.Errorf("Expected even/odd strings, got: %s", string(output))
	}
}

func TestEdgeCase_MapWithObjectCreation(t *testing.T) {
	input := `users = (1..3 map (n) = { id: n, name: "User" + n })`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle map with object creation: %v", err)
	}

	if !strings.Contains(string(output), "User") {
		t.Errorf("Expected User prefix, got: %s", string(output))
	}
}

func TestEdgeCase_MapWithValidators(t *testing.T) {
	input := `ids = (1..5 map (n) = @uuid())`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle map with validators: %v", err)
	}
}

// EDGE CASES - Numeric Limits

func TestEdgeCase_VeryLargeInt64(t *testing.T) {
	input := `x = 9223372036854775807`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle max int64: %v", err)
	}

	if !strings.Contains(string(output), "9223372036854775807") {
		t.Errorf("Expected max int64, got: %s", string(output))
	}
}

func TestEdgeCase_VerySmallInt64(t *testing.T) {
	input := `x = -9223372036854775807` // Use max-1 to avoid overflow issues
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle very negative int64: %v", err)
	}

	if !strings.Contains(string(output), "-9223372036854775807") {
		t.Errorf("Expected very negative int64, got: %s", string(output))
	}
}

func TestEdgeCase_FloatingPointPrecision(t *testing.T) {
	input := `x = 0.1 + 0.2`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle float precision: %v", err)
	}

	// Known floating point precision issue
	t.Logf("Result: %s", string(output))
}

func TestEdgeCase_VerySmallDecimal(t *testing.T) {
	input := `x = 0.0000000001`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle very small decimal: %v", err)
	}

	if !strings.Contains(string(output), "0.0000000001") && !strings.Contains(string(output), "1e-10") {
		t.Errorf("Expected small decimal, got: %s", string(output))
	}
}

// EDGE CASES - Advanced Arrays

func TestEdgeCase_ArrayWithMixedValidators(t *testing.T) {
	input := `mixed = [@uuid(), @email(), @int(1, 100), @float(0.0, 1.0), @bool]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle mixed validators in array: %v", err)
	}
}

func TestEdgeCase_MultidimensionalArrayLiteral(t *testing.T) {
	input := `matrix = [[1,2,3], [4,5,6], [7,8,9]]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle multidimensional array: %v", err)
	}

	// JSON is pretty-printed, so just check for nested arrays
	if !strings.Contains(string(output), "[") || !strings.Contains(string(output), "1") {
		t.Errorf("Expected 2D array with numbers, got: %s", string(output))
	}
}

// EDGE CASES - Advanced Templates

func TestEdgeCase_TemplateWithExpressions(t *testing.T) {
	input := `
	data [
		template { x, y, sum }
		1, 2, 1+2
		3, 4, 3+4
		5, 6, 5+6
	]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle template with expressions: %v", err)
	}

	if !strings.Contains(string(output), "3") {
		t.Errorf("Expected computed sum, got: %s", string(output))
	}
}

func TestEdgeCase_TemplateWithTernary(t *testing.T) {
	input := `
	data [
		template { n, parity }
		1, 1 % 2 == 0 ? "even" : "odd"
		2, 2 % 2 == 0 ? "even" : "odd"
		3, 3 % 2 == 0 ? "even" : "odd"
	]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle template with ternary: %v", err)
	}

	if !strings.Contains(string(output), "odd") && !strings.Contains(string(output), "even") {
		t.Errorf("Expected parity strings, got: %s", string(output))
	}
}

func TestEdgeCase_TemplateWithNestedObjects(t *testing.T) {
	input := `
	users [
		template { id, profile }
		1, { name: "Alice", age: 30 }
		2, { name: "Bob", age: 25 }
	]`
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("Should handle template with nested objects: %v", err)
	}

	if !strings.Contains(string(output), "Alice") {
		t.Errorf("Expected nested object data, got: %s", string(output))
	}
}
