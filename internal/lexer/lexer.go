package lexer

import (
	"fmt"
	ie "jsson/internal/errors"
	"jsson/internal/token"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
	line         int
	column       int
	errors       []string
	SourceFile   string
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0, errors: []string{}}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, width := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += width
	}
	l.column++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.ASSIGN, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.LT, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.GT, string(l.ch))
		}
	case '?':
		tok = l.newToken(token.QUESTION, string(l.ch))
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DECLARE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.COLON, string(l.ch))
		}
	case ',':
		tok = l.newToken(token.COMMA, string(l.ch))
	case '{':
		tok = l.newToken(token.LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(token.RBRACE, string(l.ch))
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.ch))
	case '(':
		tok = l.newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.ch))
	case '+':
		tok = l.newToken(token.PLUS, string(l.ch))
	case '-':
		tok = l.newToken(token.MINUS, string(l.ch))
	case '/':
		tok = l.newToken(token.SLASH, string(l.ch))
	case '*':
		tok = l.newToken(token.ASTERISK, string(l.ch))
	case '%':
		tok = l.newToken(token.MODULO, string(l.ch))
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LAND, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LOR, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column}
		} else {
			msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		}
	case '"':
		// Check for triple-quoted raw string
		if l.peekChar() == '"' {
			l.readChar() // consume second "
			if l.peekChar() == '"' {
				l.readChar() // consume third "
				// This is a raw string """..."""
				lit, ok := l.readRawString()
				tok.Line = l.line
				tok.Column = l.column
				if !ok {
					msg := l.lexErrMsg(ie.UnterminatedString())
					l.errors = append(l.errors, msg)
					tok = l.newToken(token.ILLEGAL, msg)
				} else {
					tok.Type = token.RAWSTRING
					tok.Literal = lit
				}
				return tok
			}
			// Only two quotes, backtrack
			// Actually, this is an empty string followed by another quote
			// Let's handle it as empty string
			tok.Type = token.STRING
			tok.Literal = ""
			tok.Line = l.line
			tok.Column = l.column
			l.readChar() // consume the second quote to be ready for next token
			return tok
		}
		// Regular string
		lit, ok := l.readString()
		tok.Line = l.line
		tok.Column = l.column
		if !ok {
			msg := l.lexErrMsg(ie.UnterminatedString())
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		} else {
			tok.Type = token.STRING
			tok.Literal = lit
		}
		return tok
	case '`':
		// Check for triple-backtick raw string
		if l.peekChar() == '`' {
			l.readChar() // consume second `
			if l.peekChar() == '`' {
				l.readChar() // consume third `
				// This is a raw string ```...```
				lit, ok := l.readTripleBacktickString()
				tok.Line = l.line
				tok.Column = l.column
				if !ok {
					msg := l.lexErrMsg(ie.UnterminatedString())
					l.errors = append(l.errors, msg)
					tok = l.newToken(token.ILLEGAL, msg)
				} else {
					tok.Type = token.RAWSTRING
					tok.Literal = lit
				}
				return tok
			}
			// Only two backticks, treat as empty template string
			tok.Type = token.TEMPLATESTR
			tok.Literal = ""
			tok.Line = l.line
			tok.Column = l.column
			return tok
		}
		// Single backtick = template string with interpolation
		lit, ok := l.readTemplateString()
		tok.Line = l.line
		tok.Column = l.column
		if !ok {
			msg := l.lexErrMsg(ie.UnterminatedString())
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		} else {
			tok.Type = token.TEMPLATESTR
			tok.Literal = lit
		}
		return tok
	case '.':
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.RANGE, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = l.newToken(token.DOT, string(l.ch))
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			// Check if the number contains a decimal point
			if containsDot(tok.Literal) {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			msg := l.lexErrMsg(ie.IllegalCharacter(l.ch))
			l.errors = append(l.errors, msg)
			tok = l.newToken(token.ILLEGAL, msg)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	// Build number literal from runes to avoid slicing/index off-by-one issues
	var runes []rune
	for isDigit(l.ch) {
		runes = append(runes, l.ch)
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		runes = append(runes, l.ch)
		l.readChar() // consume .
		for isDigit(l.ch) {
			runes = append(runes, l.ch)
			l.readChar()
		}
	}
	return string(runes)
}

