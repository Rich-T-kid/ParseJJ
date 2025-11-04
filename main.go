package main

import (
	"os"
	"time"
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
	s, err := newSource(_exampleAlbumsPath)
	if err != nil {
		panic(err)
	}
	defer s.F.Close()
	for {
		buf := make([]byte, 100)
		n, err := s.F.Read(buf)
		if err != nil {
			break
		}
		println(string(buf[:n]))
		time.Sleep(time.Second * 2)
	}

}
