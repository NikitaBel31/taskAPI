package memory

import (
	"context"
	"sync"

	"taskapi/internal/domain"
	"taskapi/internal/repository"
)

type Repo struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

func New() *Repo {
	return &Repo{
		tasks: make(map[string]domain.Task),
	}
}

func (r *Repo) Create(ctx context.Context, t domain.Task) (domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID] = t
	return t, nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (domain.Task, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	return t, ok, nil
}

func (r *Repo) List(ctx context.Context, f repository.Filter) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		if f.Status != nil && t.Status != *f.Status {
			continue
		}
		out = append(out, t)
	}
	return out, nil
}
