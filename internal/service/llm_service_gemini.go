package service

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"

	"github.com/Edgarmontenegro123/speakeasy-backend/internal/model"
)

// geminiLLMService generates tutor replies using the Gemini API.
type geminiLLMService struct {
	client *genai.Client
	model  string
}

// NewGeminiLLMService creates an LLMService backed by the given Gemini model
// (e.g. "gemini-2.5-flash"). Google periodically retires older model IDs, so
// the model is caller-supplied rather than hardcoded — see GEMINI_MODEL in
// internal/config for how it is configured.
func NewGeminiLLMService(ctx context.Context, apiKey, modelName string) (LLMService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating Gemini client: %w", err)
	}

	return &geminiLLMService{client: client, model: modelName}, nil
}

func (s *geminiLLMService) GenerateReply(topic model.Topic, history []model.Message, userMessage string) (string, error) {
	prompt := buildTutorPrompt(topic, history, userMessage)

	resp, err := s.client.Models.GenerateContent(context.Background(), s.model, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("generating tutor reply: %w", err)
	}

	reply := strings.TrimSpace(resp.Text())
	if reply == "" {
		return "", fmt.Errorf("gemini returned an empty reply")
	}

	return reply, nil
}

// buildTutorPrompt assembles a prompt that grounds the model in the active
// topic, adapts its language to the topic's CEFR level, and includes the
// conversation so far so the tutor can respond in context.
func buildTutorPrompt(topic model.Topic, history []model.Message, userMessage string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "You are a friendly English tutor helping a %s-level student practise the topic %q (%s). ",
		topic.Level, topic.Title, topic.Description)
	b.WriteString(levelGuidance(topic.Level))
	b.WriteString(" Keep your reply short (1-3 sentences), gently correct any mistakes, and end with a question that keeps the student talking.\n\n")

	for _, msg := range history {
		switch msg.Role {
		case model.RoleUser:
			b.WriteString("Student: " + msg.Content + "\n")
		case model.RoleAssistant:
			b.WriteString("Tutor: " + msg.Content + "\n")
		}
	}

	b.WriteString("Student: " + userMessage + "\nTutor:")

	return b.String()
}

func levelGuidance(level model.Level) string {
	switch level {
	case model.LevelA1, model.LevelA2:
		return "Use very simple vocabulary, short sentences, and the present tense where possible."
	case model.LevelB1, model.LevelB2:
		return "Use everyday vocabulary and moderately complex sentences."
	default:
		return "Use natural, idiomatic English and nuanced vocabulary."
	}
}
