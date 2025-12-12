package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + Literals
	IDENT       = "IDENT"       // user, name, age
	INT         = "INT"         // 123
	FLOAT       = "FLOAT"       // 123.45
	STRING      = "STRING"      // "hello"
	RAWSTRING   = "RAWSTRING"   // """raw text"""
	TEMPLATESTR = "TEMPLATESTR" // `template ${var}`

	// Operators
	ASSIGN   = "="
	DECLARE  = ":=" // Variable declaration
	COLON    = ":"
	QUESTION = "?"
	EQ       = "=="
	NEQ      = "!="
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="
	RANGE    = ".."
	DOT      = "."
	PLUS     = "+"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	MODULO   = "%"
	LAND     = "&&" // Logical AND
	LOR      = "||" // Logical OR

	// Delimiters
	COMMA    = ","
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	LPAREN   = "("
	RPAREN   = ")"
	AT       = "@" // Preset directive prefix

	// Keywords
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	TEMPLATE = "TEMPLATE"
	MAP      = "MAP"
	INCLUDE  = "INCLUDE"
	STEP     = "STEP"
	PRESET   = "PRESET" // @preset directive
	USE      = "USE"    // @use directive

	// Boolean literals extras
	YES = "YES"
	NO  = "NO"
	ON  = "ON"
	OFF = "OFF"

	// Validators
	UUID     = "UUID"
	EMAIL    = "EMAIL"
	URL      = "URL"
	IPV4     = "IPV4"
	IPV6     = "IPV6"
	FILEPATH = "FILEPATH"
	DATE     = "DATE"
	DATETIME = "DATETIME"
	REGEX    = "REGEX"
	VINT     = "VINT"   // @int(min, max)
	VFLOAT   = "VFLOAT" // @float(min, max)
	VBOOL    = "VBOOL"  // @bool
)

var keywords = map[string]TokenType{
	"true":     TRUE,
	"false":    FALSE,
	"template": TEMPLATE,
	"map":      MAP,
	"include":  INCLUDE,
	"step":     STEP,
	"preset":   PRESET,
	"use":      USE,

	// Boolean literals extras
	"yes": YES,
	"no":  NO,
	"on":  ON,
	"off": OFF,

	// Validators
	"uuid":     UUID,
	"email":    EMAIL,
	"url":      URL,
	"ipv4":     IPV4,
	"ipv6":     IPV6,
	"filepath": FILEPATH,
	"date":     DATE,
	"datetime": DATETIME,
	"regex":    REGEX,
	"int":      VINT,
	"float":    VFLOAT,
	"bool":     VBOOL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
