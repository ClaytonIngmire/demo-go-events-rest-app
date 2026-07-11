package server

import (
	"net/http"

	"backend-events-api/internal/handler"
)

func NewRouter(
	eventHandler *handler.EventHandler,
	metricsHandler *handler.MetricsHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// Assuming Go 1.22+ for method-based routing
	mux.Handle("POST /events", eventHandler)
	mux.Handle("GET /metrics", metricsHandler)

	return mux
}
