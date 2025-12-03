package transpiler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jsson/internal/lexer"
	"jsson/internal/parser"
)

func TestTranspileRangeAndFlatten(t *testing.T) {
	input := "ports = [ 8080..8085 ]\nsteps = [ 0..10 step 2 ]"
	l := lexer.New(input)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	// baseDir can be current dir for this test
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

	ports, ok := root["ports"].([]interface{})
	if !ok || len(ports) != 6 {
		t.Fatalf("ports not expanded correctly, got=%v", root["ports"])
	}

	steps, ok := root["steps"].([]interface{})
	if !ok || len(steps) != 6 {
		t.Fatalf("steps not expanded correctly, got=%v", root["steps"])
	}
}

func TestIncludeCyclicDetection(t *testing.T) {
	dir, err := os.MkdirTemp("", "jsson_test_cycle")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	aPath := filepath.Join(dir, "a.jsson")
	bPath := filepath.Join(dir, "b.jsson")

	// a includes b, b includes a
	if err := os.WriteFile(aPath, []byte("include \"b.jsson\"\n"), 0644); err != nil {
		t.Fatalf("could not write a.jsson: %v", err)
	}
	if err := os.WriteFile(bPath, []byte("include \"a.jsson\"\n"), 0644); err != nil {
		t.Fatalf("could not write b.jsson: %v", err)
	}

	mainInput := "include \"a.jsson\""
	l := lexer.New(mainInput)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(prog, dir, "keep", "")
	_, err = tr.Transpile()
	if err == nil {
		t.Fatalf("expected cyclic include error, got nil")
	}
	if !strings.Contains(err.Error(), "cyclic include detected") {
		t.Fatalf("unexpected error message, got: %v", err)
	}
}

func TestIncludeResolvesRelative(t *testing.T) {
	// create temp dir with included file
	dir, err := os.MkdirTemp("", "jsson_test")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	incPath := filepath.Join(dir, "included.jsson")
	if err := os.WriteFile(incPath, []byte("db = { host = \"localhost\" }\n"), 0644); err != nil {
		t.Fatalf("could not write include file: %v", err)
	}

	mainInput := "include \"included.jsson\""
	l := lexer.New(mainInput)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(prog, dir, "keep", "")
	out, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpile error: %v", err)
	}

	var root map[string]interface{}
	if err := json.Unmarshal(out, &root); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}

	if _, ok := root["db"]; !ok {
		t.Fatalf("included data not found in transpiled output: %v", root)
	}
}

func TestIncludeDoesNotOverwriteExisting(t *testing.T) {
	// included defines a=1, main defines a=2 before include; include should not overwrite
	dir, err := os.MkdirTemp("", "jsson_test")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	incPath := filepath.Join(dir, "incl.jsson")
	if err := os.WriteFile(incPath, []byte("a = 1\n"), 0644); err != nil {
		t.Fatalf("could not write include file: %v", err)
	}

	mainInput := "a = 2\ninclude \"incl.jsson\""
	l := lexer.New(mainInput)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(prog, dir, "keep", "")
	out, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpile error: %v", err)
	}

	var root map[string]interface{}
	if err := json.Unmarshal(out, &root); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}

	if val, ok := root["a"]; !ok {
		t.Fatalf("key 'a' missing: %v", root)
	} else {
		// JSON numbers are float64 when unmarshaled
		if num, ok := val.(float64); !ok || num != 2 {
			t.Fatalf("expected a=2, got %v", val)
		}
	}
}

func TestIncludeMergeOverwriteMode(t *testing.T) {
	dir, err := os.MkdirTemp("", "jsson_test_merge_overwrite")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	incPath := filepath.Join(dir, "incl.jsson")
	if err := os.WriteFile(incPath, []byte("a = 1\n"), 0644); err != nil {
		t.Fatalf("could not write include file: %v", err)
	}

	mainInput := "a = 2\ninclude \"incl.jsson\""
	l := lexer.New(mainInput)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(prog, dir, "overwrite", "")
	out, err := tr.Transpile()
	if err != nil {
		t.Fatalf("transpile error: %v", err)
	}

	var root map[string]interface{}
	if err := json.Unmarshal(out, &root); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}

	if val, ok := root["a"]; !ok {
		t.Fatalf("key 'a' missing: %v", root)
	} else {
		if num, ok := val.(float64); !ok || num != 1 {
			t.Fatalf("expected a=1 (overwritten), got %v", val)
		}
	}
}

func TestIncludeMergeErrorMode(t *testing.T) {
	dir, err := os.MkdirTemp("", "jsson_test_merge_error")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	incPath := filepath.Join(dir, "incl.jsson")
	if err := os.WriteFile(incPath, []byte("a = 1\n"), 0644); err != nil {
		t.Fatalf("could not write include file: %v", err)
	}

	mainInput := "a = 2\ninclude \"incl.jsson\""
	l := lexer.New(mainInput)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	tr := New(prog, dir, "error", "")
	_, err = tr.Transpile()
	if err == nil {
		t.Fatalf("expected merge conflict error, got nil")
	}
	if !strings.Contains(err.Error(), "include merge conflict") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
