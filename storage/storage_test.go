package storage

import (
	"context"
	"os"
	"testing"

	"todo-app/todo"
)

func TestSaveAndLoadTodos(t *testing.T) {
	// Arrange
	ctx := context.Background()
	tmpFile := "test_todos.json"
	fs := NewFileStore(tmpFile)
	defer fs.Close()

	todos := []todo.Item{
		{Description: "test1", Status: todo.NotStarted},
		{Description: "test2", Status: todo.Completed},
	}

	// Act
	if err := fs.SaveTodos(ctx, todos); err != nil {
		t.Fatalf("SaveTodos failed: %v", err)
	}

	loaded, err := fs.LoadTodos(ctx)
	if err != nil {
		t.Fatalf("LoadTodos failed: %v", err)
	}

	// Assert
	if len(loaded) != 2 ||
		loaded[0].Description != "test1" ||
		loaded[0].Status != todo.NotStarted ||
		loaded[1].Description != "test2" ||
		loaded[1].Status != todo.Completed {

		t.Errorf("Loaded todos do not match: %+v", loaded)
	}

	os.Remove(tmpFile)
}
