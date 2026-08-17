package router

import (
	"net/http"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/handler"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", handler.Health)

	return middleware.CORS(mux)
}
