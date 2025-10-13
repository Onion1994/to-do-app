package todo

import (
	"fmt"
	"strings"
)

func PrintTodos(todos []Item) {
	for _, element := range todos {
		fmt.Printf("%s: %s\n", element.Description, element.Status)
	}
}

func AddNewItem(todos []Item, desc string) []Item {
	lowerCaseDesc := strings.ToLower(desc)
	for _, item := range todos {
		if item.Description == lowerCaseDesc {
			return todos
		}
	}
	return append(todos, Item{Description: lowerCaseDesc, Status: NotStarted})
}

func RemoveItem(todos []Item, desc string) []Item {
	var updatedTodos []Item
	for _, element := range todos {
		if element.Description != strings.ToLower(desc) {
			updatedTodos = append(updatedTodos, element)
		}
	}
	return updatedTodos
}

func UpdateStatus(todos []Item, desc string, status string) {
	if IsValidStatus(status) {
		for i := range todos {
			if todos[i].Description == strings.ToLower(desc) {
				todos[i].Status = strings.ToLower(status)
			}
		}
	}
}

func UpdateDesc(todos []Item, oldDesc string, newDesc string) {
	for i := range todos {
		if todos[i].Description == strings.ToLower(oldDesc) {
			todos[i].Description = strings.ToLower(newDesc)
		}
	}
}
