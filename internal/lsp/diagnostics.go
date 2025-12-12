package lsp

import (
	"fmt"
	"strings"
)

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
	Source   string `json:"source"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

const (
	SeverityError   = 1
	SeverityWarning = 2
	SeverityInfo    = 3
	SeverityHint    = 4
)

func (s *Server) publishDiagnostics(uri, content string) error {
	errors := s.parseDocument(content)
	diagnostics := []Diagnostic{}

	for _, err := range errors {
		diag := s.errorToDiagnostic(err, content)
		diagnostics = append(diagnostics, diag)
	}

	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]interface{}{
			"uri":         uri,
			"diagnostics": diagnostics,
		},
	}

	return s.writeMessage(notification)
}

// errorToDiagnostic converts a parser error to an LSP diagnostic
func (s *Server) errorToDiagnostic(errMsg string, content string) Diagnostic {
	line, col := s.extractErrorPosition(errMsg, content)

	// Extract just the error message (remove the "Syntax wizard: file:line:col — " prefix)
	message := errMsg
	if idx := strings.Index(errMsg, " — "); idx != -1 {
		message = errMsg[idx+3:]
	}

	return Diagnostic{
		Range: Range{
			Start: Position{Line: line, Character: col},
			End:   Position{Line: line, Character: col + 5}, // Highlight a few characters
		},
		Severity: SeverityError,
		Message:  message,
		Source:   "jsson",
	}
}

// extractErrorPosition extracts line and column from error message
func (s *Server) extractErrorPosition(errMsg, content string) (int, int) {
	// Format: "Syntax wizard: file:line:col — message" or "Syntax wizard: line:col — message"
	// Example: "Syntax wizard: test.jsson:20:1 — expected '}' - wizard can't find the closing bracket"
	
	parts := strings.Split(errMsg, " ")
	for _, part := range parts {
		// Look for pattern like "20:1" or "file.jsson:20:1"
		if strings.Contains(part, ":") {
			colonParts := strings.Split(part, ":")
			if len(colonParts) >= 2 {
				lineStr := colonParts[len(colonParts)-2]
				colStr := colonParts[len(colonParts)-1]

				var line, col int
				if _, err := fmt.Sscanf(lineStr, "%d", &line); err == nil {
					if _, err := fmt.Sscanf(colStr, "%d", &col); err == nil {
						// LSP uses 0-based indexing
						return line - 1, col - 1
					}
				}
			}
		}
	}

	// Format: "error at line X: message" or similar
	var line int
	if _, err := fmt.Sscanf(errMsg, "error at line %d", &line); err == nil {
		return line - 1, 0 // LSP uses 0-based indexing
	}

	// For now, return first line
	return 0, 0
}
