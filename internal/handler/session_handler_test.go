package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/service"
)

type stubSessionService struct {
	session model.Session
	message model.Message
	err     error
}

func (s stubSessionService) CreateSession(topicID string) (model.Session, error) {
	return s.session, s.err
}

func (s stubSessionService) PostMessage(sessionID, content string) (model.Message, error) {
	return s.message, s.err
}

func TestSessionHandler_CreateSession(t *testing.T) {
	want := model.Session{
		ID:        "session-1",
		TopicID:   "topic-ordering-coffee",
		Status:    model.SessionStatusActive,
		CreatedAt: time.Now().UTC(),
		Messages:  []model.Message{},
	}
	h := NewSessionHandler(stubSessionService{session: want})

	body, _ := json.Marshal(map[string]string{"topic_id": "topic-ordering-coffee"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateSession(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var got model.Session
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.ID != want.ID || got.TopicID != want.TopicID {
		t.Errorf("expected session %+v, got %+v", want, got)
	}
}

func TestSessionHandler_CreateSession_MissingTopicID(t *testing.T) {
	h := NewSessionHandler(stubSessionService{})

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateSession(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSessionHandler_CreateSession_TopicNotFound(t *testing.T) {
	h := NewSessionHandler(stubSessionService{err: service.ErrTopicNotFound})

	body, _ := json.Marshal(map[string]string{"topic_id": "unknown-topic"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateSession(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestSessionHandler_PostMessage(t *testing.T) {
	want := model.Message{
		ID:        "message-1",
		SessionID: "session-1",
		Sender:    model.SenderAI,
		Content:   "Great start! Can you tell me more about that?",
		AudioURL:  "/api/v1/audio/message-1.mp3",
		CreatedAt: time.Now().UTC(),
	}
	h := NewSessionHandler(stubSessionService{message: want})

	body, _ := json.Marshal(map[string]string{"content": "Hello, I want to practise English."})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/session-1/messages", bytes.NewReader(body))
	req.SetPathValue("id", "session-1")
	rec := httptest.NewRecorder()

	h.PostMessage(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var got model.Message
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Sender != model.SenderAI || got.Content != want.Content {
		t.Errorf("expected message %+v, got %+v", want, got)
	}

	if got.AudioURL == "" {
		t.Error("expected a non-empty audio_url in the tutor's response")
	}
}

func TestSessionHandler_PostMessage_EmptyContent(t *testing.T) {
	h := NewSessionHandler(stubSessionService{})

	body, _ := json.Marshal(map[string]string{"content": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/session-1/messages", bytes.NewReader(body))
	req.SetPathValue("id", "session-1")
	rec := httptest.NewRecorder()

	h.PostMessage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSessionHandler_PostMessage_SessionNotFound(t *testing.T) {
	h := NewSessionHandler(stubSessionService{err: service.ErrSessionNotFound})

	body, _ := json.Marshal(map[string]string{"content": "Hello"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/unknown/messages", bytes.NewReader(body))
	req.SetPathValue("id", "unknown")
	rec := httptest.NewRecorder()

	h.PostMessage(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
