package parser

import (
	"jsson/internal/ast"
	"jsson/internal/lexer"
	"testing"
)

// Test boolean literals extras (yes/no/on/off)
func TestBooleanLiteralsExtras(t *testing.T) {
	tests := []struct {
		input    string
		key      string
		expected bool
	}{
		{"test { enabled = yes }", "enabled", true},
		{"test { disabled = no }", "disabled", false},
		{"test { active = on }", "active", true},
		{"test { inactive = off }", "inactive", false},
		{"test { traditional = true }", "traditional", true},
		{"test { old = false }", "old", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Fatalf("parser errors for '%s': %v", tt.input, p.Errors())
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement for '%s'. got=%d",
				tt.input, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
		if !ok {
			t.Fatalf("stmt not *ast.AssignmentStatement for '%s'. got=%T",
				tt.input, program.Statements[0])
		}

		if stmt.Name.Value != "test" {
			t.Fatalf("statement name wrong for '%s'. expected=test got=%s",
				tt.input, stmt.Name.Value)
		}

		objLit, ok := stmt.Value.(*ast.ObjectLiteral)
		if !ok {
			t.Fatalf("value not *ast.ObjectLiteral for '%s'. got=%T",
				tt.input, stmt.Value)
		}

		value, exists := objLit.Properties[tt.key]
		if !exists {
			t.Fatalf("key '%s' not found in object for '%s'", tt.key, tt.input)
		}

		boolLit, ok := value.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("value not *ast.BooleanLiteral for '%s'. got=%T",
				tt.input, value)
		}

		if boolLit.Value != tt.expected {
			t.Errorf("boolean value wrong for '%s'. expected=%v got=%v",
				tt.input, tt.expected, boolLit.Value)
		}
	}
}

// Test validator expressions
func TestValidatorExpressions(t *testing.T) {
	tests := []struct {
		input        string
		key          string
		expectedType string
	}{
		{"test { id = @uuid }", "id", "uuid"},
		{"test { email = @email }", "email", "email"},
		{"test { website = @url }", "website", "url"},
		{"test { ip = @ipv4 }", "ip", "ipv4"},
		{"test { ipv6 = @ipv6 }", "ipv6", "ipv6"},
		{"test { path = @filepath }", "path", "filepath"},
		{"test { created = @date }", "created", "date"},
		{"test { timestamp = @datetime }", "timestamp", "datetime"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Fatalf("parser errors for '%s': %v", tt.input, p.Errors())
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement for '%s'. got=%d",
				tt.input, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
		if !ok {
			t.Fatalf("stmt not *ast.AssignmentStatement for '%s'. got=%T",
				tt.input, program.Statements[0])
		}

		objLit, ok := stmt.Value.(*ast.ObjectLiteral)
		if !ok {
			t.Fatalf("value not *ast.ObjectLiteral for '%s'. got=%T",
				tt.input, stmt.Value)
		}

		value, exists := objLit.Properties[tt.key]
		if !exists {
			t.Fatalf("key '%s' not found in object for '%s'", tt.key, tt.input)
		}

		validator, ok := value.(*ast.ValidatorExpression)
		if !ok {
			t.Fatalf("value not *ast.ValidatorExpression for '%s'. got=%T",
				tt.input, value)
		}

		if validator.Type != tt.expectedType {
			t.Errorf("validator type wrong for '%s'. expected=%s got=%s",
				tt.input, tt.expectedType, validator.Type)
		}
	}
}

// Test keywords as property names
func TestKeywordsAsPropertyNames(t *testing.T) {
	input := `
test {
  uuid = "value1"
  email = "value2"
  ipv4 = "value3"
  date = "value4"
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

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	obj, ok := stmt.Value.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("value not *ast.ObjectLiteral. got=%T", stmt.Value)
	}

	expectedKeys := []string{"uuid", "email", "ipv4", "date"}
	if len(obj.Keys) != len(expectedKeys) {
		t.Fatalf("object keys count wrong. expected=%d got=%d",
			len(expectedKeys), len(obj.Keys))
	}

	for i, key := range expectedKeys {
		if obj.Keys[i] != key {
			t.Errorf("key %d wrong. expected=%s got=%s", i, key, obj.Keys[i])
		}

		if _, exists := obj.Properties[key]; !exists {
			t.Errorf("property '%s' not found in object", key)
		}
	}
}

// Test regex validator with pattern
func TestRegexValidatorWithPattern(t *testing.T) {
	input := `cpf = @regex("^\\d{3}\\.\\d{3}\\.\\d{3}-\\d{2}$")`

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

	validator, ok := stmt.Value.(*ast.ValidatorExpression)
	if !ok {
		t.Fatalf("value not *ast.ValidatorExpression. got=%T", stmt.Value)
	}

	if validator.Type != "regex" {
		t.Errorf("validator type wrong. expected=regex got=%s", validator.Type)
	}

	expectedPattern := `^\d{3}\.\d{3}\.\d{3}-\d{2}$`
	if validator.Pattern != expectedPattern {
		t.Errorf("validator pattern wrong. expected=%s got=%s",
			expectedPattern, validator.Pattern)
	}
}

// Test mixed boolean literals and validators in object
func TestMixedBooleanAndValidators(t *testing.T) {
	input := `
config {
  enabled = yes
  debug = on
  id = @uuid
  email = @email
  active = true
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

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	obj, ok := stmt.Value.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("value not *ast.ObjectLiteral. got=%T", stmt.Value)
	}

	// Check enabled = yes
	enabledVal, ok := obj.Properties["enabled"].(*ast.BooleanLiteral)
	if !ok || !enabledVal.Value {
		t.Errorf("enabled property should be BooleanLiteral with value true")
	}

	// Check debug = on
	debugVal, ok := obj.Properties["debug"].(*ast.BooleanLiteral)
	if !ok || !debugVal.Value {
		t.Errorf("debug property should be BooleanLiteral with value true")
	}

	// Check id = @uuid
	idVal, ok := obj.Properties["id"].(*ast.ValidatorExpression)
	if !ok || idVal.Type != "uuid" {
		t.Errorf("id property should be ValidatorExpression with type uuid")
	}

	// Check email = @email
	emailVal, ok := obj.Properties["email"].(*ast.ValidatorExpression)
	if !ok || emailVal.Type != "email" {
		t.Errorf("email property should be ValidatorExpression with type email")
	}
}
