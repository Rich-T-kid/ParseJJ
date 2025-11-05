package main

import (
	"fmt"
	"json-parser/parser"
	"os"
	"sync"
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

	//r := strings.NewReader(`{"name":"Bob","age":30,"active":true,"address":null}`)
	wg := sync.WaitGroup{}
	wg.Add(4)
	albumChan := make(chan any, 100)
	postChan := make(chan any, 100)
	userChan := make(chan any, 100)
	todoChan := make(chan any, 100)
	go consumeData(_exampleAlbumsPath, albumChan, &wg)
	go consumeData(_examplePostsPath, postChan, &wg)
	go consumeData(_exampleUsersPath, userChan, &wg)
	go consumeData(_exampleTodosPath, todoChan, &wg)

	for {
		select {
		case album, ok := <-albumChan:
			if !ok {
				albumChan = nil
			} else {
				fmt.Println("Album ->", album)
			}
		case post, ok := <-postChan:
			if !ok {
				postChan = nil
			} else {
				fmt.Println("Post ->", post)
			}
		case user, ok := <-userChan:
			if !ok {
				userChan = nil
			} else {
				fmt.Println("User ->", user)
			}
		case todo, ok := <-todoChan:
			if !ok {
				todoChan = nil
			} else {
				fmt.Println("Todo ->", todo)
			}
		}
		if albumChan == nil && postChan == nil && userChan == nil && todoChan == nil {
			break
		}
	}
	wg.Wait()
	fmt.Println("================ Done!! ================")
}

func consumeData(fileName string, value chan any, wg *sync.WaitGroup) any {
	s, err := newSource(fileName)
	if err != nil {
		return err
	}
	v, ok := parser.BasicParase(s.F).([]any)
	if !ok {
		return nil
	}
	for _, item := range v {
		value <- item
	}
	close(value)
	wg.Done()
	return nil
}
