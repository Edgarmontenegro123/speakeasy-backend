package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAudioHandler_GetAudio(t *testing.T) {
	h := NewAudioHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/audio/message-1.mp3", nil)
	req.SetPathValue("filename", "message-1.mp3")
	rec := httptest.NewRecorder()

	h.GetAudio(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "audio/mpeg" {
		t.Errorf("expected Content-Type %q, got %q", "audio/mpeg", contentType)
	}

	if rec.Body.Len() == 0 {
		t.Error("expected a non-empty audio body")
	}
}
