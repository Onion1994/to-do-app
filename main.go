package main

import (
	"flag"
	"fmt"
)

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

func main() {
	var todos []TodoItem

	addFlag := flag.String("add", "", "Add a new to-do item")
	descFlag := flag.String("desc", "", "Find to-do item by description")
	updateFlag := flag.String("update", "", "Update a to-do item status")
	deleteFlag := flag.String("delete", "", "Delete a to-do item")

	flag.Parse()

	if *addFlag != "" {
		task := TodoItem{
			Description: *addFlag,
			Status:      NotStarted,
		}

		todos = append(todos, task)
		PrintTodos(todos)
	}

	if *deleteFlag != "" {
		var updatedTodos []TodoItem
		for _, element := range todos {
			if element.Description != *deleteFlag {
				updatedTodos = append(updatedTodos, element)
			}
		}

		todos = updatedTodos
		PrintTodos(todos)
	}

	if *descFlag != "" && *updateFlag != "" {
		for i := range todos {
			if todos[i].Description == *descFlag {
				todos[i].Status = Status(*updateFlag)
			}
		}

		PrintTodos(todos)
	}

}
