package dto

import "taskapi/internal/domain"

type CreateInput struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      domain.Status `json:"status"`
}
