package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"todo-app/storage"
	"todo-app/todo"
)

func newTestApp(t *testing.T) *App {
	tmpFile := "test_todos.json"
	t.Cleanup(func() { os.Remove(tmpFile) })
	fs := &storage.FileStore{Path: tmpFile}
	app := &App{FS: fs}
	return app
}

func TestCreateAndReadFlow(t *testing.T) {
	app := newTestApp(t)

	// Create
	body, _ := json.Marshal(todo.Item{Description: "wash car"})
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), traceIDKey, "test-trace-1"))
	w := httptest.NewRecorder()

	app.CreateHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	// Read
	req2 := httptest.NewRequest(http.MethodGet, "/read", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), traceIDKey, "test-trace-2"))
	w2 := httptest.NewRecorder()

	app.ReadHandler(w2, req2)
	resp2 := w2.Result()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp2.StatusCode)
	}

	var readResp TodosResponse
	if err := json.NewDecoder(resp2.Body).Decode(&readResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(readResp.Todos) == 0 || readResp.Todos[0].Description != "wash car" {
		t.Errorf("expected todo 'wash car', got %+v", readResp.Todos)
	}
}

func TestDeleteHandler(t *testing.T) {
	app := newTestApp(t)

	// Add first
	item := todo.Item{Description: "take bins out"}
	body, _ := json.Marshal(item)
	createReq := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	createReq = createReq.WithContext(context.WithValue(createReq.Context(), traceIDKey, "trace-del"))
	createW := httptest.NewRecorder()
	app.CreateHandler(createW, createReq)

	// Now delete
	delReq := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewReader(body))
	delReq = delReq.WithContext(context.WithValue(delReq.Context(), traceIDKey, "trace-del"))
	delW := httptest.NewRecorder()

	app.DeleteHandler(delW, delReq)
	resp := delW.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	// Verify deletion
	readReq := httptest.NewRequest(http.MethodGet, "/read", nil)
	readReq = readReq.WithContext(context.WithValue(readReq.Context(), traceIDKey, "trace-del-read"))
	readW := httptest.NewRecorder()
	app.ReadHandler(readW, readReq)

	var readResp TodosResponse
	if err := json.NewDecoder(readW.Result().Body).Decode(&readResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(readResp.Todos) != 0 {
		t.Errorf("expected 0 todos after deletion, got %+v", readResp.Todos)
	}
}

func TestUpdateHandler(t *testing.T) {
	app := newTestApp(t)

	// Add first
	createItem := todo.Item{Description: "take bins out"}
	createBody, _ := json.Marshal(createItem)
	createReq := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(createBody))
	createReq = createReq.WithContext(context.WithValue(createReq.Context(), traceIDKey, "trace-update"))
	createW := httptest.NewRecorder()
	app.CreateHandler(createW, createReq)

	// Now update
	updateItem := UpdateRequest{Description: "take bins out", Field: todo.UpdateFieldStatus, NewValue: todo.Completed}
	updateBody, _ := json.Marshal(updateItem)
	updateReq := httptest.NewRequest(http.MethodPatch, "/update", bytes.NewReader(updateBody))
	updateReq = updateReq.WithContext(context.WithValue(updateReq.Context(), traceIDKey, "trace-update"))
	updateW := httptest.NewRecorder()

	app.UpdateHandler(updateW, updateReq)
	resp := updateW.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	// Verify update
	readReq := httptest.NewRequest(http.MethodGet, "/read", nil)
	readReq = readReq.WithContext(context.WithValue(readReq.Context(), traceIDKey, "trace-del-update"))
	readW := httptest.NewRecorder()
	app.ReadHandler(readW, readReq)

	var readResp TodosResponse
	if err := json.NewDecoder(readW.Result().Body).Decode(&readResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	found := false
	for _, element := range readResp.Todos {
		if element.Description == "take bins out" {
			found = true
			if element.Status != todo.Completed {
				t.Errorf("expected todo to be %s, got %s", todo.Completed, element.Status)
			}
		}
	}
	if !found {
		t.Errorf("todo 'take bins out' not found after update")
	}
}
