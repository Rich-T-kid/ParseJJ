package parser

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestLexStringBasic(t *testing.T) {
	var example = []byte(`"Hello, World!"`)
	lexer := NewLexer(bytes.NewReader(example))
	token, err := lexer.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fmt.Printf("token -> %v\n", token)
	if token.Type != TokenString {
		t.Errorf("expected token type %v, got %v", TokenString, token.Type)
	}
	if token.Value != "Hello, World!" {
		t.Errorf("expected token value 'Hello, World!', got '%s'", token.Value)
	}
}
func TestLexNumberBasic(t *testing.T) {
	var example = []byte(`12345`)
	lexer := NewLexer(bytes.NewReader(example))
	token, err := lexer.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fmt.Printf("token -> %v\n", token)
	if token.Type != TokenNumber {
		t.Errorf("expected token type %v, got %v", TokenNumber, token.Type)
	}
	if token.Value != "12345" {
		t.Errorf("expected token value '12345', got '%s'", token.Value)
	}
}

func TestLexBoolBasic(t *testing.T) {
	var example = []byte(`true`)
	lexer := NewLexer(bytes.NewReader(example))
	token, err := lexer.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fmt.Printf("token -> %v\n", token)
	if token.Type != TokenTrue {
		t.Errorf("expected token type %v, got %v", TokenTrue, token.Type)
	}
	if token.Value != "true" {
		t.Errorf("expected token value 'true', got '%s'", token.Value)
	}

	var exampleFalse = []byte(`false`)
	lexerFalse := NewLexer(bytes.NewReader(exampleFalse))
	tokenFalse, err := lexerFalse.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fmt.Printf("token -> %v\n", tokenFalse)
	if tokenFalse.Type != TokenFalse {
		t.Errorf("expected token type %v, got %v", TokenFalse, tokenFalse.Type)
	}
	if tokenFalse.Value != "false" {
		t.Errorf("expected token value 'false', got '%s'", tokenFalse.Value)
	}
	var invalidBool = []byte(`trux`)
	lexerInvalid := NewLexer(bytes.NewReader(invalidBool))
	_, err = lexerInvalid.NextToken()
	fmt.Printf("This should fail ==============\n")
	fmt.Printf("Error: -> %v", err)
	if err == nil {
		t.Fatalf("expected error for invalid boolean, got nil")
	}
}

func TestLexTokenSymbols(t *testing.T) {
	cases := []struct {
		name  string
		input string
		typ   TokenType
		val   string
	}{
		{"left brace", "{", TokenLeftBrace, "{"},
		{"right brace", "}", TokenRightBrace, "}"},
		{"left bracket", "[", TokenLeftBracket, "["},
		{"right bracket", "]", TokenRightBracket, "]"},
		{"comma", ",", TokenComma, ","},
		{"colon", ":", TokenColon, ":"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(bytes.NewReader([]byte(tc.input)))
			tok, err := l.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != tc.typ {
				t.Fatalf("type mismatch: want %v got %v", tc.typ, tok.Type)
			}
			if tok.Value != tc.val {
				t.Fatalf("value mismatch: want %q got %q", tc.val, tok.Value)
			}
		})
	}
}

func TestLexNumberVariants(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"integer", "42", "42"},
		{"negative", "-7", "-7"},
		{"float", "3.14", "3.14"},
		{"small float", "0.001", "0.001"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(bytes.NewReader([]byte(tc.input)))
			tok, err := l.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != TokenNumber {
				t.Fatalf("expected TokenNumber, got %v", tok.Type)
			}
			if tok.Value != tc.want {
				t.Fatalf("number mismatch: want %q got %q", tc.want, tok.Value)
			}
		})
	}
}

func TestLexStringEdgeCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", `""`, ""},
		{"hello world", `"hello world"`, "hello world"},
		{"unterminated", `"incomplete`, "incomplete"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tc.input))
			tok, err := l.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != TokenString {
				t.Fatalf("expected TokenString, got %v", tok.Type)
			}
			if tok.Value != tc.want {
				t.Fatalf("string mismatch: want %q got %q", tc.want, tok.Value)
			}
		})
	}
}

