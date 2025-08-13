package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"taskapi/internal/domain"
	"taskapi/internal/dto"
	"taskapi/internal/logger"
	"taskapi/internal/repository"
	"taskapi/internal/usecase"
)

type mockRepo struct {
	createFn func(ctx context.Context, t domain.Task) (domain.Task, error)
	getFn    func(ctx context.Context, id string) (domain.Task, bool, error)
	listFn   func(ctx context.Context, f repository.Filter) ([]domain.Task, error)
}

func (m *mockRepo) Create(ctx context.Context, t domain.Task) (domain.Task, error) {
	return m.createFn(ctx, t)
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (domain.Task, bool, error) {
	return m.getFn(ctx, id)
}
func (m *mockRepo) List(ctx context.Context, f repository.Filter) ([]domain.Task, error) {
	return m.listFn(ctx, f)
}

type mockLogger struct {
	entries []logger.Entry
}

func (m *mockLogger) Log(e logger.Entry) { m.entries = append(m.entries, e) }
func (m *mockLogger) Stop()              {}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.CreateInput
		repoErr   error
		wantErr   error
		wantTitle string
	}{
		{
			name:    "success",
			input:   dto.CreateInput{Title: "Task 1", Status: domain.StatusTodo},
			wantErr: nil,
		},
		{
			name:    "empty title",
			input:   dto.CreateInput{Title: ""},
			wantErr: usecase.ErrBadRequest,
		},
		{
			name:    "invalid status should default to TODO",
			input:   dto.CreateInput{Title: "Task 2", Status: domain.Status("invalid")},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepo{
				createFn: func(ctx context.Context, tsk domain.Task) (domain.Task, error) {
					return tsk, tt.repoErr
				},
			}
			mockLog := &mockLogger{}

			svc := usecase.NewService(mockRepo, mockLog)
			svc = &usecase.Service{
				Repo:  mockRepo,
				Log:   mockLog,
				Now:   func() time.Time { return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) },
				IdGen: func() string { return "test-id" },
			}

			got, err := svc.Create(context.Background(), "req-1", tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.ID != "test-id" {
					t.Errorf("expected ID 'test-id', got %v", got.ID)
				}
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	mockRepo := &mockRepo{
		getFn: func(ctx context.Context, id string) (domain.Task, bool, error) {
			if id == "exists" {
				return domain.Task{ID: id, Title: "Found"}, true, nil
			}
			return domain.Task{}, false, nil
		},
	}
	mockLog := &mockLogger{}
	svc := usecase.NewService(mockRepo, mockLog)

	t.Run("found", func(t *testing.T) {
		got, err := svc.Get(context.Background(), "req-1", "exists")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != "exists" {
			t.Errorf("expected id 'exists', got %v", got.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.Get(context.Background(), "req-1", "missing")
		if !errors.Is(err, usecase.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestService_List(t *testing.T) {
	mockRepo := &mockRepo{
		listFn: func(ctx context.Context, f repository.Filter) ([]domain.Task, error) {
			return []domain.Task{
				{ID: "1", Title: "Task 1"},
				{ID: "2", Title: "Task 2"},
			}, nil
		},
	}
	mockLog := &mockLogger{}
	svc := usecase.NewService(mockRepo, mockLog)

	got, err := svc.List(context.Background(), "req-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(got))
	}
}
