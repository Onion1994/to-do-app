package main

import "fmt"

func PrintTodos(todos []TodoItem) {
	for _, element := range todos {
		fmt.Printf("%s: %s\n", element.Description, element.Status)
	}
}

func AddNewItem(todos []TodoItem, desc string) []TodoItem {
    return append(todos, TodoItem{Description: desc, Status: NotStarted})
}

func RemoveItem(todos []TodoItem, desc string) []TodoItem {
	var updatedTodos []TodoItem
		for _, element := range todos {
			if element.Description != desc {
				updatedTodos = append(updatedTodos, element)
			}
		}
	return updatedTodos
}

func UpdateStatus(todos []TodoItem, desc string, status Status) {
	for i := range todos {
			if todos[i].Description == desc {
				todos[i].Status = Status(status)
			}
		}
}

func UpdateDesc(todos []TodoItem, oldDesc string, newDesc string) {
	for i := range todos {
			if todos[i].Description == oldDesc {
				todos[i].Description = newDesc
			}
		}
}