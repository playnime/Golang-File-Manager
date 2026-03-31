package httpServer

import "lab2pepe/server/internal/core/domain"

type HTTPHandlers struct {
	repo *domain.Repository
}

func NewHTTPHandlers(repo *domain.Repository) *HTTPHandlers {
	return &HTTPHandlers{
		repo: repo,
	}
}
