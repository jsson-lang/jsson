package validator

// ValidationResult contains the result of a validation operation
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
	Format string            `json:"format,omitempty"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	Path       string `json:"path"`
	Message    string `json:"message"`
	Value      any    `json:"value,omitempty"`
	SchemaPath string `json:"schemaPath,omitempty"`
	Expected   string `json:"expected,omitempty"`
}

// Schema represents a JSON Schema for validation
type Schema struct {
	Type                 string             `json:"type,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	MinLength            *int               `json:"minLength,omitempty"`
	MaxLength            *int               `json:"maxLength,omitempty"`
	Minimum              *float64           `json:"minimum,omitempty"`
	Maximum              *float64           `json:"maximum,omitempty"`
	Pattern              string             `json:"pattern,omitempty"`
	Enum                 []any              `json:"enum,omitempty"`
	Format               string             `json:"format,omitempty"`
	JSSonFormat          string             `json:"jssonFormat,omitempty"`
	AdditionalProperties *bool              `json:"additionalProperties,omitempty"`
	MinItems             *int               `json:"minItems,omitempty"`
	MaxItems             *int               `json:"maxItems,omitempty"`
	UniqueItems          bool               `json:"uniqueItems,omitempty"`
	Const                any                `json:"const,omitempty"`
	Default              any                `json:"default,omitempty"`
	Description          string             `json:"description,omitempty"`
	Title                string             `json:"title,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"`
	AnyOf                []*Schema          `json:"anyOf,omitempty"`
	AllOf                []*Schema          `json:"allOf,omitempty"`
	Not                  *Schema            `json:"not,omitempty"`
}

// Validator is the main validation engine
type Validator struct {
	schemas     map[string]*Schema
	customRules map[string]ValidationRule
}

// ValidationRule is a custom validation function
type ValidationRule func(value any, params map[string]any) bool
