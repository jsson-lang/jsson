package lsp

import (
	"strings"
)

type CompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind"`
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
}

const (
	CompletionItemKindText     = 1
	CompletionItemKindMethod   = 2
	CompletionItemKindFunction = 3
	CompletionItemKindKeyword  = 14
	CompletionItemKindVariable = 6
	CompletionItemKindProperty = 10
	CompletionItemKindSnippet  = 15
)

func (s *Server) handleCompletion(id interface{}, params interface{}) error {
	p, ok := params.(map[string]interface{})
	if !ok {
		return s.sendError(id, -32602, "Invalid params")
	}

	textDoc, ok := p["textDocument"].(map[string]interface{})
	if !ok {
		return s.sendError(id, -32602, "Invalid textDocument")
	}

	position, ok := p["position"].(map[string]interface{})
	if !ok {
		return s.sendError(id, -32602, "Invalid position")
	}

	uri := textDoc["uri"].(string)
	line := int(position["line"].(float64))
	character := int(position["character"].(float64))

	doc, ok := s.getDocument(uri)
	if !ok {
		return s.sendError(id, -32602, "Document not found")
	}

	items := s.getCompletionItems(doc.Content, line, character)

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"isIncomplete": false,
			"items":        items,
		},
	}

	return s.writeMessage(response)
}

func (s *Server) getCompletionItems(content string, line, character int) []CompletionItem {
	items := []CompletionItem{}

	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return items
	}

	currentLine := lines[line]
	beforeCursor := ""
	if character <= len(currentLine) {
		beforeCursor = currentLine[:character]
	}

	if strings.HasSuffix(beforeCursor, "..") || strings.HasSuffix(beforeCursor, ".. ") {
		items = append(items, s.getRangeCompletions(content, line)...)
		return items
	}

	items = append(items, s.getKeywordCompletions()...)
	items = append(items, s.getSnippetCompletions()...)

	if strings.HasSuffix(strings.TrimSpace(beforeCursor), ".") {
		items = append(items, s.getPropertyCompletions()...)
	}

	items = append(items, s.getVariableCompletions(content)...)
	items = append(items, s.getScopeVariables(content, line)...)

	return items
}

func (s *Server) getScopeVariables(content string, currentLine int) []CompletionItem {
	items := []CompletionItem{}
	lines := strings.Split(content, "\n")

	if currentLine >= len(lines) {
		return items
	}

	params := make(map[string]bool)
	openParens := 0

	for i := currentLine; i >= 0; i-- {
		line := lines[i]

		for _, ch := range line {
			if ch == '(' {
				openParens++
			} else if ch == ')' {
				openParens--
			}
		}

		if strings.Contains(line, " map (") || strings.Contains(line, " zip (") {
			startIdx := strings.Index(line, " map (")
			if startIdx == -1 {
				startIdx = strings.Index(line, " zip (")
			}
			if startIdx != -1 {
				rest := line[startIdx:]
				parenStart := strings.Index(rest, "(")
				parenEnd := strings.Index(rest, ")")
				if parenStart != -1 && parenEnd != -1 && parenEnd > parenStart {
					paramsStr := rest[parenStart+1 : parenEnd]
					for _, param := range strings.Split(paramsStr, ",") {
						param = strings.TrimSpace(param)
						if param != "" && isValidIdentifier(param) {
							params[param] = true
						}
					}
				}
			}
		}

		if openParens <= 0 && i < currentLine {
			break
		}
	}

	for param := range params {
		items = append(items, CompletionItem{
			Label:  param,
			Kind:   CompletionItemKindVariable,
			Detail: "Parameter (in scope)",
		})
	}

	return items
}

