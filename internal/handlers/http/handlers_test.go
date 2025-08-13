package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"taskapi/internal/domain"
	"taskapi/internal/dto"
	httpHandler "taskapi/internal/handlers/http"
	"taskapi/internal/usecase"
)

type mockTaskService struct {
	createFn func(ctx context.Context, reqID string, in dto.CreateInput) (domain.Task, error)
	getFn    func(ctx context.Context, reqID, id string) (domain.Task, error)
	listFn   func(ctx context.Context, reqID string, status *domain.Status) ([]domain.Task, error)
}

func (m *mockTaskService) Create(ctx context.Context, reqID string, in dto.CreateInput) (domain.Task, error) {
	return m.createFn(ctx, reqID, in)
}
func (m *mockTaskService) Get(ctx context.Context, reqID, id string) (domain.Task, error) {
	return m.getFn(ctx, reqID, id)
}
func (m *mockTaskService) List(ctx context.Context, reqID string, status *domain.Status) ([]domain.Task, error) {
	return m.listFn(ctx, reqID, status)
}

func TestRouter_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		serviceRes domain.Task
		serviceErr error
		wantCode   int
	}{
		{"success", "123", domain.Task{ID: "123", Title: "Test"}, nil, http.StatusOK},
		{"not found", "456", domain.Task{}, usecase.ErrNotFound, http.StatusNotFound},
		{"internal error", "789", domain.Task{}, errors.New("db error"), http.StatusInternalServerError},
		{"missing id", "", domain.Task{}, nil, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockTaskService{
				getFn: func(ctx context.Context, reqID, id string) (domain.Task, error) {
					return tt.serviceRes, tt.serviceErr
				},
			}

			rt := httpHandler.NewRouter(svc, nil)

			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.id, nil)
			rr := httptest.NewRecorder()
			rt.Get(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected code %d, got %d", tt.wantCode, rr.Code)
			}
		})
	}
}

func TestRouter_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		serviceRes domain.Task
		serviceErr error
		wantCode   int
	}{
		{"success", dto.CreateInput{Title: "Task"}, domain.Task{ID: "1"}, nil, http.StatusCreated},
		{"bad json", "not a json", domain.Task{}, nil, http.StatusBadRequest},
		{"bad request", dto.CreateInput{}, domain.Task{}, usecase.ErrBadRequest, http.StatusBadRequest},
		{"internal error", dto.CreateInput{Title: "Task"}, domain.Task{}, errors.New("db error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockTaskService{
				createFn: func(ctx context.Context, reqID string, in dto.CreateInput) (domain.Task, error) {
					return tt.serviceRes, tt.serviceErr
				},
			}

			rt := httpHandler.NewRouter(svc, nil)

			var buf bytes.Buffer
			switch v := tt.body.(type) {
			case string:
				buf.WriteString(v)
			default:
				_ = json.NewEncoder(&buf).Encode(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/tasks", &buf)
			rr := httptest.NewRecorder()
			rt.Create(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected code %d, got %d", tt.wantCode, rr.Code)
			}
		})
	}
}

func TestRouter_GetList(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		serviceRes []domain.Task
		serviceErr error
		wantCode   int
	}{
		{"success", "", []domain.Task{{ID: "1"}}, nil, http.StatusOK},
		{"invalid status", "?status=wrong", nil, nil, http.StatusBadRequest},
		{"internal error", "", nil, errors.New("db error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockTaskService{
				listFn: func(ctx context.Context, reqID string, status *domain.Status) ([]domain.Task, error) {
					return tt.serviceRes, tt.serviceErr
				},
			}

			rt := httpHandler.NewRouter(svc, nil)

			req := httptest.NewRequest(http.MethodGet, "/tasks"+tt.query, nil)
			rr := httptest.NewRecorder()
			rt.GetList(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("expected code %d, got %d", tt.wantCode, rr.Code)
			}
		})
	}
}
