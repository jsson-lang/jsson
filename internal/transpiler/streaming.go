package transpiler

import (
	"encoding/json"
	"fmt"
	"io"
	"jsson/internal/ast"
)

// SetStreamingMode configures streaming behavior
func (t *Transpiler) SetStreamingMode(enabled bool, threshold int64) {
	t.streamingEnabled = enabled
	if threshold > 0 {
		t.streamThreshold = threshold
	}
}

// estimateRangeSize calculates the size of a range expression
func (t *Transpiler) estimateRangeSize(e *ast.RangeExpression) int64 {
	// Try to evaluate start and end as constants
	startV, err := t.evalExpression(e.Start, nil)
	if err != nil {
		return 0
	}
	endV, err := t.evalExpression(e.End, nil)
	if err != nil {
		return 0
	}

	startInt, ok1 := startV.(int64)
	endInt, ok2 := endV.(int64)
	if !ok1 || !ok2 {
		return 0
	}

	step := int64(1)
	if e.Step != nil {
		stepV, err := t.evalExpression(e.Step, nil)
		if err == nil {
			if st, ok := stepV.(int64); ok {
				step = st
			}
		}
	} else {
		if startInt > endInt {
			step = -1
		}
	}

	if step == 0 {
		return 0
	}

	var size int64
	if step > 0 {
		size = (endInt-startInt)/step + 1
	} else {
		size = (startInt-endInt)/(-step) + 1
	}

	if size < 0 {
		return 0
	}
	return size
}

// shouldUseStreaming determines if streaming should be used for an expression
func (t *Transpiler) shouldUseStreaming(expr ast.Expression) bool {
	if !t.streamingEnabled {
		return false
	}

	switch e := expr.(type) {
	case *ast.RangeExpression:
		size := t.estimateRangeSize(e)
		return size > t.streamThreshold
	case *ast.MapExpression:
		// Check if the map body is another map (nested maps)
		if _, isMap := e.Body.(*ast.MapExpression); isMap {
			return true
		}
		// Check if the source is a large range
		return t.shouldUseStreaming(e.Left)
	case *ast.ArrayTemplate:
		// Check if any rows contain large ranges
		for _, row := range e.Rows {
			for _, expr := range row {
				if t.shouldUseStreaming(expr) {
					return true
				}
			}
		}
	}
	return false
}

// StreamWriter interface for different output formats
// Allows streaming large datasets without loading everything into memory
type StreamWriter interface {
	WriteArrayStart() error
	WriteArrayItem(item interface{}) error
	WriteArrayEnd() error
	WriteObjectStart() error
	WriteObjectKey(key string) error
	WriteObjectValue(value interface{}) error
	WriteObjectEnd() error
	Flush() error
}

// JSONStreamWriter implements streaming for JSON format
type JSONStreamWriter struct {
	writer      io.Writer
	encoder     *json.Encoder
	arrayDepth  int
	objectDepth int
	itemCount   []int // Stack to track items per level
	needsComma  bool
}

// NewJSONStreamWriter creates a new JSON streaming writer
func NewJSONStreamWriter(w io.Writer) *JSONStreamWriter {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return &JSONStreamWriter{
		writer:     w,
		encoder:    encoder,
		itemCount:  make([]int, 0),
		needsComma: false,
	}
}

func (w *JSONStreamWriter) WriteArrayStart() error {
	if w.needsComma {
		if _, err := w.writer.Write([]byte(",\n")); err != nil {
			return err
		}
	}
	if _, err := w.writer.Write([]byte("[")); err != nil {
		return err
	}
	w.arrayDepth++
	w.itemCount = append(w.itemCount, 0)
	w.needsComma = false
	return nil
}

func (w *JSONStreamWriter) WriteArrayItem(item interface{}) error {
	depth := len(w.itemCount) - 1
	if depth < 0 {
		return fmt.Errorf("WriteArrayItem called outside of array context")
	}

	if w.itemCount[depth] > 0 {
		if _, err := w.writer.Write([]byte(",\n")); err != nil {
			return err
		}
	} else {
		if _, err := w.writer.Write([]byte("\n")); err != nil {
			return err
		}
	}

	// Write indentation
	indent := make([]byte, (w.arrayDepth+w.objectDepth)*2)
	for i := range indent {
		indent[i] = ' '
	}
	if _, err := w.writer.Write(indent); err != nil {
		return err
	}

	// Encode the item
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	w.itemCount[depth]++
	return nil
}

