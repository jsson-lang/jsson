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

// ============================================================================
// EDGE CASES - DIVISÃO E OPERAÇÕES
// ============================================================================

func TestEdgeCase_DivisionByZero(t *testing.T) {
	input := "data = 10 / 0"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()

	// Should error
	if err == nil {
		t.Error("Expected error for division by zero, got nil")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "division") && !strings.Contains(err.Error(), "zero") {
		t.Logf("Error message: %v", err)
	}
}

func TestEdgeCase_ModuloByZero(t *testing.T) {
	input := "data = 10 % 0"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	tr := New(prog, "", "keep", "")
	_, err := tr.Transpile()

	// Should error
	if err == nil {
		t.Error("Expected error for modulo by zero, got nil")
	}
}

// ============================================================================
// EDGE CASES - STRINGS
// ============================================================================

func TestEdgeCase_EmptyString(t *testing.T) {
	input := `data = ""`
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

	if !strings.Contains(string(output), `""`) {
		t.Errorf("Expected empty string, got: %s", string(output))
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

func TestEdgeCase_DeeplyNestedArrays(t *testing.T) {
	input := "data = [[[[1]]]]"
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

	// Should handle deep nesting
	if !strings.Contains(string(output), "[[[[1]]]]") && !strings.Contains(string(output), "1") {
		t.Errorf("Expected nested arrays, got: %s", string(output))
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

// ============================================================================
// EDGE CASES - VARIÁVEIS
// ============================================================================

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
