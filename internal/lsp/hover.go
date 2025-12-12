package lsp

import (
	"fmt"
	"strings"
)

func (s *Server) handleHover(id interface{}, params interface{}) error {
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

	hoverInfo := s.getHoverInfo(doc.Content, line, character)

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  hoverInfo,
	}

	return s.writeMessage(response)
}

// getHoverInfo returns hover information for the given position
func (s *Server) getHoverInfo(content string, line, character int) interface{} {
	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return nil
	}

	currentLine := lines[line]
	if character >= len(currentLine) {
		return nil
	}

	// Get the word at the cursor position
	word := s.getWordAtPosition(currentLine, character)
	if word == "" {
		return nil
	}

	// Check if it's a parameter in a map/zip expression
	paramInfo := s.getParameterInfo(currentLine, word, character)
	if paramInfo != "" {
		return map[string]interface{}{
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": paramInfo,
			},
		}
	}

	// Check if it's a variable declaration or usage
	varInfo := s.getVariableInfo(content, word, line)
	if varInfo != "" {
		return map[string]interface{}{
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": varInfo,
			},
		}
	}

	// Check for object/array names
	objInfo := s.getObjectInfo(content, word, line)
	if objInfo != "" {
		return map[string]interface{}{
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": objInfo,
			},
		}
	}

	// Get documentation for the word (keywords/operators)
	doc := s.getDocumentation(word)
	if doc == "" {
		return nil
	}

	return map[string]interface{}{
		"contents": map[string]interface{}{
			"kind":  "markdown",
			"value": doc,
		},
	}
}

// getParameterInfo returns information about a parameter in map/zip
func (s *Server) getParameterInfo(line, paramName string, character int) string {
	// Check if we're inside a map or zip expression
	// Pattern: (... map (param) = ...) or (... zip (param1, param2) = ...)

	// Look backwards from cursor position for "map (" or "zip ("
	beforeCursor := line[:character]

	if strings.Contains(beforeCursor, "map (") || strings.Contains(beforeCursor, "zip (") {
		// Check if the parameter is actually in the parentheses
		startParen := strings.LastIndex(beforeCursor, "(")
		if startParen != -1 {
			afterParen := line[startParen:]
			endParen := strings.Index(afterParen, ")")
			if endParen != -1 {
				paramsSection := afterParen[1:endParen]
				if strings.Contains(paramsSection, paramName) {
					if strings.Contains(beforeCursor, "map (") {
						return fmt.Sprintf("# Parameter: `%s`\n\n**Map parameter** - represents each item in the mapped collection.\n\n**Example:**\n```jsson\nitems = (1..5 map (%s) = %s * 2)\n```",
							paramName, paramName, paramName)
					} else if strings.Contains(beforeCursor, "zip (") {
						return fmt.Sprintf("# Parameter: `%s`\n\n**Zip parameter** - represents corresponding items from parallel ranges.\n\n**Example:**\n```jsson\npairs [\n  template { a, b }\n  1..3, 10..12\n]\n```",
							paramName)
					}
				}
			}
		}
	}

	return ""
} // getVariableInfo returns information about a variable
func (s *Server) getVariableInfo(content, varName string, currentLine int) string {
	lines := strings.Split(content, "\n")

	// Look for variable declaration
	for i, line := range lines {
		if strings.Contains(line, varName+" :=") || strings.Contains(line, varName+":=") {
			parts := strings.Split(line, ":=")
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				return fmt.Sprintf("# Variable: `%s`\n\n**Declared at line %d**\n\n```jsson\n%s\n```\n\n**Value:** `%s`",
					varName, i+1, strings.TrimSpace(line), value)
			}
		}
	}

	return ""
}

