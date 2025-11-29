package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputPtr := flag.String("i", "", "Input JSSON file")
	formatPtr := flag.String("f", "json", "Output format: json|yaml|toml")
	mergeMode := flag.String("include-merge", "keep", "Include merge strategy: keep|overwrite|error")
	flag.Parse()

	if *inputPtr == "" {
		fmt.Println("Please provide an input file with -i")
		os.Exit(1)
	}

	// Validate format
	format := strings.ToLower(*formatPtr)
	validFormats := map[string]bool{
		"json": true, "yaml": true, "toml": true,
		"typescript": true, "ts": true,
	}

	if !validFormats[format] {
		fmt.Printf("Invalid format: %s. Must be json, yaml, toml or typescript\n", *formatPtr)
		os.Exit(1)
	}

	// Normalize aliases
	if format == "ts" {
		format = "typescript"
	}

	data, err := ioutil.ReadFile(*inputPtr)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Resolve the absolute path of the input file and pass its directory as baseDir
	absInput, err := filepath.Abs(*inputPtr)
	if err != nil {
		fmt.Printf("Error resolving input path: %v\n", err)
		os.Exit(1)
	}
	baseDir := filepath.Dir(absInput)

	l := lexer.New(string(data))
	l.SetSourceFile(absInput)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Println("\t" + msg)
		}
		os.Exit(1)
	}

	t := transpiler.New(program, baseDir, *mergeMode, absInput)

	var output []byte
	switch format {
	case "json":
		output, err = t.Transpile()
	case "yaml":
		output, err = t.TranspileToYAML()
	case "toml":
		output, err = t.TranspileToTOML()
	case "typescript":
		output, err = t.TranspileToTypeScript()

	}

	if err != nil {
		fmt.Printf("Transpilation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
