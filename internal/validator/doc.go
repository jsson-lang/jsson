/*
JSSON Validator
===============

Validates transpiled JSSON output against schemas in multiple formats.
Supports JSON Schema, YAML Schema, and custom JSSON validation rules.

Features:
  - JSON Schema Draft 7 validation
  - YAML Schema validation
  - TOML structure validation
  - Custom JSSON rules validation
  - Multi-format schema support
  - Detailed error reporting

Usage:

	validator := validator.New()
	result := validator.Validate(transpiledData, schema, "json")

Package Structure:
  - types.go: Core types (ValidationResult, ValidationError, Schema, Validator)
  - schema.go: Schema loading and registration
  - validate.go: Main validation logic
  - formats.go: Format validators (email, uuid, url, etc.)
  - helpers.go: Utility functions
  - parsers.go: TOML and TypeScript parsers
  - quick.go: Quick validation functions
*/
package validator
