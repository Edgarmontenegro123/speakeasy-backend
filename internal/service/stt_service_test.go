package service

import "testing"

func TestSTTService_Transcribe(t *testing.T) {
	s := NewSTTService()

	transcript, err := s.Transcribe([]byte("fake audio bytes"), "audio/wav")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello, I want to practice my English today"
	if transcript != want {
		t.Errorf("expected transcript %q, got %q", want, transcript)
	}
}
