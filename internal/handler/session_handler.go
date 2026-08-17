package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/service"
)

type SessionHandler struct {
	service service.SessionService
}

func NewSessionHandler(service service.SessionService) *SessionHandler {
	return &SessionHandler{service: service}
}

type createSessionRequest struct {
	TopicID string `json:"topic_id"`
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TopicID == "" {
		writeJSONError(w, http.StatusBadRequest, "topic_id is required")
		return
	}

	session, err := h.service.CreateSession(req.TopicID)
	if err != nil {
		if errors.Is(err, service.ErrTopicNotFound) {
			writeJSONError(w, http.StatusNotFound, "topic not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

type postMessageRequest struct {
	Content string `json:"content"`
}

func (h *SessionHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")

	var req postMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return
	}

	message, err := h.service.PostMessage(sessionID, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			writeJSONError(w, http.StatusNotFound, "session not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
