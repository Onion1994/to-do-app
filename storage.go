package main

import (
	"encoding/json"
	"os"
	"log/slog"
)

var todoFile = "todos.json"

func LoadTodos() ([]TodoItem, error) {
	file, err := os.Open(todoFile)

	if err != nil {
		if os.IsNotExist(err) {
			slog.Info("Todo file does not exist, starting with empty list")
			return []TodoItem{}, nil
		}
		slog.Error("Failed to open todo file", "error", err)
		return nil, err
	}

	defer file.Close()

	var todos []TodoItem
	err = json.NewDecoder(file).Decode(&todos)
	if err != nil {
		slog.Error("Failed to decode todos", "error", err)
		return nil, err
	}

	slog.Info("Loaded todos from disk", "count", len(todos))
	return todos, nil
}

func SaveTodos(todos []TodoItem) error {
	file, err := os.Create(todoFile)
	if err != nil {
		slog.Error("Failed to encode todos", "error", err)
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(todos)
	if err != nil {
		slog.Error("Failed to encode todos", "error", err)
		return err
	}
	
	slog.Info("Saved todos to disk", "count", len(todos))
	return nil
}
