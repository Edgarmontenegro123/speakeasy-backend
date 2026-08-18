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
// file (field "audio", audio/webm, audio/ogg or audio/wav), which is run
// through the STT service to obtain a transcription. It writes a 400
// response and returns ok=false on failure.
func (h *SessionHandler) extractContent(w http.ResponseWriter, r *http.Request) (string, bool) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return h.extractContentFromAudio(w, r)
	}

	return h.extractContentFromJSON(w, r)
}

func (h *SessionHandler) extractContentFromJSON(w http.ResponseWriter, r *http.Request) (string, bool) {
	var req postMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		log.Printf("post message 400: invalid or empty JSON content (content-type=%q): %v", r.Header.Get("Content-Type"), err)
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return "", false
	}

	return req.Content, true
}

func (h *SessionHandler) extractContentFromAudio(w http.ResponseWriter, r *http.Request) (string, bool) {
	if err := r.ParseMultipartForm(maxAudioUploadSize); err != nil {
		log.Printf("post message 400: failed to parse multipart form (content-type=%q): %v", r.Header.Get("Content-Type"), err)
		writeJSONError(w, http.StatusBadRequest, "invalid multipart form data")
		return "", false
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		log.Printf("post message 400: missing or invalid %q form field: %v", "audio", err)
		writeJSONError(w, http.StatusBadRequest, "an audio file is required")
		return "", false
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !isSupportedAudioMIME(mimeType) {
		log.Printf("post message 400: unsupported audio content-type %q (filename=%q)", mimeType, header.Filename)
		writeJSONError(w, http.StatusBadRequest, "unsupported audio format, expected audio/webm, audio/ogg, or audio/wav")
		return "", false
	}

	audio, err := io.ReadAll(file)
	if err != nil {
		log.Printf("post message 400: failed to read audio file (content-type=%q): %v", mimeType, err)
		writeJSONError(w, http.StatusBadRequest, "failed to read audio file")
		return "", false
	}

	content, err := h.stt.Transcribe(audio, mimeType)
	if err != nil {
		log.Printf("failed to transcribe audio: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return "", false
	}
	if content == "" {
		log.Printf("post message 400: transcription returned empty content (content-type=%q)", mimeType)
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return "", false
	}

	return content, true
}

// isSupportedAudioMIME reports whether mimeType is an accepted audio upload
// format, tolerating parameters such as "audio/webm;codecs=opus".
func isSupportedAudioMIME(mimeType string) bool {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	return strings.HasPrefix(mimeType, "audio/webm") ||
		strings.HasPrefix(mimeType, "audio/wav") ||
		strings.HasPrefix(mimeType, "audio/ogg")
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
