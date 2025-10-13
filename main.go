package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	"todo-app/todo"
	"todo-app/todostore"

	"github.com/google/uuid"
)

type contextKey string

const traceIDKey contextKey = "traceID"

func startCLI(ctx context.Context, view bool, add, find, updateStatus, updateDesc, remove string) {
	switch {
	case view:
		todostore.GetAll(ctx)
	case add != "":
		todostore.Add(ctx, add)
	case remove != "":
		todostore.Remove(ctx, remove)
	case find != "" && updateStatus != "":
		todostore.Update(ctx, find, todo.UpdateFieldStatus, updateStatus)
	case find != "" && updateDesc != "":
		todostore.Update(ctx, find, todo.UpdateFieldDescription, updateDesc)
	default:
		slog.InfoContext(ctx, "No CLI action specified")
	}
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", CreateHandler)
	// TODO: add /get, /update, /delete handlers

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", TraceMiddleware(mux)); err != nil {
		slog.Error("Server failed", "error", err)
	}
}

func main() {
	traceID := uuid.New().String()
	ctx := context.WithValue(context.Background(), traceIDKey, traceID)
	logger := slog.Default().With("traceID", traceID)
	slog.SetDefault(logger)
	slog.InfoContext(ctx, "Application started")
	defer slog.InfoContext(ctx, "Application exited")

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
		startCLI(ctx, *viewFlag, *addFlag, *findFlag, *updateStatusFlag, *updateDescFlag, *removeFlag)
	}
}
