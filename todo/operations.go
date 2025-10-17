package todo

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrItemExists    = errors.New("item already exists")
	ErrItemNotFound  = errors.New("item not found")
	ErrInvalidStatus = errors.New("invalid status")
	ErrDuplicateDesc = errors.New("new description already exists")
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
			return todos, fmt.Errorf("%w: %s", ErrItemExists, desc)
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
		return todos, fmt.Errorf("%w: %s", ErrItemNotFound, desc)
	}

	return updatedTodos, nil
}

func UpdateStatus(todos []Item, desc, status string) error {
	if !IsValidStatus(status) {
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}

	for i := range todos {
		if todos[i].Description == strings.ToLower(desc) {
			todos[i].Status = strings.ToLower(status)
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrItemNotFound, desc)
}

func UpdateDesc(todos []Item, oldDesc string, newDesc string) error {
	found := false
	for i := range todos {
		if todos[i].Description == strings.ToLower(newDesc) {
			return fmt.Errorf("%w: %s", ErrDuplicateDesc, newDesc)
		}
	}

	for i := range todos {
		if todos[i].Description == strings.ToLower(oldDesc) {
			found = true
			todos[i].Description = strings.ToLower(newDesc)
		}
	}

	if found {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrItemNotFound, oldDesc)
}
