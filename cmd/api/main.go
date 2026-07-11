package main

import (
	"log"
	"net/http"
	"time"

	"backend-events-api/internal/handler"
	"backend-events-api/internal/repository"
	"backend-events-api/internal/server"
	"backend-events-api/internal/service"
)

func main() {
	eventMemoryRepo := repository.NewEventMemory()

	eventService := service.NewEventService(eventMemoryRepo)
	metricsService := service.NewMetricsService(eventMemoryRepo)

	eventHandler := handler.NewEventHandler(eventService)
	metricsHandler := handler.NewMetricsHandler(metricsService)

	router := server.NewRouter(eventHandler, metricsHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
