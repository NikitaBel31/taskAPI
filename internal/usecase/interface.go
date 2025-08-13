package usecase

import (
	"context"
	"taskapi/internal/domain"
	"taskapi/internal/dto"
)

type TaskService interface {
	Create(ctx context.Context, reqID string, in dto.CreateInput) (domain.Task, error)
	Get(ctx context.Context, reqID, id string) (domain.Task, error)
	List(ctx context.Context, reqID string, status *domain.Status) ([]domain.Task, error)
}
