package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/service"
)

type TopicHandler struct {
	service service.TopicService
}

func NewTopicHandler(service service.TopicService) *TopicHandler {
	return &TopicHandler{service: service}
}

func (h *TopicHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	topics := h.service.ListTopics()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}
