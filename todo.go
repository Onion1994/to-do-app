package main

type Status string

const (
	NotStarted Status = "not started"
	Started    Status = "started"
	Completed  Status = "completed"
)

type TodoItem struct {
	Description string
	Status      Status
}
