package main

import (
	"context"
	"flag"
	"log/slog"

	"github.com/google/uuid"
	"todo-app/storage"
	"todo-app/todo"
)

type contextKey string

const (
	traceIDKey contextKey = "traceID"
)

func main() {
	traceID := uuid.New().String()
	ctx := context.WithValue(context.Background(), traceIDKey, traceID)

	logger := slog.Default().With("traceID", traceID)
	slog.SetDefault(logger)

	slog.InfoContext(ctx, "Application started")
	defer func() {
		slog.InfoContext(ctx, "Application exited")
	}()

	todos, err := storage.LoadTodos(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load todos", "error", err)
		return
	}

	viewFlag := flag.Bool("view", false, "View to-do list")
	addFlag := flag.String("add", "", "Add a new to-do item")
	findFlag := flag.String("find", "", "Find to-do item by description")
	updateStatusFlag := flag.String("update-status", "", "Update a to-do item status")
	updateDescriptionFlag := flag.String("update-description", "", "Update a to-do item description")
	removeFlag := flag.String("remove", "", "Remove a to-do item")

	flag.Parse()

	switch {
	case *viewFlag:
		todo.PrintTodos(todos)
	case *addFlag != "":
		todos = todo.AddNewItem(todos, *addFlag)
		if err := storage.SaveTodos(ctx, todos); err != nil {
			slog.ErrorContext(ctx, "Failed to save after add", "error", err)
		}
	case *removeFlag != "":
		todos = todo.RemoveItem(todos, *removeFlag)
		if err := storage.SaveTodos(ctx, todos); err != nil {
			slog.ErrorContext(ctx, "Failed to save after remove", "error", err)
		}
	case *findFlag != "" && *updateStatusFlag != "":
		todo.UpdateStatus(todos, *findFlag, *updateStatusFlag)
		if err := storage.SaveTodos(ctx, todos); err != nil {
			slog.ErrorContext(ctx, "Failed to save after status update", "error", err)
		}
	case *findFlag != "" && *updateDescriptionFlag != "":
		todo.UpdateDesc(todos, *findFlag, *updateDescriptionFlag)
		if err := storage.SaveTodos(ctx, todos); err != nil {
			slog.ErrorContext(ctx, "Failed to save after description update", "error", err)
		}
	}
}
