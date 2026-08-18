package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
	"github.com/Edgarmontenegro123/speakeasy-backend/internal/service"
)

type stubSessionService struct {
	session    model.Session
	message    model.Message
	err        error
	gotContent *string
}

func (s stubSessionService) CreateSession(topicID string) (model.Session, error) {
	return s.session, s.err
}

func (s stubSessionService) PostMessage(sessionID, content string) (model.Message, error) {
	if s.gotContent != nil {
		*s.gotContent = content
	}
	return s.message, s.err
}

type stubSTTService struct {
	transcript string
	err        error
	gotMIME    *string
}

func (s stubSTTService) Transcribe(audio []byte, mimeType string) (string, error) {
	if s.gotMIME != nil {
		*s.gotMIME = mimeType
	}
	if s.err != nil {
		return "", s.err
	}
	return s.transcript, nil
}

func TestSessionHandler_CreateSession(t *testing.T) {
	want := model.Session{
		ID:        "session-1",
		TopicID:   "topic-ordering-coffee",
		Status:    model.SessionStatusActive,
		CreatedAt: time.Now().UTC(),
		Messages:  []model.Message{},
	}
	h := NewSessionHandler(stubSessionService{session: want}, stubSTTService{})

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
	h := NewSessionHandler(stubSessionService{}, stubSTTService{})

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateSession(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSessionHandler_CreateSession_TopicNotFound(t *testing.T) {
	h := NewSessionHandler(stubSessionService{err: service.ErrTopicNotFound}, stubSTTService{})

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
	h := NewSessionHandler(stubSessionService{message: want}, stubSTTService{})

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
	h := NewSessionHandler(stubSessionService{}, stubSTTService{})

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
	h := NewSessionHandler(stubSessionService{err: service.ErrSessionNotFound}, stubSTTService{})

	body, _ := json.Marshal(map[string]string{"content": "Hello"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/unknown/messages", bytes.NewReader(body))
	req.SetPathValue("id", "unknown")
	rec := httptest.NewRecorder()

	h.PostMessage(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func newMultipartAudioRequest(t *testing.T, audio []byte, mimeType string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", `form-data; name="audio"; filename="audio.rec"`)
	partHeader.Set("Content-Type", mimeType)

	part, err := writer.CreatePart(partHeader)
	if err != nil {
		t.Fatalf("failed to create form part: %v", err)
	}
	if _, err := part.Write(audio); err != nil {
		t.Fatalf("failed to write audio bytes: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/session-1/messages", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetPathValue("id", "session-1")

	return httptest.NewRecorder(), req
}

func TestSessionHandler_PostMessage_MultipartAudio(t *testing.T) {
	want := model.Message{
		ID:        "message-1",
		SessionID: "session-1",
		Sender:    model.SenderAI,
		Content:   "Great start! Can you tell me more about that?",
		AudioURL:  "/api/v1/audio/message-1.mp3",
		CreatedAt: time.Now().UTC(),
	}
	transcript := "Hello, I want to practice my English today"
	var gotContent string
	var gotMIME string

	h := NewSessionHandler(
		stubSessionService{message: want, gotContent: &gotContent},
		stubSTTService{transcript: transcript, gotMIME: &gotMIME},
	)

	rec, req := newMultipartAudioRequest(t, []byte("fake audio bytes"), "audio/wav")

	h.PostMessage(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	if gotMIME != "audio/wav" {
		t.Errorf("expected the STT service to receive mime type %q, got %q", "audio/wav", gotMIME)
	}

	if gotContent != transcript {
		t.Errorf("expected transcribed content %q, got %q", transcript, gotContent)
	}

	var got model.Message
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.AudioURL == "" {
		t.Error("expected a non-empty audio_url in the tutor's response")
	}
}

func TestSessionHandler_PostMessage_MultipartAudio_Webm(t *testing.T) {
	want := model.Message{ID: "message-1", Sender: model.SenderAI, AudioURL: "/api/v1/audio/message-1.mp3"}
	h := NewSessionHandler(stubSessionService{message: want}, stubSTTService{transcript: "Hi there"})

	rec, req := newMultipartAudioRequest(t, []byte("fake webm bytes"), "audio/webm;codecs=opus")

	h.PostMessage(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestSessionHandler_PostMessage_MultipartAudio_Ogg(t *testing.T) {
	want := model.Message{ID: "message-1", Sender: model.SenderAI, AudioURL: "/api/v1/audio/message-1.mp3"}
	h := NewSessionHandler(stubSessionService{message: want}, stubSTTService{transcript: "Hi there"})

	rec, req := newMultipartAudioRequest(t, []byte("fake ogg bytes"), "audio/ogg")

	h.PostMessage(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestSessionHandler_PostMessage_MultipartUnsupportedFormat(t *testing.T) {
	h := NewSessionHandler(stubSessionService{}, stubSTTService{transcript: "Hi there"})

	rec, req := newMultipartAudioRequest(t, []byte("fake bytes"), "video/mp4")

	h.PostMessage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestSessionHandler_PostMessage_TranscriptionError(t *testing.T) {
	h := NewSessionHandler(stubSessionService{}, stubSTTService{err: errors.New("gemini unavailable")})

	rec, req := newMultipartAudioRequest(t, []byte("fake audio bytes"), "audio/wav")

	h.PostMessage(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestSessionHandler_PostMessage_MultipartMissingFile(t *testing.T) {
	h := NewSessionHandler(stubSessionService{}, stubSTTService{})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("topic", "no audio here"); err != nil {
		t.Fatalf("failed to write field: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/session-1/messages", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetPathValue("id", "session-1")
	rec := httptest.NewRecorder()

	h.PostMessage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
