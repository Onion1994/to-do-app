package storage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"todo-app/todo"
)

func TestConcurrentReads(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/concurrent_reads.json"
	fs := NewFileStore(tmpFile)
	defer fs.Close()

	initial := []todo.Item{
		{Description: "task 1", Status: todo.NotStarted},
		{Description: "task 2", Status: todo.Started},
	}
	if err := fs.SaveTodos(context.Background(), initial); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	const numReaders = 100
	var wg sync.WaitGroup
	wg.Add(numReaders)

	errors := make(chan error, numReaders)
	readCounts := make(chan int, numReaders)

	for range numReaders {
		go func() {
			defer wg.Done()

			todos, err := fs.LoadTodos(context.Background())
			if err != nil {
				errors <- err
				return
			}

			readCounts <- len(todos)
		}()
	}

	wg.Wait()
	close(errors)
	close(readCounts)

	for err := range errors {
		t.Errorf("concurrent read failed: %v", err)
	}

	for count := range readCounts {
		if count != 2 {
			t.Errorf("expected 2 todos, got %d", count)
		}
	}
}

func TestConcurrentReadAndWrite(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/concurrent_rw.json"
	fs := NewFileStore(tmpFile)
	defer fs.Close()

	initial := []todo.Item{{Description: "initial", Status: todo.NotStarted}}
	if err := fs.SaveTodos(context.Background(), initial); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	const numOperations = 100
	var wg sync.WaitGroup
	wg.Add(numOperations * 2)

	errors := make(chan error, numOperations*2)

	for range numOperations {
		go func() {
			defer wg.Done()
			_, err := fs.LoadTodos(context.Background())
			if err != nil {
				errors <- err
			}
		}()
	}

	for i := range numOperations {
		go func() {
			defer wg.Done()
			todos := []todo.Item{{Description: fmt.Sprintf("task %d", i), Status: todo.NotStarted}}
			if err := fs.SaveTodos(context.Background(), todos); err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("operation failed: %v", err)
	}

	todos, err := fs.LoadTodos(context.Background())
	if err != nil {
		t.Fatalf("final load failed: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("expected 1 todo after all writes, got %d", len(todos))
	}
}

func TestActorShutdown(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/shutdown_test.json"
	fs := NewFileStore(tmpFile)

	todos := []todo.Item{{Description: "test", Status: todo.NotStarted}}
	if err := fs.SaveTodos(context.Background(), todos); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	done := make(chan struct{})
	go func() {
		fs.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Close() hung - actor didn't shutdown")
	}
}
