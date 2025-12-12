package transpiler

import (
	"encoding/json"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"strings"
	"testing"
)

// Test basic preset definition and usage
func TestPresetBasicUsage(t *testing.T) {
	input := `
@preset "server-defaults" {
    port = 8080
    host = "localhost"
    timeout = 30
}

dev_server = @use "server-defaults"
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	devServer, exists := result["dev_server"].(map[string]interface{})
	if !exists {
		t.Fatal("dev_server not found in output")
	}

	// Check all preset properties were applied
	if devServer["port"] != float64(8080) {
		t.Errorf("port wrong. expected=8080 got=%v", devServer["port"])
	}
	if devServer["host"] != "localhost" {
		t.Errorf("host wrong. expected=localhost got=%v", devServer["host"])
	}
	if devServer["timeout"] != float64(30) {
		t.Errorf("timeout wrong. expected=30 got=%v", devServer["timeout"])
	}
}

// Test preset with overrides
func TestPresetWithOverrides(t *testing.T) {
	input := `
@preset "server-defaults" {
    port = 8080
    host = "localhost"
    timeout = 30
    maxConnections = 100
}

prod_server = @use "server-defaults" {
    port = 443
    host = "0.0.0.0"
    timeout = 60
}
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	prodServer, exists := result["prod_server"].(map[string]interface{})
	if !exists {
		t.Fatal("prod_server not found in output")
	}

	// Check overridden values
	if prodServer["port"] != float64(443) {
		t.Errorf("port not overridden. expected=443 got=%v", prodServer["port"])
	}
	if prodServer["host"] != "0.0.0.0" {
		t.Errorf("host not overridden. expected=0.0.0.0 got=%v", prodServer["host"])
	}
	if prodServer["timeout"] != float64(60) {
		t.Errorf("timeout not overridden. expected=60 got=%v", prodServer["timeout"])
	}

	// Check non-overridden value remains
	if prodServer["maxConnections"] != float64(100) {
		t.Errorf("maxConnections should remain. expected=100 got=%v", prodServer["maxConnections"])
	}
}

// Test legacy @"preset" syntax
func TestPresetLegacySyntax(t *testing.T) {
	input := `
@preset "api-config" {
    timeout = 30
    retries = 3
}

api = @"api-config"
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	api, exists := result["api"].(map[string]interface{})
	if !exists {
		t.Fatal("api not found in output")
	}

	if api["timeout"] != float64(30) {
		t.Errorf("timeout wrong. expected=30 got=%v", api["timeout"])
	}
	if api["retries"] != float64(3) {
		t.Errorf("retries wrong. expected=3 got=%v", api["retries"])
	}
}

// Test preset in nested object
func TestPresetInNestedObject(t *testing.T) {
	input := `
@preset "logging" {
    level = "info"
    format = "json"
}

service {
    name = "my-service"
    logging = @use "logging" {
        level = "debug"
    }
}
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	service, exists := result["service"].(map[string]interface{})
	if !exists {
		t.Fatal("service not found in output")
	}

	if service["name"] != "my-service" {
		t.Errorf("name wrong. expected=my-service got=%v", service["name"])
	}

	logging, exists := service["logging"].(map[string]interface{})
	if !exists {
		t.Fatal("logging not found in service")
	}

	// Check override
	if logging["level"] != "debug" {
		t.Errorf("level not overridden. expected=debug got=%v", logging["level"])
	}
	// Check non-overridden
	if logging["format"] != "json" {
		t.Errorf("format should remain. expected=json got=%v", logging["format"])
	}
}

// Test multiple presets
func TestMultiplePresets(t *testing.T) {
	input := `
@preset "api-defaults" {
    timeout = 30
    retries = 3
}

@preset "logging" {
    level = "info"
    format = "json"
}

service {
    api = @use "api-defaults"
    logging = @use "logging"
}
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	service, exists := result["service"].(map[string]interface{})
	if !exists {
		t.Fatal("service not found in output")
	}

	api, exists := service["api"].(map[string]interface{})
	if !exists {
		t.Fatal("api not found in service")
	}
	if api["timeout"] != float64(30) || api["retries"] != float64(3) {
		t.Errorf("api preset not applied correctly. got=%v", api)
	}

	logging, exists := service["logging"].(map[string]interface{})
	if !exists {
		t.Fatal("logging not found in service")
	}
	if logging["level"] != "info" || logging["format"] != "json" {
		t.Errorf("logging preset not applied correctly. got=%v", logging)
	}
}

// Test preset transpilation to YAML format
func TestPresetWithYAMLOutput(t *testing.T) {
	input := `
@preset "config" {
    enabled = true
    count = 5
}

app = @use "config"
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")

	// First transpile to JSON to verify it works
	jsonOutput, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error (JSON): %v", err)
	}

	// Verify JSON contains the data
	var result map[string]interface{}
	if err := json.Unmarshal(jsonOutput, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	app, exists := result["app"].(map[string]interface{})
	if !exists {
		t.Fatal("app not found in JSON output")
	}

	if app["enabled"] != true || app["count"] != float64(5) {
		t.Errorf("preset values not applied correctly in JSON. got=%v", app)
	}

	// Now test YAML output
	yamlOutput, err := tr.TranspileToYAML()
	if err != nil {
		t.Fatalf("transpiler error (YAML): %v", err)
	}

	yamlStr := string(yamlOutput)

	// Verify YAML contains expected structure (may vary in format)
	if !strings.Contains(yamlStr, "app") {
		t.Errorf("YAML output missing 'app'. Full output:\n%s", yamlStr)
	}
}

// Test preset transpilation to TOML format
func TestPresetWithTOMLOutput(t *testing.T) {
	input := `
@preset "database" {
    host = "localhost"
    port = 5432
}

db = @use "database"
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")

	// First transpile to JSON to verify it works
	jsonOutput, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpiler error (JSON): %v", err)
	}

	// Verify JSON contains the data
	var result map[string]interface{}
	if err := json.Unmarshal(jsonOutput, &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	db, exists := result["db"].(map[string]interface{})
	if !exists {
		t.Fatal("db not found in JSON output")
	}

	if db["host"] != "localhost" || db["port"] != float64(5432) {
		t.Errorf("preset values not applied correctly in JSON. got=%v", db)
	}

	// Now test TOML output
	tomlOutput, err := tr.TranspileToTOML()
	if err != nil {
		t.Fatalf("transpiler error (TOML): %v", err)
	}

	tomlStr := string(tomlOutput)

	// Verify TOML contains expected structure
	if !strings.Contains(tomlStr, "db") {
		t.Errorf("TOML output missing 'db'. Full output:\n%s", tomlStr)
	}
}

// Test undefined preset reference error
func TestUndefinedPresetReference(t *testing.T) {
	input := `api = @use "undefined-preset"`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(program, "", "keep", "")
	_, err := tr.Transpile()

	// Should produce an error for undefined preset
	if err == nil {
		t.Error("expected error for undefined preset, got nil")
	} else if !strings.Contains(err.Error(), "undefined") && !strings.Contains(err.Error(), "not found") {
		t.Errorf("error message should mention undefined/not found. got: %v", err)
	}
}
