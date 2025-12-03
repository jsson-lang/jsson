package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FormatContext(sourceFile string, line, col int) string {
	if sourceFile == "" {
		return fmt.Sprintf("%d:%d", line, col)
	}

	data, err := os.ReadFile(sourceFile)
	if err != nil {
		// fallback to basename
		return fmt.Sprintf("%s:%d:%d", filepath.Base(sourceFile), line, col)
	}

	lines := strings.Split(string(data), "\n")
	idx := line - 1
	if idx < 0 || idx >= len(lines) {
		return fmt.Sprintf("%s:%d:%d", filepath.Base(sourceFile), line, col)
	}

	srcLine := lines[idx]
	// Build caret line: use rune count to handle multi-byte characters
	caretPos := 0
	for i, _ := range []rune(srcLine) {
		if i >= col-1 {
			break
		}
		caretPos++
	}

	// Use spaces for caret; if col is beyond line length put caret at end}
	if col-1 > len([]rune(srcLine)) {
		caretPos = len([]rune(srcLine))
	}

	caret := strings.Repeat(" ", caretPos) + "^"

	return fmt.Sprintf("%s:%d:%d\n%s\n%s", filepath.Base(sourceFile), line, col, srcLine, caret)
}
