package transpiler

import (
	"jsson/pkg/lexer"
	"jsson/pkg/parser"
	"testing"
)

// Benchmark large range without streaming
func BenchmarkLargeRange_NoStreaming(b *testing.B) {
	input := "data = 0..9999"

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := parser.New(l)
		prog := p.ParseProgram()

		tr := New(prog, "", "keep", "")
		tr.SetStreamingMode(false, 100000) // Disabled

		_, err := tr.Transpile()
		if err != nil {
			b.Fatalf("Transpile error: %v", err)
		}
	}
}

// Benchmark large range with streaming
func BenchmarkLargeRange_WithStreaming(b *testing.B) {
	input := "data = 0..9999"

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := parser.New(l)
		prog := p.ParseProgram()

		tr := New(prog, "", "keep", "")
		tr.SetStreamingMode(true, 1000) // Enabled with low threshold

		_, err := tr.Transpile()
		if err != nil {
			b.Fatalf("Transpile error: %v", err)
		}
	}
}

// Benchmark very large range
func BenchmarkVeryLargeRange_NoStreaming(b *testing.B) {
	input := "data = 0..99999"

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := parser.New(l)
		prog := p.ParseProgram()

		tr := New(prog, "", "keep", "")
		tr.SetStreamingMode(false, 1000000)

		_, err := tr.Transpile()
		if err != nil {
			b.Fatalf("Transpile error: %v", err)
		}
	}
}

// Benchmark map transformation
func BenchmarkMapTransform_NoStreaming(b *testing.B) {
	input := `data = (0..999 map (x) = { id = x, value = x * 2 })`

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := parser.New(l)
		prog := p.ParseProgram()

		tr := New(prog, "", "keep", "")
		tr.SetStreamingMode(false, 100000)

		_, err := tr.Transpile()
		if err != nil {
			b.Fatalf("Transpile error: %v", err)
		}
	}
}

// Benchmark template with range
func BenchmarkTemplate_WithRange(b *testing.B) {
	input := `users [
  template { id, name }
  
  0..999, "user"
]`

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := parser.New(l)
		prog := p.ParseProgram()

		tr := New(prog, "", "keep", "")
		tr.SetStreamingMode(false, 100000)

		_, err := tr.Transpile()
		if err != nil {
			b.Fatalf("Transpile error: %v", err)
		}
	}
}

// Benchmark nested map (matrix generation)
func BenchmarkNestedMap_SmallMatrix(b *testing.B) {
	input := `matrix = (0..9 map (y) = (0..9 map (x) = x * y))`

	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := parser.New(l)
		prog := p.ParseProgram()

		tr := New(prog, "", "keep", "")
		tr.SetStreamingMode(false, 100000)

		_, err := tr.Transpile()
		if err != nil {
			b.Fatalf("Transpile error: %v", err)
		}
	}
}
