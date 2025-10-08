package main

import "fmt"

func PrintTodos(todos []TodoItem) {
	for _, element := range todos {
		fmt.Printf("%s: %s\n", element.Description, element.Status)
	}
}
