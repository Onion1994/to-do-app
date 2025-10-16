package todostore

import (
	"context"
	"fmt"
	"todo-app/storage"
	"todo-app/todo"
)

func GetAll(ctx context.Context, fs *storage.FileStore) error {
	todos, err := fs.LoadTodos(ctx)
	if err != nil {
		return err
	}

	todo.PrintTodos(todos)
	return nil
}

func Add(ctx context.Context, desc string, fs *storage.FileStore) error {
	todos, err := fs.LoadTodos(ctx)
	if err != nil {
		return err
	}

	todos, err = todo.AddNewItem(todos, desc)
	if err != nil {
		return err
	}

	if err := fs.SaveTodos(ctx, todos); err != nil {
		return err
	}

	return nil
}

func Remove(ctx context.Context, desc string, fs *storage.FileStore) error {
	todos, err := fs.LoadTodos(ctx)
	if err != nil {
		return err
	}

	todos, err = todo.RemoveItem(todos, desc)
	if err != nil {
		return err
	}

	if err := fs.SaveTodos(ctx, todos); err != nil {
		return err
	}

	return nil
}

func Update(ctx context.Context, desc string, field todo.UpdateField, newValue string, fs *storage.FileStore) error {
	todos, err := fs.LoadTodos(ctx)
	if err != nil {
		return err
	}

	switch field {
	case todo.UpdateFieldDescription:
		if err := todo.UpdateDesc(todos, desc, newValue); err != nil {
			return err
		}
	case todo.UpdateFieldStatus:
		if err := todo.UpdateStatus(todos, desc, newValue); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid update field: %s. Valid fields are: %s, %s", field, todo.UpdateFieldDescription, todo.UpdateFieldStatus)
	}

	if err := fs.SaveTodos(ctx, todos); err != nil {
		return err
	}

	return nil
}
