package storage

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"todo-app/todo"
)

type loadRequest struct {
	ctx      context.Context
	response chan loadResponse
}

type loadResponse struct {
	todos []todo.Item
	err   error
}

type saveRequest struct {
	ctx      context.Context
	todos    []todo.Item
	response chan error
}

type FileStore struct {
	Path    string
	loadCh  chan loadRequest
	saveCh  chan saveRequest
	closeCh chan struct{}
}

func NewFileStore(path string) *FileStore {
	fs := &FileStore{
		Path:    path,
		loadCh:  make(chan loadRequest),
		saveCh:  make(chan saveRequest),
		closeCh: make(chan struct{}),
	}
	go fs.actor()
	return fs
}

func (fs *FileStore) actor() {
	for {
		select {
		case req := <-fs.loadCh:
			todos, err := fs.loadFromDisk(req.ctx)
			req.response <- loadResponse{todos: todos, err: err}

		case req := <-fs.saveCh:
			err := fs.saveToDisk(req.ctx, req.todos)
			req.response <- err

		case <-fs.closeCh:
			return
		}
	}
}

func (fs *FileStore) LoadTodos(ctx context.Context) ([]todo.Item, error) {
	respCh := make(chan loadResponse, 1)
	fs.loadCh <- loadRequest{ctx: ctx, response: respCh}
	resp := <-respCh
	return resp.todos, resp.err
}

func (fs *FileStore) SaveTodos(ctx context.Context, todos []todo.Item) error {
	respCh := make(chan error, 1)
	fs.saveCh <- saveRequest{ctx: ctx, todos: todos, response: respCh}
	return <-respCh
}

func (fs *FileStore) Close() {
	close(fs.closeCh)
}

func (fs *FileStore) loadFromDisk(ctx context.Context) ([]todo.Item, error) {
	file, err := os.Open(fs.Path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.InfoContext(ctx, "Todo file does not exist, starting with empty list")
			return []todo.Item{}, nil
		}
		slog.ErrorContext(ctx, "Failed to open todo file", "error", err)
		return nil, err
	}
	defer file.Close()

	var todos []todo.Item
	if err := json.NewDecoder(file).Decode(&todos); err != nil {
		slog.ErrorContext(ctx, "Failed to decode todos", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Loaded todos from disk", "count", len(todos))
	return todos, nil
}

func (fs *FileStore) saveToDisk(ctx context.Context, todos []todo.Item) error {
	file, err := os.Create(fs.Path)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create or open todo file", "error", err)
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(todos); err != nil {
		slog.ErrorContext(ctx, "Failed to encode todos", "error", err)
		return err
	}

	slog.InfoContext(ctx, "Saved todos to disk", "count", len(todos))
	return nil
}
