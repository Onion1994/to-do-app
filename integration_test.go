package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"todo-app/storage"
	"todo-app/todo"
)

func startTestServer(t *testing.T) (baseURL string, cleanup func()) {
	t.Helper()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "todos.json")
	if err := os.WriteFile(tmpFile, []byte("[]"), 0644); err != nil {
		t.Fatalf("failed to initialize temp file: %v", err)
	}

	fs := &storage.FileStore{Path: tmpFile}
	app := &App{FS: fs}

	mux := http.NewServeMux()
	mux.HandleFunc("/create", app.CreateHandler)
	mux.HandleFunc("/read", app.ReadHandler)
	mux.HandleFunc("/update", app.UpdateHandler)
	mux.HandleFunc("/delete", app.DeleteHandler)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	server := &http.Server{Handler: TraceMiddleware(mux)}
	go server.Serve(listener)

	cleanup = func() {
		server.Shutdown(context.Background())
		os.Remove(tmpFile)
	}

	return "http://" + listener.Addr().String(), cleanup
}

func doRequest(t *testing.T, client *http.Client, method string, url string, body any) *http.Response {
	t.Helper()

	var buf io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		buf = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func TestCreateAndReadFlow(t *testing.T) {
	baseURL, cleanup := startTestServer(t)
	defer cleanup()

	client := &http.Client{Timeout: time.Second}

	tests := []struct {
		name       string
		item       todo.Item
		wantStatus int
	}{
		{"valid item", todo.Item{Description: "wash car"}, http.StatusCreated},
		{"empty description", todo.Item{Description: ""}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, client, http.MethodPost, baseURL+"/create", tt.item)
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Create status = %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}

	resp := doRequest(t, client, http.MethodGet, baseURL+"/read", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Read status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var readResp TodosResponse
	if err := json.NewDecoder(resp.Body).Decode(&readResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(readResp.Todos) != 1 {
		t.Errorf("expected %d todos, got %d", 1, len(readResp.Todos))
	}
}

func TestUpdateHandlerIntegration(t *testing.T) {
	baseURL, cleanup := startTestServer(t)
	defer cleanup()
	client := &http.Client{Timeout: time.Second}

	// create item
	item := todo.Item{Description: "go to gym"}
	resp := doRequest(t, client, http.MethodPost, baseURL+"/create", item)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create status = %v, want %v", resp.StatusCode, http.StatusCreated)
	}

	tests := []struct {
		name       string
		updateReq  UpdateRequest
		wantStatus int
		checkFn    func(t *testing.T, todos []todo.Item)
	}{
		{
			"update status",
			UpdateRequest{Description: "go to gym", Field: todo.UpdateFieldStatus, NewValue: todo.Completed},
			http.StatusOK,
			func(t *testing.T, todos []todo.Item) {
				if todos[0].Status != todo.Completed {
					t.Errorf("expected status completed, got %s", todos[0].Status)
				}
			},
		},
		{
			"update description",
			UpdateRequest{Description: "go to gym", Field: todo.UpdateFieldDescription, NewValue: "go to park"},
			http.StatusOK,
			func(t *testing.T, todos []todo.Item) {
				if todos[0].Description != "go to park" {
					t.Errorf("expected desc go to park, got %s", todos[0].Description)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := doRequest(t, client, http.MethodPatch, baseURL+"/update", tt.updateReq)
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Update status = %v, want %v", resp.StatusCode, tt.wantStatus)
			}

			if resp.StatusCode == http.StatusOK {
				readResp := TodosResponse{}
				resp2 := doRequest(t, client, http.MethodGet, baseURL+"/read", nil)
				defer resp2.Body.Close()
				if err := json.NewDecoder(resp2.Body).Decode(&readResp); err != nil {
					t.Fatalf("failed to decode read response: %v", err)
				}
				tt.checkFn(t, readResp.Todos)
			}
		})
	}
}

func TestDeleteHandlerIntegration(t *testing.T) {
	baseURL, cleanup := startTestServer(t)
	defer cleanup()
	client := &http.Client{Timeout: time.Second}

	item := todo.Item{Description: "take bins out"}
	resp := doRequest(t, client, http.MethodPost, baseURL+"/create", item)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create status = %v, want %v", resp.StatusCode, http.StatusCreated)
	}

	resp = doRequest(t, client, http.MethodDelete, baseURL+"/delete", item)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Delete status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	resp = doRequest(t, client, http.MethodGet, baseURL+"/read", nil)
	defer resp.Body.Close()
	var readResp TodosResponse
	if err := json.NewDecoder(resp.Body).Decode(&readResp); err != nil {
		t.Fatalf("failed to decode read response: %v", err)
	}
	if len(readResp.Todos) != 0 {
		t.Errorf("expected 0 todos, got %d", len(readResp.Todos))
	}
}
