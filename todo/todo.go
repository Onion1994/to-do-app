package todo

import "strings"

type Item struct {
    Description string
    Status      string
}

type UpdateField string

const (
    UpdateFieldDescription UpdateField = "description"
    UpdateFieldStatus      UpdateField = "status"
)

const (
	NotStarted = "not started"
	Started    = "started"
	Completed  = "completed"
)

func IsValidStatus(s string) bool {
	switch strings.ToLower(s) {
	case NotStarted, Started, Completed:
		return true
	}
	return false
}