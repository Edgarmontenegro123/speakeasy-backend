package router

import (
	"context"
	"log"
	"net/http"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/config"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/handler"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/middleware"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/repository"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/service"
)

func New(cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	healthRepo := repository.NewHealthRepository()
	healthService := service.NewHealthService(healthRepo)
	healthHandler := handler.NewHealthHandler(healthService)

	topicRepo := repository.NewTopicRepository()
	topicService := service.NewTopicService(topicRepo)
	topicHandler := handler.NewTopicHandler(topicService)

	sessionRepo := repository.NewSessionRepository()
	ttsService := service.NewTTSService()
	sttService := service.NewSTTService()
	llmService := newLLMService(cfg)
	sessionService := service.NewSessionService(sessionRepo, topicRepo, ttsService, llmService)
	sessionHandler := handler.NewSessionHandler(sessionService, sttService)

	audioHandler := handler.NewAudioHandler()

	mux.HandleFunc("GET /api/v1/health", healthHandler.GetHealth)
	mux.HandleFunc("GET /api/v1/topics", topicHandler.ListTopics)
	mux.HandleFunc("POST /api/v1/sessions", sessionHandler.CreateSession)
	mux.HandleFunc("POST /api/v1/sessions/{id}/messages", sessionHandler.PostMessage)
	mux.HandleFunc("GET /api/v1/audio/{filename}", audioHandler.GetAudio)

	return middleware.CORS(mux)
}

// newLLMService uses Gemini when an API key is configured, falling back to
// a stub tutor reply otherwise (e.g. local development without a key).
func newLLMService(cfg config.Config) service.LLMService {
	if cfg.GeminiAPIKey == "" {
		log.Println("GEMINI_API_KEY not set, using stub tutor replies")
		return service.NewStubLLMService()
	}

	llmService, err := service.NewGeminiLLMService(context.Background(), cfg.GeminiAPIKey, cfg.GeminiModel)
	if err != nil {
		log.Printf("failed to initialise Gemini LLM service, falling back to stub replies: %v", err)
		return service.NewStubLLMService()
	}

	log.Printf("using Gemini model %q for tutor replies", cfg.GeminiModel)
	return llmService
}
