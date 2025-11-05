package parser

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func runParser(input string) any {
	return BasicParase(bytes.NewBufferString(input))
}

// ------------------------------
// Tests
// ------------------------------
func TestEmptyInput(t *testing.T) {
	var example = []byte(`[]`)
	result := BasicParase(bytes.NewBuffer(example))

	// assert it's a slice
	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}

	if len(arr) != 0 {
		t.Errorf("expected empty array, got length %d", len(arr))
	}
}

func TestTwoFieldsStringsOnly(t *testing.T) {
	var example = []byte(`{"name":"Alice","age":30}`)
	result := BasicParase(bytes.NewBuffer(example))

	obj, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}

	fmt.Println("res ->", obj)

	if len(obj) != 2 {
		t.Errorf("expected length 2, got %d", len(obj))
	}
	if obj["name"] != "Alice" {
		t.Errorf("expected name to be Alice, got %v", obj["name"])
	}
	if obj["age"] != float64(30) {
		t.Errorf("expected age to be 30, got %v", obj["age"])
	}
}

func TestMultipleStringFields(t *testing.T) {
	var example = []byte(`{"name":"Alice","city":"New York","country":"USA","occupation":"Engineer"}`)
	result := BasicParase(bytes.NewBuffer(example))

	obj, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}

	fmt.Println("res ->", obj)

	expected := map[string]string{
		"name":       "Alice",
		"city":       "New York",
		"country":    "USA",
		"occupation": "Engineer",
	}

	if len(obj) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(obj))
	}

	for k, v := range expected {
		val, exists := obj[k]
		if !exists {
			t.Errorf("missing key %s in result", k)
			continue
		}
		if val != v {
			t.Errorf("expected %s to be %s, got %v", k, v, val)
		}
	}
}

// ---------- GROUP 1: BASIC TESTS ----------
func TestBasicJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
	}{
		{
			name:  "Flat object",
			input: `{"name":"Alice","age":30}`,
			expected: map[string]any{
				"name": "Alice",
				"age":  30.0,
			},
		},
		{
			name:  "Nested object",
			input: `{"user":{"name":"Bob","city":"London"}}`,
			expected: map[string]any{
				"user": map[string]any{
					"name": "Bob",
					"city": "London",
				},
			},
		},
		{
			name:  "Array in object",
			input: `{"numbers":[1,2,3]}`,
			expected: map[string]any{
				"numbers": []any{1.0, 2.0, 3.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runParser(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("mismatch for %s:\nexpected %#v\ngot      %#v", tt.name, tt.expected, result)
			}
		})
	}
}

// ---------- GROUP 2: NESTED / ARRAYS ----------
func TestNestedJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:  "Array of objects",
			input: `[{"id":1,"item":"book"},{"id":2,"item":"pen"}]`,
			expected: []any{
				map[string]any{"id": 1.0, "item": "book"},
				map[string]any{"id": 2.0, "item": "pen"},
			},
		},
		{
			name:  "Deep nesting",
			input: `{"a":{"b":{"c":{"d":{"e":"final"}}}}}`,
			expected: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": map[string]any{
							"d": map[string]any{
								"e": "final",
							},
						},
					},
				},
			},
		},
		{
			name:  "Nested arrays",
			input: `[[1,2],[3,4],[5,[6,7]]]`,
			expected: []any{
				[]any{1.0, 2.0},
				[]any{3.0, 4.0},
				[]any{5.0, []any{6.0, 7.0}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BasicParase(bytes.NewBufferString(tt.input))
			if tt.name == "Nested arrays" {
				fmt.Printf("result: %#v\n", result)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("mismatch for %s:\nexpected %#v\ngot      %#v", tt.name, tt.expected, result)
			}
		})
	}
}

// ---------- GROUP 3: COMPLEX / REALISTIC ----------
func TestComplexJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
	}{
		{
			name: "Complex data array",
			input: `{"data":[
				{"type":"fruit","items":["apple","banana"]},
				{"type":"drink","items":["water","juice"]}
			]}`,
			expected: map[string]any{
				"data": []any{
					map[string]any{"type": "fruit", "items": []any{"apple", "banana"}},
					map[string]any{"type": "drink", "items": []any{"water", "juice"}},
				},
			},
		},
		{
			name: "Full payload",
			input: `{
				"users": [
					{"id":1,"name":"John","roles":["admin","dev"]},
					{"id":2,"name":"Jane","roles":["user"],"meta":{"verified":true,"points":42}}
				],
				"active":true,
				"stats":{"total":2,"online":1}
			}`,
			expected: map[string]any{
				"users": []any{
					map[string]any{"id": 1.0, "name": "John", "roles": []any{"admin", "dev"}},
					map[string]any{"id": 2.0, "name": "Jane", "roles": []any{"user"}, "meta": map[string]any{"verified": true, "points": 42.0}},
				},
				"active": true,
				"stats":  map[string]any{"total": 2.0, "online": 1.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runParser(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("mismatch for %s:\nexpected %#v\ngot      %#v", tt.name, tt.expected, result)
			}
		})
	}
}
