package main

import (
	"context"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	tmpFile := "test_todos.json"
	origFile := todoFile
	todoFile = tmpFile

	exitCode := m.Run()

	todoFile = origFile
	os.Remove(tmpFile)

	os.Exit(exitCode)
}

func TestSaveAndLoadTodos(t *testing.T) {
	// Arrange
	ctx := context.Background()
	todos := []TodoItem{
		{Description: "test1", Status: NotStarted},
		{Description: "test2", Status: Completed},
	}

	// Act
	if err := SaveTodos(ctx, todos); err != nil {
		t.Fatalf("SaveTodos failed: %v", err)
	}

	loaded, err := LoadTodos(ctx)
	if err != nil {
		t.Fatalf("LoadTodos failed: %v", err)
	}

	// Assert
	if len(loaded) != 2 ||
		loaded[0].Description != "test1" ||
		loaded[0].Status != NotStarted ||
		loaded[1].Description != "test2" ||
		loaded[1].Status != Completed {

		t.Errorf("Loaded todos do not match: %+v", loaded)
	}
}
