package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"todo-app/storage"
	"todo-app/todo"
	"todo-app/todostore"
)

type App struct {
	FS *storage.FileStore
}

type TodosResponse struct {
	TraceID string
	Todos   []todo.Item
}

type UpdateRequest struct {
	Description string
	Field       todo.UpdateField
	NewValue    string
}

func (a *App) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := ctx.Value(traceIDKey).(string)

	var item todo.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		slog.ErrorContext(ctx, "failed to decode request", "traceID", traceID, "error", err)
		return
	}

	slog.InfoContext(ctx, "Creating todo", "desc", item.Description, "traceID", traceID)

	if err := todostore.Add(ctx, item.Description, a.FS); err != nil {
		if errors.Is(err, todo.ErrItemExists) ||
			errors.Is(err, todo.ErrItemIsEmpty) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   err.Error(),
				"traceID": traceID,
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "failed to create item",
				"traceID": traceID,
			})
		}
		slog.ErrorContext(ctx, "failed to create item", "traceID", traceID, "error", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo created",
		"traceID": traceID,
	})
}

func (a *App) ReadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := ctx.Value(traceIDKey).(string)

	todos, err := a.FS.LoadTodos(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"traceID": traceID})
		slog.ErrorContext(ctx, "failed to fetch todo items", "traceID", traceID, "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TodosResponse{
		TraceID: traceID,
		Todos:   todos,
	})
}

func (a *App) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := ctx.Value(traceIDKey).(string)

	var request UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		slog.ErrorContext(ctx, "failed to decode request", "traceID", traceID, "error", err)
		return
	}

	slog.InfoContext(ctx, "Updating todo", "desc", request.Description)
	if err := todostore.Update(ctx, request.Description, request.Field, request.NewValue, a.FS); err != nil {
		if errors.Is(err, todostore.ErrInvalidUpdateField) ||
			errors.Is(err, todo.ErrDuplicateDesc) ||
			errors.Is(err, todo.ErrInvalidStatus) ||
			errors.Is(err, todo.ErrItemNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   err.Error(),
				"traceID": traceID,
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "failed to update item",
				"traceID": traceID,
			})
		}
		slog.ErrorContext(ctx, "failed to update item", "traceID", traceID, "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo updated",
		"traceID": traceID,
	})
}

func (a *App) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := ctx.Value(traceIDKey).(string)

	var item todo.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		slog.ErrorContext(ctx, "failed to decode request", "traceID", traceID, "error", err)
		return
	}

	slog.InfoContext(ctx, "Deleting todo", "desc", item.Description, "traceID", traceID)

	if err := todostore.Remove(ctx, item.Description, a.FS); err != nil {
		if errors.Is(err, todo.ErrItemNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "item not found",
				"traceID": traceID,
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "failed to delete item",
				"traceID": traceID,
			})
		}
		slog.ErrorContext(ctx, "failed to delete item", "traceID", traceID, "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Todo deleted",
		"traceID": traceID,
	})
}

func (a *App) ListPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := ctx.Value(traceIDKey).(string)

	todos, err := a.FS.LoadTodos(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to load todos", "traceID", traceID, "error", err)
		http.Error(w, "Failed to load todos", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/list.html")
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse template", "traceID", traceID, "error", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, todos); err != nil {
		slog.ErrorContext(ctx, "failed to execute template", "traceID", traceID, "error", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
