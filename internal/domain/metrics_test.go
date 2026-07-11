package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewEventsByType(t *testing.T) {
	m := NewEventsByType()
	if len(m) != len(validEventTypes) {
		t.Fatalf("expected %d keys, got %d", len(validEventTypes), len(m))
	}

	for et := range validEventTypes {
		k := string(et)
		v, ok := m[k]
		if !ok {
			t.Errorf("missing key %s", k)
		}
		if v != 0 {
			t.Errorf("expected zero for key %s, got %d", k, v)
		}
	}
}

func TestCalculateEventMetrics_Basic(t *testing.T) {
	ts := EventTime{Time: time.Now().Add(-time.Minute)}
	acc := uuid.New()
	u1 := uuid.New()
	u2 := uuid.New()

	e1, err := NewEvent(ts, EventNewDevice, acc, u1)
	if err != nil {
		t.Fatalf("unexpected error creating event1: %v", err)
	}
	e2, err := NewEvent(ts, EventSignIn, acc, u2)
	if err != nil {
		t.Fatalf("unexpected error creating event2: %v", err)
	}
	e3, err := NewEvent(ts, EventSignIn, acc, u1)
	if err != nil {
		t.Fatalf("unexpected error creating event3: %v", err)
	}

	events := []Event{*e1, *e2, *e3}
	metrics := CalculateEventMetrics(events)

	if metrics.TotalEvents != 3 {
		t.Fatalf("expected TotalEvents 3, got %d", metrics.TotalEvents)
	}

	expected := NewEventsByType()
	expected[string(EventNewDevice)] = 1
	expected[string(EventSignIn)] = 2

	if !reflect.DeepEqual(metrics.EventsByType, expected) {
		t.Fatalf("events by type mismatch. expected %v, got %v", expected, metrics.EventsByType)
	}

	if metrics.UniqueUsers != 2 {
		t.Fatalf("expected UniqueUsers 2, got %d", metrics.UniqueUsers)
	}
}

func TestCalculateEventMetrics_Empty(t *testing.T) {
	metrics := CalculateEventMetrics([]Event{})

	if metrics.TotalEvents != 0 {
		t.Fatalf("expected TotalEvents 0, got %d", metrics.TotalEvents)
	}

	for k, v := range metrics.EventsByType {
		if v != 0 {
			t.Fatalf("expected zero count for %s, got %d", k, v)
		}
	}

	if metrics.UniqueUsers != 0 {
		t.Fatalf("expected UniqueUsers 0, got %d", metrics.UniqueUsers)
	}
}
