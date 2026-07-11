package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid" // RFC 4122 UUIDs
)

type EventTime struct {
	time.Time
}

func (et *EventTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Try RFC3339
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		et.Time = t // Assign to the embedded field
		return nil
	}

	// Fallback to ISO8601
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}

	et.Time = t.UTC()
	return nil
}

func (et EventTime) MarshalJSON() ([]byte, error) {
	// et.Format works because of method promotion!
	formatted := "\"" + et.Format(time.RFC3339) + "\""
	return []byte(formatted), nil
}

type EventType string

const (
	EventNewDevice  EventType = "NewDevice"
	EventSignIn     EventType = "SignIn"
	EventCreateItem EventType = "CreateItem"
	EventDeleteItem EventType = "DeleteItem"
	EventViewItem   EventType = "ViewItem"
)

var validEventTypes = map[EventType]bool{
	EventNewDevice:  true,
	EventSignIn:     true,
	EventCreateItem: true,
	EventDeleteItem: true,
	EventViewItem:   true,
}

type Event struct {
	Timestamp EventTime `json:"timestamp"`
	EventType EventType `json:"event_type"`
	Account   uuid.UUID `json:"account"`
	User      uuid.UUID `json:"user"`
}

func NewEvent(timestamp EventTime, eventType EventType, account uuid.UUID, user uuid.UUID) (*Event, error) {
	if timestamp.IsZero() || timestamp.After(time.Now()) {
		return nil, errors.New("The event timestamp is not valid")
	} else if !validEventTypes[eventType] {
		return nil, errors.New("The event type provided is not valid")
	} else if account.Version() != 4 {
		return nil, errors.New("The provided account UUID is not valid")
	} else if user.Version() != 4 {
		return nil, errors.New("The provided user UUID is not valid")
	}

	return &Event{
		Timestamp: timestamp,
		EventType: eventType,
		Account:   account,
		User:      user,
	}, nil
}

func (e *Event) UnmarshalJSON(b []byte) error {
	type Alias Event // an alias is used to prevent infinite loops

	var tmp Alias

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	res, err := NewEvent(tmp.Timestamp, tmp.EventType, tmp.Account, tmp.User)
	if err != nil {
		return err
	}

	*e = *res
	return nil
}

type EventRepository interface {
	Save(Event)
	GetAllForAccountAftertTime(uuid.UUID, time.Time) []Event
}
