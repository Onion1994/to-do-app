package todo

import "strings"

const (
	NotStarted = "not started"
	Started    = "started"
	Completed  = "completed"
)

type Item struct {
	Description string
	Status      string
}

func IsValidStatus(s string) bool {
	switch strings.ToLower(s) {
	case NotStarted, Started, Completed:
		return true
	}
	return false
}
