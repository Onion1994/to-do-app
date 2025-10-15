package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"todo-app/storage"
	"todo-app/todo"
	"todo-app/todostore"
)

type TodosResponse struct {
	TraceID string
	Todos   []todo.Item
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := ctx.Value(traceIDKey).(string)

	var item todo.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		slog.ErrorContext(ctx, "Failed to decode request", "error", err)
		return
	}

	slog.InfoContext(ctx, "Creating todo", "desc", item.Description)

	if err := todostore.Add(ctx, item.Description); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"traceID": traceID,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo created",
		"traceID": traceID,
	})
}

func ReadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := ctx.Value(traceIDKey).(string)

	todos, err := storage.LoadTodos(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"traceID": traceID,
		})
		slog.ErrorContext(ctx, "Failed to fetch todo items", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TodosResponse{
		TraceID: traceID,
		Todos:   todos,
	})
}
