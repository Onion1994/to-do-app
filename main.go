package main

import (
	"flag"
)

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
