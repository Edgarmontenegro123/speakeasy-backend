package service

import "testing"

func TestSTTService_Transcribe(t *testing.T) {
	s := NewSTTService()

	transcript := s.Transcribe([]byte("fake audio bytes"))

	want := "Hello, I want to practice my English today"
	if transcript != want {
		t.Errorf("expected transcript %q, got %q", want, transcript)
	}
}
