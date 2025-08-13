package repository

import (
	"context"

	"taskapi/internal/domain"
)

type Filter struct {
	Status *domain.Status
}

// Так как таска маленькая и копирование дешевое, то передаю ее по значению
type TaskRepository interface {
	Create(ctx context.Context, t domain.Task) (domain.Task, error)
	GetByID(ctx context.Context, id string) (domain.Task, bool, error)
	List(ctx context.Context, f Filter) ([]domain.Task, error)
}
