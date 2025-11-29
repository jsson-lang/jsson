package parser

import (
	"fmt"
	"jsson/internal/ast"
	ie "jsson/internal/errors"
	"jsson/internal/lexer"
	"jsson/internal/token"
	"strconv"
	"strings"
)

const (
	_ int = iota
	LOWEST
	TERNARY     // ? :
	LOGICAL     // && ||
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SUM         // + -
	PRODUCT     // * / %
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index] or obj.prop
)

var precedences = map[token.TokenType]int{
	token.LAND:     LOGICAL,
	token.LOR:      LOGICAL,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MODULO:   PRODUCT,
	token.QUESTION: TERNARY,
	token.DOT:      INDEX,
	token.RANGE:    INDEX,
	token.MAP:      INDEX, // map binds tightly like a method call
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func (p *Parser) addError(msg string) {
	var loc string
	if p.l != nil && p.l.SourceFile != "" {
		loc = ie.FormatContext(p.l.SourceFile, p.curToken.Line, p.curToken.Column)
	} else {
		loc = fmt.Sprintf("%d:%d", p.curToken.Line, p.curToken.Column)
	}
	fun := "Syntax wizard:"
	p.errors = append(p.errors, fmt.Sprintf("%s %s — %s", fun, loc, msg))
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		// Could be Assignment (key = val), VariableDeclaration (key := val), Object (key { ... }) or ArrayTemplate (key [ ... ])
		if p.peekToken.Type == token.DECLARE {
			return p.parseVariableDeclaration()
		} else if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignment()
		} else if p.peekToken.Type == token.LBRACE {
			return p.parseObjectStatement()
		} else if p.peekToken.Type == token.LBRACKET {
			return p.parseArrayTemplateStatement()
		} else {
			return nil
		}
	default:
		if p.curToken.Type == token.INCLUDE {
			return p.parseIncludeStatement()
		}
		return nil
	}
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}

	p.nextToken() // consume include

	if p.curToken.Type != token.STRING && p.curToken.Type != token.RAWSTRING {
		p.addError(ie.IncludePathExpected())
		return nil
	}

	stmt.Path = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return stmt
}

func (p *Parser) parseAssignment() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	p.nextToken() // consume ASSIGN

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	stmt := &ast.VariableDeclaration{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	p.nextToken() // consume DECLARE (:=)

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseObjectStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseArrayTemplateStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // consume IDENT
	stmt.Value = p.parseArrayTemplate()
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.parsePrefix()
	if prefix == nil {
		return nil
	}

	for p.peekToken.Type != token.EOF && precedence < p.peekPrecedence() {
		infix := p.parseInfix(prefix)
		if infix == nil {
			return prefix
		}
		prefix = infix
	}

	return prefix
}

func (p *Parser) parsePrefix() ast.Expression {
	switch p.curToken.Type {
	case token.IDENT:
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.INT:
		return p.parseIntegerLiteral()
	case token.FLOAT:
		return p.parseFloatLiteral()
	case token.STRING, token.RAWSTRING, token.TEMPLATESTR:
		return p.parseStringLiteral()
	case token.TRUE, token.FALSE:
		return p.parseBooleanLiteral()
	case token.LPAREN:
		return p.parseGroupedExpression()
	case token.LBRACKET:
		return p.parseArrayLiteral()
	case token.LBRACE:
		return p.parseObjectLiteral()
	case token.MINUS:
		// Unary minus for negative numbers
		return p.parsePrefixExpression()
	default:
		return nil
	}
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	switch p.peekToken.Type {
	case token.PLUS, token.MINUS, token.SLASH, token.ASTERISK, token.MODULO,
		token.EQ, token.NEQ, token.LT, token.GT, token.LTE, token.GTE,
		token.LAND, token.LOR:
		p.nextToken()
		return p.parseBinaryExpression(left)
	case token.QUESTION:
		p.nextToken()
		return p.parseConditionalExpression(left)
	case token.DOT:
		p.nextToken()
		return p.parseMemberExpression(left)
	case token.RANGE:
		p.nextToken()
		return p.parseRangeExpression(left)
	case token.MAP:
		p.nextToken()
		return p.parseMapExpression(left)
	default:
		return nil
	}
}

