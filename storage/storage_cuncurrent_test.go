package storage

import (
	"context"
	"sync"
	"testing"

	"todo-app/todo"
)

func TestConcurrentReads(t *testing.T) {
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

