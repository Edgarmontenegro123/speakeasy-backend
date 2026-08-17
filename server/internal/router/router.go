package router

import (
	"net/http"

	"speakeasy-server/internal/handler"
	"speakeasy-server/internal/repository"
	"speakeasy-server/internal/service"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	healthRepo := repository.NewHealthRepository()
	healthService := service.NewHealthService(healthRepo)
	healthHandler := handler.NewHealthHandler(healthService)

	mux.HandleFunc("GET /api/health", healthHandler.GetHealth)

	return mux
}