func (w *JSONStreamWriter) WriteArrayEnd() error {
	if w.arrayDepth == 0 {
		return fmt.Errorf("WriteArrayEnd called without matching WriteArrayStart")
	}

	if len(w.itemCount) > 0 {
		w.itemCount = w.itemCount[:len(w.itemCount)-1]
	}

	if _, err := w.writer.Write([]byte("\n")); err != nil {
		return err
	}

	// Write indentation for closing bracket
	indent := make([]byte, (w.arrayDepth+w.objectDepth-1)*2)
	for i := range indent {
		indent[i] = ' '
	}
	if _, err := w.writer.Write(indent); err != nil {
		return err
	}

	if _, err := w.writer.Write([]byte("]")); err != nil {
		return err
	}

	w.arrayDepth--
	w.needsComma = true
	return nil
}

func (w *JSONStreamWriter) WriteObjectStart() error {
	if w.needsComma {
		if _, err := w.writer.Write([]byte(",\n")); err != nil {
			return err
		}
	}
	if _, err := w.writer.Write([]byte("{")); err != nil {
		return err
	}
	w.objectDepth++
	w.itemCount = append(w.itemCount, 0)
	w.needsComma = false
	return nil
}

func (w *JSONStreamWriter) WriteObjectKey(key string) error {
	depth := len(w.itemCount) - 1
	if depth < 0 {
		return fmt.Errorf("WriteObjectKey called outside of object context")
	}

	if w.itemCount[depth] > 0 {
		if _, err := w.writer.Write([]byte(",\n")); err != nil {
			return err
		}
	} else {
		if _, err := w.writer.Write([]byte("\n")); err != nil {
			return err
		}
	}

	// Write indentation
	indent := make([]byte, (w.arrayDepth+w.objectDepth)*2)
	for i := range indent {
		indent[i] = ' '
	}
	if _, err := w.writer.Write(indent); err != nil {
		return err
	}

	// Write key
	keyJSON, err := json.Marshal(key)
	if err != nil {
		return err
	}
	if _, err := w.writer.Write(keyJSON); err != nil {
		return err
	}
	if _, err := w.writer.Write([]byte(": ")); err != nil {
		return err
	}

	w.needsComma = false
	return nil
}

func (w *JSONStreamWriter) WriteObjectValue(value interface{}) error {
	depth := len(w.itemCount) - 1
	if depth < 0 {
		return fmt.Errorf("WriteObjectValue called outside of object context")
	}

	// Encode the value
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	w.itemCount[depth]++
	w.needsComma = true
	return nil
}

func (w *JSONStreamWriter) WriteObjectEnd() error {
	if w.objectDepth == 0 {
		return fmt.Errorf("WriteObjectEnd called without matching WriteObjectStart")
	}

	if len(w.itemCount) > 0 {
		w.itemCount = w.itemCount[:len(w.itemCount)-1]
	}

	if _, err := w.writer.Write([]byte("\n")); err != nil {
		return err
	}

	// Write indentation for closing brace
	indent := make([]byte, (w.arrayDepth+w.objectDepth-1)*2)
	for i := range indent {
		indent[i] = ' '
	}
	if _, err := w.writer.Write(indent); err != nil {
		return err
	}

	if _, err := w.writer.Write([]byte("}")); err != nil {
		return err
	}

	w.objectDepth--
	w.needsComma = true
	return nil
}

func (w *JSONStreamWriter) Flush() error {
	if flusher, ok := w.writer.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}

// RangeIterator allows iterating over ranges without materializing the entire slice
type RangeIterator struct {
	current int64
	end     int64
	step    int64
	done    bool
}

// NewRangeIterator creates a new range iterator
func NewRangeIterator(start, end, step int64) *RangeIterator {
	if step == 0 {
		step = 1
		if start > end {
			step = -1
		}
	}
	return &RangeIterator{
		current: start,
		end:     end,
		step:    step,
		done:    false,
	}
}

// Next returns the next value in the range and whether there are more values
func (ri *RangeIterator) Next() (int64, bool) {
	if ri.done {
		return 0, false
	}

	if ri.step > 0 && ri.current > ri.end {
		ri.done = true
		return 0, false
	}
	if ri.step < 0 && ri.current < ri.end {
		ri.done = true
		return 0, false
	}

	val := ri.current
	ri.current += ri.step
	return val, true
}

// Size returns the total number of elements in the range
func (ri *RangeIterator) Size() int64 {
	if ri.step > 0 {
		return (ri.end-ri.current)/ri.step + 1
	}
	return (ri.current-ri.end)/(-ri.step) + 1
}
