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

func AddNewItem(todos []Item, desc string) ([]Item, error) {
	lowerCaseDesc := strings.ToLower(desc)
	for _, item := range todos {
		if item.Description == lowerCaseDesc {
			return todos, fmt.Errorf("item already exists")
		}
	}

	return append(todos, Item{Description: lowerCaseDesc, Status: NotStarted}), nil
}

func RemoveItem(todos []Item, desc string) ([]Item, error) {
	var updatedTodos []Item
	for _, element := range todos {
		if element.Description != strings.ToLower(desc) {
			updatedTodos = append(updatedTodos, element)
		}
	}

	if len(updatedTodos) == len(todos) {
		return todos, fmt.Errorf("item does not exist")
	}

	return updatedTodos, nil
}

func UpdateStatus(todos []Item, desc string, status string) error {
	if IsValidStatus(status) {
		found := false
		for i := range todos {
			if todos[i].Description == strings.ToLower(desc) {
				found = true
				todos[i].Status = strings.ToLower(status)
			}
		}

		if found {
			return nil
		}

		return fmt.Errorf("item not found")
	}

	return fmt.Errorf("invalid status '%s'. Valid statuses are: %s, %s, %s", status, NotStarted, Started, Completed)
}

func UpdateDesc(todos []Item, oldDesc string, newDesc string) error {
	found := false
	for i := range todos {
		if todos[i].Description == strings.ToLower(oldDesc) {
			found = true
			todos[i].Description = strings.ToLower(newDesc)
		}
	}

	if found {
		return nil
	}

	return fmt.Errorf("item not found")
}
