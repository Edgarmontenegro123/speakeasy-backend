package service

import "testing"

func TestTTSService_GenerateAudioURL(t *testing.T) {
	s := NewTTSService()

	url := s.GenerateAudioURL("message-1")

	want := "/api/v1/audio/message-1.mp3"
	if url != want {
		t.Errorf("expected audio URL %q, got %q", want, url)
	}
}