func (l *Lexer) readString() (string, bool) {
	// consume opening quote
	l.readChar()
	var runes []rune

	for {
		if l.ch == '"' {
			// consume closing quote and return content
			l.readChar()
			return string(runes), true
		}
		if l.ch == '\\' {
			// Handle escape sequences
			l.readChar()
			switch l.ch {
			case 'n':
				runes = append(runes, '\n')
			case 't':
				runes = append(runes, '\t')
			case '"':
				runes = append(runes, '"')
			case '\\':
				runes = append(runes, '\\')
			default:
				// Unknown escape, keep literal backslash and char
				runes = append(runes, '\\')
				runes = append(runes, l.ch)
			}
			l.readChar()
			continue
		}
		if l.ch == 0 {
			// unterminated
			return "", false
		}
		runes = append(runes, l.ch)
		l.readChar()
	}
}

// readRawString reads a triple-quoted raw string ("""...""")
// It preserves ALL content literally - no escape processing
func (l *Lexer) readRawString() (string, bool) {
	// Opening """ already consumed
	l.readChar() // move past third quote
	var runes []rune

	for {
		// Check for closing """
		if l.ch == '"' && l.peekChar() == '"' {
			l.readChar() // consume first "
			if l.peekChar() == '"' {
				l.readChar() // consume second "
				l.readChar() // consume third "
				return string(runes), true
			}
			// Only two quotes, add them and continue
			runes = append(runes, '"')
			runes = append(runes, l.ch)
			l.readChar()
			continue
		}
		if l.ch == 0 {
			// unterminated
			return "", false
		}
		// Track newlines for line counting
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		runes = append(runes, l.ch)
		l.readChar()
	}
}

// readTemplateString reads a backtick template string (`...`)
// Preserves content including ${...} for later interpolation parsing
func (l *Lexer) readTemplateString() (string, bool) {
	// consume opening backtick
	l.readChar()
	var runes []rune

	for {
		if l.ch == '`' {
			// consume closing backtick and return content
			l.readChar()
			return string(runes), true
		}
		if l.ch == 0 {
			// unterminated
			return "", false
		}
		// Track newlines for line counting
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		runes = append(runes, l.ch)
		l.readChar()
	}
}

// readTripleBacktickString reads a triple-backtick raw string (```...```)
// Preserves ALL content literally - no interpolation
func (l *Lexer) readTripleBacktickString() (string, bool) {
	// Opening ``` already consumed
	l.readChar() // move past third backtick
	var runes []rune

	for {
		// Check for closing ```
		if l.ch == '`' && l.peekChar() == '`' {
			l.readChar() // consume first `
			if l.peekChar() == '`' {
				l.readChar() // consume second `
				l.readChar() // consume third `
				return string(runes), true
			}
			// Only two backticks, add them and continue
			runes = append(runes, '`')
			runes = append(runes, l.ch)
			l.readChar()
			continue
		}
		if l.ch == 0 {
			// unterminated
			return "", false
		}
		// Track newlines for line counting
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		runes = append(runes, l.ch)
		l.readChar()
	}
}

func (l *Lexer) lexErrf(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	if l.SourceFile != "" {
		ctx := ie.FormatContext(l.SourceFile, l.line, l.column)
		return fmt.Sprintf("Lex goblin: %s — %s", ctx, msg)
	}
	ctx := fmt.Sprintf("%d:%d", l.line, l.column)
	return fmt.Sprintf("Lex goblin: %s — %s", ctx, msg)
}

// lexErrMsg formats an already-formatted error message with context
func (l *Lexer) lexErrMsg(msg string) string {
	if l.SourceFile != "" {
		ctx := ie.FormatContext(l.SourceFile, l.line, l.column)
		return fmt.Sprintf("Lex goblin: %s — %s", ctx, msg)
	}
	ctx := fmt.Sprintf("%d:%d", l.line, l.column)
	return fmt.Sprintf("Lex goblin: %s — %s", ctx, msg)
}

func (l *Lexer) Errors() []string {
	return l.errors
}

func (l *Lexer) SetSourceFile(path string) {
	l.SourceFile = path
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || (l.ch == '/' && l.peekChar() == '/') {
		if l.ch == '/' && l.peekChar() == '/' {
			l.skipComment()
			continue
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.readChar()
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func containsDot(s string) bool {
	for _, ch := range s {
		if ch == '.' {
			return true
		}
	}
	return false
}
