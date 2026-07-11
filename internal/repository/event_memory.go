package repository

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"backend-events-api/internal/domain"
)

type EventMemory struct {
	mu     sync.RWMutex
	events map[uuid.UUID][]domain.Event
}

func NewEventMemory() *EventMemory {
	return &EventMemory{
		events: make(map[uuid.UUID][]domain.Event),
	}
}

func (r *EventMemory) Save(e domain.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	account := e.Account

	r.events[account] = append(r.events[account], e)
}

func (r *EventMemory) GetAllForAccountAftertTime(account uuid.UUID, since time.Time) []domain.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []domain.Event
	for _, e := range r.events[account] {
		if e.Timestamp.After(since) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
