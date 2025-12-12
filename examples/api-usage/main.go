package main

import (
	"fmt"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
)

// Example of using JSSON as a library
func main() {
	// JSSON v0.0.6 source with new features
	source := `
// Boolean literals
settings {
    debug = yes
    production = no
    ssl_enabled = on
}

// Validators - auto-generate data
user {
    id = @uuid
    email = @email
    website = @url
    created_at = @datetime
}

// Presets for reusable configs
@preset "api-defaults" {
    timeout = 30
    retries = 3
    cache = true
}

dev_api = @use "api-defaults"

prod_api = @use "api-defaults" {
    timeout = 60
    retries = 5
}

// Ranges and maps
ports = (8080..8090)
`

	// Parse
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parsing errors
	if len(p.Errors()) > 0 {
		fmt.Println("❌ Parser errors:")
		for _, err := range p.Errors() {
			fmt.Println("  ", err)
		}
		return
	}

	fmt.Println("✅ Parsing successful!")

	// Transpile to JSON
	tr := transpiler.New(program, "", "keep", "")
	output, err := tr.Transpile()
	if err != nil {
		fmt.Printf("❌ Transpile error: %v\n", err)
		return
	}

	fmt.Println("\n📄 JSON Output:")
	fmt.Println(string(output))

	// Transpile to YAML
	yamlOutput, err := tr.TranspileToYAML()
	if err != nil {
		fmt.Printf("❌ YAML transpile error: %v\n", err)
		return
	}

	fmt.Println("\n📄 YAML Output:")
	fmt.Println(string(yamlOutput))

	// Transpile to TOML
	tomlOutput, err := tr.TranspileToTOML()
	if err != nil {
		fmt.Printf("❌ TOML transpile error: %v\n", err)
		return
	}

	fmt.Println("\n📄 TOML Output:")
	fmt.Println(string(tomlOutput))

	fmt.Println("\n✨ All transpilations successful!")
}
