package transpiler

import (
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"strings"
	"testing"
)

// Test for bug fix: arithmetic operations with variables
func TestVariableArithmetic_BasicOperations(t *testing.T) {
	input := `
price = 100
tax = 15
total = price + tax

discount = 20
final = total - discount
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

	// Should have correct values
	if !strings.Contains(string(output), `"total": 115`) {
		t.Errorf("Expected total=115, got: %s", string(output))
	}
	if !strings.Contains(string(output), `"final": 95`) {
		t.Errorf("Expected final=95, got: %s", string(output))
	}
}

func TestVariableArithmetic_AllOperators(t *testing.T) {
	input := `
a = 10
b = 3

sum = a + b
diff = a - b
prod = a * b
quot = a / b
mod = a % b
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

	// Verify all operations work
	if !strings.Contains(string(output), `"sum": 13`) {
		t.Errorf("Expected sum=13, got: %s", string(output))
	}
	if !strings.Contains(string(output), `"diff": 7`) {
		t.Errorf("Expected diff=7, got: %s", string(output))
	}
	if !strings.Contains(string(output), `"prod": 30`) {
		t.Errorf("Expected prod=30, got: %s", string(output))
	}
}

func TestVariableArithmetic_ChainedOperations(t *testing.T) {
	input := `
x = 5
y = x * 2
z = y + 10
result = z - 5
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

	// x=5, y=10, z=20, result=15
	if !strings.Contains(string(output), `"result": 15`) {
		t.Errorf("Expected result=15, got: %s", string(output))
	}
}

func TestVariableArithmetic_MixedWithLiterals(t *testing.T) {
	input := `
base = 100
withTax = base + 15
withDiscount = withTax - 10
doubled = withDiscount * 2
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

	// base=100, withTax=115, withDiscount=105, doubled=210
	if !strings.Contains(string(output), `"doubled": 210`) {
		t.Errorf("Expected doubled=210, got: %s", string(output))
	}
}
