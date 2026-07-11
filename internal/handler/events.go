package handler

import (
	"encoding/json"
	"net/http"

	"backend-events-api/internal/domain"
	"backend-events-api/internal/service"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{service: svc}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var event domain.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	h.service.ProcessEvent(r.Context(), event)

	w.WriteHeader(http.StatusOK)
}
