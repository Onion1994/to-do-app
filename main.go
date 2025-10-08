package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const todoFile = "todos.json"

type Status string

const (
	NotStarted Status = "not started"
	Started    Status = "started"
	Completed  Status = "completed"
)

type TodoItem struct {
	Description string
	Status      Status
}

func PrintTodos(todos []TodoItem) {
	for _, element := range todos {
		fmt.Printf("%s: %s\n", element.Description, element.Status)
	}
}

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

func main() {
	todos, _ := LoadTodos()

	addFlag := flag.String("add", "", "Add a new to-do item")
	findFlag := flag.String("find", "", "Find to-do item by description")
	updateFlag := flag.String("update", "", "Update a to-do item status")
	removeFlag := flag.String("remove", "", "Remove a to-do item")

	flag.Parse()

	if *addFlag != "" {
		task := TodoItem{
			Description: *addFlag,
			Status:      NotStarted,
		}

		todos = append(todos, task)
		SaveTodos(todos)
	}

	if *removeFlag != "" {
		var updatedTodos []TodoItem
		for _, element := range todos {
			if element.Description != *removeFlag {
				updatedTodos = append(updatedTodos, element)
			}
		}

		todos = updatedTodos
		SaveTodos(todos)
	}

	if *findFlag != "" && *updateFlag != "" {
		for i := range todos {
			if todos[i].Description == *findFlag {
				todos[i].Status = Status(*updateFlag)
			}
		}

		SaveTodos(todos)
	}
}
