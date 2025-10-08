package main

import (
	"encoding/json"
	"os"
)

var todoFile = "todos.json"

func LoadTodos() ([]TodoItem, error) {
	file, err := os.Open(todoFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []TodoItem{}, nil
		}
		return nil, err
	}
	defer file.Close()
	var todos []TodoItem
	err = json.NewDecoder(file).Decode(&todos)
	return todos, err
}

func SaveTodos(todos []TodoItem) error {
	file, err := os.Create(todoFile)
	if err != nil {
		return err
	}
	defer file.Close()
	PrintTodos(todos)
	return json.NewEncoder(file).Encode(todos)
}
