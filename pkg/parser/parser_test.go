package parser

import (
	"jsson/internal/ast"
	"jsson/internal/lexer"
	"jsson/internal/token"
	"testing"
)

func TestParseRangeExpression(t *testing.T) {
	input := "p = 8080..8085"
	// Debug: dump lexer tokens
	ld := lexer.New(input)
	for {
		tok := ld.NextToken()
		t.Logf("tok: %#v", tok)
		if tok.Type == token.EOF {
			break
		}
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	re, ok := stmt.Value.(*ast.RangeExpression)
	if !ok {
		t.Fatalf("value not *ast.RangeExpression. got=%T", stmt.Value)
	}

	start, ok := re.Start.(*ast.IntegerLiteral)
	if !ok || start.Value != 8080 {
		t.Fatalf("range start wrong. expected=8080 got=%v", re.Start)
	}
	end, ok := re.End.(*ast.IntegerLiteral)
	if !ok || end.Value != 8085 {
		t.Fatalf("range end wrong. expected=8085 got=%v", re.End)
	}
	if re.Step != nil {
		t.Fatalf("expected nil step, got=%v", re.Step)
	}
}

func TestParseRangeWithStep(t *testing.T) {
	input := "s = 0..10 step 2"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	re, ok := stmt.Value.(*ast.RangeExpression)
	if !ok {
		t.Fatalf("value not *ast.RangeExpression. got=%T", stmt.Value)
	}

	start, ok := re.Start.(*ast.IntegerLiteral)
	if !ok || start.Value != 0 {
		t.Fatalf("range start wrong. expected=0 got=%v", re.Start)
	}
	end, ok := re.End.(*ast.IntegerLiteral)
	if !ok || end.Value != 10 {
		t.Fatalf("range end wrong. expected=10 got=%v", re.End)
	}
	step, ok := re.Step.(*ast.IntegerLiteral)
	if !ok || step.Value != 2 {
		t.Fatalf("range step wrong. expected=2 got=%v", re.Step)
	}
}

func TestParseIncludeStatement(t *testing.T) {
	input := "include \"db.jsson\""
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.IncludeStatement)
	if !ok {
		t.Fatalf("stmt not *ast.IncludeStatement. got=%T", program.Statements[0])
	}
}
