package parser

import (
	"fmt"
	"io"
	"strconv"
)

// for now assume that the json is valid and we dont have to check for errors

/*

bool, for JSON booleans
float64, for JSON numbers
string, for JSON strings
[]any, for JSON arrays
map[string]any, for JSON objects
nil for JSON null
*/

// there is a finite number of possible states that can be yielded during parsing
type CURRENTSTATE byte

type ValueType int

const (
	CHAR    ValueType = iota // string value | c | incases in a quote state
	NUMBER                   // number value | 1
	BOOLEAN                  // true or false | true false
	NULL                     // null value

)

type Parser struct {
	state CURRENTSTATE
	stack [][]byte
	lexer *Lexer
}

func newParser(l *Lexer) *Parser {
	return &Parser{
		stack: make([][]byte, 0),
		lexer: l,
	}
}

func (p *Parser) push(b []byte) {
	p.stack = append(p.stack, append([]byte(nil), b...)) // copy safely
}

func (p *Parser) pop() []byte {
	if len(p.stack) == 0 {
		return nil
	}
	last := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return last
}
func (p *Parser) parse_object() map[string]any {
	return p.consume()
}
func (p *Parser) parse_array() []any {
	var arr []any
	for {
		tok, _ := p.lexer.NextToken()
		switch tok.Type {
		case TokenRightBracket, TokenEOF:
			return arr
		case TokenComma:
			continue
		case TokenString, TokenNumber, TokenTrue, TokenFalse, TokenNull:
			trueVal := parseLiteral(tok)
			fmt.Printf("Parser.parse_array(): passed in value %T parsed value %T\n", tok.Value, trueVal)
			arr = append(arr, parseLiteral(tok))
		case TokenLeftBrace:
			arr = append(arr, p.parse_object())
		case TokenLeftBracket:
			nested := p.parse_array()
			arr = append(arr, nested)

		}
	}
}
func parseLiteral(tok Token) any {
	switch tok.Type {
	case TokenTrue:
		return true
	case TokenFalse:
		return false
	case TokenNull:
		return nil
	case TokenNumber:
		if i, err := strconv.ParseFloat(tok.Value, 64); err == nil {
			return i
		}
	}
	return tok.Value
}

// {"rich":"tmp"}
func (p *Parser) consume() map[string]any {
	obj := make(map[string]any)
	for {
		tok, err := p.lexer.NextToken()
		if err != nil {
			break
		}
		switch tok.Type {
		case TokenLeftBrace: // parse as  bject
			continue
		case TokenLeftBracket: // parse as array
			p.parse_array()
		case TokenString:
			// could be a key or a value
			strVal := tok.Value
			p.push([]byte(strVal))
		case TokenColon:
			value, _ := p.lexer.NextToken()
			key := string(p.pop())

			switch value.Type {
			case TokenString, TokenNumber, TokenTrue, TokenFalse, TokenNull:
				trueVal := parseLiteral(value)
				fmt.Printf("Parser.Consume(): passed in value %T parsed value %T\n", value.Value, trueVal)
				obj[key] = trueVal
			case TokenLeftBrace:
				obj[key] = p.parse_object()
			case TokenLeftBracket:
				obj[key] = p.parse_array()
			}
		case TokenComma:
			// if its not inside and array or object just continue
			continue
		case TokenRightBrace, TokenEOF:
			return obj
		}
	}
	return obj
}

// assume all elements are at the outmost layer for now
func basicParase(data io.Reader) any {
	l := NewLexer(data)
	p := newParser(l)
	return p.parseValue()
}

func (p *Parser) parseValue() any {
	tok, _ := p.lexer.NextToken()

	switch tok.Type {
	case TokenLeftBrace:
		return p.parse_object()
	case TokenLeftBracket:
		return p.parse_array()

	case TokenString, TokenNumber, TokenTrue, TokenFalse, TokenNull:
		return parseLiteral(tok)
	default:
		return nil
	}
}
