package parser

import (
	"fmt"
	"io"
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

const (
	LEFTBRACE   CURRENTSTATE = '{' // {
	RIGHTBRACE  CURRENTSTATE = '}' // }
	LEFTSQUARE  CURRENTSTATE = '[' // [
	RIGHTSQUARE CURRENTSTATE = ']' // ]
	COMMA       CURRENTSTATE = ',' // ,
	COLON       CURRENTSTATE = ':' // :
	QUOTE       CURRENTSTATE = '"' // "
	SPACE       CURRENTSTATE = ' ' // space
)

type ValueType int

const (
	CHAR    ValueType = iota // string value | c | incases in a quote state
	NUMBER                   // number value | 1
	BOOLEAN                  // true or false | true false
	NULL                     // null value

)

type Parser struct {
	state  CURRENTSTATE
	stack  [][]byte
	reader io.Reader
}

func newParser(r io.Reader) *Parser {
	return &Parser{
		stack:  make([][]byte, 0),
		reader: r,
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
	return nil
}
func (p *Parser) parse_array() []any {
	return nil
}
func (p *Parser) parse_string() string {
	return ""
}
func (p *Parser) parse_number() float64 {
	return 0
}

// these next two can be hand coded lowkey
func (p *Parser) parse_boolean() bool {
	return false
}
func (p *Parser) parse_null() any {
	return nil
}

// {"rich":"tmp"}
func (p *Parser) consume() map[string]any {
	result := make(map[string]any)
	buffer := make([]byte, 1024)
	for {
		n, err := p.reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading:", err)
		}
		for i := 0; i < n; i++ {
			switch buffer[i] {
			case '{':
				p.parse_object()
			case '[':
				p.parse_array()
			case '"':
				p.parse_string()
			case 't', 'f':
				p.parse_boolean()
			case 'n':
				p.parse_null()
			case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
				p.parse_number()
			default:
				// check if number
			}
		}
	}
	return result
}

// assume all elements are at the outmost layer for now
func basicParase(data io.Reader) map[string]any {
	p := newParser(data)
	return p.consume()
}
