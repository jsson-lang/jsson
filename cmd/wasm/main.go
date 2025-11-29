//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"jsson/internal/transpiler"
	"syscall/js"
)

func transpile(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "No input provided",
		}
	}

	input := args[0].String()
	format := "json" // default
	if len(args) >= 2 {
		format = args[1].String()
	}

	// Lexer
	l := lexer.New(input)
	l.SetSourceFile("playground.jsson")

	// Parser
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return map[string]interface{}{
			"error": fmt.Sprintf("Parser errors: %v", p.Errors()),
		}
	}

	// Transpiler
	t := transpiler.New(program, ".", "keep", "playground.jsson")

	var output []byte
	var err error

	switch format {
	case "yaml":
		output, err = t.TranspileToYAML()
	case "toml":
		output, err = t.TranspileToTOML()
	case "typescript", "ts":
		output, err = t.TranspileToTypeScript()
	default:
		output, err = t.Transpile()
	}

	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Transpilation error: %v", err),
		}
	}

	return map[string]interface{}{
		"output": string(output),
	}
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("transpileJSSON", js.FuncOf(transpile))
	<-c
}
