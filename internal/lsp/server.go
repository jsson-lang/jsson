package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"jsson/internal/lexer"
	"jsson/internal/parser"
)

type Server struct {
	reader    io.Reader
	writer    io.Writer
	mu        sync.Mutex
	documents map[string]*Document
}

type Document struct {
	URI     string
	Content string
	Version int
}

func NewServer(reader io.Reader, writer io.Writer) *Server {
	return &Server{
		reader:    reader,
		writer:    writer,
		documents: make(map[string]*Document),
	}
}

func (s *Server) Start() error {
	log.Println("JSSON Language Server starting...")

	for {
		msg, err := s.readMessage()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Printf("Error reading message: %v", err)
			return err
		}

		if err := s.handleMessage(msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

func (s *Server) readMessage() (map[string]interface{}, error) {
	reader := bufio.NewReader(s.reader)
	var headers = make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil || contentLength == 0 {
		return nil, fmt.Errorf("invalid Content-Length: %s", contentLengthStr)
	}

	content := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, content); err != nil {
		return nil, err
	}

	var msg map[string]interface{}
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *Server) writeMessage(msg interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	if _, err := s.writer.Write([]byte(header)); err != nil {
		return err
	}

	if _, err := s.writer.Write(content); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleMessage(msg map[string]interface{}) error {
	method, ok := msg["method"].(string)
	if !ok {
		return nil
	}

	id := msg["id"]
	params := msg["params"]

	log.Printf("Received method: %s", method)

	switch method {
	case "initialize":
		return s.handleInitialize(id, params)
	case "initialized":
		return nil
	case "textDocument/didOpen":
		return s.handleDidOpen(params)
	case "textDocument/didChange":
		return s.handleDidChange(params)
	case "textDocument/didClose":
		return s.handleDidClose(params)
	case "textDocument/completion":
		return s.handleCompletion(id, params)
	case "textDocument/hover":
		return s.handleHover(id, params)
	case "textDocument/semanticTokens/full":
		return s.handleSemanticTokensFull(id, params)
	case "shutdown":
		return s.handleShutdown(id)
	case "exit":
		return io.EOF
	default:
		log.Printf("Unhandled method: %s", method)
		return nil
	}
}

func (s *Server) handleInitialize(id interface{}, params interface{}) error {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"capabilities": map[string]interface{}{
				"textDocumentSync": map[string]interface{}{
					"openClose": true,
					"change":    1,
				},
				"completionProvider": map[string]interface{}{
					"triggerCharacters": []string{".", ":", "="},
				},
				"hoverProvider": true,
				"semanticTokensProvider": map[string]interface{}{
					"legend": map[string]interface{}{
						"tokenTypes": []string{
							"namespace", "type", "class", "enum", "interface",
							"struct", "typeParameter", "parameter", "variable", "property",
							"enumMember", "event", "function", "method", "macro",
							"keyword", "modifier", "comment", "string", "number",
							"regexp", "operator",
						},
						"tokenModifiers": []string{},
					},
					"full": true,
				},
			},
			"serverInfo": map[string]interface{}{
				"name":    "jsson-lsp",
				"version": "0.0.6",
			},
		},
	}

	return s.writeMessage(response)
}

// handleDidOpen handles textDocument/didOpen notification
func (s *Server) handleDidOpen(params interface{}) error {
	p, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid params")
	}

	textDoc, ok := p["textDocument"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid textDocument")
	}

	uri := textDoc["uri"].(string)
	text := textDoc["text"].(string)
	version := int(textDoc["version"].(float64))

	s.mu.Lock()
	s.documents[uri] = &Document{
		URI:     uri,
		Content: text,
		Version: version,
	}
	s.mu.Unlock()

	return s.publishDiagnostics(uri, text)
}

// handleDidChange handles textDocument/didChange notification
func (s *Server) handleDidChange(params interface{}) error {
	p, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid params")
	}

	textDoc, ok := p["textDocument"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid textDocument")
	}

	uri := textDoc["uri"].(string)
	version := int(textDoc["version"].(float64))

	changes, ok := p["contentChanges"].([]interface{})
	if !ok || len(changes) == 0 {
		return fmt.Errorf("invalid contentChanges")
	}

	change := changes[0].(map[string]interface{})
	text := change["text"].(string)

	s.mu.Lock()
	s.documents[uri] = &Document{
		URI:     uri,
		Content: text,
		Version: version,
	}
	s.mu.Unlock()

	return s.publishDiagnostics(uri, text)
}

// handleDidClose handles textDocument/didClose notification
func (s *Server) handleDidClose(params interface{}) error {
	p, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid params")
	}

	textDoc, ok := p["textDocument"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid textDocument")
	}

	uri := textDoc["uri"].(string)

	s.mu.Lock()
	delete(s.documents, uri)
	s.mu.Unlock()

	return nil
}

// handleShutdown handles the shutdown request
func (s *Server) handleShutdown(id interface{}) error {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  nil,
	}

	return s.writeMessage(response)
}

// getDocument retrieves a document from the cache
func (s *Server) getDocument(uri string) (*Document, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.documents[uri]
	return doc, ok
}

// parseDocument parses a JSSON document and returns any errors
func (s *Server) parseDocument(content string) []string {
	l := lexer.New(content)
	p := parser.New(l)

	_ = p.ParseProgram()

	return p.Errors()
} // Context is required by some functions but not used
var _ = context.Background()
