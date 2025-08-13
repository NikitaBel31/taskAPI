package mapper

import (
	"taskapi/internal/domain"
	"taskapi/internal/dto"
	"time"
)

func ToDomainTask(in dto.CreateInput, id string, now time.Time) domain.Task {
	return domain.Task{
		ID:          id,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
