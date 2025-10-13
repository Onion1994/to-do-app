package todostore

import (
    "context"
	"fmt"
    "todo-app/storage"
    "todo-app/todo"
)

func GetAll(ctx context.Context,) error {
	todos, err := storage.LoadTodos(ctx)
	if err != nil {
		return err
	}

	todo.PrintTodos(todos)

	return nil
}

func Add(ctx context.Context, desc string) error {
	todos, err := storage.LoadTodos(ctx)
	if err != nil {
		return err
	}

	todos = todo.AddNewItem(todos, desc)

	if err := storage.SaveTodos(ctx, todos); err != nil {
		return err
	}

	return nil
}

func Remove(ctx context.Context, desc string) error {
	todos, err := storage.LoadTodos(ctx)
	if err != nil {
		return err
	}

	todos = todo.RemoveItem(todos, desc)

	if err := storage.SaveTodos(ctx, todos); err != nil {
		return err
	}

	return nil
}

func Update(ctx context.Context, desc string, field todo.UpdateField, newValue string) error {
	todos, err := storage.LoadTodos(ctx)
	if err != nil {
		return err
	}

	switch field {
	case todo.UpdateFieldDescription:
		todo.UpdateDesc(todos, desc, newValue)
	case todo.UpdateFieldStatus:
		todo.UpdateStatus(todos, desc, newValue)
	default:
		return fmt.Errorf("invalid update field: %s", field)
	}

	if err := storage.SaveTodos(ctx, todos); err != nil {
		return err
	}

	return nil
}