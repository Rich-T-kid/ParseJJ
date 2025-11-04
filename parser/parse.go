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
	state    CURRENTSTATE
	stack    [][]byte
	inQuote  bool
	inNumber bool
	setKV    bool
}

func newParser() *Parser {
	return &Parser{
		stack: make([][]byte, 0),
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

// assume all elements are at the outmost layer for now
func basicParase(data io.Reader) map[string]any {
	result := make(map[string]any)
	buffer := make([]byte, 1024)
	p := newParser()
	for {
		n, err := data.Read(buffer)
		if err == io.EOF {
			return result
		}
		// increment our selves
		for i := 0; i < n; i++ {
			switch buffer[i] {
			case byte(LEFTBRACE):
				p.state = LEFTBRACE
				// LEFTBRACE
			case byte(RIGHTBRACE):
				p.state = RIGHTBRACE
				// RIGHTBRACE
			case byte(QUOTE):
				p.state = QUOTE
				// read until next quote
				for j := i + 1; j < n; j++ {
					if buffer[j] == byte(QUOTE) {
						// found end quote
						content := append([]byte(nil), buffer[i+1:j]...)
						if p.setKV {
							// if this is the second part of a key:value pair write this to result
							key := p.pop()
							fmt.Printf("setting %s to %s", key, content)
							result[string(key)] = string(content)

							p.setKV = false
						} else {
							// otherwise this is a key so store it in the stack
							p.push(content)
						}
						i = j // move i to j
						break
					}
				}
				// QUOTE
			case byte(COLON):
				p.state = COLON
				p.setKV = true
			case byte(COMMA):
				p.state = COMMA
			}

		}
	}

}
