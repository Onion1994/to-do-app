package main

import (
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
	todos := []TodoItem{
		{Description: "Test1", Status: NotStarted},
		{Description: "Test2", Status: Completed},
	}

	// Act
	if err := SaveTodos(todos); err != nil {
		t.Fatalf("SaveTodos failed: %v", err)
	}

	loaded, err := LoadTodos()
	if err != nil {
		t.Fatalf("LoadTodos failed: %v", err)
	}

	// Assert
	if len(loaded) != 2 ||
		loaded[0].Description != "Test1" ||
		loaded[0].Status != NotStarted ||
		loaded[1].Description != "Test2" ||
		loaded[1].Status != Completed {

		t.Errorf("Loaded todos do not match: %+v", loaded)
	}
}
