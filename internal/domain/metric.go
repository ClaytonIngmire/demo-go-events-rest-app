package domain

type EventMetrics struct {
	TotalEvents  int            `json:"total_events"`
	EventsByType map[string]int `json:"events_by_type"`
	UniqueUsers  int            `json:"unique_users"`
}

func NewEventsByType() map[string]int {
	eventsByType := make(map[string]int, len(validEventTypes))

	for et := range validEventTypes {
		eventsByType[string(et)] = 0
	}

	return eventsByType
}

func CalculateEventMetrics(events []Event) EventMetrics {
	metrics := EventMetrics{
		TotalEvents:  len(events),
		EventsByType: NewEventsByType(),
	}

	uniqueUsers := make(map[string]struct{})

	for _, e := range events {
		metrics.EventsByType[string(e.EventType)]++
		uniqueUsers[e.User.String()] = struct{}{}
	}

	metrics.UniqueUsers = len(uniqueUsers)
	return metrics
}