// getObjectInfo returns information about an object or array
func (s *Server) getObjectInfo(content, name string, currentLine int) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for object definition: name {
		if strings.HasPrefix(trimmed, name+" {") || strings.HasPrefix(trimmed, name+"{") {
			return fmt.Sprintf("# Object: `%s`\n\n**Defined at line %d**\n\n```jsson\n%s\n```",
				name, i+1, strings.TrimSpace(line))
		}

		// Check for array definition: name [
		if strings.HasPrefix(trimmed, name+" [") || strings.HasPrefix(trimmed, name+"[") {
			return fmt.Sprintf("# Array: `%s`\n\n**Defined at line %d**\n\n```jsson\n%s\n```",
				name, i+1, strings.TrimSpace(line))
		}

		// Check for assignment: name =
		if strings.HasPrefix(trimmed, name+" =") || strings.HasPrefix(trimmed, name+"=") {
			return fmt.Sprintf("# Property: `%s`\n\n**Defined at line %d**\n\n```jsson\n%s\n```",
				name, i+1, strings.TrimSpace(line))
		}
	}

	return ""
} // getWordAtPosition extracts the word at the given character position
func (s *Server) getWordAtPosition(line string, character int) string {
	if character >= len(line) {
		return ""
	}

	// Find the start of the word
	start := character
	for start > 0 && isIdentifierChar(line[start-1]) {
		start--
	}

	// Find the end of the word
	end := character
	for end < len(line) && isIdentifierChar(line[end]) {
		end++
	}

	if start >= end {
		return ""
	}

	return line[start:end]
}

