package service

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const transcriptionInstruction = "Transcribe the following audio recording exactly as spoken, in English. Return only the transcript text, with no extra commentary or punctuation added."

// geminiSTTService transcribes audio using Gemini's native audio
// understanding: the recording is sent to the model as inline data alongside
// a transcription instruction, rather than using a dedicated speech API.
type geminiSTTService struct {
	client *genai.Client
	model  string
}

func NewGeminiSTTService(ctx context.Context, apiKey, modelName string) (STTService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	return &geminiSTTService{client: client, model: modelName}, nil
}

func (s *geminiSTTService) Transcribe(audio []byte, mimeType string) (string, error) {
	content := buildTranscriptionContent(audio, mimeType)

	resp, err := s.client.Models.GenerateContent(context.Background(), s.model, []*genai.Content{content}, nil)
	if err != nil {
		return "", fmt.Errorf("transcribing audio: %w", err)
	}

	transcript := strings.TrimSpace(resp.Text())
	if transcript == "" {
		return "", fmt.Errorf("gemini returned an empty transcript")
	}

	return transcript, nil
}

func buildTranscriptionContent(audio []byte, mimeType string) *genai.Content {
	return genai.NewContentFromParts([]*genai.Part{
		genai.NewPartFromText(transcriptionInstruction),
		genai.NewPartFromBytes(audio, mimeType),
	}, genai.RoleUser)
}
