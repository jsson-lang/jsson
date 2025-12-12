package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// HTTP Server Tests
// ============================================================================

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if healthResp.Status != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", healthResp.Status)
	}
	if healthResp.Service != "jsson" {
		t.Errorf("expected service 'jsson', got '%s'", healthResp.Service)
	}
	if healthResp.JssonVersion != Version {
		t.Errorf("expected jsson version '%s', got '%s'", Version, healthResp.JssonVersion)
	}
}

func TestVersionEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()

	versionHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var versionResp VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if versionResp.JssonVersion != Version {
		t.Errorf("expected jsson version '%s', got '%s'", Version, versionResp.JssonVersion)
	}
	if versionResp.ServerVersion != ServerVersion {
		t.Errorf("expected server version '%s', got '%s'", ServerVersion, versionResp.ServerVersion)
	}
}

func TestTranspileEndpoint(t *testing.T) {
	reqBody := TranspileRequest{
		Source: `app { name = "test" port = 8080 }`,
		Format: "json",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transpile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	transpileHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var transpileResp TranspileResponse
	if err := json.NewDecoder(resp.Body).Decode(&transpileResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !transpileResp.Success {
		t.Errorf("expected success, got errors: %v", transpileResp.Errors)
	}

	if transpileResp.Output == nil {
		t.Error("expected output, got nil")
	}
}

func TestTranspileEndpointYAML(t *testing.T) {
	reqBody := TranspileRequest{
		Source: `config { debug = true }`,
		Format: "yaml",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transpile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	transpileHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var transpileResp TranspileResponse
	if err := json.NewDecoder(resp.Body).Decode(&transpileResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !transpileResp.Success {
		t.Errorf("expected success, got errors: %v", transpileResp.Errors)
	}

	if transpileResp.Format != "yaml" {
		t.Errorf("expected format 'yaml', got '%s'", transpileResp.Format)
	}

	if transpileResp.OutputRaw == "" {
		t.Error("expected raw YAML output")
	}
}

func TestValidateEndpoint(t *testing.T) {
	reqBody := ValidateRequest{
		Source: `app { name = "test" }`,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	validateHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var validateResp ValidateResponse
	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !validateResp.Valid {
		t.Errorf("expected valid JSSON, got errors: %v", validateResp.Errors)
	}
}

func TestValidateEndpointInvalidSyntax(t *testing.T) {
	reqBody := ValidateRequest{
		Source: `app { name = `,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/validate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	validateHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var validateResp ValidateResponse
	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if validateResp.Valid {
		t.Error("expected invalid JSSON, got valid")
	}

	if len(validateResp.Errors) == 0 {
		t.Error("expected errors, got none")
	}
}

func TestValidateWithSchemaEndpoint(t *testing.T) {
	reqBody := ValidateWithSchemaRequest{
		Source: `app { name = "test" port = 8080 }`,
		Schema: `{
			"type": "object",
			"properties": {
				"app": {
					"type": "object",
					"properties": {
						"name": {"type": "string"},
						"port": {"type": "integer"}
					}
				}
			}
		}`,
		OutputFormat: "json",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/validate-schema", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	validateWithSchemaHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var schemaResp ValidateWithSchemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&schemaResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !schemaResp.Valid {
		t.Errorf("expected valid, got errors: %v", schemaResp.Errors)
	}

	if schemaResp.TranspiledData == nil {
		t.Error("expected transpiled data, got nil")
	}
}

func TestValidateWithSchemaEndpointFailure(t *testing.T) {
	reqBody := ValidateWithSchemaRequest{
		Source: `app { name = "test" port = -1 }`,
		Schema: `{
			"type": "object",
			"properties": {
				"app": {
					"type": "object",
					"properties": {
						"name": {"type": "string"},
						"port": {"type": "integer", "minimum": 1}
					}
				}
			}
		}`,
		OutputFormat: "json",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/validate-schema", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	validateWithSchemaHandler(w, req)

	resp := w.Result()

	var schemaResp ValidateWithSchemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&schemaResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if schemaResp.Valid {
		t.Error("expected validation to fail for port=-1")
	}

	if len(schemaResp.Errors) == 0 {
		t.Error("expected validation errors, got none")
	}
}

func TestCORSMiddleware(t *testing.T) {
	serverCORS = true

	handler := corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()

	corsHeader := resp.Header.Get("Access-Control-Allow-Origin")
	if corsHeader != "*" {
		t.Errorf("expected CORS header '*', got '%s'", corsHeader)
	}
}

func TestCORSOptionsRequest(t *testing.T) {
	serverCORS = true

	handler := corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for OPTIONS request")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for OPTIONS, got %d", resp.StatusCode)
	}
}
