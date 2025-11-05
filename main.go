package main

import (
	"fmt"
	"json-parser/parser"
	"os"
	"strings"
)

const basePath = "test_data"
const _exampleAlbumsPath = basePath + "/example_albums.json"
const _examplePostsPath = basePath + "/example_posts.json"
const _exampleUsersPath = basePath + "/example_users.json"
const _exampleTodosPath = basePath + "/example_todos.json"

type source struct {
	F *os.File
}

func newSource(filePath string) (*source, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &source{F: f}, nil
}
func main() {

	r := strings.NewReader(`{"name":"Bob","age":30,"active":true,"address":null}`)
	tokens, _ := parser.GenTokens(r)
	for i, tok := range tokens {
		fmt.Printf("pos: %d tokens Type: %v Value: %v\n", i, tok.Type, tok.Value)
	}

}
