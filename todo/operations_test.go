package todo

import (
	"testing"
)

func TestAddNewItem(t *testing.T) {
	tests := []struct {
		name        string
		input       []Item
		desc        string
		wantDesc    string
		wantStatus  string
		wantErr     bool
		expectedLen int
	}{
		{"normal add", []Item{}, "test", "test", NotStarted, false, 1},
		{"normalises to lowercase", []Item{}, "TEst", "test", NotStarted, false, 1},
		{"duplicate same case", []Item{{Description: "test", Status: NotStarted}}, "test", "test", NotStarted, true, 1},
		{"duplicate different case", []Item{{Description: "test", Status: NotStarted}}, "TEst", "test", NotStarted, true, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddNewItem(tt.input, tt.desc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}

			if len(got) != tt.expectedLen {
				t.Errorf("expected length %d, got %d", tt.expectedLen, len(got))
			}
			if len(got) > 0 && got[0].Description != tt.wantDesc {
				t.Errorf("expected description %q, got %q", tt.wantDesc, got[0].Description)
			}
			if len(got) > 0 && got[0].Status != tt.wantStatus {
				t.Errorf("expected status %q, got %q", tt.wantStatus, got[0].Status)
			}
		})
	}
}

func TestRemoveItem(t *testing.T) {
	tests := []struct {
		name        string
		input       []Item
		desc        string
		wantErr     bool
		expectedLen int
	}{
		{"remove existing", []Item{{Description: "test"}}, "test", false, 0},
		{"remove ingores case", []Item{{Description: "test"}}, "TEst", false, 0},
		{"remove absent", []Item{{Description: "test1", Status: NotStarted}}, "test2", true, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveItem(tt.input, tt.desc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}

			if len(got) != tt.expectedLen {
				t.Errorf("expected length %d, got %d", tt.expectedLen, len(got))
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	tests := []struct {
		name       string
		todos      []Item
		targetDesc string
		newStatus  string
		wantStatus string
		wantErr    bool
	}{
		{"valid update", []Item{{Description: "test", Status: NotStarted}}, "test", Started, Started, false},
		{"case-insensitive match", []Item{{Description: "test", Status: NotStarted}}, "TEST", Started, Started, false},
		{"invalid status", []Item{{Description: "test", Status: NotStarted}}, "test", "invalid status", NotStarted, true},
		{"absent item", []Item{{"test", NotStarted}}, "nope", Completed, NotStarted, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateStatus(tt.todos, tt.targetDesc, tt.newStatus)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}

			if tt.todos[0].Status != tt.wantStatus {
				t.Errorf("expected status %q, got %q", tt.wantStatus, tt.todos[0].Status)
			}
		})
	}
}

func TestUpdateDesc(t *testing.T) {
	tests := []struct {
		name       string
		todos      []Item
		targetDesc string
		newDesc    string
		wantDesc   string
		wantErr    bool
	}{
		{"valid update", []Item{{"test1", NotStarted}}, "test1", "test2", "test2", false},
		{"case-insensitive update", []Item{{"test1", NotStarted}}, "TeSt1", "TEst2", "test2", false},
		{"absent item", []Item{{"test1", NotStarted}}, "test2", "test3", "test1", true},
		{"duplicate new desc", []Item{{"test1", NotStarted}, {"test2", NotStarted}}, "test1", "test2", "test1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateDesc(tt.todos, tt.targetDesc, tt.newDesc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}

			if tt.todos[0].Description != tt.wantDesc {
				t.Errorf("expected description %q, got %q", tt.wantDesc, tt.todos[0].Description)
			}
		})
	}
}
