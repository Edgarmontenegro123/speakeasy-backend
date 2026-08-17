package router

import (
	"net/http"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/handler"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/middleware"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/service"
)

func New() http.Handler {
	mux := http.NewServeMux()

	healthRepo := repository.NewHealthRepository()
	healthService := service.NewHealthService(healthRepo)
	healthHandler := handler.NewHealthHandler(healthService)

	mux.HandleFunc("GET /api/v1/health", healthHandler.GetHealth)

	return middleware.CORS(mux)
}