func (p *Parser) parseRangeExpression(left ast.Expression) ast.Expression {
	expr := &ast.RangeExpression{Token: p.curToken, Start: left}

	precedence := p.curPrecedence()
	// move to the token after '..'
	p.nextToken()
	// parse end expression
	expr.End = p.parseExpression(precedence)

	// If there's a step clause after end
	if p.peekToken.Type == token.STEP {
		p.nextToken() // move to STEP
		p.nextToken() // move to step value
		if p.curToken.Type == token.INT {
			expr.Step = p.parseIntegerLiteral()
		} else {
			// try parsing any expression as step
			expr.Step = p.parseExpression(LOWEST)
		}
	}

	return expr
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	expr := &ast.MemberExpression{Token: p.curToken, Left: left}
	p.nextToken() // consume .

	if p.curToken.Type != token.IDENT {
		p.addError(ie.ExpectedIdentifierAfterDot())
		return nil
	}

	expr.Property = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return expr
}

func (p *Parser) parseConditionalExpression(condition ast.Expression) ast.Expression {
	expr := &ast.ConditionalExpression{
		Token:     p.curToken,
		Condition: condition,
	}

	p.nextToken() // consume ?
	expr.Consequence = p.parseExpression(TERNARY)

	if p.peekToken.Type != token.COLON {
		p.addError(ie.MissingColonInTernary())
		return nil
	}

	p.nextToken() // move to :
	p.nextToken() // consume :
	expr.Alternative = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RPAREN {
		p.addError(ie.MissingClosingParen())
		return nil
	}
	p.nextToken()
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// Handle unary minus for negative numbers
	expr := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     &ast.IntegerLiteral{Token: p.curToken, Value: 0}, // 0 - value
	}

	p.nextToken() // consume MINUS
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseArrayTemplate() ast.Expression {
	at := &ast.ArrayTemplate{Token: p.curToken}
	p.nextToken() // consume [

	// Check if there's a template definition
	hasTemplate := p.curToken.Type == token.TEMPLATE

	if hasTemplate {
		p.nextToken() // consume template
		at.Template = p.parseObjectLiteral().(*ast.ObjectLiteral)
		p.nextToken() // consume }
	}

	// Check for map clause
	if p.curToken.Type == token.MAP {
		at.Map = p.parseMapClause()

		// If no template was defined, create an implicit one based on map parameter
		if !hasTemplate && at.Map != nil {
			// Create a template with a single field matching the map parameter name
			at.Template = &ast.ObjectLiteral{
				Token:      at.Map.Token,
				Properties: make(map[string]ast.Expression),
				Keys:       []string{at.Map.Param.Value},
			}
			// The property value is just the identifier itself
			at.Template.Properties[at.Map.Param.Value] = at.Map.Param
		}
	}

	// If still no template, this is an error case
	if at.Template == nil {
		p.addError("array must have either 'template' definition or 'map' clause")
		return at
	}

	at.Rows = [][]ast.Expression{}
	expectedCols := len(at.Template.Keys)

	for p.curToken.Type != token.RBRACKET && p.curToken.Type != token.EOF {
		// Skip any stray closing braces that may remain after nested object parsing
		for p.curToken.Type == token.RBRACE {
			p.nextToken()
		}
		row := []ast.Expression{}
		for i := 0; i < expectedCols; i++ {
			if p.curToken.Type == token.COMMA {
				p.nextToken()
			}
			if p.curToken.Type == token.RBRACKET {
				break
			}
			expr := p.parseExpression(LOWEST)
			if expr != nil {
				row = append(row, expr)
			}
			p.nextToken()
		}
		if len(row) > 0 {
			at.Rows = append(at.Rows, row)
		}
		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	return at
}

func (p *Parser) parseMapClause() *ast.MapClause {
	mc := &ast.MapClause{Token: p.curToken}
	p.nextToken() // consume map
	p.nextToken() // consume (
	mc.Param = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // consume param
	p.nextToken() // consume )
	p.nextToken() // consume =
	mc.Body = p.parseObjectLiteral().(*ast.ObjectLiteral)
	p.nextToken() // consume }
	return mc
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(ie.IntegerTooSpicy(p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	isRaw := p.curToken.Type == token.RAWSTRING
	isTemplate := p.curToken.Type == token.TEMPLATESTR
	value := p.curToken.Literal

	// Check for interpolations in template strings (${...})
	if isTemplate && strings.Contains(value, "${") {
		return p.parseTemplateString(value)
	}

	// Check for interpolations in raw strings ({...}) - old syntax, deprecated
	if isRaw && strings.Contains(value, "{") {
		return p.parseInterpolatedString(value)
	}

	return &ast.StringLiteral{
		Token: p.curToken,
		Value: value,
		IsRaw: isRaw || isTemplate, // Both are raw (no escape processing)
	}
}

// parseTemplateString parses a template string with ${var} interpolations
func (p *Parser) parseTemplateString(content string) ast.Expression {
	interp := &ast.InterpolatedString{
		Token: p.curToken,
		Parts: []interface{}{},
	}

	var currentText strings.Builder
	i := 0

	for i < len(content) {
		if i < len(content)-1 && content[i] == '$' && content[i+1] == '{' {
			// Save any accumulated text
			if currentText.Len() > 0 {
				interp.Parts = append(interp.Parts, currentText.String())
				currentText.Reset()
			}

			// Find matching }
			depth := 1
			start := i + 2 // skip ${
			i += 2
			for i < len(content) && depth > 0 {
				if content[i] == '{' {
					depth++
				} else if content[i] == '}' {
					depth--
				}
				i++
			}

			if depth == 0 {
				// Parse the expression inside ${}
				exprText := content[start:i]
				exprLexer := lexer.New(exprText)
				exprParser := New(exprLexer)
				expr := exprParser.parseExpression(LOWEST)

				// Check if parsing was successful
				if expr != nil && len(exprParser.Errors()) == 0 {
					interp.Parts = append(interp.Parts, expr)
				} else {
					// Failed to parse, treat as literal text
					currentText.WriteString("${")
					currentText.WriteString(exprText)
					currentText.WriteString("}")
				}
			} else {
				// Unmatched ${, treat as literal
				currentText.WriteString(content[start-2 : i])
			}
		} else {
			currentText.WriteByte(content[i])
			i++
		}
	}

	// Add any remaining text
	if currentText.Len() > 0 {
		interp.Parts = append(interp.Parts, currentText.String())
	}

	return interp
}

// parseInterpolatedString parses a raw string with {var} interpolations
func (p *Parser) parseInterpolatedString(content string) ast.Expression {
	interp := &ast.InterpolatedString{
		Token: p.curToken,
		Parts: []interface{}{},
	}

	var currentText strings.Builder
	i := 0

	for i < len(content) {
		if content[i] == '{' {
			// Save any accumulated text
			if currentText.Len() > 0 {
				interp.Parts = append(interp.Parts, currentText.String())
				currentText.Reset()
			}

			// Find matching }
			depth := 1
			start := i + 1
			i++
			for i < len(content) && depth > 0 {
				if content[i] == '{' {
					depth++
				} else if content[i] == '}' {
					depth--
				}
				i++
			}

			if depth == 0 {
				// Parse the expression inside {}
				exprText := content[start:i]
				exprLexer := lexer.New(exprText)
				exprParser := New(exprLexer)
				expr := exprParser.parseExpression(LOWEST)

				if expr != nil && len(exprParser.Errors()) == 0 {
					interp.Parts = append(interp.Parts, expr)
				} else {
					// Failed to parse, treat as literal text
					currentText.WriteString("{")
					currentText.WriteString(exprText)
					currentText.WriteString("}")
				}
			} else {
				// Unmatched {, treat as literal
				currentText.WriteString(content[start-1 : i])
			}
		} else {
			currentText.WriteByte(content[i])
			i++
		}
	}

	// Add any remaining text
	if currentText.Len() > 0 {
		interp.Parts = append(interp.Parts, currentText.String())
	}

	return interp
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == token.TRUE}
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.curToken}
	obj.Properties = make(map[string]ast.Expression)
	obj.Keys = []string{}
	obj.Declarations = []*ast.VariableDeclaration{} // Initialize declarations

	p.nextToken() // consume {

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		if p.curToken.Type != token.IDENT {
			p.nextToken()
			continue
		}

		key := p.curToken.Literal
		p.nextToken() // consume key

		// Check if it's a variable declaration (:=) or property assignment (=)
		if p.curToken.Type == token.DECLARE {
			// Variable declaration: key := value
			p.nextToken() // consume :=
			val := p.parseExpression(LOWEST)
			decl := &ast.VariableDeclaration{
				Token: p.curToken,
				Name:  &ast.Identifier{Value: key},
				Value: val,
			}
			obj.Declarations = append(obj.Declarations, decl)
			p.nextToken() // consume value
		} else if p.curToken.Type == token.ASSIGN || p.curToken.Type == token.COLON {
			// Property assignment: key = value
			obj.Keys = append(obj.Keys, key)
			p.nextToken() // consume = or :
			val := p.parseExpression(LOWEST)
			obj.Properties[key] = val
			p.nextToken() // consume value
		} else if p.curToken.Type == token.LBRACE {
			obj.Keys = append(obj.Keys, key)
			val := p.parseExpression(LOWEST)
			obj.Properties[key] = val
			p.nextToken()
		} else if p.curToken.Type == token.LBRACKET {
			// Support arrays as object property values
			obj.Keys = append(obj.Keys, key)
			val := p.parseArrayLiteral()
			obj.Properties[key] = val
			p.nextToken()
		} else {
			obj.Keys = append(obj.Keys, key)
			obj.Properties[key] = nil
		}

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	if p.curToken.Type != token.RBRACE {
		p.addError(ie.MissingClosingBrace())
	}

	return obj
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = []ast.Expression{}

	p.nextToken() // consume [

	for p.curToken.Type != token.RBRACKET && p.curToken.Type != token.EOF {
		elem := p.parseExpression(LOWEST)
		if elem != nil {
			array.Elements = append(array.Elements, elem)
		}
		p.nextToken()

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	return array
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseMapExpression(left ast.Expression) ast.Expression {
	expression := &ast.MapExpression{Token: p.curToken, Left: left}

	// Expect '('
	if p.peekToken.Type != token.LPAREN {
		p.addError(ie.ExpectedToken(token.LPAREN, p.peekToken.Literal))
		return nil
	}
	p.nextToken() // consume map, now cur is (

	// Expect Identifier (iterator variable)
	if p.peekToken.Type != token.IDENT {
		p.addError(ie.ExpectedToken(token.IDENT, p.peekToken.Literal))
		return nil
	}
	p.nextToken() // consume (, now cur is IDENT
	expression.Iterator = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expect ')'
	if p.peekToken.Type != token.RPAREN {
		p.addError(ie.ExpectedToken(token.RPAREN, p.peekToken.Literal))
		return nil
	}
	p.nextToken() // consume IDENT, now cur is )

	// Expect '='
	if p.peekToken.Type != token.ASSIGN {
		p.addError(ie.ExpectedToken(token.ASSIGN, p.peekToken.Literal))
		return nil
	}
	p.nextToken() // consume ), now cur is =

	p.nextToken() // consume =, now cur is start of expression

	// Parse Body
	expression.Body = p.parseExpression(LOWEST)

	return expression
}
