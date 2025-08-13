package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"taskapi/internal/dto"
	"taskapi/internal/usecase/mapper"
	"taskapi/internal/usecase/validation"
	"time"

	"taskapi/internal/domain"
	"taskapi/internal/logger"
	"taskapi/internal/repository"
)

var (
	ErrNotFound   = errors.New("TASK NOT FOUND")
	ErrBadRequest = errors.New("BAD REQUEST")
)

const (
	EventTaskCreated = "task_created"
	EventTaskRead    = "task_read"
	EventTaskList    = "task_list"
)

type Logger interface {
	Log(logger.Entry)
}

type Service struct {
	Repo  repository.TaskRepository
	Log   Logger
	Now   func() time.Time
	IdGen func() string
}

func NewService(repo repository.TaskRepository, log Logger) *Service {
	return &Service{
		Repo: repo,
		Log:  log,
		Now:  func() time.Time { return time.Now().UTC() },
		IdGen: func() string {
			var b [16]byte
			_, _ = rand.Read(b[:])
			return hex.EncodeToString(b[:])
		},
	}
}

func (s *Service) Create(ctx context.Context, reqID string, in dto.CreateInput) (domain.Task, error) {
	if in.Title == "" {
		return domain.Task{}, ErrBadRequest
	}
	if !validation.IsValidStatus(in.Status) {
		in.Status = domain.StatusTodo
	}

	now := s.Now()
	t := mapper.ToDomainTask(in, s.IdGen(), now)
	out, err := s.Repo.Create(ctx, t)
	s.Log.Log(logger.Entry{
		Time:      now,
		Event:     EventTaskCreated,
		RequestID: reqID,
		Data: map[string]any{
			"id":     t.ID,
			"status": t.Status,
		},
		Error: validation.ErrString(err),
	})
	return out, err
}

func (s *Service) Get(ctx context.Context, reqID, id string) (domain.Task, error) {
	t, ok, err := s.Repo.GetByID(ctx, id)
	if err == nil && !ok {
		err = ErrNotFound
	}
	s.Log.Log(logger.Entry{
		Time:      s.Now(),
		Event:     EventTaskRead,
		RequestID: reqID,
		Data:      map[string]any{"id": id},
		Error:     validation.ErrString(err),
	})
	return t, err
}

func (s *Service) List(ctx context.Context, reqID string, status *domain.Status) ([]domain.Task, error) {
	tasks, err := s.Repo.List(ctx, repository.Filter{Status: status})
	s.Log.Log(logger.Entry{
		Time:      s.Now(),
		Event:     EventTaskList,
		RequestID: reqID,
		Data: map[string]any{
			"status": validation.StatusString(status),
			"count":  len(tasks),
		},
		Error: validation.ErrString(err),
	})
	return tasks, err
}
