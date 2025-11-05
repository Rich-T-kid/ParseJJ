package parser

import (
	"bytes"
	"testing"
)

func TestEmptyInput(t *testing.T) {
	var example = []byte(`[]`)
	result := basicParase(bytes.NewBuffer(example))
	if len(result) != 0 {
		t.Errorf("expected length 0, got %d", len(result))
	}
}

/*
func TestTwoFieldsStringsOnly(t *testing.T) {
	var example = []byte(`{"name":"Alice","age":"30"}`)
	result := basicParase(bytes.NewBuffer(example))
	fmt.Println("res ->", result)
	if len(result) != 2 {
		t.Errorf("expected length 2, got %d", len(result))
	}
	if result["name"] != "Alice" {
		t.Errorf("expected name to be Alice, got %v", result["name"])
	}
	if result["age"] != "30" {
		t.Errorf("expected age to be 30, got %v", result["age"])
	}
}
func TestMultipleStringFields(t *testing.T) {
	var example = []byte(`{"name":"Alice","city":"New York","country":"USA","occupation":"Engineer"}`)
	result := basicParase(bytes.NewBuffer(example))
	fmt.Println("res ->", result)

	expected := map[string]string{
		"name":       "Alice",
		"city":       "New York",
		"country":    "USA",
		"occupation": "Engineer",
	}

	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("expected %s to be %s, got %v", k, v, result[k])
		}
	}
}
*/
