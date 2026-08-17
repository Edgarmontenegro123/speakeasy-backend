package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/service"
)

const maxAudioUploadSize = 32 << 20 // 32 MiB

type SessionHandler struct {
	service service.SessionService
	stt     service.STTService
}

func NewSessionHandler(sessionService service.SessionService, sttService service.STTService) *SessionHandler {
	return &SessionHandler{service: sessionService, stt: sttService}
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

	content, ok := h.extractContent(w, r)
	if !ok {
		return
	}

	message, err := h.service.PostMessage(sessionID, content)
	if err != nil {
		if errors.Is(err, service.ErrSessionNotFound) {
			writeJSONError(w, http.StatusNotFound, "session not found")
			return
		}
		log.Printf("failed to post message for session %q: %v", sessionID, err)
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// extractContent resolves the message content from either a JSON body
// ({"content": "..."}) or a multipart/form-data upload carrying an audio
// file (field "file"), which is run through the STT service to obtain a
// transcription. It writes a 400 response and returns ok=false on failure.
func (h *SessionHandler) extractContent(w http.ResponseWriter, r *http.Request) (string, bool) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return h.extractContentFromAudio(w, r)
	}

	return h.extractContentFromJSON(w, r)
}

func (h *SessionHandler) extractContentFromJSON(w http.ResponseWriter, r *http.Request) (string, bool) {
	var req postMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return "", false
	}

	return req.Content, true
}

func (h *SessionHandler) extractContentFromAudio(w http.ResponseWriter, r *http.Request) (string, bool) {
	if err := r.ParseMultipartForm(maxAudioUploadSize); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid multipart form data")
		return "", false
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "an audio file is required")
		return "", false
	}
	defer file.Close()

	audio, err := io.ReadAll(file)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to read audio file")
		return "", false
	}

	content := h.stt.Transcribe(audio)
	if content == "" {
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return "", false
	}

	return content, true
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
