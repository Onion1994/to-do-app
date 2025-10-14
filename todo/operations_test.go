package todo

import (
	"testing"
)

func TestAddNewItem(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos = AddNewItem(todos, "test")

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
	todos = AddNewItem(todos, "TeSt")

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
	todos = AddNewItem(todos, "test")
	todos = AddNewItem(todos, "test")
	todos = AddNewItem(todos, "TEST")

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
	todos = AddNewItem(todos, "test1")
	todos = AddNewItem(todos, "test2")
	todos = RemoveItem(todos, "test1")

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
	todos = AddNewItem(todos, "test1")
	todos = AddNewItem(todos, "test2")
	todos = RemoveItem(todos, "TeSt1")

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test2" {
		t.Errorf("todos do not match: %+v", todos)
	}
}

func TestRemovingAbsentDoesNothing(t *testing.T) {
	// Arrange
	var todos []Item

	// Act
	todos = AddNewItem(todos, "test1")
	todos = RemoveItem(todos, "test2")

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
	todos = AddNewItem(todos, "test")
	UpdateStatus(todos, "test", Completed)

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
	todos = AddNewItem(todos, "test")
	UpdateStatus(todos, "TeSt", "COmpleTED")

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
	todos = AddNewItem(todos, "test")
	UpdateStatus(todos, "test", "banana")

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
	todos = AddNewItem(todos, "test1")
	UpdateDesc(todos, "test1", "test2")

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
	todos = AddNewItem(todos, "test1")
	UpdateDesc(todos, "TeSt1", "TEst2")

	// Assert
	expected := 1
	actual := len(todos)

	if actual != expected ||
		todos[0].Description != "test2" ||
		todos[0].Status != NotStarted {
		t.Errorf("todos do not match: %+v", todos)
	}
}