func (s *Server) getRangeCompletions(content string, currentLine int) []CompletionItem {
	items := []CompletionItem{}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		// Look for variable declarations (name := value)
		if strings.Contains(line, ":=") {
			parts := strings.Split(line, ":=")
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Check if it's likely a number
				if varName != "" && isValidIdentifier(varName) {
					// Check if value is numeric or contains numbers
					isNumeric := false
					for _, ch := range value {
						if ch >= '0' && ch <= '9' {
							isNumeric = true
							break
						}
					}

					if isNumeric {
						items = append(items, CompletionItem{
							Label:  varName,
							Kind:   CompletionItemKindVariable,
							Detail: "Range end: " + value,
						})
					}
				}
			}
		}
	}

	// Add common range patterns
	items = append(items, CompletionItem{
		Label:      "10",
		Kind:       CompletionItemKindText,
		Detail:     "Range: 0..10",
		InsertText: "10",
	})

	items = append(items, CompletionItem{
		Label:      "100",
		Kind:       CompletionItemKindText,
		Detail:     "Range: 0..100",
		InsertText: "100",
	})

	return items
} // getKeywordCompletions returns keyword completions
func (s *Server) getKeywordCompletions() []CompletionItem {
	keywords := []struct {
		label  string
		detail string
		doc    string
	}{
		{"template", "Template definition", "Define a template for structured data"},
		{"map", "Map transformation", "Transform data with a map function"},
		{"zip", "Zip ranges", "Combine multiple ranges into tuples"},
		{"include", "Include file", "Include another JSSON file"},
		{"step", "Range step", "Define step size for ranges"},
		{"@preset", "Preset definition", "Define a reusable preset configuration"},
		{"true", "Boolean true", "Boolean true value"},
		{"false", "Boolean false", "Boolean false value"},
		{"yes", "Boolean true", "Boolean true (alternative syntax)"},
		{"no", "Boolean false", "Boolean false (alternative syntax)"},
		{"on", "Boolean true", "Boolean true (alternative syntax)"},
		{"off", "Boolean false", "Boolean false (alternative syntax)"},
		{"null", "Null value", "Null value"},
		{"@uuid", "UUID validator", "Validates/generates UUID"},
		{"@email", "Email validator", "Validates/generates email"},
		{"@url", "URL validator", "Validates/generates URL"},
		{"@ipv4", "IPv4 validator", "Validates/generates IPv4"},
		{"@ipv6", "IPv6 validator", "Validates/generates IPv6"},
		{"@filepath", "File path validator", "Validates file path"},
		{"@date", "Date validator", "Validates/generates date"},
		{"@datetime", "DateTime validator", "Validates/generates datetime"},
		{"@regex", "Regex validator", "Validates with regex pattern"},
		{"@int", "Integer validator", "Generates random integer with min/max range"},
		{"@float", "Float validator", "Generates random float with min/max range"},
		{"@bool", "Boolean validator", "Generates random boolean value"},
		{"and", "Logical AND", "Logical AND operator"},
		{"or", "Logical OR", "Logical OR operator"},
		{"not", "Logical NOT", "Logical NOT operator"},
	}

	items := make([]CompletionItem, 0, len(keywords))
	for _, kw := range keywords {
		items = append(items, CompletionItem{
			Label:         kw.label,
			Kind:          CompletionItemKindKeyword,
			Detail:        kw.detail,
			Documentation: kw.doc,
		})
	}

	return items
} // getSnippetCompletions returns snippet completions
func (s *Server) getSnippetCompletions() []CompletionItem {
	return []CompletionItem{
		{
			Label:         "template",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Template array",
			Documentation: "Create an array with template definition",
			InsertText:    "[\n  template { $1 }\n\n  $2\n]",
		},
		{
			Label:         "map",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Map transformation",
			Documentation: "Transform data with map",
			InsertText:    "($1 map ($2) = $3)",
		},
		{
			Label:         "variable",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Variable declaration",
			Documentation: "Declare a variable with :=",
			InsertText:    "$1 := $2",
		},
		{
			Label:         "range",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Numeric range",
			Documentation: "Create a numeric range",
			InsertText:    "$1..$2",
		},
		{
			Label:         "ternary",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Ternary operator",
			Documentation: "Conditional expression",
			InsertText:    "$1 ? $2 : $3",
		},
		{
			Label:         "preset",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Preset definition",
			Documentation: "Define a reusable preset",
			InsertText:    "@preset \"$1\" {\n  $2\n}",
		},
		{
			Label:         "use preset",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Use preset",
			Documentation: "Apply a preset configuration",
			InsertText:    "@\"$1\" {\n  $2\n}",
		},
		{
			Label:         "object",
			Kind:          CompletionItemKindSnippet,
			Detail:        "Object definition",
			Documentation: "Create an object",
			InsertText:    "$1 {\n  $2\n}",
		},
	}
}

// getPropertyCompletions returns property completions (for obj.prop)
func (s *Server) getPropertyCompletions() []CompletionItem {
	// Common properties that might be used
	return []CompletionItem{
		{Label: "id", Kind: CompletionItemKindProperty},
		{Label: "name", Kind: CompletionItemKindProperty},
		{Label: "value", Kind: CompletionItemKindProperty},
		{Label: "type", Kind: CompletionItemKindProperty},
		{Label: "age", Kind: CompletionItemKindProperty},
		{Label: "email", Kind: CompletionItemKindProperty},
	}
}

// getVariableCompletions extracts variable declarations from the document
func (s *Server) getVariableCompletions(content string) []CompletionItem {
	items := []CompletionItem{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// Look for variable declarations (name := value)
		if strings.Contains(line, ":=") {
			parts := strings.Split(line, ":=")
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				if varName != "" && isValidIdentifier(varName) {
					items = append(items, CompletionItem{
						Label:  varName,
						Kind:   CompletionItemKindVariable,
						Detail: "Variable",
					})
				}
			}
		}
	}

	return items
}

// isValidIdentifier checks if a string is a valid identifier
func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}

	// First character must be letter or underscore
	if !((s[0] >= 'a' && s[0] <= 'z') || (s[0] >= 'A' && s[0] <= 'Z') || s[0] == '_') {
		return false
	}

	// Rest can be letters, digits, or underscores
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}

	return true
}

// sendError sends an error response
func (s *Server) sendError(id interface{}, code int, message string) error {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}

	return s.writeMessage(response)
}
