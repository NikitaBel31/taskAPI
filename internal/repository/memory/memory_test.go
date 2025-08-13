package memory_test

import (
	"context"
	"taskapi/internal/domain"
	"taskapi/internal/repository"
	"taskapi/internal/repository/memory"
	"testing"
)

func TestRepo_CreateAndGetByID(t *testing.T) {
	tests := []struct {
		name      string
		task      domain.Task
		wantFound bool
		wantErr   bool
		searchID  string
	}{
		{
			name:      "create and find task",
			task:      domain.Task{ID: "1", Title: "Test Task", Status: domain.StatusTodo},
			wantFound: true,
			searchID:  "1",
		},
		{
			name:      "search non-existing task",
			task:      domain.Task{ID: "2", Title: "Another Task", Status: domain.StatusDone},
			wantFound: false,
			searchID:  "not-exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := memory.New()

			_, err := repo.Create(ctx, tt.task)
			if err != nil {
				t.Fatalf("unexpected error on Create: %v", err)
			}

			got, ok, err := repo.GetByID(ctx, tt.searchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
			if ok != tt.wantFound {
				t.Errorf("expected found=%v, got %v", tt.wantFound, ok)
			}
			if ok && got.ID != tt.task.ID {
				t.Errorf("expected ID %q, got %q", tt.task.ID, got.ID)
			}
		})
	}
}

func TestRepo_List(t *testing.T) {
	ctx := context.Background()
	repo := memory.New()

	tasks := []domain.Task{
		{ID: "1", Title: "Todo Task", Status: domain.StatusTodo},
		{ID: "2", Title: "In Progress Task", Status: domain.StatusInProgress},
		{ID: "3", Title: "Done Task", Status: domain.StatusDone},
	}

	for _, task := range tasks {
		if _, err := repo.Create(ctx, task); err != nil {
			t.Fatalf("failed to create task: %v", err)
		}
	}

	tests := []struct {
		name      string
		filter    repository.Filter
		wantCount int
	}{
		{
			name:      "list all tasks",
			filter:    repository.Filter{},
			wantCount: 3,
		},
		{
			name:      "filter by StatusTodo",
			filter:    repository.Filter{Status: ptrStatus(domain.StatusTodo)},
			wantCount: 1,
		},
		{
			name:      "filter by StatusDone",
			filter:    repository.Filter{Status: ptrStatus(domain.StatusDone)},
			wantCount: 1,
		},
		{
			name:      "filter with no matches",
			filter:    repository.Filter{Status: ptrStatus(domain.Status("unknown"))},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(ctx, tt.filter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("expected %d tasks, got %d", tt.wantCount, len(got))
			}
		})
	}
}

func ptrStatus(s domain.Status) *domain.Status {
	return &s
}
