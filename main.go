package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	"todo-app/storage"
	"todo-app/todo"
	"todo-app/todostore"

	"github.com/google/uuid"
)

type contextKey string

const traceIDKey contextKey = "traceID"

func startCLI(view bool, add, find, updateStatus, updateDesc, remove string) {
	traceID := uuid.New().String()
	ctx := context.WithValue(context.Background(), traceIDKey, traceID)
	fs := &storage.FileStore{Path: "todos.json"}

	switch {
	case view:
		todostore.GetAll(ctx, fs)
	case add != "":
		slog.InfoContext(ctx, "Creating todo", "desc", add, "traceID", traceID)
		if err := todostore.Add(ctx, add, fs); err != nil {
			slog.ErrorContext(ctx, "failed to add item", "traceID", traceID, "error", err)
		}
	case remove != "":
		slog.InfoContext(ctx, "Deleting todo", "desc", remove, "traceID", traceID)
		if err := todostore.Remove(ctx, remove, fs); err != nil {
			slog.ErrorContext(ctx, "failed to remove item", "traceID", traceID, "error", err)
		}
	case find != "" && updateStatus != "":
		slog.InfoContext(ctx, "Updating todo", "desc", updateStatus)
		if err := todostore.Update(ctx, find, todo.UpdateFieldDescription, updateStatus, fs); err != nil {
			slog.ErrorContext(ctx, "failed to update item", "traceID", traceID, "error", err)
		}
	case find != "" && updateDesc != "":
		slog.InfoContext(ctx, "Updating todo", "desc", updateStatus)
		if err := todostore.Update(ctx, find, todo.UpdateFieldDescription, updateDesc, fs); err != nil {
			slog.ErrorContext(ctx, "failed to update item", "traceID", traceID, "error", err)
		}
	default:
		slog.InfoContext(ctx, "No CLI action specified")
	}
}

func startServer() {
	fs := &storage.FileStore{Path: "todos.json"}
	app := &App{FS: fs}

	mux := http.NewServeMux()
	mux.HandleFunc("/create", app.CreateHandler)
	mux.HandleFunc("/read", app.ReadHandler)
	mux.HandleFunc("/update", app.UpdateHandler)
	mux.HandleFunc("/delete", app.DeleteHandler)

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", TraceMiddleware(mux)); err != nil {
		slog.Error("Server failed", "error", err)
	}
}

func main() {
	modeFlag := flag.String("mode", "cli", "Choose mode: cli or server")
	viewFlag := flag.Bool("view", false, "View to-do list")
	addFlag := flag.String("add", "", "Add a new to-do item")
	findFlag := flag.String("find", "", "Find a to-do item by description")
	updateStatusFlag := flag.String("update-status", "", "Update a to-do item status")
	updateDescFlag := flag.String("update-description", "", "Update a to-do item description")
	removeFlag := flag.String("remove", "", "Remove a to-do item")

	flag.Parse()

	if *modeFlag == "server" {
		startServer()
	} else {
		startCLI(*viewFlag, *addFlag, *findFlag, *updateStatusFlag, *updateDescFlag, *removeFlag)
	}
}
