package todo

import (
	"testing"
)

func TestAddNewItem(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}

	// Assert
	if len(todos) != 1 ||
		todos[0].Description != "test" ||
		todos[0].Status != "not started" {

		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestNewItemsAreNormalisedToLowerCase(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "teSt")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}

	// Assert
	if len(todos) != 1 ||
		todos[0].Description != "test" ||
		todos[0].Status != NotStarted {

		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestItemsCannotBeDuplicated(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}
	todos, err = AddNewItem(todos, "test")
	if err == nil {
		t.Fatal("expected error for duplicate item, got nil")
	}
	todos, err = AddNewItem(todos, "TEST")
	if err == nil {
		t.Fatal("expected error for duplicate item, got nil")
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected {
		t.Errorf("todos should be %+v but are %+v", expected, actual)
	}
}

func TestRemoveItem(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test1")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}
	todos, err = AddNewItem(todos, "test2")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}
	todos, err = RemoveItem(todos, "test1")
	if err != nil {
		t.Fatalf("RemoveItem failed: %v", err)
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test2" {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestRemoveItemIgnoresCase(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test1")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}
	todos, err = AddNewItem(todos, "test2")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}
	todos, err = RemoveItem(todos, "tEsT1")
	if err != nil {
		t.Fatalf("RemoveItem failed: %v", err)
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test2" {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestRemovingAbsentItem(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test1")
	if err != nil {
		t.Fatalf("AddNewItem failed: %v", err)
	}
	todos, err = RemoveItem(todos, "test2")
	if err == nil {
		t.Fatal("expected error for removing absent item, got nil")
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test1" {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestUpdateStatus(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test")
	if err != nil {
		t.Errorf("AddNewItem failed: %v", err)
	}

	if err = UpdateStatus(todos, "test", Completed); err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test" ||
		todos[0].Status != Completed {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestUpdateStatusIgnoresCase(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test")
	if err != nil {
		t.Errorf("AddNewItem failed: %v", err)
	}

	if err = UpdateStatus(todos, "TeSt", "COmpleTED"); err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test" ||
		todos[0].Status != Completed {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestUpdateStatusOnlyIfValidStatus(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test")
	if err != nil {
		t.Errorf("AddNewItem failed: %v", err)
	}

	if err = UpdateStatus(todos, "test", "banana"); err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test" ||
		todos[0].Status != NotStarted {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestUpdateDesc(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test1")
	if err != nil {
		t.Errorf("AddNewItem failed: %v", err)
	}

	if err = UpdateDesc(todos, "test1", "test2"); err != nil {
		t.Fatalf("UpdateDesc failed: %v", err)
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test2" ||
		todos[0].Status != NotStarted {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestUpdateDescIgnoresCase(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test1")
	if err != nil {
		t.Errorf("AddNewItem failed: %v", err)
	}

	if err = UpdateDesc(todos, "TeSt1", "TEst2"); err != nil {
		t.Fatalf("UpdateDesc failed: %v", err)
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test2" ||
		todos[0].Status != NotStarted {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestUpdateDescAbsentItem(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test1")
	if err != nil {
		t.Errorf("AddNewItem failed: %v", err)
	}

	if err = UpdateDesc(todos, "test2", "test3"); err == nil {
		t.Fatal("expected error for updating absent item, got nil")
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test1" ||
		todos[0].Status != NotStarted {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestUpdateStatusAbsentItem(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos, err := AddNewItem(todos, "test1")
	if err != nil {
		t.Errorf("AddNewItem failed: %v", err)
	}

	if err = UpdateStatus(todos, "test2", "completed"); err == nil {
		t.Fatal("expected error for updating absent item, got nil")
	}

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test1" ||
		todos[0].Status != NotStarted {
		t.Errorf("todos do not match: %+v", todos)
	}
}
