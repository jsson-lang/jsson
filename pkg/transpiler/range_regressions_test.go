package transpiler

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"jsson/pkg/lexer"
	"jsson/pkg/parser"
)

func TestMultipleRangesFlatten(t *testing.T) {
	input := "points = [ 0..2, 10..12 ]"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tDir, _ := os.Getwd()
	tr := New(prog, tDir, "keep", "")
	out, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpile error: %v", err)
	}

	var root map[string]interface{}
	if err := json.Unmarshal(out, &root); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}

	pts, ok := root["points"].([]interface{})
	if !ok || len(pts) != 6 {
		t.Fatalf("points not flattened correctly, got=%v", root["points"])
	}
}

func TestNestedMapOverRange(t *testing.T) {
	// outer range 0..2, inner map over 5..6 producing g+x
	input := "nested = (0..2 map (g) = (5..6 map (x) = g + x))"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tDir, _ := os.Getwd()
	tr := New(prog, tDir, "keep", "")
	out, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpile error: %v", err)
	}

	var root map[string]interface{}
	if err := json.Unmarshal(out, &root); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}

	nested, ok := root["nested"].([]interface{})
	if !ok || len(nested) != 3 {
		t.Fatalf("unexpected nested result, got=%v", root["nested"])
	}

	// each inner element should be an array of length 2
	for i, el := range nested {
		arr, ok := el.([]interface{})
		if !ok || len(arr) != 2 {
			t.Fatalf("nested[%d] unexpected: %v", i, el)
		}
	}
}

func TestBinaryOpOnRangeReturnsError(t *testing.T) {
	input := "bad = (0..2) + 5"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tDir, _ := os.Getwd()
	tr := New(prog, tDir, "keep", "")
	_, err := tr.Transpile()
	if err == nil {
		t.Fatalf("expected error when applying binary op to range, got nil")
	}
	if !strings.Contains(err.Error(), "cannot apply operator") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
