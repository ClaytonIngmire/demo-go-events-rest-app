package service

import (
	"context"

	"backend-events-api/internal/domain"
)

type EventService struct {
	repo domain.EventRepository
}

func NewEventService(repo domain.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) ProcessEvent(ctx context.Context, event domain.Event) {
	s.repo.Save(event)
}
