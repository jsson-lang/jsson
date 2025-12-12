package parser

import (
	"jsson/internal/ast"
	"jsson/internal/lexer"
	"testing"
)

// Test @preset statement definition
func TestPresetStatement(t *testing.T) {
	input := `
@preset "api-defaults" {
    timeout = 30
    retries = 3
    cache = true
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.PresetStatement)
	if !ok {
		t.Fatalf("stmt not *ast.PresetStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "api-defaults" {
		t.Errorf("preset name wrong. expected=api-defaults got=%s", stmt.Name.Value)
	}

	if stmt.Body == nil {
		t.Fatal("preset body is nil")
	}

	expectedProps := map[string]interface{}{
		"timeout": int64(30),
		"retries": int64(3),
		"cache":   true,
	}

	for key, expectedVal := range expectedProps {
		prop, exists := stmt.Body.Properties[key]
		if !exists {
			t.Errorf("property '%s' not found in preset body", key)
			continue
		}

		switch expected := expectedVal.(type) {
		case int64:
			intLit, ok := prop.(*ast.IntegerLiteral)
			if !ok {
				t.Errorf("property '%s' not IntegerLiteral. got=%T", key, prop)
				continue
			}
			if intLit.Value != expected {
				t.Errorf("property '%s' value wrong. expected=%d got=%d", key, expected, intLit.Value)
			}
		case bool:
			boolLit, ok := prop.(*ast.BooleanLiteral)
			if !ok {
				t.Errorf("property '%s' not BooleanLiteral. got=%T", key, prop)
				continue
			}
			if boolLit.Value != expected {
				t.Errorf("property '%s' value wrong. expected=%v got=%v", key, expected, boolLit.Value)
			}
		}
	}
}

// Test preset reference with @use syntax
func TestPresetReferenceUse(t *testing.T) {
	input := `dev_api = @use "api-defaults"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "dev_api" {
		t.Errorf("assignment name wrong. expected=dev_api got=%s", stmt.Name.Value)
	}

	ref, ok := stmt.Value.(*ast.PresetReference)
	if !ok {
		t.Fatalf("value not *ast.PresetReference. got=%T", stmt.Value)
	}

	if ref.Name.Value != "api-defaults" {
		t.Errorf("preset reference name wrong. expected=api-defaults got=%s", ref.Name.Value)
	}

	if ref.Overrides != nil {
		t.Errorf("expected no overrides, got=%v", ref.Overrides)
	}
}

// Test preset reference with legacy @ syntax
func TestPresetReferenceLegacy(t *testing.T) {
	input := `dev_api = @"api-defaults"`

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

	ref, ok := stmt.Value.(*ast.PresetReference)
	if !ok {
		t.Fatalf("value not *ast.PresetReference. got=%T", stmt.Value)
	}

	if ref.Name.Value != "api-defaults" {
		t.Errorf("preset reference name wrong. expected=api-defaults got=%s", ref.Name.Value)
	}
}

// Test preset reference with overrides
func TestPresetReferenceWithOverrides(t *testing.T) {
	input := `
prod_api = @use "api-defaults" {
    timeout = 60
    retries = 5
}
`
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

	ref, ok := stmt.Value.(*ast.PresetReference)
	if !ok {
		t.Fatalf("value not *ast.PresetReference. got=%T", stmt.Value)
	}

	if ref.Name.Value != "api-defaults" {
		t.Errorf("preset reference name wrong. expected=api-defaults got=%s", ref.Name.Value)
	}

	if ref.Overrides == nil {
		t.Fatal("expected overrides, got nil")
	}

	// Check timeout override
	timeout, exists := ref.Overrides.Properties["timeout"]
	if !exists {
		t.Error("timeout override not found")
	} else {
		intLit, ok := timeout.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("timeout not IntegerLiteral. got=%T", timeout)
		} else if intLit.Value != 60 {
			t.Errorf("timeout value wrong. expected=60 got=%d", intLit.Value)
		}
	}

	// Check retries override
	retries, exists := ref.Overrides.Properties["retries"]
	if !exists {
		t.Error("retries override not found")
	} else {
		intLit, ok := retries.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("retries not IntegerLiteral. got=%T", retries)
		} else if intLit.Value != 5 {
			t.Errorf("retries value wrong. expected=5 got=%d", intLit.Value)
		}
	}
}

// Test preset in nested object
func TestPresetInNestedObject(t *testing.T) {
	input := `
service {
    name = "my-service"
    api = @use "api-defaults"
}
`
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

	objLit, ok := stmt.Value.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("value not *ast.ObjectLiteral. got=%T", stmt.Value)
	}

	// Check name property
	name, exists := objLit.Properties["name"]
	if !exists {
		t.Fatal("name property not found")
	}
	strLit, ok := name.(*ast.StringLiteral)
	if !ok || strLit.Value != "my-service" {
		t.Errorf("name property wrong. expected=my-service got=%v", name)
	}

	// Check api property is preset reference
	api, exists := objLit.Properties["api"]
	if !exists {
		t.Fatal("api property not found")
	}
	ref, ok := api.(*ast.PresetReference)
	if !ok {
		t.Fatalf("api property not *ast.PresetReference. got=%T", api)
	}
	if ref.Name.Value != "api-defaults" {
		t.Errorf("preset name wrong. expected=api-defaults got=%s", ref.Name.Value)
	}
}

// Test multiple presets in one file
func TestMultiplePresets(t *testing.T) {
	input := `
@preset "api-defaults" {
    timeout = 30
}

@preset "logging" {
    level = "info"
    format = "json"
}

dev_api = @use "api-defaults"
logging = @use "logging"
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements. got=%d", len(program.Statements))
	}

	// Check first preset
	preset1, ok := program.Statements[0].(*ast.PresetStatement)
	if !ok {
		t.Fatalf("statement 0 not *ast.PresetStatement. got=%T", program.Statements[0])
	}
	if preset1.Name.Value != "api-defaults" {
		t.Errorf("preset 1 name wrong. expected=api-defaults got=%s", preset1.Name.Value)
	}

	// Check second preset
	preset2, ok := program.Statements[1].(*ast.PresetStatement)
	if !ok {
		t.Fatalf("statement 1 not *ast.PresetStatement. got=%T", program.Statements[1])
	}
	if preset2.Name.Value != "logging" {
		t.Errorf("preset 2 name wrong. expected=logging got=%s", preset2.Name.Value)
	}

	// Check first usage
	usage1, ok := program.Statements[2].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("statement 2 not *ast.AssignmentStatement. got=%T", program.Statements[2])
	}
	ref1, ok := usage1.Value.(*ast.PresetReference)
	if !ok || ref1.Name.Value != "api-defaults" {
		t.Errorf("usage 1 wrong. expected @use api-defaults")
	}

	// Check second usage
	usage2, ok := program.Statements[3].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("statement 3 not *ast.AssignmentStatement. got=%T", program.Statements[3])
	}
	ref2, ok := usage2.Value.(*ast.PresetReference)
	if !ok || ref2.Name.Value != "logging" {
		t.Errorf("usage 2 wrong. expected @use logging")
	}
}
