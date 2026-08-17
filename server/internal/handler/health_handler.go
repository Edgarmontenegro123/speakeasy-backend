package handler

import (
	"encoding/json"
	"net/http"

	"speakeasy-server/internal/service"
)

type HealthHandler struct {
	service service.HealthService
}

func NewHealthHandler(service service.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	status := h.service.CheckHealth()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
