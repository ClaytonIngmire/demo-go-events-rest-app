package service

import (
	"time"

	"github.com/google/uuid"

	"backend-events-api/internal/domain"
)

type MetricsService struct {
	repo domain.EventRepository
}

func NewMetricsService(repo domain.EventRepository) *MetricsService {
	return &MetricsService{
		repo: repo,
	}
}

func (s *MetricsService) GetEventMetrics(account uuid.UUID) domain.EventMetrics {
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	events := s.repo.GetAllForAccountAftertTime(account, twentyFourHoursAgo)

	return domain.CalculateEventMetrics(events)
}