func TestLexNullAndUnexpected(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		l := NewLexer(strings.NewReader("null"))
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tok.Type != TokenNull {
			t.Fatalf("expected TokenNull, got %v", tok.Type)
		}
		if tok.Value != "null" {
			t.Fatalf("null value mismatch: want %q got %q", "null", tok.Value)
		}
	})

	t.Run("unexpected char", func(t *testing.T) {
		l := NewLexer(strings.NewReader("@"))
		_, err := l.NextToken()
		if err == nil {
			t.Fatalf("expected error for unexpected char, got nil")
		}
		if !strings.Contains(err.Error(), "unexpected character") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}

func TestGenTokens_SimpleObject(t *testing.T) {
	input := `{"name":"Alice","city":"Paris"}`
	expected := []Token{
		{Type: TokenLeftBrace, Value: "{"},
		{Type: TokenString, Value: "name"},
		{Type: TokenColon, Value: ":"},
		{Type: TokenString, Value: "Alice"},
		{Type: TokenComma, Value: ","},
		{Type: TokenString, Value: "city"},
		{Type: TokenColon, Value: ":"},
		{Type: TokenString, Value: "Paris"},
		{Type: TokenRightBrace, Value: "}"},
	}

	r := strings.NewReader(input)
	tokens, err := GenTokens(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, tok := range tokens {
		exp := expected[i]
		if tok.Type != exp.Type || tok.Value != exp.Value {
			t.Errorf("token %d mismatch: got {Type:%v, Value:%q}, expected {Type:%v, Value:%q}",
				i, tok.Type, tok.Value, exp.Type, exp.Value)
		}
	}
}
func TestGenTokens_NestedObject(t *testing.T) {
	input := `{"person":{"name":"Alice","city":"Paris"}}`

	expected := []Token{
		{Type: TokenLeftBrace, Value: "{"},
		{Type: TokenString, Value: "person"},
		{Type: TokenColon, Value: ":"},
		{Type: TokenLeftBrace, Value: "{"},
		{Type: TokenString, Value: "name"},
		{Type: TokenColon, Value: ":"},
		{Type: TokenString, Value: "Alice"},
		{Type: TokenComma, Value: ","},
		{Type: TokenString, Value: "city"},
		{Type: TokenColon, Value: ":"},
		{Type: TokenString, Value: "Paris"},
		{Type: TokenRightBrace, Value: "}"},
		{Type: TokenRightBrace, Value: "}"},
	}

	r := strings.NewReader(input)
	tokens, err := GenTokens(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, tok := range tokens {
		exp := expected[i]
		if tok.Type != exp.Type || tok.Value != exp.Value {
			t.Errorf("token %d mismatch:\n got      {Type:%v, Value:%q}\n expected {Type:%v, Value:%q}",
				i, tok.Type, tok.Value, exp.Type, exp.Value)
		}
	}
}

func TestGenTokens_Primitives(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Array of numbers", `[1,2,3,4]`},
		{"Array of strings", `["a","b","c"]`},
		{"Object with array value", `{"nums":[1,2,3]}`},
		{"Array of objects", `[{"id":1},{"id":2}]`},
		{"Nested arrays", `[["a","b"],["c","d"]]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			tokens, err := GenTokens(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) == 0 {
				t.Fatalf("no tokens returned for %q", tt.input)
			}

			// Print token stream for quick inspection
			for i, tok := range tokens {
				fmt.Printf("[%s] %d: Type=%v, Value=%q\n", tt.name, i, tok.Type, tok.Value)
			}
		})
	}
}

func TestGenTokens_Nested(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Mixed types", `{"name":"Bob","age":30,"active":true,"address":null}`},
		{"Whitespace handling", `{
			"a" : "b" ,
			"c" : 123
		}`},
		{"Deep nesting", `{"a":{"b":{"c":{"d":1}}}}`},
		{"Negative and float numbers", `{"temp":-12.5,"offset":0.001}`},
		{"Boolean and null in array", `[true,false,null]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			tokens, err := GenTokens(r)
			if err != nil {
				t.Fatalf("test: %s ,unexpected error: %v", tt.name, err)
			}
			if len(tokens) == 0 {
				t.Fatalf("no tokens returned for %q", tt.input)
			}

			// Print token stream for visual validation
			for i, tok := range tokens {
				fmt.Printf("[%s] %d: Type=%v, Value=%q\n", tt.name, i, tok.Type, tok.Value)
			}
		})
	}
}
