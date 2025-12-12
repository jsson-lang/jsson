// Package validator re-exports jsson/internal/validator for public use.
// This package provides schema validation for transpiled JSSON output.
package validator

import "jsson/internal/validator"

// Re-export types
type (
	Validator        = validator.Validator
	ValidationResult = validator.ValidationResult
	ValidationError  = validator.ValidationError
	Schema           = validator.Schema
	ValidationRule   = validator.ValidationRule
)

// Re-export constructor
var New = validator.New

// Re-export quick validation functions
var (
	QuickValidate     = validator.QuickValidate
	QuickValidateJSON = validator.QuickValidateJSON
	QuickValidateYAML = validator.QuickValidateYAML
)

// Re-export inline validators
var (
	ValidateUUID     = validator.ValidateUUID
	ValidateEmail    = validator.ValidateEmail
	ValidateURL      = validator.ValidateURL
	ValidateIPv4     = validator.ValidateIPv4
	ValidateIPv6     = validator.ValidateIPv6
	ValidateFilePath = validator.ValidateFilePath
	ValidateDate     = validator.ValidateDate
	ValidateDateTime = validator.ValidateDateTime
	ValidateRegex    = validator.ValidateRegex
)

// Re-export format validators
var (
	ValidateFormat  = validator.ValidateFormat
	IsValidEmail    = validator.IsValidEmail
	IsValidURL      = validator.IsValidURL
	IsValidUUID     = validator.IsValidUUID
	IsValidIPv4     = validator.IsValidIPv4
	IsValidIPv6     = validator.IsValidIPv6
	IsValidDate     = validator.IsValidDate
	IsValidDateTime = validator.IsValidDateTime
	IsValidSemVer   = validator.IsValidSemVer
	IsValidHexColor = validator.IsValidHexColor
	IsValidPort     = validator.IsValidPort
)
