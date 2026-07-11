package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"backend-events-api/internal/service"
)

type MetricsHandler struct {
	service *service.MetricsService
}

func NewMetricsHandler(svc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{service: svc}
}

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accountStr := r.URL.Query().Get("account")
	if accountStr == "" {
		http.Error(w, "Missing 'account' query parameter", http.StatusBadRequest)
		return
	}

	account, err := uuid.Parse(accountStr)
	if err != nil || account.Version() != 4 {
		http.Error(w, "Invalid 'account' UUID format", http.StatusBadRequest)
		return
	}

	metrics := h.service.GetEventMetrics(account)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
