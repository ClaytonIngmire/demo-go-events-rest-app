package repository_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"backend-events-api/internal/domain"
	"backend-events-api/internal/repository"
)

func TestEventMemorySaveAndGetAllForAccountAfterTime(t *testing.T) {
	repo := repository.NewEventMemory()

	accountOne := uuid.New()
	accountTwo := uuid.New()
	user := uuid.New()

	baseTime := time.Now().Add(-2 * time.Hour)
	olderEvent := newTestEvent(t, baseTime.Add(-30*time.Minute), accountOne, user)
	newerEvent := newTestEvent(t, baseTime.Add(30*time.Minute), accountOne, user)
	otherAccountEvent := newTestEvent(t, baseTime.Add(45*time.Minute), accountTwo, user)

	repo.Save(olderEvent)
	repo.Save(newerEvent)
	repo.Save(otherAccountEvent)

	events := repo.GetAllForAccountAftertTime(accountOne, baseTime)

	if len(events) != 1 {
		t.Fatalf("expected 1 event for account after time, got %d", len(events))
	}

	if events[0].Timestamp != newerEvent.Timestamp {
		t.Fatalf("expected the newer event to be returned, got timestamp %v", events[0].Timestamp)
	}

	if events[0].Account != accountOne {
		t.Fatalf("expected the event to belong to account %s, got %s", accountOne, events[0].Account)
	}
}

func TestEventMemoryReturnsNilSliceForUnknownAccount(t *testing.T) {
	repo := repository.NewEventMemory()
	account := uuid.New()
	user := uuid.New()

	event := newTestEvent(t, time.Now().Add(-10*time.Minute), account, user)
	repo.Save(event)

	events := repo.GetAllForAccountAftertTime(uuid.New(), time.Now().Add(-1*time.Hour))

	if events != nil {
		t.Fatal("expected a nil slice for an unknown account")
	}

	if len(events) != 0 {
		t.Fatalf("expected no events for an unknown account, got %d", len(events))
	}
}

func newTestEvent(t *testing.T, timestamp time.Time, account uuid.UUID, user uuid.UUID) domain.Event {
	t.Helper()

	event, err := domain.NewEvent(domain.EventTime{Time: timestamp}, domain.EventNewDevice, account, user)
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	return *event
}
