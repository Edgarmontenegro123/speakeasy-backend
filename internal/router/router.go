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

	topicRepo := repository.NewTopicRepository()
	topicService := service.NewTopicService(topicRepo)
	topicHandler := handler.NewTopicHandler(topicService)

	sessionRepo := repository.NewSessionRepository()
	sessionService := service.NewSessionService(sessionRepo, topicRepo)
	sessionHandler := handler.NewSessionHandler(sessionService)

	mux.HandleFunc("GET /api/v1/health", healthHandler.GetHealth)
	mux.HandleFunc("GET /api/v1/topics", topicHandler.ListTopics)
	mux.HandleFunc("POST /api/v1/sessions", sessionHandler.CreateSession)
	mux.HandleFunc("POST /api/v1/sessions/{id}/messages", sessionHandler.PostMessage)

	return middleware.CORS(mux)
}
