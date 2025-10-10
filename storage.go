package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
)

var todoFile = "todos.json"

func LoadTodos(ctx context.Context) ([]TodoItem, error) {
	file, err := os.Open(todoFile)

	if err != nil {
		if os.IsNotExist(err) {
			slog.InfoContext(ctx, "Todo file does not exist, starting with empty list")
			return []TodoItem{}, nil
		}
		slog.ErrorContext(ctx, "Failed to open todo file", "error", err)
		return nil, err
	}

	defer file.Close()

	var todos []TodoItem
	err = json.NewDecoder(file).Decode(&todos)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to decode todos", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Loaded todos from disk", "count", len(todos))
	return todos, nil
}

func SaveTodos(ctx context.Context, todos []TodoItem) error {
	file, err := os.Create(todoFile)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to encode todos", "error", err)
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(todos)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to encode todos", "error", err)
		return err
	}

	slog.InfoContext(ctx, "Saved todos to disk", "count", len(todos))
	return nil
}