// isIdentifierChar checks if a character is valid in an identifier
func isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// getDocumentation returns documentation for a keyword or symbol
func (s *Server) getDocumentation(word string) string {
	docs := map[string]string{
		"template": "# template\n\nDefines a template for structured data in arrays.\n\n**Example:**\n```jsson\nusers [\n  template { name, age }\n  \n  John, 25\n  Jane, 30\n]\n```",

		"map": "# map\n\nTransforms data by applying a function to each element.\n\n**Example:**\n```jsson\nitems = (1..5 map (x) = x * 2)\n// Result: [2, 4, 6, 8, 10]\n```",

		"zip": "# zip\n\nCombines multiple ranges into tuples.\n\n**Example:**\n```jsson\npairs [\n  template { a, b }\n  \n  1..3, 10..12\n]\n// Result: [{a:1, b:10}, {a:2, b:11}, {a:3, b:12}]\n```",

		"preset": "# @preset\n\nDefines a reusable configuration preset.\n\n**Example:**\n```jsson\n@preset \"base\" {\n  enabled = true\n  timeout = 5000\n}\n\nconfig = @\"base\" {\n  name = \"My Config\"\n}\n```",

		"include": "# include\n\nIncludes another JSSON file.\n\n**Example:**\n```jsson\ninclude \"config.jsson\"\n```",

		"step": "# step\n\nDefines the step size for a range.\n\n**Example:**\n```jsson\nnumbers = 0..10 step 2\n// Result: [0, 2, 4, 6, 8, 10]\n```",

		"true": "# true\n\nBoolean true value.",

		"false": "# false\n\nBoolean false value.",

		"yes": "# yes\n\nBoolean true value (alternative syntax).\n\nEquivalent to `true`.",

		"no": "# no\n\nBoolean false value (alternative syntax).\n\nEquivalent to `false`.",

		"on": "# on\n\nBoolean true value (alternative syntax).\n\nEquivalent to `true`. Commonly used for flags and settings.",

		"off": "# off\n\nBoolean false value (alternative syntax).\n\nEquivalent to `false`. Commonly used for flags and settings.",

		"@uuid": "# @uuid\n\nGenerates UUID v4 (Universally Unique Identifier) values.\n\n**Format:** `xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx` (36 characters)\n\n**Use Cases:**\n- Unique user/entity identifiers\n- Database primary keys\n- API request tracking\n- Session IDs\n\n**Example:**\n```jsson\nusers [\n  template { name }\n  map (u) = {\n    id = @uuid\n    name = u.name\n  }\n  \"Alice\", \"Bob\"\n]\n```\n\n**Output:** `\"7f3e8c2a-4d1b-4e9f-a3c5-8b2d7e6f1a9c\"`",

		"@email": "# @email\n\nGenerates valid email addresses.\n\n**Format:** `user{random}@example.com`\n\n**Use Cases:**\n- User registration data\n- Contact information\n- Notification recipients\n- Test fixtures\n\n**Example:**\n```jsson\ncontacts [\n  template { name }\n  map (c) = {\n    name = c.name\n    email = @email\n    verified = @bool\n  }\n  \"John\", \"Jane\"\n]\n```\n\n**Output:** `\"user8234@example.com\"`",

		"@url": "# @url\n\nGenerates valid HTTP/HTTPS URLs.\n\n**Format:** `https://example.com`\n\n**Use Cases:**\n- API endpoints\n- Webhook URLs\n- External resource links\n- Configuration URLs\n\n**Example:**\n```jsson\nservices [\n  template { name }\n  map (s) = {\n    name = s.name\n    endpoint = @url\n    active = @bool\n  }\n  \"auth\", \"payment\"\n]\n```\n\n**Output:** `\"https://example.com\"`",

		"@ipv4": "# @ipv4\n\nGenerates IPv4 addresses (32-bit).\n\n**Format:** `xxx.xxx.xxx.xxx` (0-255 per octet)\n\n**Use Cases:**\n- Server IP configurations\n- Network device addresses\n- Firewall rules\n- Load balancer targets\n\n**Example:**\n```jsson\nservers [\n  template { name, port }\n  map (s) = {\n    name = s.name\n    ip = @ipv4\n    port = s.port\n  }\n  \"web-01\", 8080\n  \"web-02\", 8081\n]\n```\n\n**Output:** `\"192.168.142.73\"`",

		"@ipv6": "# @ipv6\n\nGenerates IPv6 addresses (128-bit).\n\n**Format:** `xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx` (hexadecimal)\n\n**Use Cases:**\n- Modern network configurations\n- Cloud infrastructure\n- IoT device addressing\n- Dual-stack environments\n\n**Example:**\n```jsson\ninfra {\n  ipv6 = @ipv6\n  ipv4 = @ipv4\n}\n```\n\n**Output:** `\"2001:0db8:85a3:0000:0000:8a2e:0370:7334\"`",

		"@filepath": "# @filepath\n\nGenerates valid file system paths.\n\n**Format:** `/path/to/file.txt`\n\n**Use Cases:**\n- Log file locations\n- Configuration file paths\n- Resource file references\n- Output destinations\n\n**Example:**\n```jsson\nloggers [\n  template { name, level }\n  map (l) = {\n    name = l.name\n    level = l.level\n    output = @filepath\n  }\n  \"app\", \"info\"\n  \"error\", \"error\"\n]\n```\n\n**Output:** `\"/path/to/file.txt\"`",

		"@date": "# @date\n\nGenerates dates in ISO 8601 format.\n\n**Format:** `YYYY-MM-DD`\n\n**Use Cases:**\n- Event dates\n- Scheduling data\n- Historical records\n- Expiration dates\n\n**Example:**\n```jsson\nevents [\n  template { name }\n  map (e) = {\n    name = e.name\n    date = @date\n    active = @bool\n  }\n  \"Launch\", \"Release\"\n]\n```\n\n**Output:** `\"2025-12-11\"` (current date)",

		"@datetime": "# @datetime\n\nGenerates timestamps in RFC3339/ISO 8601 format.\n\n**Format:** `YYYY-MM-DDTHH:MM:SSZ`\n\n**Use Cases:**\n- Event timestamps\n- Log entries\n- Audit trails\n- API responses\n\n**Example:**\n```jsson\nlogs [\n  template { level, message }\n  map (l) = {\n    timestamp = @datetime\n    level = l.level\n    message = l.message\n  }\n  \"info\", \"Started\"\n  \"error\", \"Failed\"\n]\n```\n\n**Output:** `\"2025-12-11T14:23:45Z\"` (current time)",

		"@regex": "# @regex\n\nGenerates values matching a custom regex pattern.\n\n**Syntax:** `@regex(\"pattern\")`\n\n**Common Patterns:**\n- CPF: `^\\\\d{3}\\\\.\\\\d{3}\\\\.\\\\d{3}-\\\\d{2}$`\n- Phone: `^\\\\(\\\\d{2}\\\\)\\\\s\\\\d{4,5}-\\\\d{4}$`\n- Zip Code: `^\\\\d{5}-\\\\d{3}$`\n\n**Example:**\n```jsson\nusers [\n  template { name }\n  map (u) = {\n    name = u.name\n    cpf = @regex(\"^\\\\d{3}\\\\.\\\\d{3}\\\\.\\\\d{3}-\\\\d{2}$\")\n    phone = @regex(\"^\\\\(\\\\d{2}\\\\)\\\\s\\\\d{5}-\\\\d{4}$\")\n  }\n  \"Alice\", \"Bob\"\n]\n```\n\n**Output:** Placeholder text matching the pattern",

		"@int": "# @int\n\nGenerates random integer values within a specified range (inclusive).\n\n**Syntax:** `@int(min, max)`\n\n**Parameters:**\n- `min`: Minimum value (inclusive)\n- `max`: Maximum value (inclusive)\n\n**Use Cases:**\n- Age ranges (18-65)\n- Scores/ratings (0-100)\n- Port numbers (3000-9000)\n- Quantities (1-1000)\n- IDs (1000-9999)\n\n**Example:**\n```jsson\nproducts [\n  template { name }\n  map (p) = {\n    id = @int(1000, 9999)\n    name = p.name\n    stock = @int(0, 500)\n    rating = @int(1, 5)\n  }\n  \"Laptop\", \"Phone\"\n]\n```\n\n**Output:** Random integer like `42`, `1337`, `8765`",

		"@float": "# @float\n\nGenerates random floating-point values within a specified range.\n\n**Syntax:** `@float(min, max)`\n\n**Parameters:**\n- `min`: Minimum value (inclusive)\n- `max`: Maximum value (exclusive)\n\n**Use Cases:**\n- Prices/amounts (9.99-999.99)\n- Ratings (0.0-5.0)\n- Percentages (0.0-100.0)\n- Coordinates (lat/lng)\n- Temperatures (-50.0-50.0)\n\n**Example:**\n```jsson\nproducts [\n  template { name }\n  map (p) = {\n    name = p.name\n    price = @float(9.99, 999.99)\n    discount = @float(0.0, 0.5)\n    rating = @float(0.0, 5.0)\n  }\n  \"Laptop\", \"Mouse\"\n]\n```\n\n**Output:** Random float like `42.73`, `299.99`, `3.87`",

		"@bool": "# @bool\n\nGenerates random boolean values (50/50 chance).\n\n**Syntax:** `@bool`\n\n**Output:** `true` or `false`\n\n**Use Cases:**\n- Feature flags\n- User status (active/verified)\n- Configuration toggles\n- Conditional data\n- A/B testing\n\n**Example:**\n```jsson\nusers [\n  template { name, age }\n  map (u) = {\n    name = u.name\n    age = u.age\n    active = @bool\n    verified = @bool\n    premium = @bool\n  }\n  \"Alice\", @int(18, 65)\n  \"Bob\", @int(18, 65)\n]\n```\n\n**Output:** Randomly `true` or `false`",

		"null": "# null\n\nRepresents a null value.",

		"and": "# and\n\nLogical AND operator.\n\n**Example:**\n```jsson\nresult = x > 5 and y < 10\n```\n\nAlternatively, use `&&` operator.",

		"or": "# or\n\nLogical OR operator.\n\n**Example:**\n```jsson\nresult = x > 5 or y < 10\n```\n\nAlternatively, use `||` operator.",

		"not": "# not\n\nLogical NOT operator.\n\n**Example:**\n```jsson\nresult = not isActive\n```\n\nAlternatively, use `!` operator.",
	}

	if doc, ok := docs[word]; ok {
		return doc
	}

	// Check for operators
	operatorDocs := map[string]string{
		":=": "# Variable Declaration\n\nDeclares a variable.\n\n**Example:**\n```jsson\napi_url := \"https://api.example.com\"\n```",

		"..": "# Range Operator\n\nCreates a numeric range.\n\n**Example:**\n```jsson\nnumbers = 1..10\n// Result: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]\n```",

		"?": "# Ternary Operator\n\nConditional expression.\n\n**Example:**\n```jsson\nstatus = age >= 18 ? \"adult\" : \"minor\"\n```",
	}

	if doc, ok := operatorDocs[word]; ok {
		return doc
	}

	return ""
}
