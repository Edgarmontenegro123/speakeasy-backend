package service

import (
	"bytes"
	"strings"
	"testing"
)

func TestBuildTranscriptionContent(t *testing.T) {
	audio := []byte("fake wav bytes")

	content := buildTranscriptionContent(audio, "audio/wav")

	if content.Role != "user" {
		t.Errorf("expected role %q, got %q", "user", content.Role)
	}

	if len(content.Parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(content.Parts))
	}

	if !strings.Contains(content.Parts[0].Text, "Transcribe") {
		t.Errorf("expected the first part to be a transcription instruction, got %q", content.Parts[0].Text)
	}

	inlineData := content.Parts[1].InlineData
	if inlineData == nil {
		t.Fatal("expected the second part to carry inline audio data")
	}

	if inlineData.MIMEType != "audio/wav" {
		t.Errorf("expected mime type %q, got %q", "audio/wav", inlineData.MIMEType)
	}

	if !bytes.Equal(inlineData.Data, audio) {
		t.Errorf("expected inline data to match the input audio bytes")
	}
}
