package parser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type TokenType int

const (
	TokenLeftBrace TokenType = iota
	TokenRightBrace
	TokenString
	TokenNumber
	TokenTrue
	TokenFalse
	TokenNull
	TokenComma
	TokenColon
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
}
type Lexer struct {
	r   *bufio.Reader
	pos int
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

func (l *Lexer) next() (byte, error) {
	b, err := l.r.ReadByte()
	if err != nil {
		return 0, err
	}
	l.pos++
	return b, nil
}

func (l *Lexer) unread() {
	_ = l.r.UnreadByte()
	l.pos--
}

func (l *Lexer) NextToken() (Token, error) {
	for {
		char, err := l.next()
		if err == io.EOF {
			return Token{Type: TokenEOF}, nil
		}
		if unicode.IsSpace(rune(char)) {
			continue
		}
		switch char {
		case '{':
			return Token{Type: TokenLeftBrace, Value: "{"}, nil
		case '}':
			return Token{Type: TokenRightBrace, Value: "}"}, nil
		case '[':
			return Token{Type: TokenLeftBrace, Value: "["}, nil
		case ']':
			return Token{Type: TokenRightBrace, Value: "]"}, nil
		case ',':
			return Token{Type: TokenComma, Value: ","}, nil
		case ':':
			return Token{Type: TokenColon, Value: ":"}, nil
		case '"':
			return l.lexString()
		case 't', 'f':
			l.unread()
			return l.lexBoolean()
		case 'n':
			l.unread()
			return l.lexNull()
		default:
			if unicode.IsDigit(rune(char)) || char == '-' {
				l.unread()
				return l.lexNumber()
			}
			return Token{}, fmt.Errorf("unexpected character: %c", char)
		}

	}
}

// "input text here"
func (l *Lexer) lexString() (Token, error) {
	var str []byte
	for {
		char, err := l.next()
		if err == io.EOF || char == '"' {
			break
		}
		// keep appending char as long as theres chars and theres no enclosing quote
		str = append(str, char)
	}
	return Token{Type: TokenString, Value: string(str)}, nil
}
func (l *Lexer) lexNumber() (Token, error) {
	var strInt []byte
	for {
		char, err := l.next()
		if err != nil {
			break
		}
		if !unicode.IsDigit(rune(char)) && char != '.' && char != '-' {
			l.unread()
			break
		}
		strInt = append(strInt, char)
	}
	return Token{Type: TokenNumber, Value: string(strInt)}, nil
}
func (l *Lexer) lexBoolean() (Token, error) {
	for {
		char, err := l.next()
		if err != nil {
			break
		}
		switch char {
		case 't': // true
			// we know exactly what is next so we can hardcode it
			for i := range "rue" {
				b, err := l.next()
				if err != nil {
					return Token{}, fmt.Errorf("unexpected end of input in boolean")
				}
				if b != "rue"[i] {
					return Token{}, fmt.Errorf("invalid boolean value")
				}
			}
			return Token{Type: TokenTrue, Value: "true"}, nil
		case 'f':
			for i := range "alse" {
				b, err := l.next()
				if err != nil {
					return Token{}, fmt.Errorf("unexpected end of input in boolean")
				}
				if b != "alse"[i] {
					return Token{}, fmt.Errorf("invalid boolean value")
				}
			}
			return Token{Type: TokenFalse, Value: "false"}, nil
		default:
			return Token{}, io.ErrUnexpectedEOF
		}
	}
	return Token{}, nil
}
func (l *Lexer) lexNull() (Token, error) {
	// consume the leading 'n'
	_, err := l.next()
	if err != nil {
		return Token{}, err
	}
	for i := range "ull" {
		b, err := l.next()
		if err != nil {
			return Token{}, fmt.Errorf("unexpected end of input in null")
		}
		if b != "ull"[i] {
			return Token{}, fmt.Errorf("invalid null value")
		}
	}
	return Token{Type: TokenNull, Value: "null"}, nil
}
func GenTokens(r io.Reader) ([]Token, error) {
	lexer := NewLexer(r)
	var tokens []Token
	for {
		token, err := lexer.NextToken()
		if err != nil {
			return nil, err
		}
		fmt.Printf(
			"Token=%v", token,
		)
		if token.Type == TokenEOF {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
