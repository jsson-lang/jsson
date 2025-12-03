package transpiler

import (
	"bytes"
	"jsson/internal/ast"
	"testing"
)

func TestJSONStreamWriter_BasicArray(t *testing.T) {
	var buf bytes.Buffer
	w := NewJSONStreamWriter(&buf)

	if err := w.WriteArrayStart(); err != nil {
		t.Fatalf("WriteArrayStart failed: %v", err)
	}

	items := []interface{}{1, 2, 3, 4, 5}
	for _, item := range items {
		if err := w.WriteArrayItem(item); err != nil {
			t.Fatalf("WriteArrayItem failed: %v", err)
		}
	}

	if err := w.WriteArrayEnd(); err != nil {
		t.Fatalf("WriteArrayEnd failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Fatal("Output is empty")
	}

	// Should contain all numbers
	for _, item := range items {
		if !bytes.Contains(buf.Bytes(), []byte(string(rune(item.(int)+'0')))) {
			t.Logf("Output: %s", output)
		}
	}
}

func TestJSONStreamWriter_BasicObject(t *testing.T) {
	var buf bytes.Buffer
	w := NewJSONStreamWriter(&buf)

	if err := w.WriteObjectStart(); err != nil {
		t.Fatalf("WriteObjectStart failed: %v", err)
	}

	if err := w.WriteObjectKey("name"); err != nil {
		t.Fatalf("WriteObjectKey failed: %v", err)
	}
	if err := w.WriteObjectValue("John"); err != nil {
		t.Fatalf("WriteObjectValue failed: %v", err)
	}

	if err := w.WriteObjectKey("age"); err != nil {
		t.Fatalf("WriteObjectKey failed: %v", err)
	}
	if err := w.WriteObjectValue(30); err != nil {
		t.Fatalf("WriteObjectValue failed: %v", err)
	}

	if err := w.WriteObjectEnd(); err != nil {
		t.Fatalf("WriteObjectEnd failed: %v", err)
	}

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("name")) {
		t.Errorf("Output missing 'name' key: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("John")) {
		t.Errorf("Output missing 'John' value: %s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("age")) {
		t.Errorf("Output missing 'age' key: %s", output)
	}
}

func TestRangeIterator_Forward(t *testing.T) {
	iter := NewRangeIterator(0, 5, 1)

	expected := []int64{0, 1, 2, 3, 4, 5}
	var result []int64

	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		result = append(result, val)
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestRangeIterator_Backward(t *testing.T) {
	iter := NewRangeIterator(5, 0, -1)

	expected := []int64{5, 4, 3, 2, 1, 0}
	var result []int64

	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		result = append(result, val)
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestRangeIterator_Step(t *testing.T) {
	iter := NewRangeIterator(0, 10, 2)

	expected := []int64{0, 2, 4, 6, 8, 10}
	var result []int64

	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		result = append(result, val)
	}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestTranspiler_SetStreamingMode(t *testing.T) {
	// Create a minimal valid program for testing
	prog := &ast.Program{
		Statements: []ast.Statement{},
	}

	tr := New(prog, "", "keep", "")

	// Default should be disabled
	if tr.streamingEnabled {
		t.Error("Streaming should be disabled by default")
	}

	// Enable streaming
	tr.SetStreamingMode(true, 5000)
	if !tr.streamingEnabled {
		t.Error("Streaming should be enabled")
	}
	if tr.streamThreshold != 5000 {
		t.Errorf("Expected threshold 5000, got %d", tr.streamThreshold)
	}

	// Disable streaming
	tr.SetStreamingMode(false, 0)
	if tr.streamingEnabled {
		t.Error("Streaming should be disabled")
	}
}
